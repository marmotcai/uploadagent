FROM marmotcai/golang AS building

MAINTAINER marmotcai "marmotcai@163.com"
RUN yum install -y gcc-c++

ENV APP_NAME=uploadagent
ENV APP_SOURCE_DIR=$GOPATH/src/github.com/marmotcai/${APP_NAME}
ENV APP_GIT_URL=github.com/marmotcai/uploadagent
ENV OUTPUT_PATH=${APP_SOURCE_DIR}/output
ENV OUTPUT_PACKETS=/root/${APP_NAME}.tar.gz

RUN gopm get -g -v ${APP_GIT_URL}
WORKDIR ${APP_SOURCE_DIR}

ENV ENTRYPOINT_FILE=${APP_SOURCE_DIR}/build.sh
RUN chmod +x ${ENTRYPOINT_FILE} && \
    ${ENTRYPOINT_FILE}

FROM marmotcai/centos-base AS uploadagent

ENV APP_NAME=uploadagent
ENV UA_PATH=/root/ua
RUN mkdir -p $UA_PATH

COPY --from=building ${OUTPUT_PACKETS} ${UA_PATH}/

WORKDIR ${UA_PATH}
RUN tar xvf ${APP_NAME}.tar.gz
RUN rm -f ${APP_NAME}.tar.gz

RUN chmod +x ./ua && ./ua

RUN yum install -y mediainfo
