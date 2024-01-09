#!/usr/bin/env bash

# set -x

cmd_dir="cmd/notifier"
binary=event-notifier

if ! [ -x "$(command -v go)" ]; then
    echo -e "go has to be installed"
    exit 1
fi

if ! [ -x "$(command -v git)" ]; then
    echo -e "git has to be installed"
    exit 1
fi

cd ${cmd_dir} && \
    go build -v \
        -ldflags="-X main.appVersion=$(git describe --tags --long --dirty)" \
        -o ${binary}
