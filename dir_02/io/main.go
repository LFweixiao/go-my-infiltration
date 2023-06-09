package main

import (
	"fmt"
	"log"
	"os"
)

// 标准输入stdin 读取io.Reader
type FooReader struct{}

// 标准输入  stdin 读取数据
func (FooReader *FooReader) Read(b []byte) (int, error) {
	fmt.Print("in > ")
	return os.Stdin.Read(b)
}

type FooWriter struct {
}

func (FooReader *FooWriter) Writer(b []byte) (int, error) {
	fmt.Print("out > ")
	return os.Stdin.Write(b)
}

// io演示
func main() {

	// 实例化 reader writer
	var (
		reader FooReader
		writer FooWriter
	)

	// 创建缓冲区
	input := make([]byte, 2048)

	// 读取缓冲
	s, err := reader.Read(input)
	if err != nil {
		log.Fatalf("读取失败")
	}
	fmt.Printf(" 读取数据 %d 结尾 \n", s)

	s, err = writer.Writer(input)
	if err != nil {
		log.Fatalf("写入失败")
	}
	fmt.Printf(" 写入数据 %d 结尾 \n", s)
}
