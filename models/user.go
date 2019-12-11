package models

import "net/http"

type User struct {
	Id int            `json:"id"`
	Email string      `json:"email"`
	Login string      `json:"login"`
	Fullname string   `json:"fullname"`
	Password string   `json:"password"`
	AccVerified bool  `json:"acc_verified"`
}

type Tokens struct {
	Access *http.Cookie
	Refresh *http.Cookie
}