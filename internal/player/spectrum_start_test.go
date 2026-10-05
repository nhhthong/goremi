// Tests of how the player starts mpv with and without the spectrum: the command line, the message level it asks for, and the filter mpv ends up with.
package player

import (
	"context"
	"encoding/json"
	"net"
	"reflect"
	"strings"
	"testing"
)

// START: TestMpvArgsWithSpectrum

func TestMpvArgsWithSpectrum(t *testing.T) {
	t.Setenv("GOREMI_MPV_AO", "")
	args := mpvArgs("/tmp/x.sock", true)
	want := []string{"--no-config", "--idle=yes", "--no-terminal", "--input-ipc-server=/tmp/x.sock", "--af=" + SpectrumFilter(), "--msg-level=all=error,ffmpeg=v"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("args = %q, want %q", args, want)
	}
}

// END: TestMpvArgsWithSpectrum

// START: TestMpvArgsAudioOutput

func TestMpvArgsAudioOutput(t *testing.T) {
	t.Setenv("GOREMI_MPV_AO", "null")
	if args := mpvArgs("/tmp/x.sock", false); !reflect.DeepEqual(args[len(args)-1:], []string{"--ao=null"}) {
		t.Fatalf("args = %q, want --ao=null last", args)
	}
}

// END: TestMpvArgsAudioOutput

// START: TestMpvArgsWithoutSpectrum

func TestMpvArgsWithoutSpectrum(t *testing.T) {
	for _, a := range mpvArgs("/tmp/x.sock", false) {
		if strings.HasPrefix(a, "--af") || strings.HasPrefix(a, "--msg-level") {
			t.Fatalf("argument %q with the spectrum off", a)
		}
	}
}

// END: TestMpvArgsWithoutSpectrum

// START: TestSubscribeLevel

// subscribed runs subscribe on a player whose mpv is a fake that answers every command, and returns the commands it received.
func subscribed(t *testing.T, spectrum bool, log func(string, ...any)) [][]any {
	t.Helper()
	ours, theirs := net.Pipe()
	p := &Player{conn: ours, events: make(chan Event, 16), pending: map[int]chan answer{}, log: log}
	go p.read()
	got := make(chan [][]any, 1)
	go func() {
		dec := json.NewDecoder(theirs)
		var cmds [][]any
		for {
			var req struct {
				Command   []any `json:"command"`
				RequestID int   `json:"request_id"`
			}
			if dec.Decode(&req) != nil {
				break
			}
			cmds = append(cmds, req.Command)
			reply, _ := json.Marshal(map[string]any{"request_id": req.RequestID, "error": "success"})
			theirs.Write(append(reply, '\n'))
			if len(cmds) == 1 {
				break
			}
		}
		got <- cmds
	}()
	if err := p.subscribe(spectrum); err != nil {
		t.Fatal(err)
	}
	cmds := <-got
	ours.Close()
	theirs.Close()
	return cmds
}

func TestSubscribeLevel(t *testing.T) {
	if got := subscribed(t, true, nil); !reflect.DeepEqual(got, [][]any{{"request_log_messages", "v"}}) {
		t.Fatalf("with the spectrum: %v, want request_log_messages v", got)
	}
	if got := subscribed(t, false, func(string, ...any) {}); !reflect.DeepEqual(got, [][]any{{"request_log_messages", "warn"}}) {
		t.Fatalf("with a log only: %v, want request_log_messages warn", got)
	}
}

// END: TestSubscribeLevel

// START: TestStartWithSpectrumSetsFilter

func TestStartWithSpectrumSetsFilter(t *testing.T) {
	p, err := StartWith(context.Background(), Options{Spectrum: true})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	af, err := p.Property("af")
	if err != nil {
		t.Fatal(err)
	}
	if list, ok := af.([]any); !ok || len(list) != 1 {
		t.Fatalf("af = %#v, want a list of one filter", af)
	}
}

// END: TestStartWithSpectrumSetsFilter

// START: TestStartWithoutSpectrumHasNoFilter

func TestStartWithoutSpectrumHasNoFilter(t *testing.T) {
	p, err := StartWith(context.Background(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	af, err := p.Property("af")
	if err != nil {
		t.Fatal(err)
	}
	if list, ok := af.([]any); !ok || len(list) != 0 {
		t.Fatalf("af = %#v, want an empty list", af)
	}
}

// END: TestStartWithoutSpectrumHasNoFilter
