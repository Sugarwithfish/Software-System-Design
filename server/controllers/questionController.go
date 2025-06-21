package controllers

import (
	"github.com/Sugarwithfish/Software-System-Design/models/request"
	"github.com/Sugarwithfish/Software-System-Design/server"
	"github.com/Sugarwithfish/Software-System-Design/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Summary      获取题目列表
// @Description  根据题型、关键词、页码和页大小获取题目列表
// @Accept       json
// @Produce      json
// @Param        type     query int    false "题型：0-全部 1-单选 2-多选 3-编程"
// @Param        keyword  query string false "关键词"
// @Param        pageNum  query int    false "页码"
// @Param        pageSize query int    false "页大小"
// @Success      200 {object} response.Response "成功 (Code: 0), 参数错误 (Code: -1), 内部错误 (Code: -3)"
func ListQuestionHandler() func(c *gin.Context) {
	return func(c *gin.Context) {

		questionType, err := strconv.Atoi(c.Query("type"))
		if err != nil {
			utils.ParamError(c, err.Error())
			return
		}
		keyword := c.Query("keyword")
		pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
		if err != nil {
			utils.ParamError(c, err.Error())
			return
		}
		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		if err != nil {
			utils.ParamError(c, err.Error())
			return
		}
		listQuestionParams := request.ListQuestionParam{
			Type:       questionType,
			Keyword:    keyword,
			PageNumber: pageNum,
			PageSize:   pageSize,
		}
		questions, err := server.ListQuestion(listQuestionParams)
		if err != nil {
			utils.FailWithMsg(c, utils.ERROR_INTERNAL, err.Error())
			return
		}
		utils.Success(c, questions)
	}
}

// @Summary      删除题目
// @Description  根据题目ID列表批量删除题目
// @Accept       json
// @Produce      json
// @Param        request body request.DeleteQuestion true "题目ID"
// @Success      200 {object} response.Response "成功 (Code: 0), 参数错误 (Code: -1), 内部错误 (Code: -3)"
// @Router       /api/v1/questions [delete]
func DeleteQuestionHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 绑定并校验请求体
		var req request.DeleteQuestion
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.FailWithMsg(c, utils.ERROR_PARAM, err.Error())
			return
		}
		var ids []int64

		for _, idStr := range req.IDs {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				utils.ParamError(c, err.Error())
				return
			}
			ids = append(ids, id)
		}

		// 调用服务层执行删除逻辑
		deletedCount, err := server.DeleteQuestion(ids) // 建议将 service.DeleteQuestion 重命名为 service.DeleteQuestions
		if err != nil {
			utils.FailWithMsg(c, utils.ERROR_INTERNAL, err.Error())
			return
		}

		// 成功响应，返回实际删除的题目数量
		utils.Success(c, deletedCount)
	}
}

// @Summary      获取题目详情
// @Description  根据题目ID获取题目详细信息
// @Accept       json
// @Produce      json
// @Param        id path int true "题目ID"
// @Success      200 {object} response.Response "成功 (Code: 0), 参数错误 (Code: -1), 内部错误 (Code: -3)"
// @Router       /api/v1/questions/{id} [get]
func GetQuestionHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 获取请求中的题目 id
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			utils.ParamError(c, err.Error())
			return
		}
		// 调用 service 层查询
		questionDetail, err := server.GetQuestion(id)
		if err != nil {
			utils.FailWithMsg(c, utils.ERROR_INTERNAL, err.Error())
			return
		}
		if questionDetail == nil {
			utils.FailWithMsg(c, utils.ERROR_NOT_FOUND, "未找到问题")
			return
		}
		utils.Success(c, questionDetail)
	}
}

