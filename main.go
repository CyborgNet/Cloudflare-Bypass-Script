package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println(os.Args[0] + " [타겟] [시간] [슬롯] (프록시 리스트 - 클플 x)")
		return
	}
	loader := new(Loader)
	loader.target = os.Args[1]
	duration, _ := strconv.Atoi(os.Args[2])
	loader.duration = time.Second * time.Duration(duration)
	slots, _ := strconv.Atoi(os.Args[3])
	loader.slots = slots
	var proxies []string
	if len(os.Args) > 4 {
		file, err := os.Open(os.Args[4])
		if err != nil {
			fmt.Printf("failed opening file: %s", err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			proxies = append(proxies, scanner.Text())
		}
	}
	loader.proxylist = proxies
	loader.Run()
}
