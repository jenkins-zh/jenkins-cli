#!/bin/sh

set -e

BIN_PATH=/usr/bin/jcli

# If we run make directly, any files created on the bind mount
# will have awkward ownership.  So we switch to a user with the
# same user and group IDs as source directory.  We have to set a
# few things up so that sudo works without complaining later on.
${BIN_PATH} job build ${JOB_NAME} \
    -b --url https://xxxx.com \
    --config-load false \
    --wait -l \
    --logger-level info $*