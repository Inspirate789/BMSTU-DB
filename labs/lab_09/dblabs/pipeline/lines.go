package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type serverData struct {
	conn net.Conn
	id   int
	text string
}

func connectDB(_ interface{}) (interface{}, error) {
	dialer := &net.Dialer{}
	connect, err := dialer.Dial("tcp", "127.0.0.1:3302")

	if err == nil {
		scanner := bufio.NewScanner(connect)
		for i := 0; i < 4; i++ {
			scanner.Scan()
		}
	}

	return serverData{conn: connect}, err
}

func insertTerm(arg interface{}) (interface{}, error) {
	sd, ok := arg.(serverData)
	if !ok {
		return sd, errors.New("invalid type conversion (serverData)")
	}

	_, err := sd.conn.Write([]byte("0 1\n"))
	if err != nil {
		return sd, err
	}

	scanner := bufio.NewScanner(sd.conn)
	scanner.Scan()
	scanner.Scan()
	msg := scanner.Text()
	if !strings.Contains(msg, " successfully") {
		return sd, errors.New(msg)
	}
	scanner.Scan()
	sd.id, err = strconv.Atoi(scanner.Text())

	return sd, err
}

func selectTerm(arg interface{}) (interface{}, error) {
	sd, ok := arg.(serverData)
	if !ok {
		return nil, errors.New("invalid type conversion (serverData)")
	}

	_, err := sd.conn.Write([]byte(fmt.Sprintf("0 2 %d\n", sd.id)))
	if err != nil {
		return sd, err
	}

	scanner := bufio.NewScanner(sd.conn)
	scanner.Scan()
	scanner.Scan()
	msg := scanner.Text()
	if !strings.Contains(msg, " successfully") {
		return sd, errors.New(msg)
	}
	scanner.Scan()
	sd.text = scanner.Text()

	return sd, nil
}

func deleteTerm(arg interface{}) (interface{}, error) {
	sd, ok := arg.(serverData)
	if !ok {
		return nil, errors.New("invalid type conversion (serverData)")
	}

	_, err := sd.conn.Write([]byte(fmt.Sprintf("0 3 %d\n", sd.id)))
	if err != nil {
		return sd, err
	}

	scanner := bufio.NewScanner(sd.conn)
	scanner.Scan()
	scanner.Scan()
	msg := scanner.Text()
	if !strings.Contains(msg, " successfully") {
		return sd, errors.New(msg)
	}
	scanner.Scan()
	text := scanner.Text()

	if text != sd.text {
		return sd, errors.New("incorrect delete")
	}

	return sd, nil
}

func disconnectDB(arg interface{}) (interface{}, error) {
	sd, ok := arg.(serverData)
	if !ok {
		return nil, errors.New("invalid type conversion")
	}

	return nil, sd.conn.Close()
}
