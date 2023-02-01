package main

import (
	"bufio"
	"dblabs/database"
	"dblabs/server/handlers"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO: create good UI
func sendGreeting(conn net.Conn) error {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Welcome to %s, friend from %s\n", conn.LocalAddr(), conn.RemoteAddr()))
	builder.WriteString("Commands: \n")
	builder.WriteString("menu - get menu\n")
	builder.WriteString("exit - quit from server\n")
	_, err := conn.Write([]byte(builder.String()))

	return err
}

func sendMenu(conn net.Conn) error {
	builder := strings.Builder{}
	builder.WriteString("Menu: \n")

	labNumbers := make([]int, 0, len(handlers.Handlers))
	for n := range handlers.Handlers {
		labNumbers = append(labNumbers, n)
	}
	sort.Ints(labNumbers)

	for _, labNumber := range labNumbers {
		builder.WriteString(fmt.Sprintf("lab %d: \n", labNumber))
		for i, curHandler := range handlers.Handlers[labNumber] {
			builder.WriteString(fmt.Sprintf("\t%d - ", i+1))
			builder.WriteString(curHandler.Title())
			builder.WriteRune('\n')
		}
	}
	builder.WriteString("menu - get menu\n")
	builder.WriteString("exit - quit from server\n")
	builder.WriteString("\nTo execute the command, enter the lab number, the query number, and the arguments required by the query, separated by a space\n")
	_, err := conn.Write([]byte(builder.String()))

	return err
}

func commandProcessing(db *database.Database, cmd string) (res string, err error) {
	defer func() {
		if recErr := recover(); recErr != nil {
			res = ""
			err = errors.New(fmt.Sprintf("work failed with panic: %v\n", recErr))
			return
		}
	}()
	tokens := strings.Split(cmd, " ")

	labNumber, err := strconv.Atoi(tokens[0])
	if err != nil {
		return "", errors.New("unrecognized lab")
	}

	if _, inMap := handlers.Handlers[labNumber]; !inMap {
		return "", errors.New("incorrect lab number")
	}

	cmdNumber, err := strconv.Atoi(tokens[1])
	if err != nil {
		return "", errors.New("unrecognized command")
	}

	if cmdNumber < 1 || cmdNumber > len(handlers.Handlers[labNumber]) {
		return "", errors.New("incorrect command number")
	}

	return handlers.Handlers[labNumber][cmdNumber-1].Execute(db, tokens[2:])
}

func reply(conn net.Conn, msg string, db *database.Database) error {
	_, err := conn.Write([]byte("Request received, processing in progress\n"))
	if err != nil {
		return errors.New(fmt.Sprintf("cannot write to connection with %s: %v", conn.RemoteAddr(), err))
	}

	startTime := time.Now().UnixNano() / 1e6
	res, err := commandProcessing(db, msg)
	endTime := time.Now().UnixNano() / 1e6

	if err != nil {
		log.Printf(fmt.Sprintf("Error while processing request from %s: %v\n", conn.RemoteAddr(), err))
		_, err = conn.Write([]byte(fmt.Sprintf("%v\n", err)))
		if err != nil {
			return errors.New(fmt.Sprintf("cannot write to connection with %s: %v", conn.RemoteAddr(), err))
		}
	} else {
		log.Printf("Request from %s processed (%d ms), response sent\n", conn.RemoteAddr(), endTime-startTime)
		_, err = conn.Write([]byte("Request was processed successfully. Result: \n" + res + "\n"))
		if err != nil {
			return errors.New(fmt.Sprintf("cannot write to connection with %s: %v", conn.RemoteAddr(), err))
		}
	}

	return nil
}

func messageRoutine(conn net.Conn, db *database.Database) error {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Printf("Received from %s: %s\n", conn.RemoteAddr(), text)
		if text == "menu" {
			err := sendMenu(conn)
			if err != nil {
				return err
			}
			continue
		} else if text == "exit" {
			break
		}

		err := reply(conn, text, db)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.New(fmt.Sprintf("cannot scan connection with %s: %v", conn.RemoteAddr(), err))
	}

	return nil
}

func handleConnection(conn net.Conn, db *database.Database) {
	defer conn.Close()

	log.Printf("Connected to %s\n", conn.RemoteAddr())

	err := sendGreeting(conn)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	} else {
		err = messageRoutine(conn, db)
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
	}

	log.Printf("Closing connection with %s\n", conn.RemoteAddr())
}

/*
nc 127.0.0.1 3302
*/

// TODO: add input commands (at least for normal server termination, not by ctrl+c)
func main() {
	data, err := os.ReadFile("server/conf/AuthDB.json")
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	var dbRequest database.AuthDB
	err = json.Unmarshal(data, &dbRequest)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	log.Println("Connecting to database...")
	db := database.Database{ConnTimeout: 10 * time.Second}
	if err = db.Connect(dbRequest); err != nil {
		log.Fatalf("%v\n", err)
	}
	defer db.Disconnect()
	log.Println("Complete")

	//log.Println("Init database...")
	//if err = db.Init(); err != nil {
	//	log.Fatalf("%v\n", err)
	//}
	//log.Println("Complete")

	l, err := net.Listen("tcp", "0.0.0.0:3302")
	if err != nil {
		log.Printf("Cannot listen: %v\n", err)
	}
	defer l.Close()

	wg := sync.WaitGroup{}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Cannot accept: %v\n", err)
			break
		}

		wg.Add(1)
		go func() {
			handleConnection(conn, &db)
			wg.Done()
		}()
	}

	wg.Wait()
}
