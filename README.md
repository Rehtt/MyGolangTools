# MyGolangTools

## FileServer

一个简单的文件服务

<img src="https://github.com/Rehtt/MyGolangTools/blob/master/img/fileServer.png?raw=true">

### Logs:

每天自动将下载次数高于平均值的文件移动到A盘，将低于平均值的文件移动到B盘。
  
显示两个盘的文件及下载。


## GetLiveBilibiliVideoURL

获取哔哩哔哩直播视频的地址


### shell ffmpeg自动侦测下载
指定时间循环侦测地址并自动下载

将live.sh与live.go编译的文件放在同一目录
<code>./live.sh "直播地址"</code>

## SimpleFileServer
简单的一个文件服务器。
默认监听地址：0.0.0.0:8080
<pre>
-host 监听地址
-port 监听端口
-path 文件地址
</pre>
