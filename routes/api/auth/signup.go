//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package auth

import (
	"github.com/OGFris/itoplist-backend/database"
	"github.com/OGFris/itoplist-backend/utils"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"regexp"
)

var r = regexp.MustCompile("(?i)^[.a-z0-9_-]+$")

func Signup(ctx *fasthttp.RequestCtx) {
	user := &database.User{}
	for _, v := range []string{"username", "first_name", "last_name", "email", "password"} {
		if !ctx.PostArgs().Has(v) {
			ctx.Error("not all required parameters were provided, "+v+" is missing!", http.StatusBadRequest)

			return
		}
	}

	if !govalidator.IsEmail(string(ctx.PostArgs().Peek("email"))) {
		ctx.Error("invalid email", http.StatusBadRequest)

		return
	}

	user.FirstName = string(ctx.PostArgs().Peek("first_name"))
	user.LastName = string(ctx.PostArgs().Peek("last_name"))
	user.Username = string(ctx.PostArgs().Peek("username"))
	user.Email = string(ctx.PostArgs().Peek("email"))
	user.Password = string(ctx.PostArgs().Peek("password"))

	checkUser := &database.User{}

	// check if username and/or email are already used
	if database.Instance.Where("username = ?", user.Username).Or("email = ?", user.Email).Find(&checkUser).Error != gorm.ErrRecordNotFound {
		ctx.Error("username and/or email are already used", http.StatusBadRequest)

		return
	}

	// check username size
	if len(user.Username) > 20 || len(user.Username) < 4 {
		ctx.Error("bad username size", http.StatusBadRequest)

		return
	}

	// check if username matches the regex
	if !r.MatchString(user.Username) {
		ctx.Error("bad username", http.StatusBadRequest)

		return
	}

	var err error

	// encrypt the password
	user.Password, err = utils.Encrypt(user.Password)
	if err != nil {
		ctx.Error("unexpected error occurred. Error 001", http.StatusInternalServerError)

		return
	}

	// TODO: do the email validation part.
	user.Validated = 0
	user.ValidationToken = utils.GenerateToken()

	// Create the user and check for any error.
	if errs := database.Instance.Create(user).GetErrors(); len(errs) != 0 {
		ctx.Error("Internal database error. Error 002", http.StatusInternalServerError)
		log.Fatalln(errs)

		return
	}

	utils.WriteJson(&ctx.Response, user)
}
