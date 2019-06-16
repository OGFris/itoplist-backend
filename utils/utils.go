//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package utils

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"math/rand"
	"testing"
	"time"
)

type FormError struct {
	Message string `json:"message"`
}

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func WriteJson(w *fasthttp.Response, data interface{}) {
	w.Header.Set("Content-Type", "application/json")
	PanicError(jsoniter.NewEncoder(w.BodyWriter()).Encode(data))
}

func AssertEq(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatal(fmt.Sprintf("%v != %v", a, b))
	}
}

// GenerateToken returns a 16 chars long string of random characters.
func GenerateToken() string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	tokenBytes := make([]rune, 16)
	for i := range tokenBytes {
		tokenBytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(tokenBytes)
}
