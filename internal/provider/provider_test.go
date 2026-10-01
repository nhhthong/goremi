// Tests for the Provider interface contract.
package provider

import (
	"reflect"
	"sort"
	"testing"
)

// START: TestYouTubeImplementsProvider

func TestYouTubeImplementsProvider(t *testing.T) {
	var p Provider = &YouTubeProvider{}
	if p == nil {
		t.Fatal("provider is nil")
	}
}

// END: TestYouTubeImplementsProvider

// START: TestProviderMethodSet

func TestProviderMethodSet(t *testing.T) {
	typ := reflect.TypeOf((*Provider)(nil)).Elem()
	var got []string
	for i := 0; i < typ.NumMethod(); i++ {
		got = append(got, typ.Method(i).Name)
	}
	sort.Strings(got)
	want := []string{"Details", "Resolve", "Search"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("methods = %v, want %v", got, want)
	}
}

// END: TestProviderMethodSet
