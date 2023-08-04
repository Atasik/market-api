package handler

import (
	"context"
	"errors"
	"fmt"
	"market/pkg/model"
	"market/pkg/service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type OptionsKey string

const (
	OptionsContextKey OptionsKey = "query_options"
)

const (
	authorizationHeader = "Authorization"
	defaultSortField    = "created_at"
	defaultLimit        = 25
	maxLimit            = 50
	defaultOffset       = 0
)

var (
	ErrNoQuery = errors.New("no query")
)

func accessLog(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("access log middleware")
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Since(start),
		)
	})
}

func auth(s service.User, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware", r.URL.Path)

		header := r.Header.Get(authorizationHeader)
		if header == "" {
			newErrorResponse(w, "empty auth header", http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			newErrorResponse(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			newErrorResponse(w, "empty token", http.StatusUnauthorized)
			return
		}

		token, err := s.CheckToken(headerParts[1])
		if err != nil {
			newErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := service.ContextWithToken(r.Context(), token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("panic middleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// я знаю, что это коряво, но лень рефакторить пока
func query(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("queryMiddleware", r.URL.Path)

		sortBy := strings.ToLower(r.URL.Query().Get("sort_by"))
		sortOrder := strings.ToUpper(r.URL.Query().Get("sort_order"))

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = defaultLimit
		}

		var offset int
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			offset = defaultOffset
		}

		if sortBy == "" {
			sortBy = defaultSortField
		}

		if sortOrder == "" {
			sortOrder = model.DESCENDING
		}

		if limit <= 0 {
			limit = defaultLimit
		}

		if limit > maxLimit {
			limit = maxLimit
		}

		if page <= 0 {
			offset = defaultOffset
		} else {
			offset = (page - 1) * limit
		}

		options := &Options{
			SortBy:    sortBy,
			SortOrder: sortOrder,
			Limit:     limit,
			Offset:    offset,
		}
		ctx := contextWithOptions(r.Context(), options)
		next(w, r.WithContext(ctx))
	})
}

type Options struct {
	SortBy, SortOrder string
	Limit             int
	Offset            int
}

func optionsFromContext(ctx context.Context) (*Options, error) {
	options, ok := ctx.Value(OptionsContextKey).(*Options)
	if !ok || options == nil {
		return nil, ErrNoQuery
	}
	return options, nil
}

func contextWithOptions(ctx context.Context, opts *Options) context.Context {
	return context.WithValue(ctx, OptionsContextKey, opts)
}
