#!/bin/bash
export GOROOT=/usr/lib/google-golang
export GOPATH=$PWD
pushd .
cd src/app
# XXX: This is stupid.
cp -r ../localize/data/ .
gcloud --project sbb-status-4f4eb app deploy --quiet
popd

