package chat_api

import (
	"blogx_server/common/res"
	"blogx_server/utils/jwt"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var onlineMap = map[uint]map[string]*websocket.Conn{}

func (ChatApi) ChatView(c *gin.Context) {
	claims, err := jwt.ParseTokenByGin(c)
	if err != nil || claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("ws 升级失败 %s", err)
		return
	}

	userID := claims.UserID
	addr := conn.RemoteAddr().String()
	addrMap, ok := onlineMap[userID]
	if !ok {
		onlineMap[userID] = map[string]*websocket.Conn{
			addr: conn,
		}
	} else {
		_, ok := addrMap[addr]
		if !ok {
			addrMap[addr] = conn
		}
	}
	fmt.Println("进入后", onlineMap)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Errorf("读取消息失败 %s", err)
			break
		}
		fmt.Println("收到的消息", msg, msgType)

		err = conn.WriteMessage(1, []byte("你好"))
		if err != nil {
			logrus.Errorf("发动消息失败 %s", err)
			break
		}
	}

	defer conn.Close()

	// 关闭
	addrMap2, ok2 := onlineMap[userID]
	if ok2 {
		_, ok := addrMap2[addr]
		if ok {
			delete(addrMap2, addr)
		}
		if len(addrMap2) == 0 {
			delete(onlineMap, userID)
		}
	}
	fmt.Println("关闭", onlineMap)
}
