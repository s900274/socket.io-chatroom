package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nkovacs/go-socket.io"
)

var room = make(map[string]int) //房间的详情

type MsgStruct struct {
	Username  string `json:"username"`
	Message	  string `json:"msg"`
	Uuid 	  string `json:"uuid"`
	RoomId	  string `json:"roomId"`
	Type	  string `json:"type"`
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	_ = server.On("connection", func(so socketio.Socket) {

		log.Println("on connection")

		roomId := joinRoom(room)

		room[roomId]++

		_ = so.Join(roomId)

		msgStruct := &MsgStruct{}
		msgStruct.RoomId = roomId
		msgStruct.Type	 = "connect"

		respStr, _ := json.Marshal(msgStruct)
		so.Emit("chat message", string(respStr))

		log.Println("room list", room)
		_ = so.On("chat message", func(msg string) {

			log.Println(so.Id())
			log.Println(so.Rooms())

			msgStruct := &MsgStruct{}
			json.Unmarshal([]byte(msg), msgStruct)

			log.Println("chat message", msg)
			log.Println(msgStruct)
			log.Println("emit:", so.Emit("chat message", msg))
			switch msgStruct.Type {
			case "message":
				_ = so.BroadcastTo(msgStruct.RoomId, "chat message", msg)
			case "leave":
				log.Println("room list", room)
				//room[msgStruct.RoomId]--
				if room[msgStruct.RoomId] == 1 {
					delete(room, msgStruct.RoomId)
				} else {
					room[msgStruct.RoomId]--
				}
				so.Leave(msgStruct.RoomId)
			}
		})
		_ = so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	_ = server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("/Users/WeiWenWang/Documents/GoProject/socket-chatroom/public")))
	log.Println("Serving at localhost:12345...")
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func joinRoom(roomList map[string]int) (string){

	for k, v := range roomList {

		if v == 1 {
			return k
		}
	}

	return "Room_"+strconv.FormatInt(time.Now().UnixNano(), 10)
	}