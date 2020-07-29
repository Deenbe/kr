#!/bin/bash

set -ex

BINARY=""
if [[ "$(uname)" = "Linux" ]]; then
  BINARY="kr_linux_amd64"
fi

if [[ "$(uname)" = "Darwin" ]]; then
  BINARY="kr_darwin_amd64"
fi

if [[ "$BINARY" = "" ]]; then
    echo "OS $(uname) is not supported"
    exit 1
fi

curl -o /usr/local/bin/kr -L "https://github.com/Deenbe/kr/releases/latest/download/$BINARY" 
chmod +x /usr/local/bin/kr

