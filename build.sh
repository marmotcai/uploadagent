#!/bin/bash

echo "golang build..."

APP_SOURCE_DIR=$GOPATH/src/github.com/marmotcai/${APP_NAME}

SOURCE_URL=http://git.atoml.com/caijun/${APP_NAME}
ENV GIT_URL=github.com/marmotcai/uploadagent

if [ ! -d "${APP_SOURCE_DIR}" ]; then
  echo "git clone ${GIT_URL} to ${APP_SOURCE_DIR}" 
  go get -v ${GIT_URL}
fi

cd ${APP_SOURCE_DIR}
git pull

OUTPUT_PATH=${APP_SOURCE_DIR}/output

echo "go build ${APP_SOURCE_DIR} to ${OUTPUT_PATH}/ua"
go build -o ${OUTPUT_PATH}/ua

cd ${OUTPUT_PATH}
tar -zcvf ${OUTPUT_PACKETS} ./*

