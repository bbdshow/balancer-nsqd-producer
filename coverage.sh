#!/bin/bash

set -e
for d in $(go list ./... | grep -E 'balancer-nsqd-producer$|algorithm$'); do
    go test -v -covermode=count -coverprofile=coverage.out $d
done

for d in $(go list ./... | grep -E 'balancer-nsqd-producer$|algorithm$'); do
   go test -race -covermode=atomic -coverprofile=coverage.txt 
done

