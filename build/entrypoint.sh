#!/bin/sh

if [ "${1:0:1}" = "-" ]; then
  set -- hkserver "$@"
fi

exec "$@"