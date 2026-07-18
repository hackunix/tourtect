#!/bin/sh
set -eu

tourtect-migrate
if [ "${TOURTECT_SEED_DEMO:-false}" = "true" ]; then
  tourtect-seed
fi
exec tourtect-api
