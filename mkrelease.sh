#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -o demon/linux_amd64/demon
CGO_ENABLED=0 GOOS=darwin go build -o demon/darwin_amd64/demon
CGO_ENABLED=0 GOOS=windows go build -o demon/windows_amd64/demon.exe

zip -r demon.zip demon
