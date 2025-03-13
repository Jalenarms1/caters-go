package types

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Uid string  `json:"uid"`
	Exp float64 `json:"exp"`
	jwt.RegisteredClaims
}

type UserCtxKey string
type ShopSlugCtxKey string

var AuthKey UserCtxKey = "foodgo-auth"
var UrlSlugKey ShopSlugCtxKey = "foodgo-url-slug"

type Error struct {
	Err        error
	ReturnCode int
}

type ErrHandlerFunc func(w http.ResponseWriter, r *http.Request) *Error
