package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func escanear(wg *sync.WaitGroup, puertos chan int, ip string, tipoCon string) {
	for {
		puerto := <-puertos
		ipParsed, _ := net.LookupHost(ip)
		con, err := net.DialTimeout(tipoCon, fmt.Sprintf("%s:%d", ipParsed[0], puerto), time.Second)
		if err == nil {
			fmt.Printf("Puerto %d abierto.\n", puerto)
			con.Close()
		}
		wg.Done()
	}
}

func main() {
	var wg sync.WaitGroup
	var escaneres int
	var tipoCon string
	puertos := make(chan int)

	if len(os.Args) < 3 {
		fmt.Println("Uso: netscan IP PuertoMinimo-PuertoMaximo [Tipo de conexion] [SubprocesosACrear]")
		os.Exit(0)
	}

	rPuerto := strings.Split(os.Args[2], "-")
	maxPuerto, _ := strconv.Atoi(rPuerto[1])

	if len(os.Args) > 3 {
		tipoCon = os.Args[3]
	} else {
		tipoCon = "tcp"
	}
	if len(os.Args) > 4 {
		escaneres, _ = strconv.Atoi(os.Args[4])
	} else {
		escaneres = 50
	}

	for i := 0; i <= escaneres; i++ {
		go escanear(&wg, puertos, os.Args[1], tipoCon)
	}
	for i, _ := strconv.Atoi(rPuerto[0]); i <= maxPuerto; i++ {
		puertos <- i
		wg.Add(1)
	}
	wg.Wait()
}
