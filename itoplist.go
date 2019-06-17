//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package main

import (
	"flag"
	"github.com/OGFris/itoplist-backend/database"
	"github.com/OGFris/itoplist-backend/routes/api/article"
	"github.com/OGFris/itoplist-backend/utils"
	"github.com/buaazp/fasthttprouter"
	_ "github.com/joho/godotenv/autoload"
	"github.com/valyala/fasthttp"
	"os"
)

func main() {
	database.Init()

	seed := flag.Bool("seed", false, "database seeder")
	if *seed {
		// Seed the database with fake info

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := fasthttprouter.New()

	router.POST("/api/article", article.Create)

	s := &fasthttp.Server{
		Handler:          router.Handler,
		DisableKeepalive: true,
	}

	utils.PanicError(s.ListenAndServe(":" + port))
}
