// Tests for the mpv player: it starts, answers, and fails cleanly without mpv.
package player

import (
	"testing"
	"time"
)

// START: TestStartAnswersGetProperty

func TestStartAnswersGetProperty(t *testing.T) {
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer p.Close()
	got, err := p.Property("idle-active")
	if err != nil || got != true {
		t.Fatalf("Property(idle-active) = %v, %v, want true, nil", got, err)
	}
}

// END: TestStartAnswersGetProperty

// START: TestStartStaysIdle

func TestStartStaysIdle(t *testing.T) {
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer p.Close()
	time.Sleep(time.Second)
	got, err := p.Property("idle-active")
	if err != nil || got != true {
		t.Fatalf("Property(idle-active) after 1s = %v, %v, want true, nil", got, err)
	}
}

// END: TestStartStaysIdle

// START: TestStartWithoutMpv

func TestStartWithoutMpv(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	p, err := Start()
	if p != nil || err == nil {
		t.Fatalf("Start() = %v, %v, want nil and an error", p, err)
	}
}

// END: TestStartWithoutMpv
