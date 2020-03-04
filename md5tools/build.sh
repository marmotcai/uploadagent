#!/usr/bin/env bash

go build  -o md5tools main.go
tar -zcvf ${OUTPUT_PACKETS} md5tools
echo ${OUTPUT_PACKETS}
./md5tools