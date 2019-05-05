#!/usr/bin/env bash

PROJECT="wxpay_sdk"

OLD_CWD=`pwd -P`

DIR_PATH=`dirname $0`
PROJECT_ROOT="${OLD_CWD}/${DIR_PATH}"

FILE_DATE=`date "+%Y%m%d_%H:%M:%S"`

FILE_NAME="${PROJECT}.${FILE_DATE}.tar.gz"

cd ${PROJECT_ROOT}
if [[ ! -d "pack" ]]
then
    mkdir -p "pack"
fi

tar --exclude="pack" --exclude=".git" -czvf ./pack/${FILE_NAME} .
if [[ $? -ne 0 ]]
then
    echo "Packing ${FILE_NAME} failure."
    exit 1
fi

cd ${OLD_CWD}
echo "Packing ${FILE_NAME} success."
