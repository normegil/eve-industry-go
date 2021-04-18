#!/bin/sh
GOLANG_CI_VERSION="v1.39.0"
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $GOLANG_CI_VERSION