//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package database

import "time"

// Article will be stored on Elasticsearch.
type Article struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Hidden      bool      `json:"hidden"`
	Date        time.Time `json:"date"`
}
