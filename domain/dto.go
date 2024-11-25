package domain

type ValidateRequest struct {
	FriendID int    `json:"friend_id"`
	Code     string `json:"code"`
}
