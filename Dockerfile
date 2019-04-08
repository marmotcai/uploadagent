FROM marmotcai/golang AS base

MAINTAINER marmotcai "marmotcai@163.com"

RUN gopm get -g golang.org/x/sys && \
    gopm get -g golang.org/x/text && \
    gopm get -g golang.org/x/time && \
    gopm get -g golang.org/x/crypto && \
    gopm get -g golang.org/x/net

RUN gopm get -g -v gopkg.in/urfave/cli.v1 && \
    gopm get -g -v gopkg.in/yaml.v2 && \
    gopm get -g -v github.com/aliyun/aliyun-oss-go-sdk && \
    gopm get -g -v github.com/aws/aws-sdk-go && \
    gopm get -g -v github.com/spf13/viper && \
    gopm get -g -v github.com/astaxie/beego && \
    gopm get -g -v github.com/bramvdbogaerde/go-scp && \
    gopm get -g -v github.com/fatih/color && \
    gopm get -g -v github.com/secsy/goftp && \
    gopm get -g -v github.com/yanyiwu/gosimhash

RUN yum install -y gcc-c++

FROM base AS building

ENV APP_NAME=uploadagent
ENV APP_SOURCE_DIR=$GOPATH/src/github.com/marmotcai/${APP_NAME}
ENV SOURCE_URL=http://git.atoml.com/caijun/${APP_NAME}
ENV GIT_URL=http://git.atoml.com/caijun/${APP_NAME}.git

RUN git clone ${GIT_URL} ${APP_SOURCE_DIR} 
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



