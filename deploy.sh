#!/bin/bash
export GOROOT=/usr/lib/google-golang
export GOPATH=/usr/local/google/home/dan/.go/
gcloud --project sbb-status-4f4eb app deploy --quiet

