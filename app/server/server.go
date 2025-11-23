package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"sysken-pay-api/app/config"
	"sysken-pay-api/app/infra/query"
	"sysken-pay-api/app/ui/api/health"
	api_item "sysken-pay-api/app/ui/api/item"
	api_user "sysken-pay-api/app/ui/api/user"
	"sysken-pay-api/app/usecase/item"
	"sysken-pay-api/app/usecase/user"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// TODO エンドポイントとUI層との接続
// usecase層とdomain層の生成
// CROSの設定
// サーバーの立ち上げ

const (
	requestTimeout    = 60 * time.Second
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 10 * time.Second
)

func Run(db *sql.DB) error {
	addr := ":" + strconv.Itoa(config.Port())

	// Repository
	userRepo := query.NewUserProfileRepository(db)
	itemRepo := query.NewItemRepository(db)

	// UseCase
	registerUserUseCase := user.NewRegisterUserUseCase(userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(userRepo)
	registerItemUseCase := item.NewRegisterItemUseCase(itemRepo)
	updateItemUseCase := item.NewUpdateItemUseCase(itemRepo)
	findItemByJanCodeUseCase := item.NewFindItemByJanCodeUseCase(itemRepo)
	getAllItemsUseCase := item.NewGetAllItemsUseCase(itemRepo)

	// Handler
	userHandler := api_user.NewUserHandler(registerUserUseCase, updateUserUseCase)
	itemHandler := api_item.NewItemHandler(registerItemUseCase, updateItemUseCase, findItemByJanCodeUseCase, getAllItemsUseCase)

	// ルーターの設定
	r := chi.NewRouter()

	// ミドルウェアの設定
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(newCORS().Handler)

	// リクエストタイムアウトをコンテキストに設定
	// リクエストがタイムアウトした場合、ctx.Done()を通じて通知し、以降の処理を停止する
	r.Use(middleware.Timeout(requestTimeout))

	// v1 エンドポイント
	r.Route("/v1", func(r chi.Router) {
		// ユーザー関連
		r.Route("/user", func(r chi.Router) {
			r.Post("/", userHandler.RegisterUser)
			r.Patch("/{user_id}", userHandler.UpdateUser)
		})

		// 商品関連
		r.Route("/item", func(r chi.Router) {
			r.Post("/", itemHandler.ResisterItem)
			r.Patch("/", itemHandler.UpdateItem)
			r.Get("/{jan_code}", itemHandler.GetItemByJanCode)
			r.Get("/", itemHandler.GetAllItems)
		})

		// ヘルスチェック
		r.Route("/health", func(r chi.Router) {
			r.Get("/", health.Check)
		})
	})

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	l, err := net.Listen("tcp", addr)
	slog.Info("Serve on 127.0.0.1", "addr", addr)
	if err != nil {
		slog.Error("failed to listen", "err", err)
	}

	go func() {
		if err = srv.Serve(l); err != nil {
			slog.Error("failed to serve", "err", err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", "err", err)
	}

	return nil
}

func newCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
	})
}
