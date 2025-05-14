#!/usr/bin/env bash

if [[ $# -lt 1 || $1 -lt 1 ]]; then
	echo "usage: $(basename $0) <count>"
	exit 1
fi

PROJECT_ROOT=$( cd "$( dirname "${BASH_SOURCE[0]:-$0}" )/.." && pwd )
HOOKS_PATH=$PROJECT_ROOT/.aagh
HOOK_NAME="post-checkout"
HOOK_SCRIPT=$HOOKS_PATH/$HOOK_NAME
SPECIFIED_COUNT=$1

cleanup() {
	rm $HOOKS_PATH/_/$HOOK_NAME
	rm -rf $HOOK_SCRIPT
}

create() {
	WAIT=$((1+RANDOM % 3))
	echo 'echo "$0" slept for ' "$WAIT" >> $HOOK_SCRIPT$1
	echo "sleep $WAIT" >> $HOOK_SCRIPT$1
}

trap cleanup EXIT
go build -o $HOOKS_PATH/_/$HOOK_NAME $PROJECT_ROOT/cmd/runner/main.go

if [[ $SPECIFIED_COUNT -gt 1 ]]; then
	mkdir -p $HOOK_SCRIPT

	for ((i=1; i <= $SPECIFIED_COUNT; i++)); do
		for type in lint format test; do
			create "/$i-$type"
		done
	done
else
	create ""
fi

go run $PROJECT_ROOT/cmd/cli/main.go run $HOOK_NAME
