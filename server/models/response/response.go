package response

type AIGenQuestion struct {
	Title   string `json:"title" gorm:"type:text"`
	Content string `json:"content" gorm:"type:text"`
	Options string `json:"options,omitempty" gorm:"type:text"`
	Answer  string `json:"answer,omitempty" gorm:"type:text"`
}

type CodingQuestion struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ChoiceQuestion struct {
	Title   string   `json:"title"`
	Options []string `json:"options"`
	Answer  []int    `json:"answer"`
}
