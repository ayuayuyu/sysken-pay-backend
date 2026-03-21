package balance

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	apierrors "sysken-pay-api/app/ui/api/pkg/errors"
	"sysken-pay-api/app/usecase/balance"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	GetBalance(w http.ResponseWriter, r *http.Request)
	GetPurchaseHistories(w http.ResponseWriter, r *http.Request)
}

func NewBalanceHandler(
	getBalanceUseCase balance.GetBalanceUseCase,
	getPurchaseHistoriesUseCase balance.GetPurchaseHistoriesUseCase,
) Handler {
	return &balanceHandlerImpl{
		getBalanceUseCase:           getBalanceUseCase,
		getPurchaseHistoriesUseCase: getPurchaseHistoriesUseCase,
	}
}

var _ Handler = (*balanceHandlerImpl)(nil)

type balanceHandlerImpl struct {
	getBalanceUseCase           balance.GetBalanceUseCase
	getPurchaseHistoriesUseCase balance.GetPurchaseHistoriesUseCase
}

func (h *balanceHandlerImpl) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		log.Printf("user_id is missing in URL")
		apierrors.RespondError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	ctx := r.Context()
	b, err := h.getBalanceUseCase.GetBalance(ctx, userID)
	if err != nil {
		log.Printf("Failed to get balance: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := toGetBalanceResponse(b)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Failed to encode response: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *balanceHandlerImpl) GetPurchaseHistories(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		log.Printf("user_id is missing in URL")
		apierrors.RespondError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	page := 1
	perPage := 20
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v >= 1 {
			page = v
		}
	}
	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if v, err := strconv.Atoi(pp); err == nil && v >= 1 {
			perPage = v
		}
	}

	ctx := r.Context()
	histories, totalAmount, pageInfo, err := h.getPurchaseHistoriesUseCase.GetPurchaseHistories(ctx, userID, page, perPage)
	if err != nil {
		log.Printf("Failed to get purchase histories: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := toGetPurchaseHistoriesResponse(histories, totalAmount, pageInfo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Failed to encode response: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
	}
}
