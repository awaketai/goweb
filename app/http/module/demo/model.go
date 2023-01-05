package demo

import (
	"database/sql"
	"time"
)

type UserModel struct {
	UserId int    `json:"userModel"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
}

type User struct {
	ID           uint
	Name         string
	Email        sql.NullString
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
