FROM marmotcai/golang AS building

MAINTAINER marmotcai "marmotcai@163.com"

ENV WORK_DIR=/root
ENV APP_NAME=md5tools
ENV APP_SOURCE_DIR=$GOPATH/src/github.com/marmotcai/${APP_NAME}
ENV OUTPUT_PATH=${WORK_DIR}/output
ENV OUTPUT_PACKETS=${OUTPUT_PATH}/${APP_NAME}.tar.gz

VOLUME $OUTPUT_PATH

RUN mkdir -p APP_SOURCE_DIR
WORKDIR ${APP_SOURCE_DIR}

COPY main.go .
COPY build.sh .
RUN chmod +x build.sh
# RUN build.sh
# RUN go build  -o md5tools main.go
# RUN tar -zcvf ${OUTPUT_PACKETS} md5tools
CMD ["./build.sh"]
