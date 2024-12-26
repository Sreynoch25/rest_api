package models

import "time"

type User struct {
	ID           int        `json:"id"`
	LastName     string     `json:"last_name" validate:"required"`  
	FirstName    string     `json:"first_name" validate:"required"` 
	UserName     string     `json:"user_name" validate:"required"`
	LoginID      string     `json:"login_id" validate:"required"`
	Email        string     `json:"email" validate:"required,email"` 
	Password     string     `json:"password,omitempty" validate:"required,min=6"`
	RoleName     string     `json:"role_name" validate:"required"`
	RoleID       int        `json:"role_id"` 
	IsAdmin      bool       `json:"is_admin"`
	LoginSession *string    `json:"login_session"`
	LastLogin    *time.Time `json:"last_login"`
	CurrencyID   *int       `json:"currency_id"`
	LanguageID   *int       `json:"language_id"`
	StatusID     int        `json:"status_id"`
	Order        *int       `json:"order"`
	CreatedBy    int        `json:"created_by"` 
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedBy    int        `json:"updated_by"` 
	UpdatedAt    time.Time  `json:"updated_at"` 
	DeletedBy    *int       `json:"deleted_by"`
	DeletedAt    *time.Time `json:"deleted_at"`
}


type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
	ExpiredAt time.Time `json:"expired_at"`
}