# tssh

## golang 实现的ssh 工具

### 安装 

#### 下载安装 

[下载tssh](https://github.com/luanruisong/tssh/releases/download/v1.0.0/tssh)

#### 一键安装

```shell
sudo wget -O /usr/local/bin/tssh https://github.com/luanruisong/tssh/releases/download/v1.0.0/tssh
```

## 设置环境变量

```shell
    export TSSH_HOME=/Users/user/work/ssh_config/
```

## 查看帮助

```shell
$ tssh -h
Usage of tssh:
  -P int
    	set port in (-a|-s) (default 22)
  -a string
    	add config {user@host}
  -c string
    	connect config host {name}
  -d string
    	del config {name}
  -e	evn info
  -k string
    	set private_key path in (-a|-s)
  -l	config list
  -n string
    	set name in (-a|-s)
  -p string
    	set password in (-a|-s)
  -s string
    	set config {user@host}
```

## 相关操作

### 添加一个链接配置

#### 采用密钥模式

```shell
$ tssh -a user@host -k /Users/user/.ssh/id_rsa -n name
```

#### 采用密码模式

```shell
$ tssh -a user@host -p 123456 -n pname
```

#### 覆盖一个链接配置

```shell
$ tssh -a user@host -k /Users/user/.ssh/id_rsa -n name
$ tssh -s user@host -p 123456 -n pname
```

### 查看现有链接

```shell
$ tssh -l
No              name                ip      user                pwd                    key_path      port                 save_at
 1              name              host      user                        /Users/user/.ssh/id_rsa        22     2021-03-30 18:38:28
 2             pname              host      user             123456                                    22     2021-03-30 18:38:37
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

在windows下 这行代码会出现 panic

```go
    termWidth, termHeight, err := terminal.GetSize(fd)
    if err != nil {
        panic(err)
    }
```

翻阅了很多文档，目前还是无法解决,目前已修复为fmt打印，看起来舒服了点

![panic](https://blog-img.luanruisong.com/blog/img/20210330183152.png)