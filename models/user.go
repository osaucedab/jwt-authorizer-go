package models

import (
	"errors"
	"fmt"
	"jwt-authorizer/utils/token"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

func (u *User) SaveUser() (*User, error) {

	if err := u.HashPassword(u.Password); err != nil {
		return &User{}, err
	}

	err := DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	user.Password = string(bytes)
	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (string, error) {

	var err error

	u := User{}

	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}


func GetUserById(uid uint) (User, error){
	var user User
	fmt.Println("user id")
	fmt.Println(uid)

	if err := DB.First(&user, uid).Error; err != nil{
		fmt.Println("user")
		fmt.Println(user)
		fmt.Println(err)
		return user, errors.New("User not found")
	}

	user.PrepareGive()

	return user, nil
}

func (user *User) PrepareGive() {
	user.Password = ""
}