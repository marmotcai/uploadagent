FROM marmotcai/golang AS building

MAINTAINER marmotcai "marmotcai@163.com"
RUN yum install -y gcc-c++

ENV WORK_DIR=/root
ENV APP_NAME=uploadagent
ENV APP_GIT_URL=github.com/marmotcai/${APP_NAME}
ENV APP_SOURCE_DIR=$GOPATH/src/github.com/marmotcai/${APP_NAME}
ENV OUTPUT_PATH=${APP_SOURCE_DIR}/output
ENV OUTPUT_PACKETS=${WORK_DIR}/${APP_NAME}.tar.gz

RUN gopm get -g -v ${APP_GIT_URL}
WORKDIR ${APP_SOURCE_DIR}

ENV ENTRYPOINT_FILE=${APP_SOURCE_DIR}/build.sh
RUN chmod +x ${ENTRYPOINT_FILE} && \
    ${ENTRYPOINT_FILE}

FROM marmotcai/centos-base

ENV WORK_DIR=/root
ENV APP_NAME=uploadagent
ENV OUTPUT_PATH=${WORK_DIR}/${APP_NAME}
ENV OUTPUT_PACKETS=${WORK_DIR}/${APP_NAME}.tar.gz

COPY --from=building ${OUTPUT_PACKETS} ${OUTPUT_PACKETS}

RUN mkdir -p ${OUTPUT_PATH}
WORKDIR ${OUTPUT_PATH}
RUN tar xvf ${OUTPUT_PACKETS} . && rm -f ${OUTPUT_PACKETS}

RUN chmod +x ./${APP_NAME} && ./${APP_NAME}

RUN yum install -y mediainfo

RUN localedef -c -f UTF-8 -i zh_CN zh_CN.utf8
ENV LC_ALL zh_CN.utf8
