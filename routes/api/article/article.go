//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package article

import (
	"context"
	"github.com/OGFris/itoplist-backend/database"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

func Create(ctx *fasthttp.RequestCtx) {
	allowed, _, user := database.RequireAuth(ctx, database.ModeratorRole)
	if !allowed {
		ctx.Error("please signin to do this", http.StatusUnauthorized)

		return
	}

	a := database.Article{
		Title:       string(ctx.FormValue("title")),
		Description: string(ctx.FormValue("description")),
		Content:     string(ctx.FormValue("content")),
		Hidden:      ctx.Request.PostArgs().GetBool("hidden"),
		Date:        time.Now(),
		AuthorId:    int(user.ID()),
		Type:        ctx.Request.PostArgs().GetUintOrZero("type"),
	}

	r, err := database.Elastic.Index().Index("articles").OpType("_doc").BodyJson(a).Do(context.Background())
	if err != nil {
		ctx.Error("received an error while creating article on database", http.StatusInternalServerError)

		return
	}

	ctx.Success("text/plain", []byte(r.Result))
}

func Update(ctx *fasthttp.RequestCtx) {
	allowed, _, _ := database.RequireAuth(ctx, database.ModeratorRole)
	if !allowed {
		ctx.Error("please signin to do this", http.StatusUnauthorized)

		return
	}

	articleId := string(ctx.Request.PostArgs().Peek("article_id"))
	title := string(ctx.Request.PostArgs().Peek("title"))
	description := string(ctx.Request.PostArgs().Peek("description"))
	content := string(ctx.Request.PostArgs().Peek("content"))
	hidden := ctx.Request.PostArgs().GetBool("hidden")

	r, err := database.Elastic.Update().Index("articles").Id(articleId).Doc(map[string]interface{}{
		"title":       title,
		"description": description,
		"content":     content,
		"hidden":      hidden,
	}).Do(context.Background())
	if err != nil {
		ctx.Error("received an error while creating article on database", http.StatusInternalServerError)

		return
	}

	ctx.Success("text/plain", []byte(r.Result))
}

func Latest(ctx *fasthttp.RequestCtx) {
	r, err := database.Elastic.Search().Index("articles").Sort("date", true).Do(context.Background())
	if err != nil {
		ctx.Error("received an error while search for new articles", http.StatusInternalServerError)

		return
	}

	bytes, err := jsoniter.Marshal(r.Hits.Hits)
	if err != nil {
		ctx.Error("received an error while processing articles data", http.StatusInternalServerError)

		return
	}

	ctx.Success("application/json", bytes)
}
