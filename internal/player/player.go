// Playback through mpv: start the process, talk to it over its IPC endpoint.
package player

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"sync"
	"time"
)

// START: Player

// Event is what the player reports about playback.
type Event string

const (
	Ended  Event = "ended"
	Failed Event = "failed"
	Done   Event = "done"
)

// answer is mpv's reply to one command.
type answer struct {
	data any
	err  string
}

// Player is one running mpv and the connection to its IPC endpoint.
type Player struct {
	cmd     *exec.Cmd
	conn    net.Conn
	cleanup func()
	exited  chan struct{} // closed when the mpv process has ended
	events  chan Event
	mu      sync.Mutex // guards the fields below and serialises writes to conn
	nextID  int
	pending map[int]chan answer
	dead    bool // the connection to mpv is gone
}

// END: Player

// START: Start

// Start runs mpv idle with its own IPC endpoint and connects to it. It returns an error, never panics, when mpv is not on PATH.
func Start() (*Player, error) { return StartContext(context.Background()) }

// StartContext is Start that gives up when ctx ends: it stops its mpv and removes the private directory.
func StartContext(ctx context.Context) (*Player, error) {
	path, err := exec.LookPath("mpv")
	if err != nil {
		return nil, err
	}
	endpoint, cleanup, err := newEndpoint()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(path, "--no-config", "--idle=yes", "--no-terminal", "--input-ipc-server="+endpoint)
	if err := cmd.Start(); err != nil {
		cleanup()
		return nil, err
	}
	exited := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(exited)
	}()
	conn, err := dialRetry(ctx, endpoint, exited)
	if err != nil {
		_ = cmd.Process.Kill()
		<-exited
		cleanup()
		return nil, err
	}
	p := &Player{cmd: cmd, conn: conn, cleanup: cleanup, exited: exited, events: make(chan Event, 16), pending: map[int]chan answer{}}
	go p.read()
	return p, nil
}

// dialRetry waits for mpv to create its endpoint, and gives up at once when mpv has exited or ctx has ended.
func dialRetry(ctx context.Context, endpoint string, exited <-chan struct{}) (net.Conn, error) {
	var err error
	for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); {
		select {
		case <-exited:
			return nil, errors.New("mpv exited before its endpoint opened")
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		var conn net.Conn
		if conn, err = dial(endpoint); err == nil {
			return conn, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
	return nil, fmt.Errorf("mpv endpoint did not open: %w", err)
}

// END: Start

// START: read

// read hands every answer of mpv to the command waiting for it; when the connection ends it reports Done and fails the waiting commands.
func (p *Player) read() {
	r := bufio.NewReader(p.conn)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			break
		}
		var m struct {
			Data      any    `json:"data"`
			Error     string `json:"error"`
			RequestID int    `json:"request_id"`
			Event     string `json:"event"`
			Reason    string `json:"reason"`
		}
		if json.Unmarshal(line, &m) != nil {
			continue
		}
		if m.Event == "end-file" && m.Reason == "eof" {
			p.emit(Ended)
		}
		if m.Event == "end-file" && m.Reason == "error" {
			p.emit(Failed)
		}
		if m.RequestID == 0 {
			continue // an event, not an answer
		}
		p.mu.Lock()
		ch := p.pending[m.RequestID]
		delete(p.pending, m.RequestID)
		p.mu.Unlock()
		if ch != nil {
			ch <- answer{data: m.Data, err: m.Error}
		}
	}
	p.mu.Lock()
	p.dead = true
	for id, ch := range p.pending {
		ch <- answer{err: "mpv exited"}
		delete(p.pending, id)
	}
	p.mu.Unlock()
	p.emit(Done)
	close(p.events)
}

// emit queues an event; one nobody reads and the queue cannot hold is dropped, so the reader never blocks.
func (p *Player) emit(e Event) {
	select {
	case p.events <- e:
	default:
	}
}

// Events reports what happens to playback; Done comes when mpv has exited.
func (p *Player) Events() <-chan Event { return p.events }

// END: read

// START: command

// command sends one mpv command and returns the data of its answer.
func (p *Player) command(args ...any) (any, error) {
	ch := make(chan answer, 1)
	p.mu.Lock()
	if p.dead {
		p.mu.Unlock()
		return nil, errors.New("mpv exited")
	}
	p.nextID++
	id := p.nextID
	p.pending[id] = ch
	req, _ := json.Marshal(map[string]any{"command": args, "request_id": id})
	_, err := p.conn.Write(append(req, '\n'))
	if err != nil {
		delete(p.pending, id)
	}
	p.mu.Unlock()
	if err != nil {
		return nil, err
	}
	a := <-ch
	if a.err != "success" {
		return nil, fmt.Errorf("mpv: %s", a.err)
	}
	return a.data, nil
}

// END: command

// START: Property

// Property asks mpv for the value of a property, such as "idle-active".
func (p *Player) Property(name string) (any, error) { return p.command("get_property", name) }

// END: Property

// START: Play

// Play loads the URL, replacing what plays now.
func (p *Player) Play(url string) error {
	_, err := p.command("loadfile", url, "replace")
	return err
}

// Paused tells whether playback is paused; false when mpv does not answer.
func (p *Player) Paused() bool {
	v, err := p.Property("pause")
	b, _ := v.(bool)
	return err == nil && b
}

// seconds reads a property that mpv gives in seconds; 0 when mpv does not answer.
func (p *Player) seconds(name string) time.Duration {
	v, err := p.Property(name)
	s, _ := v.(float64)
	if err != nil {
		return 0
	}
	return time.Duration(s * float64(time.Second))
}

// Duration is the length of the loaded file; 0 when nothing is loaded.
func (p *Player) Duration() time.Duration { return p.seconds("duration") }

// TogglePause pauses playback, or resumes it when it is paused.
func (p *Player) TogglePause() error {
	_, err := p.command("cycle", "pause")
	return err
}

// Seek moves the position by the given seconds, back when negative.
func (p *Player) Seek(seconds float64) error {
	_, err := p.command("seek", seconds)
	return err
}

// Position is how far playback has come; 0 when nothing plays.
func (p *Player) Position() time.Duration { return p.seconds("time-pos") }

// END: Play

// START: Close

// Close stops mpv and removes the private endpoint directory.
func (p *Player) Close() {
	_ = p.conn.Close()
	if p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
	}
	<-p.exited
	p.cleanup()
}

// END: Close
