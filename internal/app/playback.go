// Playing a track: resolve its stream URL, hand the URL to the player, and show the artist line.
package app

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"goremi/internal/player"
	"goremi/internal/provider"
	"goremi/internal/ui"
)

// START: Player interface

// Player plays a stream URL and answers the controls; the mpv player and the test fake both fit it.
type Player interface {
	Play(url string) error
	TogglePause() error
	Seek(seconds float64) error
	Events() <-chan player.Event
	Position() time.Duration
	Duration() time.Duration
	Paused() bool
}

// seekStep is the seconds one seek key moves, as the spec decides.
const seekStep = 10

// WithPlayer returns a copy that plays through p.
func (m Model) WithPlayer(p Player) Model {
	m.player, m.plays = p, &playQueue{}
	return m
}

// END: Player interface

// START: playTrack

// loadingArtist is the artist line while the details of the track load.
const loadingArtist = "Loading…"

// playQueue runs Play calls one at a time, in the order they were queued, off the screen loop.
type playQueue struct {
	once sync.Once
	jobs chan func()
}

// submit queues f without waiting and returns the channel its result comes on.
func (q *playQueue) submit(f func() error) <-chan error {
	q.once.Do(func() {
		q.jobs = make(chan func(), 64)
		go func() {
			for job := range q.jobs {
				job()
			}
		}()
	})
	done := make(chan error, 1)
	q.jobs <- func() { done <- f() }
	return done
}

// playedMsg carries the outcome of Play for a track.
type playedMsg struct {
	track provider.Track
	err   error
}

// resolvedMsg carries the stream URL of a track, or why it could not be resolved.
type resolvedMsg struct {
	track provider.Track
	url   string
	err   error
}

// detailsMsg carries the details of the playing track, or why they could not be loaded.
type detailsMsg struct {
	track provider.Track
	err   error
}

// resolveCmd asks the provider for the stream URL of the track.
func resolveCmd(p provider.Provider, track provider.Track) tea.Cmd {
	return func() tea.Msg {
		url, err := p.Resolve(track)
		return resolvedMsg{track: track, url: url, err: err}
	}
}

// detailsCmd asks the provider for the details of the track.
func detailsCmd(p provider.Provider, track provider.Track) tea.Cmd {
	return func() tea.Msg {
		d, err := p.Details(track)
		return detailsMsg{track: d, err: err}
	}
}

// startPlay begins a play request: the artist line reads Loading… until the details arrive.
func (m Model) startPlay(track provider.Track) (Model, tea.Cmd) {
	m.playing, m.artist, m.paused, m.elapsed, m.total = track, loadingArtist, false, 0, 0
	resolve, log := resolveCmd(m.provider, track), m.log
	play := func() tea.Msg { // the play line names the track for the yt-dlp and mpv lines that follow
		if log != nil {
			log("play %s %q", track.ID, track.Title)
		}
		return resolve()
	}
	m.ended = false
	return m, play
}

// frameMsg asks the app to read the band levels, smooth the bars and draw a frame of the spectrum.
type frameMsg struct{}

// frameEvery is the interval between two frames of the spectrum: 30 frames per second, as the spec decides.
const frameEvery = time.Second / 30

// frameCmd waits one frame and then gives a frameMsg.
func frameCmd() tea.Cmd {
	return tea.Tick(frameEvery, func(time.Time) tea.Msg { return frameMsg{} })
}

// spectrumSource is a Player that also reports the band levels of the spectrum.
type spectrumSource interface {
	Bands() [ui.SpectrumBars]float64
	Spectrum() bool
}

