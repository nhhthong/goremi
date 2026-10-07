// Search input: the query text the user types.
package ui

import tea "charm.land/bubbletea/v2"

// START: SearchInput

// SearchInput holds the typed query.
type SearchInput struct {
	value string
}

func (s SearchInput) Value() string { return s.value }

// WithValue returns the input holding v.
func (s SearchInput) WithValue(v string) SearchInput {
	s.value = v
	return s
}

// Update appends the printable text of a key press and removes the last character on Backspace;
// keys with Ctrl or Alt are not text.
func (s SearchInput) Update(k tea.KeyPressMsg) SearchInput {
	if k.Code == tea.KeyBackspace {
		if r := []rune(s.value); len(r) > 0 {
			s.value = string(r[:len(r)-1])
		}
		return s
	}
	if k.Text != "" && k.Mod&(tea.ModCtrl|tea.ModAlt) == 0 {
		s.value += k.Text
	}
	return s
}

// END: SearchInput
