package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", os.Getenv("SERVER_ADDR"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader := bufio.NewReader(conn)
	prefix, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	prefix = strings.TrimRight(prefix, "\r\n")

	fmt.Printf("prefix: %s", string(prefix))
	answer := findAnswer(string(prefix))
	fmt.Printf("answer: %s\n", answer)
	_, err = conn.Write([]byte(answer + "\n"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	quote, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("quote: %s\n", quote)
}

func findAnswer(prefix string) string {
	startTime := time.Now()
	n := uint64(0)
	var s string
	for {
		data := make([]byte, 8)
		binary.LittleEndian.PutUint64(data, uint64(n))

		s = fmt.Sprintf("%s:%s", prefix, base64.StdEncoding.EncodeToString(data))

		sum := sha256.Sum256([]byte(s))
		if sum[0] == 0 && sum[1] == 0 && (sum[2]&0xf0) == 0 {
			fmt.Printf("%x\n", sum)
			fmt.Printf("n: %d\n", n)
			break
		}
		n++
	}
	finishTime := time.Now()
	fmt.Printf("Answer found in %v\n", finishTime.Sub(startTime))
	return s
}
