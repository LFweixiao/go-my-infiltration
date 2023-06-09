package main

import (
	"fmt"
	"net"
	"sync"
)

// main TCP 扫描器
func main() {
	conWithScanWorker1024()
}

// conWithScanWorker1024 线程池管理并发
func conWithScanWorker1024() {
	ports := make(chan int, 100)
	var wg sync.WaitGroup
	for i := 0; i < cap(ports); i++ {
		go worker(ports, &wg) // ports 是通信通道
	}
	for i := 1; i <= 1024; i++ {
		wg.Add(1)
		ports <- i
	}
	wg.Wait()
	close(ports)
}

// 工作区 线程池
func worker(ports chan int, wg *sync.WaitGroup) {
	for p := range ports {
		address := fmt.Sprintf("127.0.0.1:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
		} else {
			fmt.Printf("打开端口：%d \n", p) // 此线程阻塞等待通道数据
			conn.Close()
		}
		wg.Done()
	}
}

// conWithScan1024 并发等待扫描完成
// WaitGroup 计数器来等待多线程结束
func conWithScan1024() {
	// WaitGroup 通过计数的方式在归 0 前阻止线程的结束
	var wg sync.WaitGroup
	for i := 1; i <= 1024; i++ {
		wg.Add(1)
		go func(number int) {
			defer wg.Done()
			//address := fmt.Sprintf("scanme.nmap.org:%d", number)
			address := fmt.Sprintf("127.0.0.1:%d", number)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				// 端口关闭或过滤
				return
			}
			conn.Close()
			fmt.Printf("%d 打开端口", number)
		}(i)
	}
	wg.Wait()
}

// conScan1024 并发扫描1024
// 主线程过早完成推出
func conScan1024() {
	for i := 1; i <= 1024; i++ {
		go func(number int) {
			address := fmt.Sprintf("scanme.nmap.org:%d", number)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				// 端口关闭或过滤
				return
			}
			conn.Close()
			fmt.Printf("%d 打开端口", number)
		}(i)
	}
}

// scan1024 扫描1024个端口
// 单线程很慢
func scan1024() {
	for i := 1; i <= 1024; i++ {
		address := fmt.Sprintf("scanme.nmap.org:%d", i)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			// 端口关闭或过滤
			continue
		}
		conn.Close()
		fmt.Printf("%d 打开端口", i)
	}
}

// scanOne 扫描一个端口
func scanOne() {
	_, err := net.Dial("tcp", "scanme.nmap.org:80")
	if err == nil {
		fmt.Println("链接成功")
	} else {
		fmt.Println(err)
	}
}
