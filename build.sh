#!/bin/bash
SYNO=0
# -------

usage() { echo "Usage: ./build.sh [--synology]" 1>&2; exit 1; }

while :; do
    case $1 in
        -h|-\?|--help)   # Call a "show_help" function to display a synopsis, then exit.
            usage
            exit
            ;;
        --syno|--synology)
            SYNO=1
            ;;
        --)              # End of all options.
            shift
            break
            ;;
        -?*)
            printf 'WARN: Unknown option (ignored): %s\n' "$1" >&2
            ;;
        *)               # Default case: If no more options then break out of the loop.
            break
      esac

      shift
done

if [[ "$SYNO" != "0" ]]; then
    echo "[BUILD] Cross-compile synology version"
    GOOS=linux GOARCH=amd64 go build
else
    echo "[BUILD] Compile standard version"
    go build
fi