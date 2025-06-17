package controllers

import (
	"github.com/Sugarwithfish/Software-System-Design/models/dto"
	"github.com/Sugarwithfish/Software-System-Design/utils"
	"github.com/gin-gonic/gin"
)

func TestControllers(c *gin.Context) {

	var req dto.TestControllersRequestModel

	if err := c.ShouldBind(&req); err != nil {
		utils.ParamError(c, "参数错误："+err.Error())
		return
	}

	//utils.Success(c, "testing_controller")
	utils.Success(c, req)

}
