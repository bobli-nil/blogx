package chat_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum/chat_msg_type"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/focus_service"
	"blogx_server/utils/jwt"
	"blogx_server/utils/xss"
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

var OnlineMap = map[uint]map[string]*websocket.Conn{}

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
	addrMap, ok := OnlineMap[userID]
	if !ok {
		OnlineMap[userID] = map[string]*websocket.Conn{
			addr: conn,
		}
	} else {
		_, ok := addrMap[addr]
		if !ok {
			addrMap[addr] = conn
		}
	}
	fmt.Println("进入后", OnlineMap)

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

		// 具体的消息类型要做处理
		switch req.MsgType {
		case chat_msg_type.TextMsgType:
			if req.Msg.TextMsg == nil || req.Msg.TextMsg.Content == "" {
				res.SendConnFailWithMsg("文本消息内容为空", conn)
				continue
			}
		case chat_msg_type.ImageMsgType:
			if req.Msg.ImageMsg == nil || req.Msg.ImageMsg.Src == "" {
				res.SendConnFailWithMsg("图片消息内容为空", conn)
				continue
			}
		case chat_msg_type.MarkdownMsgType:
			if req.Msg.MarkdownMsg == nil || req.Msg.MarkdownMsg.Content == "" {
				res.SendConnFailWithMsg("Markdown消息内容为空", conn)
				continue
			}
			content, err := xss.XssFilter(req.Msg.MarkdownMsg.Content)
			if err != nil {
				res.SendConnFailWithMsg("Markdown内容XSS处理失败", conn)
				continue
			}
			req.Msg.MarkdownMsg.Content = content
		default:
			res.SendConnFailWithMsg("不支持的消息类型", conn)
			continue
		}

		// 判断发送人和接收人关系
		relation := focus_service.CalcUserRelationship(userID, revUser.ID)
		fmt.Println("好友关系", relation.String())
		switch relation {
		case relationship_enum.RelationStranger:
			var revUserMsgConf models.UserMessageConfModel
			if err := global.DB.Take(&revUserMsgConf, "user_id = ?", revUser.ID).Error; err != nil {
				res.SendConnFailWithMsg("接收人隐私设置不存在", conn)
				continue
			}
			if !revUserMsgConf.OpenPrivateChat {
				res.SendConnFailWithMsg("对方未开启陌生人私聊", conn)
				continue
			}
		case relationship_enum.RelationFriend:
		case relationship_enum.RelationFocus, relationship_enum.RelationFans:
			// 今天对方没有回复，那么你就只能发一条
			var chatList []models.ChatModel
			global.DB.
				Find(&chatList, "date(created_at) = date(now()) and ((send_user_id = ? and rev_user_id = ?) or (send_user_id = ? and rev_user_id = ?))", userID, revUser.ID, revUser.ID, userID)
			var sendChatCount, revUserCount int
			for _, model := range chatList {
				if model.SendUserID == userID {
					sendChatCount++
				}
				if model.RevUserID == userID {
					revUserCount++
				}
			}
			fmt.Printf("chatList %v \n", chatList)
			fmt.Printf("%d %d \n", sendChatCount, revUserCount)
			if sendChatCount > 0 && revUserCount == 0 {
				res.SendConnFailWithMsg("对方未回复，今天只能发送一条消息", conn)
				continue
			}
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
		res.SendWsMsg(OnlineMap, req.RevUserID, data)
		// 给自己也发一份
		data.IsMe = true
		res.SendConnOkWithMsg(data, conn)
	}

	defer conn.Close()

	// 关闭
	addrMap2, ok2 := OnlineMap[userID]
	if ok2 {
		_, ok := addrMap2[addr]
		if ok {
			delete(addrMap2, addr)
		}
		if len(addrMap2) == 0 {
			delete(OnlineMap, userID)
		}
	}
	fmt.Println("关闭", OnlineMap)
}
