#!/bin/sh

# start lingshuchain
cd /data
/usr/local/bin/LingShuChain -c /data/config.yaml >>nohup.out &

# start front
cp -r /data/sdk/* /front/conf/
cd /front && bash start.sh

# keep container running
tail -f /docker-start.sh

