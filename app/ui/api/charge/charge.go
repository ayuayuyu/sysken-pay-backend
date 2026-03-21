package charge

import (
	"encoding/json"
	"log"
	"net/http"
	apierrors "sysken-pay-api/app/ui/api/pkg/errors"
	"sysken-pay-api/app/usecase/charge"

	"github.com/go-chi/chi/v5"
)

//TODO APIリクエストからデータを整形してユースケースに渡すための構造体を作成する

type Handler interface {
	ChargeAmount(w http.ResponseWriter, r *http.Request)
	ChargeCancel(w http.ResponseWriter, r *http.Request)
}

func NewChargeHandler(chargeAmountUseCase charge.ChargeAmountUseCase, chargeCancelUseCase charge.ChargeCancelUseCase) Handler {
	return &chargeHandlerImpl{
		chargeAmountUseCase: chargeAmountUseCase,
		chargeCancelUseCase: chargeCancelUseCase,
	}
}

var _ Handler = (*chargeHandlerImpl)(nil)

type chargeHandlerImpl struct {
	chargeAmountUseCase charge.ChargeAmountUseCase
	chargeCancelUseCase charge.ChargeCancelUseCase
}

func (h *chargeHandlerImpl) ChargeAmount(w http.ResponseWriter, r *http.Request) {

	var req PostChargeRequest
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
	//ユースケースの呼び出し
	chargedAmount, err := h.chargeAmountUseCase.ChargeAmount(ctx, userID, req.Amount)
	if err != nil {
		log.Printf("Failed to charge amount: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//レスポンスの作成
	res := toPostChargeResponse(chargedAmount)

	w.Header().Set("Content	-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Failed to encode response: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *chargeHandlerImpl) ChargeCancel(w http.ResponseWriter, r *http.Request) {

	var req PostChargeCancelRequest
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
	//ユースケースの呼び出し
	canceledCharge, err := h.chargeCancelUseCase.ChargeCancel(ctx, userID, req.Amount)
	if err != nil {
		log.Printf("Failed to cancel charge: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//レスポンスの作成
	res := toPostChargeCancelResponse(canceledCharge)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Failed to encode response: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
