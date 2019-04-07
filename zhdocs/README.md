# clustersh

# 简介

一个简单的命令行工具，用于指定的集群上执行shell脚本，
不需要提前预装任何东西。并且可以根据
集群中每台计算机的操作系统类型选择相应的脚本进行执行。

可以将其简单地理解为`sh`命令的集群版本。只不过
`clustersh shname`中`shname`省略了`.sh`后缀名，
并且可以根据集群中不同的操作系统类型自动选择执行
`shname_centos.sh`还是执行`shname_ubuntu.sh`，
如果没有该操作系统的执行文件则直接执行`shname.sh`文件。

# 场景

基本上所有的集群管理工具都需要提前在集群每台电脑上安装相应的软件，
虽然这些工具功能强大且使用方便，
但是集群规模很大且操作系统类型不一的话，这个安装过程就会很麻烦，
就好像早上去上班的路上，要先走500m才能到达地铁站
，`clustersh`主要就是用于解决集群搭建的这"500m"的问题。
所以`clustersh`尽可能地简单易用。

当然，如果你只是希望在指定的集群中批量执行shell脚本，这个工具也很适合你。
工具支持直接指定一个ip范围，能让配置过程轻松不少。

# Get Started

`clustersh`的使用非常简单。只需任选一台与集群网络联通的linux机器
，在其上按照如下步骤操作.

假设我们的任务是给集群内所有机器统一安装一个nfs客户端，集群内有centos机器和
ubuntu机器

#### 下载clustersh

去[下载地址]()下载clustersh的二进制文件，
然后将其移动到linux的PATH路径下：

```bash
mv clustersh /usr/local/bin
```

尝试运行一下命令：

```bash
clustersh --help
```

可以看到相关的帮助信息


#### 准备一个文件夹

之后准备一个文件夹(假设是`~/clustershtest`)：

```bash
mkdir ~/clustershtest
cd ~/clustershtest
```


#### 配置ips

在文件夹下创建一个名叫ips的文件:

```bash
touch ips
```

然后在里面配置上集群中所有机器的ip，
假设我的集群中有5台机器，分别是10.10.108.23,10.10.108.71,
10.10.108.72,10.10.108.73,10.10.108.90。
于是我们可以如下配置ips：

```bash
10.10.108.23
10.10.108.71-73
10.10.108.90
```

这里我们使用了`71-73`直接指定了一个范围的ip来简化配置，`clustersh`
目前只支持在ip地址的第四段使用范围指定。

默认情况下配置文件名叫做ips，如果你不想让它叫做ips的花，可以在后面
执行`clustersh`命令是使用`--ips`指定。

#### 编写shell脚本

在文件夹下写如下两个脚本：

 - `nfs_centos.sh`,用于在centos机器上安装nfs-client

```bash
#!/bin/sh

yum install -y  nfs-utils
```

 - `nfs_ubuntu.sh`,用于在ubuntu上安装nfs-ubuntu
 
```bash
#!/bin/sh

apt install -y nfs-common
```

在开始下一步之前，你最好确保你写的
所有shell脚本在对应操作系统上都测试通过。

#### 执行clustersh

最后在文件夹下执行如下命令即可：

```bash
clustersh nfs -U root -P xxxxxx
```

`clustersh`会寻找当前目录下的`ips`文件，将其中
的ip地址读出，依次使用命令行提供的用户名和密码
（这里的用户名为`root`，密码为`xxxxxx`）登陆
这些ip。（在实践中，集群大多有统一的用户名和密码，
所以这里就使用统一的用户名与密码登陆集群了）

`nfs`是shell脚本的**简称**，`clustersh`会根据服务器
的操作系统类型将其扩充为`nfs_操作系统类型.sh`，
如果`nfs_操作系统类型.sh`文件不存在的话则扩充为`nfs.sh`.

