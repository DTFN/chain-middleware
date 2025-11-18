#! /bin/bash

# 打印所有设备的输入和发送流量
# tail -n +3 /proc/net/dev | awk '{print $1""$2":"$10}'

#
start_time=`date +%s%N`
start_rcv=$(tail -n +3 /proc/net/dev | awk '{sum+=$2}END{print sum}')
start_send=$(tail -n +3 /proc/net/dev | awk '{sum+=$10}END{print sum}')
sleep 1
end_time=`date +%s%N`
end_rcv=$(tail -n +3 /proc/net/dev | awk '{sum+=$2}END{print sum}')
end_send=$(tail -n +3 /proc/net/dev | awk '{sum+=$10}END{print sum}')

# 计算网速
net_in=$[ (end_rcv - start_rcv) * 1000000000 / (end_time - start_time) ]
net_out=$[ (end_send - start_send) * 1000000000 / (end_time - start_time)  ]

# json格式输出
echo $net_in # 入网
echo $net_out # 出网