// onFrame moves the bars one frame: towards the band levels while audio plays, towards random idle heights when paused or after the last track ended, towards nothing when mpv runs without the spectrum.
// The frame tick goes on while a track is on; with none (a failed play) it ends.
func (m Model) onFrame() (Model, tea.Cmd) {
	if m.artist == "" {
		m.framing = false
		return m, nil
	}
	var target [ui.SpectrumBars]int
	src, ok := m.player.(spectrumSource)
	switch {
	case m.paused || m.ended:
		intn := m.intn
		if intn == nil {
			intn = rand.Intn
		}
		target = ui.IdleHeights(intn)
	case ok && src.Spectrum():
		for i, db := range src.Bands() {
			target[i] = ui.SpectrumHeight(db)
		}
	}
	for i := range m.bars {
		m.bars[i] = ui.SmoothHeight(m.bars[i], target[i])
	}
	return m, frameCmd()
}

// playResolved queues the URL of a resolved track for the player and returns at once: the outcome comes back as a playedMsg. A failed Resolve plays nothing and says so under Search:.
func (m Model) playResolved(r resolvedMsg) (Model, tea.Cmd) {
	if r.err != nil {
		m.notice, m.artist = fmt.Sprintf("Cannot play %q.", r.track.Title), ""
		return m, nil
	}
	m.notice = ""
	if m.player == nil {
		return m, detailsCmd(m.provider, r.track)
	}
	player, url := m.player, r.url
	done := m.plays.submit(func() error { return player.Play(url) })
	return m, func() tea.Msg { return playedMsg{track: r.track, err: <-done} }
}

// played handles the outcome of Play: a failure says so under Search:, a success loads the details; the outcome for a track that is not the playing one any more is dropped.
func (m Model) played(p playedMsg) (Model, tea.Cmd) {
	if p.track.ID != m.playing.ID {
		return m, nil
	}
	if p.err != nil {
		m.notice, m.artist = playMessage(p.err, p.track.Title), ""
		return m, nil
	}
	return m, detailsCmd(m.provider, p.track)
}

// playMessage is the line under Search: for a track mpv could not play: mpv missing, or any other failure.
func playMessage(err error, title string) string {
	if errors.Is(err, exec.ErrNotFound) {
		return "mpv not found. Install mpv and try again."
	}
	return fmt.Sprintf("Cannot play %q.", title)
}

// showDetails puts the artist of the playing track on the artist line; a failed load falls back to the artist the search gave. The first details of a play also start the tick that refreshes the bar and the clock, and the frame tick of the spectrum when it shows.
func (m Model) showDetails(d detailsMsg) (Model, tea.Cmd) {
	if d.track.ID != m.playing.ID {
		return m, nil // a track played since: this answer is stale
	}
	m.artist = d.track.Artist
	if d.err != nil {
		m.artist = m.playing.Artist
	} else if d.track.Duration > 0 {
		m.playing.Duration = d.track.Duration
	}
	var cmds []tea.Cmd
	if !m.ticking {
		m.ticking = true
		cmds = append(cmds, tickCmd())
	}
	if m.spectrum && !m.framing { // the frame tick of the spectrum starts with the first details of a play
		m.framing = true
		cmds = append(cmds, frameCmd())
	}
	return m, tea.Batch(cmds...)
}

// END: playTrack

// START: mpvPlayer

// mpvPlayer is the Player the program uses: it starts mpv at the first Play, so a user who only searches runs no mpv, and it stops mpv in Close.
type mpvPlayer struct {
	mu     sync.Mutex
	p      *player.Player
	ev     chan player.Event // lives from the start, so the app can wait on it before mpv runs
	closed bool
	ctx    context.Context // ends at Close, to interrupt a start in progress
	cancel context.CancelFunc
	log    func(format string, a ...any) // gets the warn and error messages of mpv; nil means none
	// spectrum starts mpv with the band filter (the config key show_spectrum); start is player.StartWith unless a test sets another.
	spectrum bool
	start    func(ctx context.Context, opts player.Options) (*player.Player, error)
}

// errPlayerClosed is what Play returns once Close has run.
var errPlayerClosed = errors.New("player closed")

// newMpvPlayer makes a player that has not started mpv yet.
func newMpvPlayer() *mpvPlayer {
	ctx, cancel := context.WithCancel(context.Background())
	return &mpvPlayer{ev: make(chan player.Event, 16), ctx: ctx, cancel: cancel}
}