比如`clustersh`登陆到一台centos服务器后，发现
操作系统是centos，于是就会尝试寻找`nfs_centos.sh`，
如果有的话就执行它，没有的话则执行`nfs.sh`

如果在集群中还有更多的操作系统类型，请以如下格式命名脚本：

```bash
简称_操作系统类型.sh
```

`clustersh`当前支持识别的操作系统类型有：

|操作系统类型|
|----|
|centos|
|rhel|
|aliyun|
|fedora|
|debian|
|ubuntu|
|raspbian|

你也可以再提供一个`简称.sh`用于在**操作系统类型无法识别**或者
是**没有提供针对该种操作系统的脚本**时执行。

如果你写的脚本对所有操作系统都通用的话，你直接给一个`简称.sh`即可。

#### 查看输出

虽然shell脚本在相应的操作系统上都测试通过，
但是在集群中运行时还是有可能因为一些
莫名奇妙的原因（比如磁盘空间不足，DNS配置错误等等）
失败，`clustersh`在运行时会打印每台机器运行的成败情况，
对于少数失败的机器，最好手动登陆上去完成操作。


![clustersh fail](https://user-images.githubusercontent.com/23725000/55680669-3780ff00-594f-11e9-9781-e1859403af9d.png)


比如从上面的输出中看到`10.10.108.41`因为某些原因没能成功
执行脚本，最好手工登陆上去操作，不过这种情况属于少数，
并不会花费太多的精力。


[案例源代码](https://github.com/DQinYuan/clustersh/tree/master/examples/nfs)


# 另一个简单的案例

假设我懒得在集群中搭建的DNS服务器，我希望他们互相之间
仅仅通过本地的hosts文件来互相识别主机名，这个时候就需要统一更新
集群中机器的hosts文件，使用`clustersh`可以轻松完成任务。

假设集群中只有四台机器：10.10.108.91,10.10.108.92,10.10.108.93,
10.10.108.94

[案例源代码]()

#### 新建文件夹

新建一个专门的工作目录：

```bash
mkdir unihosts
cd unihosts
```

#### 编辑ips文件

在文件夹下编辑`ips`文件，令其内容为：

```bash
10.10.108.91-94
```

#### 编写统一的hosts文件

集群中统一的hosts文件，我称之为`unihosts`

```bash
touch unihosts
```

编辑其内容为：

```bash
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
10.10.108.91 h91
10.10.108.92 h92
10.10.108.93 h93
10.10.108.94 h94
```

#### 编写shell脚本

shell脚本的主要功能就是使用`unihosts`覆盖掉
集群服务器上的`/etc/hosts`文件

编辑`unihosts.sh`如下：

```bash
#!/bin/sh

\cp -f ./unihosts /etc/hosts
```

这个脚本对所有操作系统通用，所以我不需要
像之前的案例一样给出针对操作系统的脚本了。


#### 执行clustersh

```bash
clustersh unihosts -U root -P xxxxxx
```

这个任务比较简单，应该很快就能执行完毕。


#### 总结

这个案例就是想说明你的shell脚本里是可以使用
当前目录及子目录中的任意文件的，因为当前目录
及子目录的所有文件都会被我发送到集群中去。

![clustersh summary](https://user-images.githubusercontent.com/23725000/55681918-c8130b80-595e-11e9-8123-bfc554844551.png)


# 参数介绍

通过`cluster --help`参数即可查看所有的参数及其含义

| 全称        | 简写    |  含义  | 默认值| 
| --------   | -----   | ---- |---- |
|  --username       | -U      |   用于登陆集群服务器的用户名    |root|
| --password        | -P      |   用于登陆集群服务器的密码    | root|
| --ips        | -I      |   用于指定集群ip配置文件    | ips|
|--timeout|-T| 用于指定ssh连接的超时时间，格式为数字加单位，比如"10s"| 10s|
|--verbose|-V|尽可能地打印信息,包括shell脚本在集群执行时的全部输出，执行命令时加上该标志（不需要参数）即表示开启|不开启|









