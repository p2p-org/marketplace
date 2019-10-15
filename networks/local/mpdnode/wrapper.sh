#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/mpd/${BINARY:-mpd}
ID=${ID:-0}
LOG=${LOG:-mpd.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'mpd' E.g.: -e BINARY=mpd_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export MPDHOME="/mpd/node${ID}/mpd"

if [ -d "$(dirname "${MPDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${MPDHOME}" "$@" | tee "${MPDHOME}/${LOG}"
else
  "${BINARY}" --home "${MPDHOME}" "$@"
fi

``