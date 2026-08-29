package chat_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum/chat_msg_type"
	"blogx_server/utils/jwt"
	"encoding/json"
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

type ChatRequest struct {
	RevUserID uint                  `json:"revUserID"`
	MsgType   chat_msg_type.MsgType `json:"msgType"` // 1文本 2图片 3md
	Msg       ctype.ChatMsg         `json:"msg"`
}
type ChatResponse struct {
	ChatListResponse
}

var onlineMap = map[uint]map[string]*websocket.Conn{}

func (ChatApi) ChatView(c *gin.Context) {
	claims, err := jwt.ParseTokenByGin(c)
	if err != nil || claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}

	// 发送者用户信息
	userID := claims.UserID
	var user models.UserModel
	if err := global.DB.Take(&user, userID).Error; err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("ws 升级失败 %s", err)
		return
	}

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
		// 读取消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Errorf("读取消息失败 %s", err)
			break
		}

		// 获取请求参数
		var req ChatRequest
		err = json.Unmarshal(msg, &req)
		if err != nil {
			res.SendConnFailWithMsg("参数错误", conn)
			continue
		}

		// 判断接收人在不在
		var revUser models.UserModel
		if err := global.DB.Take(&revUser, req.RevUserID).Error; err != nil {
			res.SendConnFailWithMsg("接收人不存在", conn)
			continue
		}

		// 消息入库
		model := models.ChatModel{
			SendUserID: claims.UserID,
			RevUserID:  req.RevUserID,
			MsgType:    req.MsgType,
			Msg:        req.Msg,
		}
		if err := global.DB.Create(&model).Error; err != nil {
			res.SendConnFailWithMsg("消息发送失败", conn)
			continue
		}

		// 发送消息给接收人revUserID
		data := ChatResponse{
			ChatListResponse: ChatListResponse{
				ChatModel:        model,
				SendUserNickname: user.Nickname,
				SendUserAvatar:   user.Avatar,
				RevUserNickname:  revUser.Nickname,
				RevUserAvatar:    revUser.Avatar,
			},
		}
		res.SendWsMsg(onlineMap, req.RevUserID, data)
		// 给自己也发一份
		data.IsMe = true
		res.SendConnOkWithMsg(data, conn)
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
