#!/bin/bash

set -e

echo "mode: count" > coverage.out

for d in $(go list ./... | grep -E 'balancer-nsqd-producer$|algorithm$'); do
    go test -v -covermode=count -coverprofile=profile.out $d
    if [ -f profile.out ]; then
        cat profile.out | grep -v "mode:" >> coverage.out
        rm profile.out
    fi
done

echo "mode: goveralls" > profile.out
$GOPATH/bin/goveralls -coverprofile="profile.out" \
                          -service=travis-ci    \
                          -ignore="example" || true
rm profile.out