#!/usr/bin/env bash

set -eox pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

function cleanup()
{
    transmission-remote localhost -t 1 -r
    transmission-remote localhost --exit; killall -9 transmission-daemon
    exit 0
}
trap cleanup EXIT
trap cleanup SIGINT

mkdir -p ./mnt
transmission-daemon -w ./mnt --logfile ./transmission.log --log-level debug
transmission-remote localhost -a ${SCRIPT_DIR}/debian-12.9.0-amd64-netinst.iso.torrent
tail -f ./transmission.log