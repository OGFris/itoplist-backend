//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package article

import (
	"context"
	"github.com/OGFris/itoplist-backend/database"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

func Create(ctx *fasthttp.RequestCtx) {
	if !database.RequireAuth(ctx) {
		ctx.Error("please signin to do this", http.StatusUnauthorized)

		return
	}

	hidden := ctx.Request.PostArgs().GetBool("hidden")

	a := database.Article{
		Title:       string(ctx.FormValue("title")),
		Description: string(ctx.FormValue("title")),
		Content:     string(ctx.FormValue("title")),
		Hidden:      hidden,
		Date:        time.Now(),
	}

	r, err := database.Elastic.Index().Index("articles").OpType("_doc").BodyJson(a).Do(context.Background())
	if err != nil {
		ctx.Error("received an error while creating article on database", http.StatusInternalServerError)

		return
	}

	ctx.Success("text/plain", []byte(r.Id))
}
