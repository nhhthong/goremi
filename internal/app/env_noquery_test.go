// Tests that the app asks the terminal nothing at start: terminal image support is dropped with the artwork.
package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
)

// START: TestEnvMsgSendsNoDeviceAttributesRequest

func TestEnvMsgSendsNoDeviceAttributesRequest(t *testing.T) {
	_, cmd := New(fakeProvider{}).Update(tea.EnvMsg(uv.Environ{"TERM=xterm-256color", "KITTY_WINDOW_ID=1"}))
	for _, msg := range cmdMsgs(cmd) {
		if raw, ok := msg.(tea.RawMsg); ok && raw.Msg == ansi.RequestPrimaryDeviceAttributes {
			t.Fatal("the environment message made the app ask the terminal for its primary device attributes")
		}
	}
}

// END: TestEnvMsgSendsNoDeviceAttributesRequest

// START: TestAttributesReplyIsIgnored

func TestAttributesReplyIsIgnored(t *testing.T) {
	m := New(fakeProvider{})
	before := m.View().Content
	next, cmd := m.Update(uv.PrimaryDeviceAttributesEvent{62, 4})
	if cmd != nil {
		t.Fatal("the reply made the app start a command")
	}
	if got := next.(Model).View().Content; got != before {
		t.Fatalf("the reply changed the view:\n%q\nwant\n%q", got, before)
	}
}

// END: TestAttributesReplyIsIgnored
