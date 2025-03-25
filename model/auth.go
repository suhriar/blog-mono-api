package model

type UserAuth struct {
	ID       int64  `json:"ud"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
