package res

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

func SendConnFailWithMsg(msg string, conn *websocket.Conn) {
	byteData, _ := json.Marshal(Response{
		Code: FailServiceCode,
		Msg:  msg,
	})
	conn.WriteMessage(websocket.TextMessage, byteData)
}

func SendConnOkWithMsg(data any, conn *websocket.Conn) {
	byteData, _ := json.Marshal(Response{
		Code: SuccessCode,
		Msg:  SuccessCode.String(),
		Data: data,
	})
	conn.WriteMessage(websocket.TextMessage, byteData)
}

func SendWsMsg(onLineMap map[uint]map[string]*websocket.Conn, userID uint, data any) {
	addrMap, ok := onLineMap[userID]
	if ok {
		byteData, _ := json.Marshal(Response{
			Code: SuccessCode,
			Msg:  SuccessCode.String(),
			Data: data,
		})
		for _, w := range addrMap {
			w.WriteMessage(websocket.TextMessage, byteData)
		}
	}
}
