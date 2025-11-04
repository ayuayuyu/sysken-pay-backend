package server

import (
	"database/sql"
	"log"
	"net/http"
	"sysken-pay-api/app/infra/query"
	api_user "sysken-pay-api/app/ui/api/user"
	"sysken-pay-api/app/usecase/user"
)

// TODO エンドポイントとUI層との接続
// usecase層とdomain層の生成
// CROSの設定
// サーバーの立ち上げ
func Run(db *sql.DB) error {

	// Repository
	userRepo := query.NewUserProfileRepository(db)

	// UseCase
	registerUserUseCase := user.NewRegisterUserUseCase(userRepo)

	// Handler
	userHandler := api_user.NewUserHandler(registerUserUseCase)

	//ルーターの設定
	mux := http.NewServeMux()

	// エンドポイントの設定
	mux.HandleFunc("POST /user", userHandler.RegisterUser)

	// CORSの設定
	handler := CORSMiddleware(mux)

	log.Println("Starting server on :8080")
	// サーバーの起動
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return err
	}

	return nil
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
