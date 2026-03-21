package purchase

import (
	"encoding/json"
	"log"
	"net/http"
	domainPurchase "sysken-pay-api/app/domain/object/purchase"
	apierrors "sysken-pay-api/app/ui/api/pkg/errors"
	"sysken-pay-api/app/usecase/purchase"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	CreatePurchase(w http.ResponseWriter, r *http.Request)
	CancelPurchase(w http.ResponseWriter, r *http.Request)
}

func NewPurchaseHandler(createPurchaseUseCase purchase.CreatePurchaseUseCase, cancelPurchaseUseCase purchase.CancelPurchaseUseCase) Handler {
	return &purchaseHandlerImpl{
		createPurchaseUseCase: createPurchaseUseCase,
		cancelPurchaseUseCase: cancelPurchaseUseCase,
	}
}

var _ Handler = (*purchaseHandlerImpl)(nil)

type purchaseHandlerImpl struct {
	createPurchaseUseCase purchase.CreatePurchaseUseCase
	cancelPurchaseUseCase purchase.CancelPurchaseUseCase
}

func (h *purchaseHandlerImpl) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	var req PostPurchaseRequest
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

	ctx := r.Context()

	// Itemsの変換
	var purchaseItems []domainPurchase.PurchaseItem
	for _, item := range req.Items {
		pi, err := domainPurchase.NewPurchaseItem(item.ItemId, item.Quantity)
		if err != nil {
			log.Printf("Failed to create purchase item: %v", err)
			apierrors.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		purchaseItems = append(purchaseItems, pi)
	}

	//ユースケースの呼び出し
	createdPurchase, err := h.createPurchaseUseCase.CreatePurchase(ctx, userID, purchaseItems)
	if err != nil {
		log.Printf("Failed to create purchase: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//レスポンスの作成
	res := toPostPurchaseResponse(createdPurchase)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *purchaseHandlerImpl) CancelPurchase(w http.ResponseWriter, r *http.Request) {
	var req PostPurchaseCancelRequest
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

	ctx := r.Context()

	// Itemsの変換
	var purchaseItems []domainPurchase.PurchaseItem
	for _, item := range req.Items {
		pi, err := domainPurchase.NewPurchaseItem(item.ItemId, item.Quantity)
		if err != nil {
			log.Printf("Failed to create purchase item: %v", err)
			apierrors.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		purchaseItems = append(purchaseItems, pi)
	}

	//ユースケースの呼び出し
	canceledPurchase, err := h.cancelPurchaseUseCase.CancelPurchase(ctx, userID, purchaseItems)
	if err != nil {
		log.Printf("Failed to cancel purchase: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//レスポンスの作成
	res := toPostPurchaseResponse(canceledPurchase)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
