package global_notification_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"fmt"

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
	if err := global.DB.Take(&model, "title = ?", cr.Title).Error; err == nil {
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
				if model.IsDelete {
					deletedIDList = append(deletedIDList, model.NotificationID)
				}
				if model.IsRead {
					readMsgMap[model.NotificationID] = true
				}
			}
		}
		if len(deletedIDList) > 0 {
			query.Where("id not in ?", deletedIDList)
		}
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
		Debug:    true,
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

func (GlobalNotificationApi) RemoveAdminView(c *gin.Context) {
	cr := middleware.GetBind[models.DeleteRequest](c)

	var list []models.GlobalNotificationModel
	global.DB.Find(&list, "id in ?", cr.IDList)
	if len(list) > 0 {
		global.DB.Delete(&list)
	}

	msg := fmt.Sprintf("删除%d条全局消息，成功删除%d条", len(cr.IDList), len(list))
	res.OkWithMsg(msg, c)
}

type UserMsgActionRequest struct {
	ID   uint `json:"id" binding:"required"`
	Type int8 `json:"type" binding:"required,oneof=1 2"` // 1读取 2删除
}

func (GlobalNotificationApi) UserMsgActionView(c *gin.Context) {
	cr := middleware.GetBind[UserMsgActionRequest](c)
	claims := jwt.GetClaims(c)

	var msg models.GlobalNotificationModel
	if err := global.DB.Take(&msg, "id = ?", cr.ID).Error; err != nil {
		res.FailWithMsg("消息不存在", c)
		return
	}

	var ugnm models.UserGlobalNotificationModel
	err := global.DB.Take(&ugnm, "notification_id = ? and user_id = ?", cr.ID, claims.UserID).Error
	if err != nil {
		model := models.UserGlobalNotificationModel{
			NotificationID: cr.ID,
			UserID:         claims.UserID,
		}
		if cr.Type == 1 {
			model.IsRead = true
		} else {
			model.IsDelete = true
		}
		global.DB.Create(&model)
		res.OkWithMsg("消息读取成功", c)
		return
	}
	if ugnm.IsDelete {
		res.OkWithMsg("消息已删除", c)
		return
	}
	if cr.Type == 1 {
		global.DB.Model(&ugnm).Update("is_read", true)
		res.OkWithMsg("消息读取成功", c)
		return
	}
	if cr.Type == 2 {
		global.DB.Model(&ugnm).Update("is_delete", true)
		res.OkWithMsg("消息删除成功", c)
		return
	}

}
