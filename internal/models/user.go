package models

import (
	"golang.org/x/crypto/bcrypt"
)

// @Description User credentials for login or registration
// @Example {"username":"Wladim1r","password":"12345qwerty"}
type Request struct {
	Username string `json:"username" binding:"required" example:"Wladim1r"`
	Password string `json:"password" binding:"required" example:"12345qwerty"`
}

// @Description User data model
// @Example {"id":1,"username":"Wladim1r"}
type User struct {
	ID       uint   `json:"id"       gorm:"primarykey"                                    example:"1"`
	Username string `json:"username" gorm:"unique"                                        example:"Wladim1r" binding:"required"`
	Password string `json:"-"`
	Books    []Book `json:"-"        gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)

	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// @Description User response including ID, username and quantity of books
// @Example {"id":1,"username":"Wladim1r","total":10}
type UserResponse struct {
	ID       uint   `json:"id"       example:"1"`
	Username string `json:"username" example:"Wladim1r"`
	Total    int    `json:"total"    example:"10"`
}

// @Description Response containing user array
// @Example {"users":[{"id":1,"username":"Wladim1r","total":10},{"id":2,"username":"Ivan","total":13}]}
type UsersResponse struct {
	Users []UserResponse `json:"users"`
}
