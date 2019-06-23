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
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/objx"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
)

var facebook common.Provider

func init() {
	var err error

	facebook, err = gomniauth.Provider("facebook")
	utils.PanicError(err)
}

func FacebookLogin(ctx *fasthttp.RequestCtx) {
	state := gomniauth.NewState("after", "success")
	authUrl, err := facebook.GetBeginAuthURL(state, nil)
	if err != nil {
		ctx.Error("unexpected error", http.StatusInternalServerError)

		return
	}

	ctx.Redirect(authUrl, http.StatusFound)
}

func FacebookCallback(ctx *fasthttp.RequestCtx) {
	omap, err := objx.FromURLQuery(ctx.QueryArgs().String())
	if err != nil {
		ctx.Error("unexpected error", http.StatusInternalServerError)

		return
	}

	c, err := facebook.CompleteAuth(omap)
	if err != nil {
		ctx.Error("unexpected error", http.StatusInternalServerError)

		return
	}

	u, userErr := facebook.GetUser(c)
	if userErr != nil {
		ctx.Error("unexpected error", http.StatusInternalServerError)

		return
	}

	exist := &database.FacebookUser{}
	if gorm.IsRecordNotFoundError(database.Instance.Where("auth_code = ?", u.AuthCode()).Find(exist).Error) {
		user := &database.FacebookUser{
			AuthCode: u.AuthCode(),
			Name:     u.Name(),
			Nickname: u.Nickname(),
			Email:    u.Email(),
			// TODO: Scan the data object and detect whether the user has a phone only or an email only or both.
		}

		database.Instance.Create(user)
	}

	cookie := fasthttp.AcquireCookie()

	cookie.SetKey("jwt")
	jwt, err := exist.JWT()
	utils.PanicError(err)
	cookie.SetValue(jwt)

	ctx.Response.Header.SetCookie(cookie)
	ctx.Redirect(os.Getenv("URL"), http.StatusFound)
}
