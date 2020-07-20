package main

import (
	"encoding/base64"
	"fmt"
	"os/exec"
)

func getAnswer(pageurl, userAgent string) string {
	answer, err := exec.Command("node", "solver.js", pageurl, base64.StdEncoding.EncodeToString([]byte(userAgent))).Output()
	if err != nil {
		fmt.Println(err)
	}
	return string(answer)
}
