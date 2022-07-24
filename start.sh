if [ $# == 1 ] ; then
echo "参数正确"
else
echo '参数错误'
echo '必选参数:'
echo '   1 string 类型 执行文件名'
echo './start.sh goods-api_main'
exit 1;
fi
name=$1
chmod +x ./$name
#重启，如果已经存在则关闭重启
if pgrep -x $name > /dev/null #查看一下进程有没有启动
then
  #如果已经启动 
  echo "${name} is running"
  echo "shutting down ${name}"
  #ps -a拿到进程信息，|grep 过滤我们想要的service进程信息，然后用awk将进程号拿出来
  #前面这一行代码就是拿到service进程信息，xargs就是将进程杀掉，不要用kill -9,会强杀
  #我们还要做优雅退出，注销服务。
  if ps -a | grep $name | awk '{print $1}' | xargs kill $1
    then
      echo "starting ${name}"
      # 后台启动
      ./$name > /dev/null 2>&1 &
      echo "start ${name} success"
  fi
else
 echo "starting ${srv_name}"
  ./$srv_name > /dev/null 2>&1 &
  echo "start ${srv_name} success"
fi