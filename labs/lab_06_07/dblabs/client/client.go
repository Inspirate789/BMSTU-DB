package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func readRoutine(ch chan string, r io.Reader, cancel context.CancelFunc) {
	scanner := bufio.NewScanner(r)

	for {
		if scanner.Scan() {
			ch <- scanner.Text()
		} else {
			cancel()
			return
		}
	}
}

func readServer(ctx context.Context, cancel context.CancelFunc, conn net.Conn) {
	ch := make(chan string)
	go readRoutine(ch, conn, cancel)
OUTER:
	for {
		select {
		case str := <-ch:
			log.Printf("From server: %s\n", str)
		case <-ctx.Done():
			break OUTER
		}
	}
	log.Printf("Finished readServer\n")
}

func writeServer(ctx context.Context, cancel context.CancelFunc, conn net.Conn) {
	ch := make(chan string)
	go readRoutine(ch, os.Stdin, cancel)
OUTER:
	for {
		select {
		case str := <-ch:
			log.Printf("To server: %v\n", str)
			conn.Write([]byte(fmt.Sprintf("%s\n", str)))

			if str == "exit" {
				cancel()
			}
		case <-ctx.Done():
			break OUTER
		}

	}
	log.Printf("Finished writeServer")
}

func main() {
	dialer := &net.Dialer{}
	ctx, cancel := context.WithCancel(context.Background()) // context.WithTimeout(context.Background(), 5*time.Minute)

	conn, err := dialer.DialContext(ctx, "tcp", "127.0.0.1:3302")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer conn.Close()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		readServer(ctx, cancel, conn)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		writeServer(ctx, cancel, conn)
		wg.Done()
	}()

	//time.Sleep(1 * time.Minute)
	//cancel()
	wg.Wait()
}