// Events reports what mpv tells; nothing comes before the first Play.
func (m *mpvPlayer) Events() <-chan player.Event { return m.ev }

// forward copies the events of mpv to the stable channel; one the queue cannot hold is dropped.
func (m *mpvPlayer) forward(p *player.Player) {
	for e := range p.Events() {
		select {
		case m.ev <- e:
		default:
		}
	}
}

// Play starts mpv when it has not started yet, then plays the URL. The lock is not held while mpv starts, so Close can interrupt the start.
func (m *mpvPlayer) Play(url string) error {
	m.mu.Lock()
	p, closed := m.p, m.closed
	m.mu.Unlock()
	if closed {
		return errPlayerClosed
	}
	if p == nil {
		start := m.start
		if start == nil {
			start = player.StartWith
		}
		started, err := start(m.ctx, player.Options{Log: m.log, Spectrum: m.spectrum})
		if err != nil {
			return err
		}
		m.mu.Lock()
		if m.closed {
			m.mu.Unlock()
			started.Close()
			return errPlayerClosed
		}
		m.p, p = started, started
		m.mu.Unlock()
		go m.forward(started)
	}
	return p.Play(url)
}

// Bands is the level of each band in dB; -60 for all of them before the first Play.
func (m *mpvPlayer) Bands() [ui.SpectrumBars]float64 {
	if p := m.started(); p != nil {
		return p.Bands()
	}
	var out [ui.SpectrumBars]float64
	for i := range out {
		out[i] = -60
	}
	return out
}

// Spectrum tells whether the running mpv has the band filter.
func (m *mpvPlayer) Spectrum() bool {
	p := m.started()
	return p != nil && p.Spectrum()
}

// started returns the running player, or nil before the first Play.
func (m *mpvPlayer) started() *player.Player {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.p
}

// TogglePause does nothing before the first Play.
func (m *mpvPlayer) TogglePause() error {
	if p := m.started(); p != nil {
		return p.TogglePause()
	}
	return nil
}

// Seek does nothing before the first Play.
func (m *mpvPlayer) Seek(seconds float64) error {
	if p := m.started(); p != nil {
		return p.Seek(seconds)
	}
	return nil
}

// Position is 0 before the first Play.
func (m *mpvPlayer) Position() time.Duration {
	if p := m.started(); p != nil {
		return p.Position()
	}
	return 0
}

// Paused is false before the first Play.
func (m *mpvPlayer) Paused() bool {
	if p := m.started(); p != nil {
		return p.Paused()
	}
	return false
}

// Duration is 0 before the first Play.
func (m *mpvPlayer) Duration() time.Duration {
	if p := m.started(); p != nil {
		return p.Duration()
	}
	return 0
}

// Close stops mpv when it started, and interrupts a start in progress; it does not wait for that start to end.
func (m *mpvPlayer) Close() {
	m.mu.Lock()
	m.closed = true
	p := m.p
	m.mu.Unlock()
	m.cancel()
	if p != nil {
		p.Close()
	}
}

// END: mpvPlayer

// START: playback keys

// pauseKey and the seek keys act on the playing track while the list has focus.
func (m Model) pauseKey() Model {
	if m.player != nil && m.player.TogglePause() == nil {
		m.paused = !m.paused
	}
	return m
}

func (m Model) seekKey(seconds float64) {
	if m.player != nil {
		_ = m.player.Seek(seconds)
	}
}

// END: playback keys

// START: player events

// playerEventMsg carries one event of the player.
type playerEventMsg struct{ e player.Event }

// waitEvent waits for the next event of the player; a closed channel ends the waiting.
func waitEvent(ch <-chan player.Event) tea.Cmd {
	if ch == nil {
		return nil
	}
	return func() tea.Msg {
		e, ok := <-ch
		if !ok {
			return nil
		}
		return playerEventMsg{e: e}
	}
}

