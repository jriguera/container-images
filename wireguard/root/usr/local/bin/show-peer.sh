#!/usr/bin/with-contenv bash
# shellcheck shell=bash

if [[ ! $# -gt 0 ]]; then
  echo "You need to specify which peers to show"
  exit 0
fi

for i in "$@"
do
  [[ "${i}" =~ ^[0-9]+$ ]] && PEER_ID="peer${i}" || PEER_ID="peer_${i//[^[:alnum:]_-]/}"

  if grep -q "# ${PEER_ID}" /config/wg_confs/wg0.conf
  then
    echo "PEER ${i} QR code:"
    qrencode -t ansiutf8 < /config/${PEER_ID}/${PEER_ID}.conf
  else
    echo "PEER ${i} is not active"
  fi
done
