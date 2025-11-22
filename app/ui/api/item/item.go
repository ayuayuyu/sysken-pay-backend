package item

import (
	"encoding/json"
	"log"
	"net/http"
	apierrors "sysken-pay-api/app/ui/api/pkg/errors"
	"sysken-pay-api/app/usecase/item"

	"github.com/go-chi/chi/v5"
)

//TODO APIリクエストからデータを整形してユースケースに渡すための構造体を作成する

type Handler interface {
	ResisterItem(w http.ResponseWriter, r *http.Request)
	UpdateItem(w http.ResponseWriter, r *http.Request)
	GetItemByJanCode(w http.ResponseWriter, r *http.Request)
	GetAllItems(w http.ResponseWriter, r *http.Request)
}

func NewItemHandler(registerItemUseCase item.RegisterItemUseCase, updateItemUseCase item.UpdateItemUseCase, findItemByJanCodeUseCase item.FindItemByJanCodeUseCase, getAllItemsUseCase item.GetAllItemsUseCase) Handler {
	return &itemHandlerImpl{
		registerItemUseCase:      registerItemUseCase,
		updateItemUseCase:        updateItemUseCase,
		findItemByJanCodeUseCase: findItemByJanCodeUseCase,
		getAllItemsUseCase:       getAllItemsUseCase,
	}
}

var _ Handler = (*itemHandlerImpl)(nil)

type itemHandlerImpl struct {
	registerItemUseCase      item.RegisterItemUseCase
	updateItemUseCase        item.UpdateItemUseCase
	findItemByJanCodeUseCase item.FindItemByJanCodeUseCase
	getAllItemsUseCase       item.GetAllItemsUseCase
}

func (h *itemHandlerImpl) ResisterItem(w http.ResponseWriter, r *http.Request) {

	var req PostItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	//ユースケースの呼び出し
	createdItem, err := h.registerItemUseCase.RegisterItem(ctx, req.JanCode, req.ItemName, req.Price)
	if err != nil {
		log.Printf("Failed to register item: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//レスポンスの作成
	res := toPostItemResponse(createdItem)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (h *itemHandlerImpl) UpdateItem(w http.ResponseWriter, r *http.Request) {

	var req PostItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	//ユースケースの呼び出し
	updatedItem, err := h.updateItemUseCase.UpdateItem(ctx, req.JanCode, req.ItemName, req.Price)
	if err != nil {
		log.Printf("Failed to update item: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//レスポンスの作成
	res := toUpdateItemResponse(updatedItem)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (h *itemHandlerImpl) GetItemByJanCode(w http.ResponseWriter, r *http.Request) {

	janCode := chi.URLParam(r, "jan_code")
	ctx := r.Context()
	//ユースケースの呼び出し
	foundItem, err := h.findItemByJanCodeUseCase.GetItemByJanCode(ctx, janCode)
	if err != nil {
		log.Printf("Failed to find item by jan code: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//レスポンスの作成
	res := toGetItemResponse(foundItem)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (h *itemHandlerImpl) GetAllItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//ユースケースの呼び出し
	allItems, err := h.getAllItemsUseCase.GetAllItems(ctx)
	if err != nil {
		log.Printf("Failed to get all items: %v", err)
		apierrors.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//レスポンスの作成
	res := toGetAllItemsResponse(allItems)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		apierrors.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
}
