// The version the build stamps into the binary.
package app

// START: Version

// Version is the release shown in the badge; `make build` and the release script set it with -ldflags "-X goremi/internal/app.Version=...", and an unstamped build is "dev".
var Version = "dev"

// END: Version