// @Summary      创建题目
// @Description  创建新的题目（单选/多选/编程题）
// @Accept       json
// @Produce      json
// @Param        request body request.QuestionCreate true "题目信息"
// @Success      200 {object} response.Response "成功 (Code: 0), 参数错误 (Code: -1), 业务错误 (Code: -2), 内部错误 (Code: -3)"
// @Router       /api/v1/questions [post]
func SaveQuestionHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var question request.QuestionCreate
		err := c.ShouldBind(&question)
		if err != nil {
			utils.ParamError(c, err.Error())
			return
		}

		// 调用 Service 层创建题目
		err = server.CreateQuestion(&question)
		if err != nil {
			utils.FailWithMsg(c, utils.ERROR_INTERNAL, err.Error())
		} else {
			// 返回成功响应
			utils.Success(c, nil)
		}
	}
}

// @Summary      批量创建题目
// @Description  批量创建新的题目
// @Accept       json
// @Produce      json
// @Param        request body []request.QuestionCreate true "题目列表"
// @Success      200 {object} response.Response "成功 (Code: 0), 参数错误 (Code: -1), 内部错误 (Code: -3)"
// @Router       /api/v1/questions/batch [post]
func BatchSaveQuestionHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var questions []request.QuestionCreate
		if err := c.ShouldBindJSON(&questions); err != nil {
			utils.ParamError(c, err.Error())
			return
		}

		// 将 apiRequest.QuestionCreate 转换为 []*models.Question
		for _, questionCreate := range questions {
			err := server.CreateQuestion(&questionCreate)
			if err != nil {
				utils.FailWithMsg(c, utils.ERROR_INTERNAL, err.Error())
				return
			}
		}
		utils.Success(c, nil)
	}
}

// @Summary      编辑题目
// @Description  编辑现有题目信息
// @Accept       json
// @Produce      json
// @Param        request body request.QuestionEdit true "题目编辑信息"
// @Success      200 {object} response.Response "成功 (Code: 0), 参数错误 (Code: -1), 业务错误 (Code: -2)"
// @Router       /api/v1/questions/{id} [put]
func EditQuestionHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var questionEdit request.QuestionEdit

		if err := c.ShouldBind(&questionEdit); err != nil {
			utils.ParamError(c, err.Error())
			return
		}

		if err := server.EditQuestion(&questionEdit); err != nil {
			utils.FailWithMsg(c, utils.ERROR_INTERNAL, err.Error())
		} else {
			utils.Success(c, nil)
		}
	}
}

// @Summary     AI生成题目
// @Description 根据请求参数生成题目
// @Accept      json
// @Produce     json
// @Param       request body request.AIGenQuestion true "生成题目参数"
// @Success     200 {object} response.Response "成功 (Code: 0), 参数错误 (Code: -1), 内部错误 (Code: -3)"
// @Router      /api/v1/questions/ai-generate [post]
func AIGenQuestionHandler(cfg *server.LLMConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 获取请求参数
		var req request.AIGenQuestion
		err := c.ShouldBindJSON(&req)
		if err != nil {
			utils.ParamError(c, err.Error())
			return
		}

		if req.Type == nil {
			utils.ParamError(c, "题目类型不能为空")
			return
		}
		if req.Amount == nil {
			utils.ParamError(c, "题目数量不能为空")
			return
		}
		if req.Language == "" {
			utils.ParamError(c, "编程语言不能为空")
			return
		}

		questionType := *req.Type
		amount := *req.Amount
		// 校验请求参数
		if questionType != 1 && questionType != 2 && questionType != 3 {
			utils.ParamError(c, "请求参数错误: 题目类型不合法")
			return
		}
		if amount < 1 || amount > 10 {
			utils.ParamError(c, "请求参数错误: 题目数量不合法")
			return
		}

		llmParams := request.AIGenQuestionParams{
			Type:     questionType,
			Language: req.Language,
			Amount:   amount,
		}

		// 调用服务层逻辑
		questions, err := server.AIGenerateQuestion(llmParams, cfg)
		if err != nil {
			utils.FailWithMsg(c, utils.ERROR_INTERNAL, err.Error())
			return
		}
		utils.Success(c, questions)
	}
}
