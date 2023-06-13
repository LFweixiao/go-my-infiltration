package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"os"
	"text/tabwriter"
)

// 完整的域猜测
// 命令行 time ./subdomain_guesser -domain microsoft.com -wordlist namelist.txt -c 1000

func main() {
	var (
		// 需要执行猜测的域
		flDomain      = flag.String("domain", "", "需要执行猜测的域")
		flWordlist    = flag.String("wordlist", "", "需要查询的单词列表")
		flWorkerCount = flag.Int("c", 100, "使用的工人数量")
		flServerAddr  = flag.String("server", "8.8.8.8:53", "DNS服务")
	)
	flag.Parse()

	if *flDomain == "" || *flWordlist == "" {
		fmt.Println("-domain and -wordlist are required")
		os.Exit(1)
	}

	var results []result
	fqdns := make(chan string, *flWorkerCount)
	gather := make(chan []result)
	tracker := make(chan empty)

	fh, err := os.Open(*flWordlist)
	if err != nil {
		panic(err)
	}

	defer fh.Close()
	scanner := bufio.NewScanner(fh)

	for i := 0; i < *flWorkerCount; i++ {
		go worker(tracker, fqdns, gather, *flServerAddr)
	}
	for scanner.Scan() {
		fqdns <- fmt.Sprintf("%s.%", scanner.Text(), *flDomain)
	}
	// 可以在此处坚持 scanner.Err()
	go func() {
		for r := range gather {
			results = append(results, r...)
		}
		var e empty
		tracker <- e
	}()

	close(fqdns)
	for i := 0; i < *flWorkerCount; i++ {
		<-tracker
	}
	close(gather)
	<-tracker

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\n", r.Hostname, r.IPAddress)
	}
	w.Flush()
}

func lookupA(fqdn, serverAddr string) ([]string, error) {
	var m dns.Msg
	var ips []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	in, err := dns.Exchange(&m, serverAddr)
	if err != nil {
		return ips, err
	}
	if len(in.Answer) < 1 {
		return ips, errors.New("没有答复")
	}
	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
			return ips, nil
		}
	}
	return ips, nil
}

func lookupCNAME(fdqn, serverAddr string) ([]string, error) {
	var m dns.Msg
	var fqdns []string
	m.SetQuestion(dns.Fqdn(fdqn), dns.TypeCNAME)
	in, err := dns.Exchange(&m, serverAddr)
	if err != nil {
		return fqdns, err
	}
	if len(in.Answer) < 1 {
		return fqdns, errors.New("没有答复")
	}
	for _, answer := range in.Answer {
		if c, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, c.Target)
		}
	}
	return fqdns, nil
}

func lookup(fqdn, serverAddr string) []result {
	var results []result
	var cfqdn = fqdn // 不要修改原始信息
	for {
		cnames, err := lookupCNAME(cfqdn, serverAddr)
		if err == nil && len(cnames) > 0 {
			cfqdn = cnames[0]
			continue // 我们必须处理下一个CNAME
		}
		ips, err := lookupA(cfqdn, serverAddr)
		if err != nil {
			break //该主机名称没有A记录
		}
		for _, ip := range ips {
			results = append(results, result{IPAddress: ip, Hostname: fqdn})
		}
		break // 我们以及处理所有的结果
	}
	return results
}

func worker(tracker chan empty, fqdns chan string, gather chan []result, serverAddr string) {
	for fqdn := range fqdns {
		results := lookup(fqdn, serverAddr)
		if len(results) > 0 {
			gather <- results
		}
	}
	var e empty
	tracker <- e
}

type empty struct {
}
type result struct {
	IPAddress string
	Hostname  string
}
