#!/bin/bash

rm -f git-env && go build -ldflags="-s -w" && mv git-env ~/bin
