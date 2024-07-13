package models

import "github.com/google/uuid"

type ResponseUser struct {
	//User ID
	UserId uuid.UUID `json:"user_id,omitempty"`
	//Login for User
	Login string `json:"login,omitempty"`
	//Password for User
	Password string `json:"password,omitempty"`
	//Roles that belong to the User
	Roles []string `json:"roles,omitempty"`
}
