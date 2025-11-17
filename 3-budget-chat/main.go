package main

import (
	"bufio"
	"fmt"
	"net"
	// "net/http"

	// "github.com/arl/statsviz"
	"github.com/charmbracelet/log"
)

const PORT = 6942

func main() {
	log := getNewLogger("main")
	log.Info("Budget Chat!")

	// mux := http.NewServeMux()
	// statsviz.Register(mux)
	// go func() {
	// 	log.Info("statsviz on http://localhost:8080/debug/statsviz")
	// 	log.Info(http.ListenAndServe("localhost:8080", mux))
	// }()

	l, err := net.Listen("tcp4", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatal(err)
	}
	addr := l.Addr().String()
	log.Infof("Listening on %s", addr)

	room := new(Room)
	room.name = "common"
	room.log = getNewLogger("room")
	for i := 1; ; i++ {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		connId := fmt.Sprintf("conn-%d", i)
		log.Infof("Received connection from %s", connId)
		go reqHandler(conn, getNewLogger(connId), room)
	}
}

func reqHandler(conn net.Conn, log *log.Logger, chatRoom *Room) {
	defer conn.Close()
	welcomeMsg := "Welcome to budgetchat! What shall I call you?"
	fmt.Fprintln(conn, "Welcome to budgetchat! What shall I call you?")
	log.Debug(welcomeMsg)
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	username := scanner.Text()
	user, err := newUser(username)
	if err != nil {
		log.Error(err)
		fmt.Fprintln(conn, err)
		log.Warn("closing connection")
		return
	}
	log.SetPrefix(fmt.Sprintf("%s (%s)", log.GetPrefix(), user.UserName))
	log.Debugf("Got username %s", user.UserName)
	chatRoom.AddUser(&user)
	presenceNotif := chatRoom.GetConnectedUsers(&user)
	fmt.Fprintln(conn, presenceNotif)
	chatRoom.NotifyMembers(&user, "entered")

	go func() {
		for msg := range user.MsgSender {
			fmt.Fprintln(conn, msg)
		}
	}()

	for scanner.Scan() {
		msg := formatMessage(user.UserName, scanner.Text())
		chatRoom.BroadcastMsg(&user, msg)
	}

	chatRoom.NotifyMembers(&user, "left")
	defer chatRoom.Delete(&user)
}
