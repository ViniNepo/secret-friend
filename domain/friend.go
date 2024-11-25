package domain

type Friend struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Description  string `json:"description"`
	Requirement  string `json:"requirement"`
	SelectFriend *int   `json:"select_friend"`
	ValidateCode string `json:"validate_code"`
	IsValid      bool   `json:"is_valid"`
}
