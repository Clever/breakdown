#!/bin/bash -x

BASE_DIR="$(dirname $0)"

DIR="${PWD##*/}"

BREAKDOWN_URL=http://localhost:8080

TMP_OUT="$(mktemp)"

RAND_COMMIT_SHA="$(cat /dev/urandom | env LC_ALL=C tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)"

$BASE_DIR/bin/breakdown-cli \
    -output $TMP_OUT \
    -pretty=f \
    Clever/$DIR $RAND_COMMIT_SHA

curl -X POST -H "Content-Type: application/json" -o - \
    -d @$TMP_OUT $BREAKDOWN_URL/v1/upload

rm $TMP_OUT

