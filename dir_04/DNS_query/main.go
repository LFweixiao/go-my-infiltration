package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"os"
)

// 检索地址
func main() {
	queryNDSGuesser()
}

// queryNDSGuesser 猜测子域
func queryNDSGuesser() {
	var (
		// 需要执行猜测的域
		flDomain      = flag.String("domain", "", "需要执行猜测的域")
		flWordlist    = flag.String("wordlist", "", "需要查询的单词列表")
		flWorkerCount = flag.Int("c", 100, "使用的工人数量")
		flServerAddr  = flag.String("server", "8.8.8.8:53", "DNS服务")
	)
	flag.Parse()
	if *flDomain == "" || *flWordlist == "" {
		fmt.Println(" -domain and -flWordList 是必须的参数")
		os.Exit(1)
	}
	fmt.Println(*flWorkerCount, *flServerAddr)
}

type result struct {
	IPAddress string
	Hostname  string
}

// queryNDS 查询域名的的 NDS
func queryNDS() {
	str := "stacktitan.com"
	str = "lfweixiao.cn"
	var msg dns.Msg
	fqdn := dns.Fqdn(str)
	msg.SetQuestion(fqdn, dns.TypeA)
	in, err := dns.Exchange(&msg, "8.8.8.8:53")
	if err != nil {
		// 可以查看堆栈来定位信息
		panic(err)
	}
	if len(in.Answer) < 1 {
		fmt.Println("没有记录 一般是域名无法解析")
	}
	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			fmt.Println(a)
		}
	}
}
