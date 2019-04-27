#!/bin/bash

git_url=${1}
output_path=${2}
app_name=${3}
arch=${4}

if  [[ ${git_url} = '' ]]; then
  echo "use: docker-build.sh git_url ./output ua"
  exit 0
fi

if  [[ ${output_path} = '' ]]; then
  output_path="./"  
fi

if  [[ ${app_name} = '' ]]; then
  curtime=`date "+%Y-%m-%d-%H-%M-%S"`
  app_name="builder_${curtime}"
fi

echo "git url : ${git_url}"
echo "output path : ${output_path}"
echo "output app filename : ${app_name}"
echo "arch : ${arch}"

case $arch in
    arm)
      docker run --rm -ti \
                 --env CGO_ENABLE=0 \
                 --env GOARCH=arm \
                 --env GOOS=linux \
                 -v $output_path:/root/output \
                 marmotcai/golang-builder build $git_url $app_name
    ;;

    *)
      docker run --rm -ti \
                 -v $output_path:/root/output \
                 marmotcai/golang-builder build $git_url $app_name
    ;;

  esac

exit 0
