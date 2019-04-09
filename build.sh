#!/bin/bash

echo "golang build..."

if [ ! -d "${APP_SOURCE_DIR}" ]; then
  echo "git clone ${APP_GIT_URL} to ${APP_SOURCE_DIR}" 
  go get -v ${APP_GIT_URL}
fi

cd ${APP_SOURCE_DIR}
git pull

echo "go build ${APP_SOURCE_DIR} to ${OUTPUT_PATH}"
go build -o ${OUTPUT_PATH}/uploadagent
cd ${OUTPUT_PATH}
tar -zcvf ${OUTPUT_PACKETS} ./*