// onPlayerEvent shows a failed track under Search:, plays the next track when one ended, and waits for the next event.
func (m Model) onPlayerEvent(ev playerEventMsg) (Model, tea.Cmd) {
	if ev.e == player.Failed {
		m.notice, m.artist = fmt.Sprintf("Cannot play %q.", m.playing.Title), ""
	}
	if ev.e == player.Ended {
		next := m.stepTrack(1)
		m.ended = next == nil // the last track of the list ended: the spectrum falls idle
		return m, tea.Batch(waitEvent(m.player.Events()), next)
	}
	return m, waitEvent(m.player.Events())
}

// END: player events

// START: step track

// stepTrack plays the track delta places from the playing one in the results list; at the ends, or when the playing track is not in the list, it does nothing.
func (m Model) stepTrack(delta int) tea.Cmd {
	tracks := m.results.Tracks()
	for i, t := range tracks {
		if t.ID == m.playing.ID && m.playing.ID != "" {
			if j := i + delta; j >= 0 && j < len(tracks) {
				return func() tea.Msg { return PlayMsg{Track: tracks[j]} }
			}
			return nil
		}
	}
	return nil
}

// END: step track

// START: progress

// tickMsg asks the app to read the progress of the player.
type tickMsg struct{}

// progressMsg carries the position and the length the player reported.
type progressMsg struct {
	elapsed, total time.Duration
	paused         bool
}

// readProgress reads the player in a command, so Update never waits for mpv.
func readProgress(p Player) tea.Cmd {
	return func() tea.Msg { return progressMsg{elapsed: p.Position(), total: p.Duration(), paused: p.Paused()} }
}

// onTick starts a reading of the player.
func (m Model) onTick() (Model, tea.Cmd) {
	if m.player == nil {
		return m, nil
	}
	return m, readProgress(m.player)
}

// onProgress keeps what the player reported (a length of 0 leaves the length the track already has; the pause state replaces the app's own flag) and asks for the next tick while a track plays; with nothing playing the tick loop ends.
func (m Model) onProgress(p progressMsg) (Model, tea.Cmd) {
	m.elapsed, m.paused = p.elapsed, p.paused
	if p.total > 0 {
		m.total = p.total
	}
	if m.artist == "" {
		m.ticking = false
		return m, nil
	}
	return m, tickCmd()
}

// tickEvery is the refresh interval of the bar and the clock, as the spec decides.
const tickEvery = time.Second

// tickCmd waits one interval and then gives a tickMsg.
func tickCmd() tea.Cmd {
	return tea.Tick(tickEvery, func(time.Time) tea.Msg { return tickMsg{} })
}

// END: progress

// START: click

// controlAt tells which control of the controls row the cell is on (0 previous, 1 seek back, 2 pause, 3 seek forward, 4 next), or -1. The row is found in the drawn view, so it holds wherever the layout puts it.
func (m Model) controlAt(x, y int) int {
	lines := strings.Split(ansi.Strip(m.View().Content), "\n")
	if y < 0 || y >= len(lines) {
		return -1
	}
	row, text := lines[y], ui.ControlsText(m.paused)
	i := strings.Index(row, text)
	if i < 0 {
		return -1
	}
	rel := x - utf8.RuneCountInString(row[:i])
	runes := []rune(text)
	if rel < 0 || rel >= len(runes) || runes[rel] == ' ' {
		return -1
	}
	return strings.Count(string(runes[:rel]), " ")
}

// click runs the action of the control under a left click.
func (m Model) click(c tea.Mouse) (Model, tea.Cmd) {
	if c.Button != tea.MouseLeft {
		return m, nil
	}
	switch m.controlAt(c.X, c.Y) {
	case 0:
		return m, m.stepTrack(-1)
	case 1:
		m.seekKey(-seekStep)
	case 2:
		m = m.pauseKey()
	case 3:
		m.seekKey(seekStep)
	case 4:
		return m, m.stepTrack(1)
	}
	return m, nil
}

// END: click
