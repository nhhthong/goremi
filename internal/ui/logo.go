// The ASCII art "Goremi" shown in the player panel; made once from resources/ascii_terminal.png.
package ui

import "strings"

// START: Logo

// Logo is the art, 40 columns wide, one line per row.
var Logo = strings.Join([]string{
	"        ▟███▙                        ▐█▌",
	"    ▄  ▐█▛▀▜█▌▗▄▄▄ ▄▄▗▄ ▗▄▄▄ ▗▄▄▄▗▄▄ ▄▄",
	"  ▗▖█  ▐█     ████▌▜███▚████▙▜██████▌▜█▌",
	" ▗▐▌█▟ ▐█ ▐██▘█▌ █▌ █▌ ▐█▄▄██▐█▌██ █▌▐█▌",
	"▄▜▐▌██ ▐█  ▜█ █▌ █▌ █▌ ▐█▀▀▀▀▐█▌▜█ █▌▐█▌",
	"▛▐▐▌▌█ ▝██▄█▛ █▙▄█▌▟█▌ ▐█▙▄▄▖▟█▌██ █▌▟█▙",
	"        ▝▀▀▀  ▝▀▀▀ ▝▀▘  ▀▀▀▀ ▝▀▘▝▘ ▀▘▀▀▀",
}, "\n")

// END: Logo
