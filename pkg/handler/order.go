package handler

import (
	"encoding/json"
	"market/pkg/model"
	"market/pkg/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// @Summary	Create order
// @Security	ApiKeyAuth
// @Tags		order
// @ID			create-order
// @Product	json
// @Success	200		{object}	getOrdersResponse
// @Failure	400,404	{object}	errorResponse
// @Failure	500		{object}	errorResponse
// @Failure	default	{object}	errorResponse
// @Router		/api/order [get]
func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	order := model.Order{
		CreatedAt:   time.Now(),
		DeliveredAt: time.Now().Add(4 * 24 * time.Hour),
	}

	lastID, err := h.Services.Order.Create(sess.UserID, order)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Order was created with id LastInsertId: %v", lastID)

	q := model.OrderQueryInput{
		QueryInput: model.QueryInput{
			Limit:     defaultLimit,
			Offset:    defaultOffset,
			SortBy:    defaultSortField,
			SortOrder: model.DESCENDING,
		},
	}

	orders, err := h.Services.Order.GetAll(sess.UserID, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetOrdersResponse(w, orders, http.StatusCreated)
}

// @Summary	Get order
// @Security	ApiKeyAuth
// @Tags		order
// @ID			get-order
// @Product	json
// @Param		orderId	path		integer	true	"ID of order to get"
// @Param   sort_by query   string false "sort by" Enums(views, price, created_at)
// @Param   sort_order query string false "sort order" Enums(asc, desc)
// @Param   limit   query int false "limit" Enums(10, 25, 50)
// @Param   page  query int false "page"
// @Success	201			{object}	model.Order
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/order/{orderId} [get]
func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	options, err := optionsFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["orderId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	selectedOrder, err := h.Services.Order.GetByID(sess.UserID, orderID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := model.ProductQueryInput{
		QueryInput: model.QueryInput{
			Limit:     options.Limit,
			Offset:    options.Offset,
			SortBy:    options.SortBy,
			SortOrder: options.SortOrder,
		},
	}

	selectedOrder.Products, err = h.Services.Order.GetProductsByOrderID(sess.UserID, orderID, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(selectedOrder)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

// @Summary	Get orders
// @Security	ApiKeyAuth
// @Tags		order
// @ID			get-orders
// @Product	json
// @Param   sort_by query   string false "sort by" Enums(created_at)
// @Param   sort_order query string false "sort order" Enums(asc, desc)
// @Param   limit   query int false "limit" Enums(10, 25, 50)
// @Param   page  query int false "page"
// @Success	200		{object}	getOrdersResponse
// @Failure	400,404	{object}	errorResponse
// @Failure	500		{object}	errorResponse
// @Failure	default	{object}	errorResponse
// @Router		/api/orders [get]
func (h *Handler) getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	options, err := optionsFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	orderQuery := model.OrderQueryInput{
		QueryInput: model.QueryInput{
			Limit:     options.Limit,
			Offset:    options.Offset,
			SortBy:    options.SortBy,
			SortOrder: options.SortOrder,
		},
	}

	orders, err := h.Services.Order.GetAll(sess.UserID, orderQuery)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetOrdersResponse(w, orders, http.StatusOK)
}
