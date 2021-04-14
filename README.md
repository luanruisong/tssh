# tssh

## golang 实现的ssh 工具

### 安装 

#### 下载安装 

下载地址 [release](https://github.com/luanruisong/tssh/releases/)

windows用户请手动下载，暂时不提供一键安装模式（~~主要是批处理脚本不会写~~）

#### Mac一键安装

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/luanruisong/tssh/master/install.sh)"
```

#### homebrew 安装

对不起，我不配

![我不配](https://blog-img.luanruisong.com/blog/img/20210330204817.png)

二次冲击 Homebrew 失败

![new](https://blog-img.luanruisong.com/blog/img/20210412120135.png)

#### 2.0 抢先体验版请自行下载编译安装

***branch 2.0***

## 设置环境变量

```shell
export TSSH_HOME=/Users/user/work/ssh_config/
```

## 查看帮助

```shell
$ tssh -h

 ______   ______     ______     __  __
/\__  _\ /\  ___\   /\  ___\   /\ \_\ \
\/_/\ \/ \ \___  \  \ \___  \  \ \  __ \
   \ \_\  \/\_____\  \/\_____\  \ \_\ \_\
    \/_/   \/_____/   \/_____/   \/_/\/_/


Usage of TSSH:

  env		get env info 				(e|-e)
  version	get version info			(v|-v)
  list 		get config list				(l|-l)
  conn		connect to alias			(c|-c)
  delete 	del config by alias			(d|-d)
  add 		add config {user@host}			(a|-a)
  save 		reset config {user@host}		(s|-s)
	  -P int
			set port in (add|save) (default 22)
	  -k string
			set private_key path in (add|save)
	  -n string
			set alias name in (add|save)
	  -p string
			set password in (add|save)
```

## 相关操作

### 添加一个链接配置

#### 采用密钥模式

```shell
$ tssh -a user@host -k /Users/user/.ssh/id_rsa -n name
```

#### 采用密码模式

**密码如含有特殊字符请使用单引号**

```shell
$ tssh -a user@host -p 123456 -n pname
```

#### 覆盖一个链接配置

```shell
$ tssh -s user@host -k /Users/user/.ssh/id_rsa -n name
$ tssh -s user@host -p 123456 -n pname
```

### 查看现有链接

```shell
$ tssh -l
No              name                ip      user               auth_mode      port                 save_at
 1              name              host      user             private_key        22     2021-03-30 18:38:28
 2             pname              host      user                password        22     2021-03-30 18:38:37
```

### 删除配置

```shell
$ tssh -d name
$ tssh -d pname
```

### 链接

```shell
tssh -c name
```


## windows 实测

~~在windows下代码会出现 panic~~

![panic](https://blog-img.luanruisong.com/blog/img/20210330183152.png)

已修复windows 问题，感谢大佬提供了一个 终端跨平台解决方案

大佬项目链接 [containerd/console](https://github.com/containerd/console)

解决问题的心路历程 -> [anwu's blog](https://luanruisong.com/post/golang/tssh/)