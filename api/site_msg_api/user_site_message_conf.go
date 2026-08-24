package site_msg_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/utils/jwt"
	"blogx_server/utils/mps"

	"github.com/gin-gonic/gin"
)

func (SiteMsgApi) UserMessageConfView(c *gin.Context) {
	claims := jwt.GetClaims(c)

	var userMessageConf models.UserMessageConfModel
	if err := global.DB.Take(&userMessageConf, "user_id = ?", claims.UserID).Error; err != nil {
		res.OkWithMsg("用户消息配置不存在", c)
		return
	}

	res.OkWithData(userMessageConf, c)
}

type UserMessageConfRequest struct {
	OpenCommentMessage bool `json:"openCommentMessage" u:"open_comment_message"` // 是否开启回复和评论
	OpenDiggMessage    bool `json:"openDiggMessage" u:"open_digg_message"`       // 是否开启赞和收藏
	OpenPrivateChat    bool `json:"openPrivateChat" u:"open_private_chat"`       // 是否开启私聊
}

func (SiteMsgApi) UserMessageConfUpdateView(c *gin.Context) {
	claims := jwt.GetClaims(c)
	cr := middleware.GetBind[UserMessageConfRequest](c)

	var userMessageConf models.UserMessageConfModel
	if err := global.DB.Take(&userMessageConf, "user_id = ?", claims.UserID).Error; err != nil {
		res.OkWithMsg("用户消息配置不存在", c)
		return
	}

	mp := mps.Struct2Map(cr, "u")
	global.DB.Model(&userMessageConf).Updates(mp)
	res.OkWithMsg("用户消息配置更新成功", c)
}
