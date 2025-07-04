package request

type ListQuestionParam struct {
	Type       int    `json:"type"`
	Keyword    string `json:"keyword"`
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
}

type QuestionCreate struct {
	Type     int    `json:"type"`                               // 题型：1-单选题  2-多选题  3-编程题
	Title    string `json:"title" gorm:"type:text"`             // 标题
	Content  string `json:"content" gorm:"type:text"`           // 内容
	Options  string `json:"options,omitempty" gorm:"type:text"` // 选项
	Answer   string `json:"answer,omitempty" gorm:"type:text"`  // 答案
	Language string `json:"language"`                           // 语言
}

type QuestionEdit struct {
	ID       int64  `json:"id,string"`                          // 题目 id
	Type     int    `json:"type"`                               // 题型：1-单选题  2-多选题  3-编程题
	Title    string `json:"title" gorm:"type:text"`             // 标题
	Content  string `json:"content" gorm:"type:text"`           // 内容
	Options  string `json:"options,omitempty" gorm:"type:text"` // 选项
	Answer   string `json:"answer,omitempty" gorm:"type:text"`  // 答案
	Language string `json:"language" gorm:"size:50"`            // 语言
}

type AIGenQuestionParams struct {
	Type     int
	Language string
	Amount   int
}

type DeleteQuestion struct {
	IDs []string `json:"ids,string"` // 需要删除的 id
}

type AIGenQuestion struct {
	Type     *int   `json:"type"`     // 题型：1-单选题  2-多选题  3-编程题
	Language string `json:"language"` // 编程语言
	Amount   *int   `json:"amount"`   // 出题数量
}
