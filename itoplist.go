//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package main

import (
	"context"
	"flag"
	"github.com/OGFris/itoplist-backend/database"
	"github.com/OGFris/itoplist-backend/routes/api/article"
	"github.com/OGFris/itoplist-backend/routes/api/auth"
	"github.com/OGFris/itoplist-backend/utils"
	"github.com/buaazp/fasthttprouter"
	"github.com/icrowley/fake"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/valyala/fasthttp"
	"os"
	"time"
)

func main() {
	database.Init()

	gomniauth.SetSecurityKey(os.Getenv("SECURITY_KEY"))
	gomniauth.WithProviders(
		facebook.New(
			os.Getenv("FB_CLIENT_ID"),
			os.Getenv("FB_CLIENT_SECRET"),
			os.Getenv("URL")+"/api/auth/facebook/callback",
		),
	)

	auth.InitProviders()

	seed := flag.Bool("seed", false, "database seeder")
	if *seed {
		Seed(10)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := fasthttprouter.New()

	router.POST("/api/article", article.Create)
	router.POST("/api/article/update", article.Update)
	router.GET("/api/articles", article.Latest)
	router.POST("/api/auth/signin", auth.Signin)
	router.POST("/api/auth/signup", auth.Signup)
	router.GET("/api/auth/facebook/login", auth.FacebookLogin)
	router.GET("/api/auth/facebook/callback", auth.FacebookCallback)

	s := &fasthttp.Server{
		Handler:          router.Handler,
		DisableKeepalive: true,
	}

	utils.PanicError(s.ListenAndServe(":" + port))
}

func Seed(n int) {
	for i := 0; i < n; i++ {
		a := database.Article{
			Title:       fake.Title(),
			Description: fake.Paragraph(),
			Content:     fake.Paragraph(),
			Hidden:      false,
			Date:        time.Now(),
			AuthorId:    fake.Day(),
			Type:        0,
		}

		_, err := database.Elastic.Index().Index("articles").Type("article").BodyJson(a).Do(context.Background())
		utils.PanicError(err)
	}
}
