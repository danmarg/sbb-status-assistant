#!/bin/bash
export GOROOT=/usr/lib/google-golang
export GOPATH=$PWD
pushd .
cd src/app
gcloud --project sbb-status-4f4eb app deploy --quiet
popd

