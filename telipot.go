package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/monaco-io/request"
	"github.com/udhos/equalfile"
)

var (
	tokenFlag   string
	chatFlag    string
	versionFlag bool
	version     string
)

func init() {
	flag.StringVar(&tokenFlag, "token", "", "Telegram Bot token")
	flag.StringVar(&chatFlag, "chat", "", "Telegram Chat")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.Parse()
}

func main() {
	var ipaddr string
	if versionFlag {
		fmt.Printf("\n     🤖  telipot %v\n\n", version)
		return
	}

	ipaddr = checkIP()

	if !compareIp() {
		println("New IP -->", ipaddr)
		sendMessage(tokenFlag, chatFlag, ipaddr)
		err := os.Rename("/tmp/newipaddr.txt", "/tmp/ipaddr.txt")
		if err != nil {
			println(err)
		}
	} else {
		println("Same IP -->", ipaddr)
	}
}

func checkIP() string {
	var url = `http://ip-api.com/json/?fields=query`
	c := request.Client{
		URL:    url,
		Method: "GET",
		Header: (map[string]string{
			"User-Agent": "ipchk - The IP Checker",
		}),
		Timeout:   time.Second * 20,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp := c.Send()
	re := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	match := re.FindStringSubmatch(resp.String())
	f, err := os.Create("/tmp/newipaddr.txt")
	f.WriteString(match[0])
	if err != nil {
		println(err)
	}
	return match[0]
}

func sendMessage(token, chatId, ipaddr string) string {
	var url = `https://api.telegram.org/bot` + token + `/sendMessage`
	c := request.Client{
		URL:    url,
		Method: "POST",
		Header: (map[string]string{
			"User-Agent":   "telipot - The Telegram Bot Checker",
			"Accept":       "application/json",
			"content-type": "application/json",
		}),
		Query: map[string]string{
			"text":                     ipaddr,
			"parse_mode":               "markdown",
			"disable_web_page_preview": "False",
			"disable_notification":     "True",
			"reply_to_message_id":      "None",
			"chat_id":                  chatId,
		},
		Timeout:   time.Second * 20,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp := c.Send()
	return resp.String()
}

func compareIp() bool {
	cmp := equalfile.New(nil, equalfile.Options{})
	equal, err := cmp.CompareFile("/tmp/ipaddr.txt", "/tmp/newipaddr.txt")
	if err != nil {
		println(err)
	}
	return equal
}
