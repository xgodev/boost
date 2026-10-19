#!/usr/bin/env bash
# Builds, vets and tests every example module: each example is a builder app
# in its own Go module, which the root `go test ./...` does not reach.
#
#   scripts/examples.sh            # every example
#   scripts/examples.sh wrapper    # the examples under wrapper/
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

failed=()
while IFS= read -r mod; do
	dir="$(dirname "$mod")"
	echo "== $dir"
	if ! (cd "$dir" && go vet ./... && go test ./...); then
		failed+=("$dir")
	fi
done < <(find "${1:-.}" -name go.mod -not -path ./go.mod -not -path '*/.solvers/*' -not -path '*/testdata/*' | sort)

if ((${#failed[@]})); then
	printf 'FAILED: %s\n' "${failed[@]}" >&2
	exit 1
fi
