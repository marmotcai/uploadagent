#!/bin/bash

cmd=${1}
param1=${2}
case $cmd in 
    commit)
      git add .
      curtime=`date "+%Y-%m-%d:%H:%M:%S"`
      git commit -m "auto commit ${curtime}"
      git push
    ;;

    pull)
      git pull
    ;;

    build)
      docker build --target uploadagent -t ${param1} .
      docker run --rm -ti ${param1} /bin/bash
    ;;

    *)
      echo "use: sh make.sh commit"
      echo "use: sh make.sh pull"
      echo "use: sh make.sh build imagename"
    ;;
esac

exit 0;
