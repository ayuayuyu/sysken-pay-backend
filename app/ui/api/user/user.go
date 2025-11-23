package user

import (
	"encoding/json"
	"log"
	"net/http"
	apierrors "sysken-pay-api/app/ui/api/pkg/errors"
	"sysken-pay-api/app/usecase/user"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TODO　APIリクエストからデータを整形してユースケースに情報を渡すものを作成する
type Handler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(registerUserUseCase user.RegisterUserUseCase, updateUserUseCase user.UpdateUserUseCase) Handler {
	return &userHandlerImpl{
		registerUserUseCase: registerUserUseCase,
		updateUserUseCase:   updateUserUseCase,
	}
}

var _ Handler = (*userHandlerImpl)(nil)

type userHandlerImpl struct {
	registerUserUseCase user.RegisterUserUseCase
	updateUserUseCase   user.UpdateUserUseCase
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

func (h *userHandlerImpl) UpdateUser(w http.ResponseWriter, r *http.Request) {
	//userRequestのパース
	var req PatchUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		log.Printf("user_id is missing in URL")
		apierrors.RespondError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	parseUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Failed to parse user_id: %v", err)
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	//ユースケースの呼び出し
	updatedUser, err := h.updateUserUseCase.UpdateUser(ctx, parseUUID, req.UserName)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//レスポンスの作成
	res := toPatchUserResponse(updatedUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
}
