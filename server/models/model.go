package models

type Question struct {
	ID       int64  `json:"id"`                                 // 题目 id
	Type     int    `json:"type"`                               // 题型：1-单选题  2-多选题  3-编程题
	Title    string `json:"title" gorm:"type:text"`             // 标题
	Content  string `json:"content" gorm:"type:text"`           // 内容
	Options  string `json:"options,omitempty" gorm:"type:text"` // 选项
	Answer   string `json:"answer,omitempty" gorm:"type:text"`  // 答案
	Language string `json:"language" gorm:"size:50"`            // 语言
}

type PageQueryResult struct {
	Records interface{} `json:"records"`
	Total   int64       `json:"total"`
}
