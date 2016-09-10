package goWebChatServer

import (
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func cleanupName(oldName string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-',
			r == '_':
			return r
		}
		return -1
	}, oldName)
}

func getIP(req *http.Request) string {

	ipSlice, ok := req.Header["X-Real-Ip"]
	var ip string
	if !ok {
		ip = req.RemoteAddr
	} else {
		ip = ipSlice[0]
	}
	index := strings.Index(ip, ":")
	if index != -1 {
		ip = ip[0:index]
	}

	return ip
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
