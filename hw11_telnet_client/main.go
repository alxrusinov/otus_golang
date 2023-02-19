package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout for request")

	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatalln("host and port is required")
	}

	args := flag.Args()

	host := args[0]
	port := args[1]

	address := net.JoinHostPort(host, port)

	tcpClient := NewTelnetClient(address,
		timeout,
		os.Stdin,
		os.Stdout)

	if err := tcpClient.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer tcpClient.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := tcpClient.Send()

		if errors.Is(err, io.EOF) {
			os.Stderr.WriteString(fmt.Sprintf("...%s\n", err))
			tcpClient.Close()
			return
		}

		if err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := tcpClient.Receive()

		if errors.Is(err, io.EOF) {
			fmt.Fprintln(os.Stdout, "Bye-bye")
			os.Stderr.WriteString(fmt.Sprintf("%s\n", ErrPeerClose))
			tcpClient.Close()
			return
		}

		if err != nil {
			log.Fatalln(err)
		}
	}()

	wg.Wait()
}
