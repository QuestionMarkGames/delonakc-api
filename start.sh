#!/bin/sh

port=8088
version=1.0
name="delonakc-api"

redisService=$(docker ps -qf name=redis)
mongodbService=$(docker ps -qf name=mongodb)

# 判断redis 容器是否启动
if [ "${redisService}" = "" ]
then
   echo ""
   echo "=====start redis ...====="
   echo ""
   docker run -d --name redis daocloud.io/library/redis:3.2.9
fi

# 判断MongoDB 容器是否启动
if [ "${mongodbService}" = "" ]
then
    echo ""
    echo "=====start mongodb ...====="
    echo ""
    docker run -d --name mongodb -p 27017:27017 -v /var/lib/mongodb:/data/db mongo:4.0.8
fi

echo ""
echo "=====start running...====="
echo ""

docker build -t="${name}:${version}" .

if [ $? -ne 0 ]
then
    echo ""
    echo "====docker build error...===="
    exit 1
fi

echo ""
echo "=====start running...====="
echo ""

docker run -d --name $name --link mongodb:mongodb --link redis:redis -p $port:$port $name:$version
