FROM marmotcai/golang AS building

MAINTAINER marmotcai "marmotcai@163.com"
RUN yum install -y gcc-c++

ENV APP_NAME=uploadagent
ENV APP_SOURCE_DIR=$GOPATH/src/github.com/marmotcai/${APP_NAME}
ENV GIT_URL=github.com/marmotcai/uploadagent

RUN gopm get -g -v ${GIT_URL}
WORKDIR ${APP_SOURCE_DIR}

ENV ENTRYPOINT_FILE=${APP_SOURCE_DIR}/build.sh
RUN chmod +x ${ENTRYPOINT_FILE} && \
    ${ENTRYPOINT_FILE}

FROM marmotcai/centos-base AS app

RUN yum install -y mediainfo

ENV APP_NAME=uploadagent
ENV OUTPUT_PACKETS=/root/output/${APP_NAME}.tar.gz

ENV UA_PATH=/root/ua

COPY --from=building ${OUTPUT_PACKETS} ${UA_PATH}/

WORKDIR ${UA_PATH}
RUN tar xvf ${APP_NAME}.tar.gz
RUN rm -f ${APP_NAME}.tar.gz

RUN chmod +x ./ua && ./ua



