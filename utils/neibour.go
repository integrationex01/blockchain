package utils

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

func IsFoundHost(host string, port uint16) bool {
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		fmt.Printf("%s %v\n", address, err)
		return false
	}
	conn.Close()
	return true
}

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func FindNeibours(myHost string, myPort, startPort, endPort uint16, startIP, endIp uint8) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)

	m := PATTERN.FindStringSubmatch(myHost)
	// log.Println(m)
	if m == nil {
		return []string{}
	}
	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])
	neibours := make([]string, 0)
	for port := startPort; port <= endPort; port++ {
		for ip := startIP; ip <= endIp; ip++ {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address && IsFoundHost(guessHost, port) {
				neibours = append(neibours, guessTarget)
			}
		}
	}
	return neibours

}

func GetHost() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "127.0.0.1"
	}
	// log.Printf("Hostname: %s\n", hostname)
	addrs, err := net.LookupHost(hostname)
	if err != nil || len(addrs) == 0 {
		return "127.0.0.1"
	}
	// log.Printf("Addrs: %v\n", addrs)
	return addrs[len(addrs)-2]
}
