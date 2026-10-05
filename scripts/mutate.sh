#!/bin/sh
# Runs gremlins on the given packages, then stops the mpv processes its mutants left behind in its work directories.
# Usage: mutate.sh --threshold N <package>...  (N is the efficacy gremlins must reach)
[ "$1" = "--threshold" ] || { echo "usage: mutate.sh --threshold N <package>..." >&2; exit 2; }
n=$2
shift 2
gremlins unleash --threshold-efficacy "$n" "$@"
code=$?
pkill -f 'gremlins-.*/[g]oremi-.*/mpv.sock' # [g]: the pattern does not match this command line itself
exit $code
