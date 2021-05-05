# tssh

## golang 实现的ssh 工具

### 安装 

#### 下载安装 

下载地址 [release](https://github.com/luanruisong/tssh/releases/)

#### homebrew 安装

```shell
$ brew install tssh
```

## 环境变量

### 手动设置
```shell
export TSSH_HOME=/Users/user/work/ssh_config/
```
### 默认设置
```shell
# 默认设置在windows环境下使用%HOMEPATH%
export TSSH_HOME=$HOME/.tssh/config
```

## 查看帮助

![help](https://blog-img.luanruisong.com/blog/img/20210414135853.gif)

## 相关操作

### 添加一个链接配置

#### 采用密码模式

![add](https://blog-img.luanruisong.com/blog/img/20210414140115.gif)

#### 指定更多参数

![addmore](https://blog-img.luanruisong.com/blog/img/20210414140311.gif)

### 查看现有链接（2.0）

![list](https://blog-img.luanruisong.com/blog/img/20210414140709.gif)

### 删除配置

![del](https://blog-img.luanruisong.com/blog/img/20210414140941.gif)

## 答谢

### 跨平台终端解决方案

主要解决win下获取终端信息

大佬项目链接 [containerd/console](https://github.com/containerd/console)

### 更加友好的交互

2.0 引入了一个有意思的新包 让我们有更加友好的交互方式

大佬项目链接 [manifoldco/promptui](https://github.com/manifoldco/promptui)

## 其他

解决问题的心路历程 -> [anwu's blog](https://luanruisong.com/post/golang/tssh/)