package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type ctxKey string

const (
	PageKey  ctxKey = "page"
	LimitKey ctxKey = "limit"
)

func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		// parse page
		page, err := strconv.Atoi(q.Get("page"))
		if err != nil || page < 1 {
			page = 1
		}
		// parse limit
		limit, err := strconv.Atoi(q.Get("limit"))
		if err != nil || limit < 1 {
			limit = 10
		}
		ctx := context.WithValue(r.Context(), PageKey, page)
		ctx = context.WithValue(ctx, LimitKey, limit)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
