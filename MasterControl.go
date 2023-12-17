package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func cliente(con net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	go io.Copy(os.Stdout, con)
	for {
		scanner.Scan()
		if scanner.Text() == "" {
			continue
		} else if string([]byte(scanner.Text())[:3]) == "put" {
			con.Write([]byte(scanner.Text()))
			con1, _ := net.Listen("tcp", ":1235")
			con2, _ := con1.Accept()
			file, _ := os.Open(string([]byte(scanner.Text())[4:]))
			io.Copy(con2, file)
			//time.Sleep(1 * time.Millisecond)
			file.Close()
			con2.Close()
			con1.Close()
			continue
		} else if string([]byte(scanner.Text())[:3]) == "get" {
			con.Write([]byte(scanner.Text()))
			con1, _ := net.Listen("tcp", ":1235")
			con2, _ := con1.Accept()
			name := strings.Split(string([]byte(scanner.Text())[4:]), "/")
			file2, _ := os.Create(name[len(name)-1:][0])
			io.Copy(file2, con2)
			file2.Close()
			con2.Close()
			con1.Close()
			continue
		} else if scanner.Text() == "close" {
			con.Write([]byte(scanner.Text() + "\n"))
			con.Close()
			break
		}
		con.Write([]byte(scanner.Text() + "\n"))
	}
}

func main() {
	con, _ := net.Listen("tcp", ":1234")
	for {
		con2, _ := con.Accept()
		fmt.Println("Cliente conectado ;)")
		go cliente(con2)
	}
}
