package db

import (
	"errors"
	"github.com/Sugarwithfish/Software-System-Design/models"
	"gorm.io/gorm"
	"log"
)

func QuestionPageQuery(questionType int, keyword string, pageNum, pageSize int) ([]*models.Question, int64, error) {
	var questions []*models.Question
	var total int64

	query := DB.Model(&models.Question{})

	// Type 为 0 则查询全部类型，不添加 type 条件
	if questionType != 0 {
		query = query.Where("type = ?", questionType)
	}

	// 关键词过滤
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		log.Printf("查询题目总数失败: %v\n", err)
		return nil, 0, err
	}

	// 指定查询结果字段
	query = query.Select("id", "type", "title", "content", "options", "answer", "language")

	// 获取分页数据
	offset := (pageNum - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&questions).Error
	if err != nil {
		log.Printf("分页查询题目失败: %v\n", err)
		return nil, 0, err
	}

	return questions, total, nil
}

// DeleteQuestion 批量删除题目，并返回删除的行数
func DeleteQuestion(ids []int64) (int64, error) {
	result := DB.Where("id IN ?", ids).Delete(&models.Question{})
	if result.Error != nil {
		log.Printf("删除题目失败（IDs: %v): %v\n", ids, result.Error)
		return -1, result.Error
	}
	return result.RowsAffected, nil
}

// GetQuestionById 根据 id 查找题目
func GetQuestionById(id int64) (*models.Question, error) {
	var question models.Question
	err := DB.First(&question, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("数据库查询失败")
	}
	return &question, nil
}

// InsertQuestion 插入单条数据
func InsertQuestion(question *models.Question) error {
	result := DB.Create(question)
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.New("插入数据库失败!")
	}
	return nil
}

// UpdateQuestionByID 更新题目
func UpdateQuestionByID(question *models.Question) error {
	result := DB.Model(&question).Updates(question)
	return result.Error
}
