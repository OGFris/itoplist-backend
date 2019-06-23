package database

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
)

type FacebookUser struct {
	Model
	AuthCode string `gorm:"Type:varchar(255);Column:auth_code;NOT NULL;unique" json:"auth_code"`
	Name     string `gorm:"Type:varchar(255);Column:name;NOT NULL" json:"name"`
	Nickname string `gorm:"Type:varchar(255);Column:nickname;NOT NULL" json:"nickname"`
	Email    string `gorm:"Type:varchar(255);Column:email;NOT NULL" json:"email"`
}

func (u *FacebookUser) ID() uint {

	return u.Model.ID
}

func (u *FacebookUser) JWT() (token string, err error) {
	key := os.Getenv("JWT_KEY")

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   u.ID,
		"type": 1,
	})

	return tkn.SignedString([]byte(key))
}

func (u *FacebookUser) ByJWT(token string) error {
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
		if authType != 1 {
			return errors.New("invalid token")
		}

		err := Instance.First(&u, &FacebookUser{
			Model: Model{
				ID: id,
			},
		}).Error

		return err
	}

	return errors.New("invalid token")
}
