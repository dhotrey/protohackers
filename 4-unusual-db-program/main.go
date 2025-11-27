package main

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
)

func main() {
	log := getNewLogger("main", log.InfoLevel)
	log.Info("Unusual database program!")

	var kv sync.Map
	kv.Store("version", "Ken's Key-Value Store 1.0")
	laddr := &net.UDPAddr{
		Port: 6942,
		IP:   net.ParseIP("0.0.0.0"),
	}
	log.Info("Listening on", "laddr", laddr)

	udpConn, err := net.ListenUDP("udp", laddr)

	if err != nil {
		log.Fatal(err)
	}

	udpBuff := make([]byte, 1024) //  protocol max message size is 1000

	for {
		n, clientAddr, readError := udpConn.ReadFromUDP(udpBuff)
		if readError != nil {
			log.Fatal("Something went wrong while reading UDP packet", "readError", readError)
		}

		go func(clientAddr net.UDPAddr, msg string, udpConn *net.UDPConn, db *sync.Map) {
			log.Info("->", "src", clientAddr, "message", msg)
			if strings.Contains(msg, "=") { // insert operation
				k, v := parseKeyValue(msg)
				log.Debug("Got key value ", "key", k, "value", v)
				if k != "version" { // making version immutable
					db.Store(k, v)
					log.Info("Insert", "key", k, "val", v)
				} else {
					log.Warn("Attempted to modify version key - ignoring")
				}
			} else { // retrieval
				key := msg
				value, exists := db.Load(key)
				log.Debug("Got value for key", "value", value, "key", key)
				if exists { // only respond when key exists
					response := fmt.Sprintf("%s=%s", key, value)
					n, err := udpConn.WriteToUDP([]byte(response), &clientAddr)
					if err != nil {
						log.Error("Error writing to client")
					}
					log.Info("replied: ", "message", response, "size", n)
				} else {
					log.Warnf("Key %s doesn't exist", key)
				}
			}
		}(*clientAddr, string(udpBuff[:n]), udpConn, &kv)
	}
}

func parseKeyValue(message string) (key string, value string) {
	splits := strings.SplitN(message, "=", 2)
	key = splits[0]
	value = splits[1]
	return key, value
}
