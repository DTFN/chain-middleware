#!/bin/sh

# start LingShuChain
cd /chain/node
cp -r /chain/node/sdk/* /front/conf/
#/usr/local/bin/LingShuChain -c /data/config.yaml >>nohup.out &
bash /launch.sh

sleep 1

# start front
cd /front && bash start.sh

# keep container running
tail -f /docker-start.sh
