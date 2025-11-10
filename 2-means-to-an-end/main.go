package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/charmbracelet/log"
)

const PORT = 6942
const MESSAGE_SIZE = 9 // bytes

func main() {
	var wg sync.WaitGroup

	log := getNewLogger("main")
	log.Info("Means to an End!")

	l, err := net.Listen("tcp4", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	addr := l.Addr().String()
	log.Infof("Listening on %s", addr)

	for i := 1; i < 12; i++ { // the tester only makes 11 connections
		wg.Add(1)
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		connId := fmt.Sprintf("conn-%d", i)
		log.Infof("Received connection from %s", connId)
		go reqHandler(conn, getNewLogger(connId), &wg)
	}
	wg.Wait()
	log.Info("Bye!")
}

func reqHandler(conn net.Conn, log *log.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	db := initDb()
	for {
		msg := make([]byte, MESSAGE_SIZE)
		n, err := io.ReadAtLeast(conn, msg, MESSAGE_SIZE)
		if err != nil {
			log.Error(err)
		}
		op, timestamp, price, err := parseMsg(msg)
		if err != nil {
			log.Error("Error while parsing message")
			log.Fatal(err)
		}

		log.Infof("<-- (%d bytes) : %x | %s | %d | %d ", n, msg, op, timestamp, price)

		switch op {
		case "I":
			if !db.Add(timestamp, price) {
				log.Warnf("Closing Connection: timestamp %d exists in db", timestamp)
				return
			}
		case "Q":
			avg := db.Query(timestamp, price) // in this case timestamp is actially mintime and price is maxtime
			buff := make([]byte, 4)
			binary.BigEndian.PutUint32(buff, uint32(avg))
			n, err := conn.Write(buff)
			if err != nil {
				log.Errorf("Error writing to conn : %v", err)
			}
			log.Infof("--> (%d bytes) : %x |  %d ", n, buff, avg)

		default:
			log.Warnf("Closing Connection: undefined op (%s)", op)
			return
		}
	}
}

func parseMsg(msg []byte) (operation string, timestamp int32, price int32, err error) {
	buf := bytes.NewReader(bytes.Clone(msg[1:5]))
	err = binary.Read(buf, binary.BigEndian, &timestamp)
	if err != nil {
		return "", 0, 0, err
	}
	buf = bytes.NewReader(bytes.Clone(msg[5:9]))
	err = binary.Read(buf, binary.BigEndian, &price)
	if err != nil {
		return "", 0, 0, err
	}
	return string(msg[0]), timestamp, price, nil
}
