package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/mps"

	"github.com/gin-gonic/gin"
)

type AdminUserInfoUpdateRequest struct {
	UserID   uint           `json:"userId" binding:"required"`
	Username *string        `json:"username" s-u:"username"`
	Nickname *string        `json:"nickname" s-u:"nickname"`
	Avatar   *string        `json:"avatar" s-u:"avatar"`
	Abstract *string        `json:"abstract" s-u:"abstract"`
	Role     *enum.RoleType `json:"role" s-u:"role"`
}

func (UserApi) AdminUserInfoUpdateView(c *gin.Context) {
	//var cr AdminUserInfoUpdateRequest
	//if err := c.ShouldBindJSON(&cr); err != nil {
	//	res.FailWithError(err, c)
	//	return
	//}
	cr := middleware.GetBind[AdminUserInfoUpdateRequest](c)

	userMap := mps.Struct2Map(cr, "s-u")

	var user models.UserModel
	if err := global.DB.Take(&user, cr.UserID).Error; err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	if err := global.DB.Model(user).Updates(userMap).Error; err != nil {
		res.FailWithMsg("用户信息修改失败", c)
		return
	}

	res.OkWithMsg("用户信息修改成功", c)

}
