#!/usr/bin/env sh

# Generate test coverage statistics for Go packages.
#
# Works around the fact that `go test -coverprofile` currently does not work
# with multiple packages, see https://code.google.com/p/go/issues/detail?id=6909
#
# Usage: script/coverage [--html|--coveralls]
#
#     --html      Additionally create HTML report and open it in browser
#

set -e

workdir=.
profile="coverage.txt"
mode=atomic

generate_cover_data() {
    for pkg in "$@"; do
        f="$workdir/$(echo $pkg | tr / -).coverprofile"
        go test -covermode="$mode" -coverprofile="$f" "$pkg"
    done

    echo "mode: $mode" >"$profile"
    grep -h -v "^mode:" "$workdir"/*.coverprofile >>"$profile"
    rm -f *.coverprofile
}

show_cover_report() {
    go tool cover -${1}="$profile"
}

generate_cover_data $(go list ./...)
show_cover_report func
case "$1" in
"")
    ;;
--html)
    show_cover_report html ;;
*)
    echo >&2 "error: invalid option: $1"; exit 1 ;;
esac
