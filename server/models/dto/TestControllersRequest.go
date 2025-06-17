package dto

type TestControllersRequestModel struct {
	Type     *int `json:"type"`
	Language *int `json:"language"`
	Number   *int `json:"number"`
}
