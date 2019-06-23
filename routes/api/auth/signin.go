//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package auth

import (
	"github.com/OGFris/itoplist-backend/database"
	"github.com/OGFris/itoplist-backend/utils"
	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

func Signin(ctx *fasthttp.RequestCtx) {
	if ctx.PostArgs().Has("username") && ctx.PostArgs().Has("password") {
		username := string(ctx.PostArgs().Peek("username"))
		password := string(ctx.PostArgs().Peek("password"))
		user := &database.User{}

		err := database.Instance.Where("username = ?", username).Or("email = ?", username).Find(&user).Error
		if err == gorm.ErrRecordNotFound {
			ctx.Error("couldn't find a user with the same email and/or password", http.StatusBadRequest)

			return
		}

		if err != nil {
			ctx.Error("Internal database error. Error 003", http.StatusInternalServerError)
			log.Fatalln(err)

			return
		}

		if utils.Compare(user.Password, password) {
			bytes, err := jsoniter.Marshal(user)
			if err != nil {
				ctx.Error("unexpected error occurred. Error 004", http.StatusInternalServerError)
				log.Fatalln(err)

				return
			}

			cookie := fasthttp.AcquireCookie()

			cookie.SetKey("jwt")
			jwt, err := user.JWT()
			utils.PanicError(err)
			cookie.SetValue(jwt)

			ctx.Response.Header.SetCookie(cookie)
			ctx.Success("application/json", bytes)
		} else {
			ctx.Error("email and/or password were not correct", http.StatusBadRequest)

			return
		}

	} else {
		ctx.Error("not all required parameters were provided", http.StatusBadRequest)
	}
}
