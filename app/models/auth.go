// models/auth.go
package models

type User struct {
	ID       string `json:"id" binding:"required"`
	Username string `json:"username" binding:"required,min=1"`
	Email    string `json:"email" binding:"required,email"`
}

type Role struct {
	Name        string   `json:"name" binding:"required,min=1"`
	Permissions []string `json:"permissions" binding:"required,min=1,dive,required"`
}

type TokenRequest struct {
	User  User   `json:"user" binding:"required"`
	Roles []Role `json:"roles" binding:"required,min=1,dive"`
}

type TokenResponse struct {
	Person Person `json:"person"`
	Token  string `json:"token"`
}
