#!/usr/bin/env bash

PROJECT="wxpay_sdk"

OLD_CWD=`pwd -P`

PROJECT_ROOT="${OLD_CWD}"

OLD_GOPATH=${GOPATH}

# setting current gopath
GOPATH="GOPATH=${PROJECT_ROOT}:${OLD_GOPATH}"
export ${GOPATH}

cd ${PROJECT_ROOT}

# 1. generate doc
function generate_doc(){
    if [[ ! -d "doc" ]]
    then
        mkdir "doc"
    fi
    touch doc/USAGE
    echo "pkg info: http://localhost:PORT/github.com/wengle520/wxpay-sdk/" >> doc/USAGE
    return 0
}

# 2. build project
function build() {
    # base on gopath, stringutil is a package name
    # go install github.com/wengle520/stringutil
    go install wxpay
    if [[ $? -ne 0 ]]
    then
        echo "go install failure."
        return 1
    fi
    echo "go install success."
    return 0
}

# 3. packing
function packing(){
    bash ./pack.sh
    if [[ $? -ne 0 ]]
    then
        return 1
    fi
    return 0
}

function main() {
    generate_doc && \
    build && \
    packing
    if [[ $? -ne 0 ]]
    then
        echo "build failure."
        return 1
    fi
    return 0
}

main
ret=$?

GOPATH="GOPATH=${OLD_GOPATH}"
export ${GOPATH}
cd ${OLD_CWD}

if [[ ${ret} -ne 0 ]]
then
    exit 1
fi
echo "build success."
