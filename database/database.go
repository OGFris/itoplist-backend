//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package database

import (
	"context"
	"errors"
	"github.com/OGFris/itoplist-backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v7"
	"net/http"
	"os"
	"time"
)

var (
	Instance *gorm.DB
	Elastic  *elastic.Client
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"Default:null" sql:"index" json:"deleted_at"`
}

func Init() {
	if Instance == nil {
		var d *gorm.DB

		sql := os.Getenv("DB_USER") + ":" +
			os.Getenv("DB_PASSWORD") + "@tcp(" +
			os.Getenv("DB_ADDRESS") + ")/" +
			os.Getenv("DB_NAME")

		d, err := gorm.Open("mysql", sql+"?charset=utf8&parseTime=True&loc=Local")
		utils.PanicError(err)
		Instance = d

		d.AutoMigrate(
			&User{},
			&FacebookUser{},
		)

	}

	if Elastic == nil {
		var err error

		Elastic, err = elastic.NewClient(
			elastic.SetHttpClient(&http.Client{Transport: &Transport{
				Username: os.Getenv("ES_USERNAME"),
				Password: os.Getenv("ES_PASSWORD"),
			},
			},
			),
			elastic.SetURL(os.Getenv("ES_URL")),
			elastic.SetSniff(false),
			elastic.SetHealthcheck(false),
			elastic.SetHealthcheckTimeoutStartup(0),
		)
		utils.PanicError(err)

		exist, err := Elastic.IndexExists("articles").Do(context.Background())
		if !exist {
			body := `
			{
				"settings" : {
					"analysis" : {
						"analyzer" : {
							"default" : {
								"tokenizer" : "standard",
									"filter" : ["asciifolding", "lowercase"]
							}
						}
					}
				}
			}`

			result, err := Elastic.CreateIndex("articles").BodyString(body).Do(context.Background())
			utils.PanicError(err)

			if !result.Acknowledged {
				panic(errors.New("acknowledged should be true but returned false"))
			}
		}
	}
}
