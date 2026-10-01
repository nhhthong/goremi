//go:build integration

// Integration tests: run the real yt-dlp over the network (build tag integration).
package provider

import (
	"strings"
	"testing"
)

// START: TestSearchReal

func TestSearchReal(t *testing.T) {
	got, err := (&YouTubeProvider{}).Search("daft punk", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) < 1 || len(got) > 10 {
		t.Fatalf("want 1 to 10 tracks, got %d", len(got))
	}
	for i, tr := range got {
		if tr.ID == "" || tr.Title == "" {
			t.Errorf("track %d = %+v, want ID and Title", i, tr)
		}
	}
}

// END: TestSearchReal

// START: TestResolveReal

func TestResolveReal(t *testing.T) {
	got, err := (&YouTubeProvider{}).Resolve(Track{ID: "dQw4w9WgXcQ"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, "http") {
		t.Fatalf("got %q, want an http URL", got)
	}
}

// END: TestResolveReal

// START: TestDetailsReal

func TestDetailsReal(t *testing.T) {
	got, err := (&YouTubeProvider{}).Details(Track{ID: "dQw4w9WgXcQ"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Artist == "" || got.Duration <= 0 || !strings.HasPrefix(got.Artwork, "http") {
		t.Fatalf("track = %+v, want Artist, Duration and an http Artwork", got)
	}
}

// END: TestDetailsReal
