//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package database

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"os"
)

type User struct {
	Model
	FirstName       string `gorm:"Type:varchar(255);Column:first_name;NOT NULL"`
	LastName        string `gorm:"Type:varchar(255);Column:last_name;NOT NULL"`
	Username        string `gorm:"Type:varchar(255);Column:username;NOT NULL;unique" json:"username" valid:"matches((?i)^[.a-z0-9_-]+$),required"`
	Email           string `gorm:"Type:varchar(255);Column:email;NOT NULL;unique" json:"email" valid:"email,required"`
	Password        string `gorm:"Type:varchar(255);Column:password;NOT NULL" json:"-" valid:"required,length(6|16384)"`
	Validated       int    `gorm:"Type:tinyint(1);Column:validated;NOT NULL;Default:0" json:"-"`
	ValidationToken string `gorm:"Type:varchar(255);Column:token;NOT NULL;unique" json:"-"`
	Role            int    `gorm:"Type:int(11);Column:role;NOT NULL;Default:0" json:"role"`
}

type UserInterface interface {
	JWT() (string, error)
	ByJWT(token string) error
	ID() uint
}

const (
	NormalRole = iota
	ModeratorRole
	AdminRole
)

func (u *User) ID() uint {

	return u.Model.ID
}

func (u *User) JWT() (string, error) {
	key := os.Getenv("JWT_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   u.ID,
		"type": 0,
	})

	return token.SignedString([]byte(key))
}

func (u *User) ByJWT(token string) error {
	key := os.Getenv("JWT_KEY")

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(key), nil
	})

	if err != nil {
		return errors.New("invalid token")
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		id := uint(claims["id"].(float64))
		authType := claims["type"].(int)
		if authType != 0 {
			return errors.New("invalid token")
		}

		err := Instance.First(&u, &User{Model: Model{ID: id}}).Error

		return err
	}

	return errors.New("invalid token")
}

func RequireAuth(ctx *fasthttp.RequestCtx, permission int) (allowed bool, authType int, u UserInterface) {
	facebookUser := &FacebookUser{}
	user := &User{}

	err := user.ByJWT(string(ctx.Request.Header.Cookie("jwt")))
	err2 := facebookUser.ByJWT(string(ctx.Request.Header.Cookie("jwt")))

	// TODO: I'll PROBABLY have to replace all these bool statements with a switch case later.
	if err != nil && err2 != nil {

		return false, authType, u
	}

	if err != nil && err2 == nil {
		authType = 1
		u = facebookUser
	} else {
		u = user
		authType = 0
	}

	if permission != 0 && authType != 0 {

		return false, authType, u
	}

	if user.Role >= permission {

		return true, authType, u
	}

	return false, authType, u
}
