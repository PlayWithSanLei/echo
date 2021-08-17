package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":6430")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

type User struct {
	ID             uint64
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

type Message struct {
	OwnerID uint64
	Content string
}

func (u *User) String() string {
	return u.Addr + ", UID:" + strconv.FormatUint(u.ID, 10) + ", Enter At:" +
		u.EnterAt.Format("2006-01-01 15:04:05+8000")
}

var (
	enteringChannel = make(chan *User)
	leavingChannel  = make(chan *User)
	messageChannel  = make(chan Message, 8)
)

// 用于记录聊天室用户 并进行消息广播
// 新用户
// 用户普通消息
// 用户离开
func broadcaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			users[user] = struct{}{}
		case user := <-leavingChannel:
			delete(users, user)
			close(user.MessageChannel)
		case msg := <-messageChannel:
			for user := range users {
				if user.ID == msg.OwnerID { //不给自己发送
					continue
				}
				user.MessageChannel <- msg.Content
			}
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	go sendMessage(conn, user.MessageChannel)

	// 发送欢迎消息
	user.MessageChannel <- "Welcome to ECHO, " + user.String()
	msg := Message{
		OwnerID: user.ID,
		Content: fmt.Sprintf("[%d] has enter", user.ID),
	}

	messageChannel <- msg
	enteringChannel <- user

	var userActive = make(chan struct{})
	go func() {
		d := time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = strconv.FormatUint(user.ID, 10) + ":" + input.Text()
		messageChannel <- msg
		// 用户活跃
		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误", err)
	}

	leavingChannel <- user
	msg.Content = fmt.Sprintf("[%d] has left", user.ID)
	messageChannel <- msg
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

var idLock sync.Mutex
var globalID uint64

func GenUserID() uint64 {
	idLock.Lock()
	defer idLock.Unlock()

	globalID++
	return globalID
}
