package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type Room struct {
	name    string
	members sync.Map
	log     *log.Logger
}

func (r *Room) AddUser(u *User) {
	r.members.Store(u, time.Now())
	r.log.Debugf("added user %s to map", u.UserName)
}

func (r *Room) Delete(u *User) {
	r.members.Delete(u)
}

// BroadcastMsg sends a the message to all members in the Room
func (r *Room) BroadcastMsg(sender *User, msg string) {
	r.members.Range(func(key, value any) bool {
		u, ok := key.(*User)
		if !ok {
			r.log.Error(ok)
		}
		if u.UserName != sender.UserName {
			r.log.Info(msg)
			u.MsgSender <- msg
		}
		return true
	})

}
func (r *Room) NotifyMembers(newUser *User, action string) {
	r.members.Range(func(key, value any) bool {
		u, ok := key.(*User)
		if !ok {
			r.log.Error(ok)
		}
		notification := fmt.Sprintf("* %s has %s the room", newUser.UserName, action)
		r.log.Info(notification)
		if u.UserName != newUser.UserName {
			u.MsgSender <- notification
		}
		return true
	})

}

func (r *Room) GetConnectedUsers(newUser *User) string {
	userChan := make(chan string)
	go func() {
		r.members.Range(func(key, value any) bool {
			u, ok := key.(*User)
			if !ok {
				r.log.Error(ok)
			}
			if u.UserName != newUser.UserName {
				userChan <- u.UserName
			}
			return true
		})
		close(userChan)
	}()

	users := []string{}
	for user := range userChan {
		users = append(users, user)
	}
	// presenceString := fmt.Sprintf("* Hi %s!, the room contains : %s", newUser.UserName, strings.Join(users, ", "))
	presenceString := fmt.Sprintf("* The room contains : %s", strings.Join(users, ", "))
	r.log.Info(presenceString)
	return presenceString
}

type User struct {
	UserName  string
	CreatedAt time.Time
	MsgSender chan string // messages written to this channel gets sent to the client. make this buffered so that if any one client is slow to read from the channel it doesn't block writes for other clients.
}

func newUser(username string) (User, error) {
	validUsername := regexp.MustCompile(`^[A-Za-z0-9]{1,16}$`)
	if validUsername.MatchString(username) {
		return User{UserName: username, CreatedAt: time.Now(), MsgSender: make(chan string, 100)}, nil
	}
	return User{}, errors.New("Closing Connection -- Invalid Username: username must be between 1-16 characters long and can only have alpha numeric characters")
}
