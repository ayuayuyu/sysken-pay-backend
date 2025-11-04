package user

import (
	"encoding/json"
	"log"
	"net/http"
	apierrors "sysken-pay-api/app/ui/api/pkg/errors"
	"sysken-pay-api/app/usecase/user"
)

// TODO　APIリクエストからデータを整形してユースケースに情報を渡すものを作成する
type Handler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(registerUserUseCase user.RegisterUserUseCase) Handler {
	return &userHandlerImpl{
		registerUserUseCase: registerUserUseCase,
	}
}

var _ Handler = (*userHandlerImpl)(nil)

type userHandlerImpl struct {
	registerUserUseCase user.RegisterUserUseCase
}

func (h *userHandlerImpl) RegisterUser(w http.ResponseWriter, r *http.Request) {
	//userRequestのパース
	var req PostUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	//ユースケースの呼び出し
	createdUser, err := h.registerUserUseCase.RegisterUser(ctx, req.UserName)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//レスポンスの作成
	res := toPostUserResponse(createdUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
}
