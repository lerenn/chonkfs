function cleanup()
{
    transmission-remote localhost -t 1 -r
    transmission-remote localhost --exit; killall -9 transmission-daemon
    exit 0
}
trap cleanup EXIT
trap cleanup SIGINT

transmission-daemon -w ./mnt --logfile ./transmission.log --log-level debug
transmission-remote localhost -a ./debian-12.9.0-amd64-netinst.iso.torrent
tail -f ./transmission.log