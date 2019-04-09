# clustersh

[![Join the chat at https://gitter.im/dqinyuan/clustersh](https://badges.gitter.im/dqinyuan/clustersh.svg)](https://gitter.im/dqinyuan/clustersh?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[中文文档](https://github.com/DQinYuan/clustersh/tree/master/zhdocs)

# Introduction

a simple command tool, which can execute shell script 
in target servers, without need to install
anything in these servers. More over, it can 
execute the script for specific Linux distribution
(like ubuntu, centos etc.) in target server.

 You can simply understand it as a "cluster version" of
 `sh` command on Linux.`clustersh shname` merely omits
 `.sh`  suffix in `shname`, and it is able to choose
 whether to execute `shname_centos.sh` or 
 `shname_ubuntu.sh` according to target server's operation
 system type. If there is no script for specific
 operation system, `shname.sh` will be executed.
 
# Scenes
 
 If you want to execute shell script together
  in plenty of
  computers (perhaps hundreds of computers) 
  which are just installed linux operation
  system, `clustersh` is very suitable for you
  

"just installed linux operation
   system" only to emphasize `clustersh` needn't
   install any thing in target servers,
   it's not necessary condition.
   
   

# Get Started

`clustersh` is very simple to use. Simply choose a Linux machine that is connected to the cluster network 
and follow the steps below.

Suppose our task is to install 
an nfs client for all machines in the cluster. There are 
centos machines and ubuntu machines in the cluster.

### Download cluster sh

Download `clustersh` binary file from
the download address(use wget as follows)
 , and then move it to 
the PATH of Linux:

```bash
wget https://github.com/DQinYuan/clustersh/releases/download/v0.1.0/clustersh
chmod a+x clustersh
mv clustersh /usr/local/bin
```

Try to execute it:

```bash
clustersh --help
```

Then you can see some help info.

### Prepare a directory

Assume the prepared directory is `~/clustershtest`:

```bash
mkdir ~/clustershtest
cd ~/clustershtest
```

### Config ips file

Create a file named ips in the directory:

```bash
touch ips
```
  
Then config ips of all the machines
in the cluster in the file, 
assuming that there are 5 machines
, 10.10.108.23,10.10.108.71, 
10.10.108.72,10.10.108.73,10.10.108.90 respectively.
So we can config as follow:

```bash
10.10.108.23
10.10.108.71-73
10.10.108.90
```

Here we use `71-73` to target a range of ip.
`clustersh` just surpport range ip in last
segment of ip address. 

default config file name is `ips`, but you
can add `--ips` param when execute `clustersh`
command to customize config file name.

### Write shell script
 
Write two scripts in the directory, as follows:

 - `nfs_centos.sh`, to install nfs-client in centos machine
 
```bash
#!/bin/sh

yum install -y  nfs-utils
``` 
 
 - `nfs_ubuntu.sh`, to install nfs-client in ubuntu machine
 
```bash
#!/bin/sh

apt install -y nfs-common
```

Before starting the next step, 
you'd better make sure that 
all the shell scripts you write
 are tested on the corresponding operating system.
 

### Execute clustersh

Final, you only need to execute command as follow:

```bash
clustersh nfs -U root -P xxxxxx
```

`clustersh` will look for `ips` in current
directory, and read ip from it. Then 
 apply username and password
provided by command line(here, username is "root"
and password is "xxxxxx") to log in them.  (
in practice, servers in a cluster often have
the same username and password, so `clustersh`
only use a pair of username and password
)

`nfs` is the **simple name** of the shell script
to be executed, it will be extended
 to `nfs_os.sh` by `clustersh`
according to target server's operation system type.
If `nfs_os.sh` not exists, `nfs.sh` will be
executed.

For example, `clustersh` log in a centos server,
discovering its os is centos. Then it will look for
`nfs_centos.sh`.If successful, `clustersh` will
execute it, else it will execute `nfs.sh`.

If there is more os type in your cluster, 
you should name specific script as:

```bash
simplename_os.sh
```
(in this case, simplename is "nfs")

`clustersh` can recognize follow os:

|Linux distribution|
|----|
|centos|
|rhel|
|aliyun|
|fedora|
|debian|
|ubuntu|
|raspbian|

You can extra provide a `simplename.sh` file
to execute when `clustersh` 
**can not recognize os type** or
**there is no specific scripts for target os**

### Check output

Although shell scripts pass tests
 on the corresponding operating system, 
 it is still possible to 
fail in some servers 
for some unpredictable reasons(no enough
space on disks, dns config err, for example).
`clustersh` will print if script is executed 
successfully on a target server in the runtime.
For a few failed machines, 
better to log in them and manually 
complete your task.

![clustersh fail](https://user-images.githubusercontent.com/23725000/55680669-3780ff00-594f-11e9-9781-e1859403af9d.png)

From the picture above, we can see `10.10.108.41`
fail to execute the shell script for some reason,
you'd better log in it and manually compete it.
Fortunately, this not often happen, so it will not
 expend too many energy.
 
[source code for this case](https://github.com/DQinYuan/clustersh/tree/master/examples/nfs)

# Another case

Suppose I am too lazy to set up a DNS server in the cluster, 
so they can only identify each other's hostnames through the local hosts file.
We need to update all the hosts files of the machines in the cluster. 
`clustersh` can help you easily complete the task.

Suppose there are only 4 machines in the cluster, 
10.10.108.91,10.10.108.92,10.10.108.93, 10.10.108.94 respectively.
(here are only four machines for the sake of convenience. In reality,
 there may be hundreds of machines.)
 
[source code for this case](https://github.com/DQinYuan/clustersh/tree/master/examples/unihosts)

Steps as follows.

### Create a directory

Create a work directory

```bash
mkdir unihosts
cd unihosts
```

### Edit ips file

Edit ips file in this directory:

```bash
10.10.108.91-94
```

### Edit uniform hosts file

The uniform hosts file in the cluster, I call it "unihosts":

```bash
touch unihosts
```

edit it as follows:

```bash
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
10.10.108.91 h91
10.10.108.92 h92
10.10.108.93 h93
10.10.108.94 h94
```

### Write shell script

The main function of the script is to overwrite `/etc/hosts`
with `unihosts`.

`unihosts.sh`:

```bash
#!/bin/sh

\cp -f ./unihosts /etc/hosts
```

The script is common to all operating system,
so we needn't to provide os specific scrips like last case.

### Execute clustersh

```bash
clustersh unihosts -U root -P xxxxxx
```

The task is relatively simple, so execution is very fast 
even if there are hundreds of machines.

### Summary

This case is to illustrate that you can use any file 
in the current directory and subdirectories in your shell script, 
because all files in the current directory and subdirectories 
will be sent to the target servers by `clustersh`.

![clustersh summary](https://user-images.githubusercontent.com/23725000/55685568-8304ce80-598a-11e9-9755-879005bab0a3.png)

# Params introduction

You can get all available param and introduction by  `clustersh --help`

| Full name        | Simeple name    | Mean   | Default| 
| --------   | -----   | ---- |---- |
|  --username       | -U      |   username for ssh log in    |root|
| --password        | -P      |   password for ssh log in    | root|
| --ips        | -I      |   config file for machines' ip in cluster    | ips|
|--timeout|-T| ssh connection timeout, for example "10s"| 10s|
|--verbose|-V| print all shell output in cluster, perhaps can help you debug your shell. To open verbose, you only need to add this flag(needn't param) |Close|

# Principle

`clustersh` is only a simple abstraction on ssh.It will first 
send all files in the working directory and its  subdirectories
to the target server, After which it chooses a suitable script 
to be executed.

`clustersh` will connect to machines with the same number of CPU cores
 at the same time.
