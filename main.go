package main

import (
	"flag"
	"fmt"

	"tssh/cmd"
	"tssh/store"
)

var (
	a = flag.String("a", "", "add config {user@host}")
	s = flag.String("s", "", "set config {user@host}")
	d = flag.String("d", "", "del config {name}")
	c = flag.String("c", "", "connect config host {name}")
	l = flag.Bool("l", false, "config list")
	e = flag.Bool("e", false, "evn info")
	v = flag.Bool("v", false, "app version")
)

func main() {

	flag.Parse()
	//查看信息
	switch {
	case *e:
		store.Env()
		return
	case *v:
		fmt.Println("version", "1.1.0")
		return
	}
	//检测环境变量
	if err := store.DefaultCheck(); err != nil {
		fmt.Println(err)
		return
	}
	//删除配置操作
	if len(*d) > 0 {
		err := store.Del(*d)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	//添加与覆盖
	if cmd.AddOrSave(a, s) {
		return
	}
	//查看与链接
	if cmd.ListAndConn(l, c) {
		return
	}

	flag.PrintDefaults()

}
