// YouTube provider: searches and resolves tracks by calling the yt-dlp CLI.
package provider

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const pageSize = 10

// START: Runner

// Runner runs yt-dlp with args and returns its standard output; tests replace it.
type Runner func(args ...string) ([]byte, error)

type YouTubeProvider struct {
	Runner Runner
}

// run uses the injected Runner, else the real yt-dlp via os/exec (no shell).
func (p *YouTubeProvider) run(args ...string) ([]byte, error) {
	if p.Runner != nil {
		return p.Runner(args...)
	}
	return exec.Command("yt-dlp", args...).Output()
}

// END: Runner

// START: Search

// Search returns up to 10 tracks of the given page. Page n asks yt-dlp for the first 10n results
// and keeps items 10(n-1)+1 to 10n; the "ytsearch" prefix keeps the query one argument.
func (p *YouTubeProvider) Search(query string, page int) ([]Track, error) {
	if page < 1 {
		page = 1
	}
	args := []string{fmt.Sprintf("ytsearch%d:%s", pageSize*page, query), "--flat-playlist"}
	if page > 1 {
		args = append(args, "-I", fmt.Sprintf("%d:%d", pageSize*(page-1)+1, pageSize*page))
	}
	out, err := p.run(append(args, "-j")...)
	if err != nil {
		return nil, err
	}
	// parse one JSON object per line, stop at the page size
	var tracks []Track
	sc := bufio.NewScanner(bytes.NewReader(out))
	sc.Buffer(nil, 1<<20)
	for sc.Scan() && len(tracks) < pageSize {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, err
		}
		tracks = append(tracks, e.track())
	}
	return tracks, sc.Err()
}

// END: Search

// START: entry mapping

// entry is the part of a yt-dlp JSON object that a Track needs.
type entry struct {
	ID         string
	Title      string
	Artist     string
	Channel    string
	Album      string
	Duration   float64
	Thumbnail  string
	Thumbnails []struct{ URL string }
}

// track maps an entry to a Track. Artist is "artist", else "channel"; Album is "album";
// Duration is "duration" in seconds; Artwork is "thumbnail", else the last (largest) of "thumbnails".
func (e entry) track() Track {
	art := e.Thumbnail
	if art == "" && len(e.Thumbnails) > 0 {
		art = e.Thumbnails[len(e.Thumbnails)-1].URL
	}
	artist := e.Artist
	if artist == "" {
		artist = e.Channel
	}
	return Track{ID: e.ID, Title: e.Title, Artist: artist, Album: e.Album, Artwork: art, Duration: time.Duration(e.Duration * float64(time.Second)), Source: "youtube"}
}

// END: entry mapping

// START: Details

// Details fills Artist, Album, Artwork and Duration from the full metadata of one track;
// a metadata value that is empty never erases a value the track already has.
func (p *YouTubeProvider) Details(track Track) (Track, error) {
	out, err := p.run("-j", "https://www.youtube.com/watch?v="+track.ID)
	if err != nil {
		return track, err
	}
	var e entry
	if err := json.Unmarshal(bytes.TrimSpace(out), &e); err != nil {
		return track, err
	}
	full := e.track()
	if full.Artist != "" {
		track.Artist = full.Artist
	}
	if full.Album != "" {
		track.Album = full.Album
	}
	if full.Artwork != "" {
		track.Artwork = full.Artwork
	}
	if full.Duration != 0 {
		track.Duration = full.Duration
	}
	return track, nil
}

// END: Details

// START: Resolve

// Resolve returns the best-audio stream URL of a track; empty output is an error.
func (p *YouTubeProvider) Resolve(track Track) (string, error) {
	out, err := p.run("-f", "bestaudio", "--print", "urls", "https://www.youtube.com/watch?v="+track.ID)
	if err != nil {
		return "", err
	}
	url := strings.TrimSpace(string(out))
	if url == "" {
		return "", errors.New("yt-dlp returned no stream URL")
	}
	return url, nil
}

// END: Resolve
