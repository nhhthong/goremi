// Tests for Search skipping the lines of yt-dlp output that are not JSON.
package provider

import "testing"

// START: TestSearchSkipsBadLine

func TestSearchSkipsBadLine(t *testing.T) {
	f := &fakeRunner{out: `{"id":"a1","title":"One"}` + "\n" + "this is not json\n" + `{"id":"c3","title":"Three"}` + "\n"}
	got, err := (&YouTubeProvider{Runner: f.run}).Search("x", 1)
	if err != nil || len(got) != 2 || got[0].ID != "a1" || got[1].ID != "c3" {
		t.Fatalf("Search = %v, %v; want the tracks a1 and c3 and no error", got, err)
	}
}

// END: TestSearchSkipsBadLine

// START: TestSearchSkipsLeadingWarning

func TestSearchSkipsLeadingWarning(t *testing.T) {
	f := &fakeRunner{out: "WARNING: no formats found\n" + `{"id":"a1","title":"One"}` + "\n" + `{"id":"b2","title":"Two"}` + "\n"}
	got, err := (&YouTubeProvider{Runner: f.run}).Search("x", 1)
	if err != nil || len(got) != 2 {
		t.Fatalf("Search = %v, %v; want 2 tracks and no error", got, err)
	}
}

// END: TestSearchSkipsLeadingWarning
