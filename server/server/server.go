package server

import (
	"encoding/json"
	"errors"
	"github.com/Sugarwithfish/Software-System-Design/db"
	"github.com/Sugarwithfish/Software-System-Design/models"
	req "github.com/Sugarwithfish/Software-System-Design/models/request"
	"github.com/Sugarwithfish/Software-System-Design/utils"
	"gorm.io/gorm"
	"log"
)

func ListQuestion(param req.ListQuestionParam) (*models.PageQueryResult, error) {
	questions, total, err := db.QuestionPageQuery(param.Type, param.Keyword, param.PageNumber, param.PageSize)

	if err != nil {
		return nil, err
	}

	var questionVOs []*models.Question
	for _, question := range questions {
		questionVOs = append(questionVOs, question)
	}
	return &models.PageQueryResult{
		Records: questionVOs,
		Total:   total,
	}, nil
}

func DeleteQuestion(ids []int64) (int64, error) {
	rowCount, err := db.DeleteQuestion(ids)
	if err != nil {
		return -1, err
	}
	return rowCount, nil
}

func GetQuestion(id int64) (*models.Question, error) {
	question, err := db.GetQuestionById(id)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return question, nil
}

func CreateQuestion(questionCreate *req.QuestionCreate) error {
	log.Println(questionCreate)
	if questionCreate.Title == "" {
		return errors.New("标题不为空")
	}
	if questionCreate.Type != 3 {
		var options []string
		err := json.Unmarshal([]byte(questionCreate.Options), &options)
		if err != nil || len(options) != 4 {
			return errors.New("选择题的选项必须为 4个合法的字符串数组！")
		}
		var answer []int
		err = json.Unmarshal([]byte(questionCreate.Answer), &answer)
		if err != nil {
			log.Println(err.Error())
			return errors.New("选择题的答案格式有误！")
		}
		if questionCreate.Type == 1 && len(answer) != 1 {
			return errors.New("单选题的答案有且仅有一个！")
		}
		if questionCreate.Type == 2 && len(answer) <= 1 {
			return errors.New("多选题的答案必须大于一个！")
		}
	}
	var question = &models.Question{
		ID:       utils.GenerateID(),
		Type:     questionCreate.Type,
		Title:    questionCreate.Title,
		Content:  questionCreate.Content,
		Options:  questionCreate.Options,
		Answer:   questionCreate.Answer,
		Language: questionCreate.Language,
	}
	if err := db.InsertQuestion(question); err != nil {
		return errors.New("数据库插入失败: " + err.Error())
	}

	return nil
}

func EditQuestion(questionEdit *req.QuestionEdit) error {
	if _, err := db.GetQuestionById(questionEdit.ID); err != nil {
		return errors.New("数据库不存在该题目")
	}

	questionToUpdate := &models.Question{
		ID:       questionEdit.ID,
		Type:     questionEdit.Type,
		Title:    questionEdit.Title,
		Content:  questionEdit.Content,
		Options:  questionEdit.Options,
		Answer:   questionEdit.Answer,
		Language: questionEdit.Language,
	}

	return db.UpdateQuestionByID(questionToUpdate)

}
