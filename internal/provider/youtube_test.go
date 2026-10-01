// Tests for the YouTube provider using a fake yt-dlp runner.
package provider

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

// START: fakeRunner

type fakeRunner struct {
	args []string
	out  string
	err  error
}

// END: fakeRunner

// START: run

func (f *fakeRunner) run(args ...string) ([]byte, error) {
	f.args = args
	return []byte(f.out), f.err
}

// END: run

// START: jsonLines

func jsonLines(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(`{"id":"id` + string(rune('a'+i)) + `","title":"Song ` + string(rune('A'+i)) + `"}` + "\n")
	}
	return b.String()
}

// END: jsonLines

// START: TestSearchArgs

func TestSearchArgs(t *testing.T) {
	f := &fakeRunner{}
	p := &YouTubeProvider{Runner: f.run}
	if _, err := p.Search("daft punk", 1); err != nil {
		t.Fatal(err)
	}
	want := []string{"ytsearch10:daft punk", "--flat-playlist", "-j"}
	if !reflect.DeepEqual(f.args, want) {
		t.Fatalf("args = %q, want %q", f.args, want)
	}
}

// END: TestSearchArgs

// START: TestSearchParse

func TestSearchParse(t *testing.T) {
	f := &fakeRunner{out: `{"id":"a1","title":"Get Lucky"}` + "\n" + `{"id":"b2","title":"One More Time"}` + "\n"}
	got, err := (&YouTubeProvider{Runner: f.run}).Search("x", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ID != "a1" || got[0].Title != "Get Lucky" || got[1].ID != "b2" || got[1].Title != "One More Time" {
		t.Fatalf("tracks = %+v", got)
	}
}

// END: TestSearchParse

// START: TestSearchCapsAtTen

func TestSearchCapsAtTen(t *testing.T) {
	f := &fakeRunner{out: jsonLines(11)}
	got, err := (&YouTubeProvider{Runner: f.run}).Search("x", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 10 {
		t.Fatalf("want 10 tracks, got %d", len(got))
	}
}

// END: TestSearchCapsAtTen

// START: TestSearchEmpty

func TestSearchEmpty(t *testing.T) {
	f := &fakeRunner{}
	got, err := (&YouTubeProvider{Runner: f.run}).Search("x", 1)
	if err != nil || len(got) != 0 {
		t.Fatalf("got %+v, %v; want no tracks and no error", got, err)
	}
}

// END: TestSearchEmpty

// START: TestSearchRunnerError

func TestSearchRunnerError(t *testing.T) {
	f := &fakeRunner{err: errors.New("boom")}
	got, err := (&YouTubeProvider{Runner: f.run}).Search("x", 1)
	if err == nil || len(got) != 0 {
		t.Fatalf("got %+v, %v; want an error and no tracks", got, err)
	}
}

// END: TestSearchRunnerError

// START: TestSearchQueryIsOneArg

func TestSearchQueryIsOneArg(t *testing.T) {
	q := "--exec touch /tmp/pwned; id"
	f := &fakeRunner{}
	if _, err := (&YouTubeProvider{Runner: f.run}).Search(q, 1); err != nil {
		t.Fatal(err)
	}
	if len(f.args) != 3 || f.args[0] != "ytsearch10:"+q || f.args[1] != "--flat-playlist" || f.args[2] != "-j" {
		t.Fatalf("args = %q", f.args)
	}
}

// END: TestSearchQueryIsOneArg

// START: TestResolveURL

func TestResolveURL(t *testing.T) {
	f := &fakeRunner{out: "https://s/a\n"}
	got, err := (&YouTubeProvider{Runner: f.run}).Resolve(Track{ID: "x"})
	if err != nil || got != "https://s/a" {
		t.Fatalf("got %q, %v", got, err)
	}
}

// END: TestResolveURL

// START: TestResolveArgs

func TestResolveArgs(t *testing.T) {
	f := &fakeRunner{out: "https://s/a\n"}
	if _, err := (&YouTubeProvider{Runner: f.run}).Resolve(Track{ID: "abc123"}); err != nil {
		t.Fatal(err)
	}
	want := []string{"-f", "bestaudio", "--print", "urls", "https://www.youtube.com/watch?v=abc123"}
	if !reflect.DeepEqual(f.args, want) {
		t.Fatalf("args = %q, want %q", f.args, want)
	}
}

// END: TestResolveArgs

// START: TestResolveRunnerError

func TestResolveRunnerError(t *testing.T) {
	f := &fakeRunner{err: errors.New("boom")}
	got, err := (&YouTubeProvider{Runner: f.run}).Resolve(Track{ID: "x"})
	if err == nil || got != "" {
		t.Fatalf("got %q, %v; want an error", got, err)
	}
}

// END: TestResolveRunnerError

// START: TestResolveEmpty

func TestResolveEmpty(t *testing.T) {
	f := &fakeRunner{}
	if got, err := (&YouTubeProvider{Runner: f.run}).Resolve(Track{ID: "x"}); err == nil {
		t.Fatalf("got %q, nil; want an error", got)
	}
}

// END: TestResolveEmpty

// START: TestSearchNoYtDlp

func TestSearchNoYtDlp(t *testing.T) {
	t.Setenv("PATH", "")
	if _, err := (&YouTubeProvider{}).Search("x", 1); err == nil {
		t.Fatal("want an error")
	}
}

// END: TestSearchNoYtDlp

// START: TestDetailsNoYtDlp

func TestDetailsNoYtDlp(t *testing.T) {
	t.Setenv("PATH", "")
	if _, err := (&YouTubeProvider{}).Details(Track{ID: "x"}); err == nil {
		t.Fatal("want an error")
	}
}

// END: TestDetailsNoYtDlp

// START: TestResolveNoYtDlp

func TestResolveNoYtDlp(t *testing.T) {
	t.Setenv("PATH", "")
	if _, err := (&YouTubeProvider{}).Resolve(Track{ID: "x"}); err == nil {
		t.Fatal("want an error")
	}
}

// END: TestResolveNoYtDlp

// START: TestSearchPage2Args

func TestSearchPage2Args(t *testing.T) {
	f := &fakeRunner{}
	if _, err := (&YouTubeProvider{Runner: f.run}).Search("daft punk", 2); err != nil {
		t.Fatal(err)
	}
	want := []string{"ytsearch20:daft punk", "--flat-playlist", "-I", "11:20", "-j"}
	if !reflect.DeepEqual(f.args, want) {
		t.Fatalf("args = %q, want %q", f.args, want)
	}
}

// END: TestSearchPage2Args

// START: TestSearchPage3Args

func TestSearchPage3Args(t *testing.T) {
	f := &fakeRunner{}
	if _, err := (&YouTubeProvider{Runner: f.run}).Search("daft punk", 3); err != nil {
		t.Fatal(err)
	}
	want := []string{"ytsearch30:daft punk", "--flat-playlist", "-I", "21:30", "-j"}
	if !reflect.DeepEqual(f.args, want) {
		t.Fatalf("args = %q, want %q", f.args, want)
	}
}

// END: TestSearchPage3Args

// START: TestEntryToTrack

func TestEntryToTrack(t *testing.T) {
	f := &fakeRunner{out: `{"id":"a1","title":"Get Lucky","thumbnails":[{"url":"small"},{"url":"big"}]}` + "\n"}
	got, err := (&YouTubeProvider{Runner: f.run}).Search("x", 1)
	if err != nil {
		t.Fatal(err)
	}
	want := []Track{{ID: "a1", Title: "Get Lucky", Artwork: "big", Source: "youtube"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tracks = %+v, want %+v", got, want)
	}
}

// END: TestEntryToTrack

// START: TestEntryThumbnailField

func TestEntryThumbnailField(t *testing.T) {
	f := &fakeRunner{out: `{"id":"a1","title":"Get Lucky","thumbnail":"t1"}` + "\n"}
	got, err := (&YouTubeProvider{Runner: f.run}).Search("x", 1)
	if err != nil || len(got) != 1 || got[0].Artwork != "t1" {
		t.Fatalf("got %+v, %v; want Artwork t1", got, err)
	}
}

// END: TestEntryThumbnailField

// START: TestEntryNoThumbnail

func TestEntryNoThumbnail(t *testing.T) {
	f := &fakeRunner{out: `{"id":"a1","title":"Get Lucky"}` + "\n"}
	got, err := (&YouTubeProvider{Runner: f.run}).Search("x", 1)
	if err != nil || len(got) != 1 || got[0].Artwork != "" {
		t.Fatalf("got %+v, %v; want an empty Artwork and no error", got, err)
	}
}

// END: TestEntryNoThumbnail

// START: entryOf

// entryOf runs Search on one JSON line and returns the single track.
func entryOf(t *testing.T, line string) Track {
	t.Helper()
	f := &fakeRunner{out: line + "\n"}
	got, err := (&YouTubeProvider{Runner: f.run}).Search("x", 1)
	if err != nil || len(got) != 1 {
		t.Fatalf("got %+v, %v; want one track", got, err)
	}
	return got[0]
}

// END: entryOf

// START: TestArtistField

func TestArtistField(t *testing.T) {
	got := entryOf(t, `{"id":"a","title":"t","artist":"Daft Punk","channel":"DaftPunkVEVO"}`)
	if got.Artist != "Daft Punk" {
		t.Fatalf("Artist = %q", got.Artist)
	}
}

// END: TestArtistField

// START: TestArtistFallsBackToChannel

func TestArtistFallsBackToChannel(t *testing.T) {
	got := entryOf(t, `{"id":"a","title":"t","channel":"Daft Punk"}`)
	if got.Artist != "Daft Punk" {
		t.Fatalf("Artist = %q", got.Artist)
	}
}

// END: TestArtistFallsBackToChannel

// START: TestArtistEmpty

func TestArtistEmpty(t *testing.T) {
	got := entryOf(t, `{"id":"a","title":"t"}`)
	if got.Artist != "" {
		t.Fatalf("Artist = %q, want empty", got.Artist)
	}
}

// END: TestArtistEmpty

// START: TestAlbumField

func TestAlbumField(t *testing.T) {
	got := entryOf(t, `{"id":"a","title":"t","album":"Random Access Memories"}`)
	if got.Album != "Random Access Memories" {
		t.Fatalf("Album = %q", got.Album)
	}
}

// END: TestAlbumField

// START: TestAlbumEmpty

func TestAlbumEmpty(t *testing.T) {
	got := entryOf(t, `{"id":"a","title":"t"}`)
	if got.Album != "" {
		t.Fatalf("Album = %q, want empty", got.Album)
	}
}

// END: TestAlbumEmpty

// START: TestDetailsArgs

func TestDetailsArgs(t *testing.T) {
	f := &fakeRunner{out: `{"id":"abc123"}`}
	if _, err := (&YouTubeProvider{Runner: f.run}).Details(Track{ID: "abc123"}); err != nil {
		t.Fatal(err)
	}
	want := []string{"-j", "https://www.youtube.com/watch?v=abc123"}
	if !reflect.DeepEqual(f.args, want) {
		t.Fatalf("args = %q, want %q", f.args, want)
	}
}

// END: TestDetailsArgs

// START: TestDetailsFills

func TestDetailsFills(t *testing.T) {
	f := &fakeRunner{out: `{"id":"abc","title":"Full Title","artist":"Daft Punk","album":"Random Access Memories","duration":248,"thumbnail":"http://t","channel":"c"}`}
	in := Track{ID: "abc", Title: "Get Lucky", Source: "youtube"}
	got, err := (&YouTubeProvider{Runner: f.run}).Details(in)
	if err != nil {
		t.Fatal(err)
	}
	want := Track{ID: "abc", Title: "Get Lucky", Artist: "Daft Punk", Album: "Random Access Memories", Artwork: "http://t", Duration: 248 * time.Second, Source: "youtube"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("track = %+v, want %+v", got, want)
	}
}

// END: TestDetailsFills

// START: TestDetailsFractionalDuration

func TestDetailsFractionalDuration(t *testing.T) {
	f := &fakeRunner{out: `{"id":"abc","duration":248.5}`}
	got, err := (&YouTubeProvider{Runner: f.run}).Details(Track{ID: "abc"})
	if err != nil || got.Duration != 248500*time.Millisecond {
		t.Fatalf("got %v, %v; want 248.5s", got.Duration, err)
	}
}

// END: TestDetailsFractionalDuration

// START: TestDetailsKeepsExisting

func TestDetailsKeepsExisting(t *testing.T) {
	f := &fakeRunner{out: `{"id":"abc","channel":"Daft Punk"}`}
	got, err := (&YouTubeProvider{Runner: f.run}).Details(Track{ID: "abc", Artwork: "old"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Artist != "Daft Punk" || got.Album != "" || got.Artwork != "old" {
		t.Fatalf("track = %+v", got)
	}
}

// END: TestDetailsKeepsExisting

// START: TestDetailsRunnerError

func TestDetailsRunnerError(t *testing.T) {
	f := &fakeRunner{err: errors.New("boom")}
	if _, err := (&YouTubeProvider{Runner: f.run}).Details(Track{ID: "abc"}); err == nil {
		t.Fatal("want an error")
	}
}

// END: TestDetailsRunnerError
