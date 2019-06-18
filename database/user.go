//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package database

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

const (
	NormalRole = iota
	ModeratorRole
	AdminRole
)
