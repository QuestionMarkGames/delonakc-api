#!/bin/sh

port=8088
version=1.0
name="delonakc-api"
mongoname="mongodb-test"
rdname="redis-test"

redisService=$(docker ps -qf name=$rdname)
mongodbService=$(docker ps -qf name=$mongoname)

# 判断redis 容器是否启动
if [ "${redisService}" = "" ]
then
   echo ""
   echo "=====start redis ...====="
   echo ""
   rdid=$(docker ps -aqf name=$rdname)
   if [ "${rdid}" != "" ]
   then
     docker stop $rdname && docker rm -f $rdname
   fi
   docker run -d --name $rdname redis:latest
fi
# 判断MongoDB 容器是否启动
if [ "${mongodbService}" = "" ]
then
    echo ""
    echo "=====start mongodb ...====="
    echo ""
    mdbid=$(docker ps -aqf name=$mongoname)
    if [ "${mdbid}" != "" ]
    then
        docker stop $mongoname && docker rm -f $mongoname
    fi
    #docker run -d --name $mongoname -p 27017:27017 delonakc/mongo:1.0.0 --config /etc/mongod.conf --smallfiles
    docker run -d --name $mongoname delonakc/mongo:version-1.1.0
fi

echo ""
echo "=====start running...====="
echo ""

# tag名称
tag="${name}-test:${version}"
label="${name}-test-${version}"
prevImageId=$(docker images -qf "label=${label}")

# 删除之前的镜像
if [ "${prevImageId}" != "" ]
then
    docker rmi -f $prevImageId
fi

docker build -t="${tag}" --label="${label}" .

if [ $? -ne 0 ]
then
    echo ""
    echo "====docker build error...===="
    exit 1
fi

echo ""
echo "=====start running...====="
echo ""

prevContainer=$(docker ps -aqf "name=${name}")

# 删除之前的容器
if [ "${prevContainer}" != "" ]
then
    docker rm -f $name
fi

# 启动容器
docker run -d --name $name --link $mongoname:mongodb --link $rdname:redis -p $port:$port $tag
