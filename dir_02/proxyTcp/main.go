package main

import (
	"io"
	"log"
	"net"
)

// 代理服务请求转发
func main() {
	listen, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalln("意料之外的错误，无法绑定的端口")
	}
	for {
		accept, err := listen.Accept()
		if err != nil {
			log.Fatalln("意料之外的错误，无法开启入口")
		}
		go handle(accept)
	}

}

func handle(src net.Conn) {
	dst, err := net.Dial("tcp", "lfweixiao.cn:80")
	defer dst.Close()
	if err != nil {
		log.Fatalf("意料之外的错误，无法绑定端口")
	}
	go func() {
		// 请求复制到目标
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()
	// 目标到输出复制回源
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}
}
