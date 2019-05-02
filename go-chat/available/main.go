package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func exists(domain string) (bool, error) {
	// 接続先のURL
	const whoisServer = "com.whois-servers.net"
	// TCPコネクション張る: WHOISの仕様に基づき43番ポートに接続
	//                      -> ポート番号だけ別にしておいた方が便利
	conn, err := net.Dial("tcp", whoisServer+":43")
	if err != nil {
		return false, err
	}
	defer conn.Close()

	// ドメイン名を投げつける
	conn.Write([]byte(domain + "\r\n"))
	// レスポンスをチェックして"no match"が含まれれば未使用
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		if strings.Contains(strings.ToLower(sc.Text()), "no match") {
			return false, nil
		}
	}

	return true, nil
}

var canUse = map[bool]string{true: "x", false: "o"}

func main() {
	sc := bufio.NewScanner(os.Stdin)

	var domain string
	var exist bool
	var err error
	for sc.Scan() {
		domain = sc.Text()
		fmt.Printf("%s -> ", domain)

		if exist, err = exists(domain); err != nil {
			panic(err)
		}
		fmt.Println(canUse[exist])

		time.Sleep(1)
	}
}
