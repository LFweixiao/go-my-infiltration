package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

// TCP 显示服务
func main() {

}

// 最佳优化 copy
func conpyL(conn net.Conn) {
	defer conn.Close()
	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatalln("意料之外的读或写错误")
	}
}

//  bufio 包含 Reader 和 Writer
func bufioL(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	// 分隔符 代表读取长度
	s, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("未知的读错误")
	}
	log.Println("读 %d bytes: %s", len(s), s)

	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString(s); err != nil {
		log.Fatalln("未出的写错误")
	}
	writer.Flush()
}

//  direct 直接使用 conn来编写
func direct() {
	listen, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("端口绑定失败")
	}
	log.Println("监听端口 20080")

	for {
		accept, err := listen.Accept()
		if err != nil {
			log.Fatalln("入口无法关联")
		}
		log.Println("关联 入口")
		go echo(accept)
	}
}

// echo  conn是一个处理函数 它仅回显接收到的数据
func echo(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 512)
	for {
		// 读取数据到缓冲区
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("链接 断开")
			break
		}
		if err != nil {
			log.Println("意料之外的异常")
			break
		}
		log.Println("接收 %d bytes : %s \n", size, string(b))

		// 通过conn.write 发送数据
		log.Println("写数据")
		if _, err := conn.Write(b[0:size]); err != nil {
			log.Fatalln("意料之外的写数据")
		}
	}
}
