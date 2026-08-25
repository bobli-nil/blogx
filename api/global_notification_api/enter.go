package global_notification_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type GlobalNotificationApi struct{}

type CreateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Icon    string `json:"icon"`
	Href    string `json:"href"`
}

func (GlobalNotificationApi) CreateView(c *gin.Context) {
	cr := middleware.GetBind[CreateRequest](c)

	var model models.GlobalNotificationModel
	if err := global.DB.Take(&model, "title = ?", cr.Title).Error; err != nil {
		res.FailWithMsg("全局消息名称重复", c)
		return
	}

	err := global.DB.Create(&models.GlobalNotificationModel{
		Title:   cr.Title,
		Content: cr.Content,
		Icon:    cr.Icon,
		Href:    cr.Href,
	}).Error
	if err != nil {
		res.FailWithMsg("创建消息失败", c)
		return
	}

	res.OkWithMsg("全局消息创建成功", c)
}

type ListRequest struct {
	common.PageInfo
	Type int8 `form:"type" binding:"required,oneof=1 2"` // 1自己可见的 2管理员
}

type ListResponse struct {
	models.GlobalNotificationModel
	IsRead bool `json:"isRead"`
}

func (GlobalNotificationApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[ListRequest](c)
	claims := jwt.GetClaims(c)

	query := global.DB.Where("")
	readMsgMap := make(map[uint]bool)

	switch cr.Type {
	case 1:
		deletedIDList := make([]uint, 0)
		var userGlobalNotificationList []models.UserGlobalNotificationModel
		global.DB.Model(&models.UserGlobalNotificationModel{}).Find(&userGlobalNotificationList, "user_id = ?", claims.UserID)
		if len(userGlobalNotificationList) > 0 {
			for _, model := range userGlobalNotificationList {
				readMsgMap[model.NotificationID] = model.IsRead
				deletedIDList = append(deletedIDList, model.NotificationID)
			}
		}
		query.Where("id not in ?", deletedIDList)
	case 2:
		if claims.Role != enum.AdminRole {
			res.FailWithMsg("权限错误", c)
			return
		}
	}

	_list, count, _ := common.ListQuery(models.GlobalNotificationModel{}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title", "content"},
		Where:    query,
	})

	list := make([]ListResponse, 0)
	for _, model := range _list {
		list = append(list, ListResponse{
			GlobalNotificationModel: model,
			IsRead:                  readMsgMap[model.ID],
		})
	}

	res.OkWithList(list, count, c)
}
