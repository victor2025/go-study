#!/bin/sh
# 递归找到进程最底层子进程并杀除.
mainName=$1
echo "=====开始====="
#grep -v可以过滤掉grep的进程，但测试时有时会出现问题，所以加上获取第一行
mainId=`ps -ef |grep ${mainName}|grep -v 'grep' |head -1|cut -c 9-15`
#也可以使用这种方法获取查出的第一个参数
#mainId=`ps -A |grep ${mainName}|awk '{print $1}'`
#去掉空格
mainId=`echo ${mainId}|sed 's/ //g'`
echo "mainId===${mainId}"
#查主进程下所有子进程 格式为main.sh(275)---children1.sh(27641)---sleep(27643)
pidLine=`pstree -p ${mainId}`
echo "pidLine===pidLine${pidLine}"
#取括号中的内容
pidLine=`echo $pidLine | awk 'BEGIN{ FS="(" ; RS=")" } NF>1 { print $NF }'`
#echo $pidLine
for pid in $pidLine
    do 
        echo "kill ${pid}"
        kill ${pid}
    done