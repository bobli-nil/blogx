package chat_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum/chat_msg_type"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/focus_service"
	"blogx_server/utils/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

type SessionListRequest struct {
	common.PageInfo
}

type SessionListResponse struct {
	UserID       uint                       `json:"userID"`
	UserNickname string                     `json:"userNickname"`
	UserAvatar   string                     `json:"userAvatar"`
	Msg          ctype.ChatMsg              `json:"msg"`
	MsgType      chat_msg_type.MsgType      `json:"msgType"`
	MsgDate      *time.Time                 `json:"MsgDate"`
	Relation     relationship_enum.Relation `json:"relation"`
}

type SessionResult struct {
	User1        uint       `gorm:"column:user1" json:"user1"`
	User2        uint       `gorm:"column:user2" json:"user2"`
	MaxTime      *time.Time `gorm:"column:max_time" json:"maxTime"`
	MsgCount     int        `gorm:"column:msg_count" json:"msgCount"`
	LatestChatID uint       `gorm:"column:latest_chat_id" json:"latestChatID"`
}

func (ChatApi) SessionListView(c *gin.Context) {
	cr := middleware.GetBind[SessionListRequest](c)
	claims := jwt.GetClaims(c)

	var deletedIDList []uint
	global.DB.
		Model(&models.UserChatActionModel{}).
		Where("user_id = ? and is_delete = ?", claims.UserID, true).
		Select("chat_id").Scan(&deletedIDList)

	args := []any{claims.UserID, claims.UserID}

	sql := `
        select m2.user1, m2.user2, m2.max_time, m2.msg_count, m1.id latest_chat_id
        from chat m1
        inner join (
            select
                least(send_user_id, rev_user_id) user1,
                greatest(send_user_id, rev_user_id) user2,
                max(created_at) max_time,
                count(*) msg_count
            from chat
            where (send_user_id = ? or rev_user_id = ?)
    `

	if len(deletedIDList) > 0 {
		sql += " and id not in ?"
		args = append(args, deletedIDList)
	}

	sql += `
            group by user1, user2
            order by max_time desc
            limit ? offset ?
        ) m2
        on least(m1.send_user_id, m1.rev_user_id) = m2.user1
        and greatest(m1.send_user_id, m1.rev_user_id) = m2.user2
        and m1.created_at = m2.max_time;
    `
	args = append(args, cr.GetLimit(), cr.GetOffset())

	_list := make([]SessionResult, 0)
	global.DB.Raw(sql, args...).Scan(&_list)

	var userIDList []uint
	var chatIDList []uint
	for _, table := range _list {
		chatIDList = append(chatIDList, table.LatestChatID)
		if table.User1 == claims.UserID {
			userIDList = append(userIDList, table.User2)
		}
		if table.User2 == claims.UserID {
			userIDList = append(userIDList, table.User1)
		}
	}

	userMap := common.ScanMapV2(models.UserModel{}, common.ScanOption{
		Where: global.DB.Where("id in ?", userIDList),
	})
	chatMap := common.ScanMapV2(models.ChatModel{}, common.ScanOption{
		Where: global.DB.Where("id in ?", chatIDList),
	})

	relationMap := focus_service.CalcUserPatchRelationship2(claims.UserID, userIDList)

	var list = make([]SessionListResponse, 0)
	for _, table := range _list {
		item := SessionListResponse{}
		if table.User1 == claims.UserID {
			item.UserID = table.User2
		}
		if table.User2 == claims.UserID {
			item.UserID = table.User1
		}
		item.UserNickname = userMap[item.UserID].Nickname
		item.UserAvatar = userMap[item.UserID].Avatar
		item.Msg = chatMap[table.LatestChatID].Msg
		item.MsgType = chatMap[table.LatestChatID].MsgType
		item.MsgDate = chatMap[table.LatestChatID].CreatedAt
		item.Relation = relationMap[item.UserID]

		list = append(list, item)
	}

	res.OkWithList(list, len(list), c)

}
