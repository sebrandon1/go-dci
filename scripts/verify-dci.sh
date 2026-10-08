#!/usr/bin/env bash

set -euo pipefail

run_smoke() {
	local command="$*"

	printf 'Verifying go-dci %s...\n' "$command"
	if ! ./go-dci "$@" > /dev/null; then
		printf 'Failed to verify go-dci %s\n' "$command" >&2
		exit 1
	fi
}

run_smoke identity
run_smoke topics --name go-dci-nightly-smoke-probe
run_smoke components --name go-dci-nightly-smoke-probe
run_smoke products
run_smoke teams --name go-dci-nightly-smoke-probe
# RemoteCI credentials are not authorized to list users (HTTP 401).
run_smoke remotecis
# The live API does not allow listing job states (HTTP 405).
run_smoke jobs --age 1
run_smoke ocpcount --age 1
