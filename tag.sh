#!/bin/sh

VERSION=$(hub tag --list | tail -n 1)
VERSION_BITS=(${VERSION//./ })
VNUM1=${VERSION_BITS[0]}
VNUM2=${VERSION_BITS[1]}
VNUM3=${VERSION_BITS[2]}
VNUM3=$((VNUM3+1))

NEW_TAG="$VNUM1.$VNUM2.$VNUM3"
echo "Updating $VERSION to $NEW_TAG"

#get current hash and see if it already has a tag
GIT_COMMIT=`git rev-parse HEAD`
NEEDS_TAG=`git describe --contains $GIT_COMMIT`

if [[ -z "${NEEDS_TAG}" ]]; then
    echo "hub release create -c -a release/jcli-darwin-amd64.tar.gz v${NEW_TAG}"
fi
