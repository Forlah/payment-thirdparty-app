package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"thirdparty-service/database"
	"thirdparty-service/environment"
	"thirdparty-service/models"
	"time"

	"github.com/go-chi/chi"
)

type HttpHandler struct {
	config       *environment.Config
	mongodbStore database.MongoDBStore
}

func NewHTTPHandler(config *environment.Config, store database.MongoDBStore) *HttpHandler {
	return &HttpHandler{config: config, mongodbStore: store}
}

func (handler *HttpHandler) responseWriter(w http.ResponseWriter, response interface{}, codes ...int) {
	statusCode := http.StatusOK
	if len(codes) > 0 {
		statusCode = codes[0]
	}
	_, ok := response.(*models.ErrorResponse)
	if ok {
		statusCode = http.StatusBadRequest
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshalling response")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		return
	}
}

func (handler *HttpHandler) PostPaymentHandler(w http.ResponseWriter, r *http.Request) {
	paymentType := r.URL.Query().Get("type")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("Error closing request body")
		}
	}()

	var payload models.PostPaymentRequestPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate account exist
	account, err := handler.mongodbStore.GetAccountByID(payload.AccountId)
	if err != nil {
		resp := models.ErrorResponse{
			ErrorMessage: "account not found",
		}
		handler.responseWriter(w, resp, http.StatusNotFound)
		return
	}

	switch paymentType {
	case "debit":
		// check balance
		if payload.Amount > float64(account.Balance) {
			resp := models.ErrorResponse{
				ErrorMessage: "insufficient funds",
			}
			handler.responseWriter(w, resp, http.StatusInternalServerError)
			return
		}

		newBalance := account.Balance - payload.Amount

		transaction := &models.Transaction{
			Reference: payload.Reference,
			AccountID: payload.AccountId,
			Amount:    payload.Amount,
			Type:      models.DEBIT,
			Status:    models.SUCCESS,
			CreatedAt: time.Now().Unix(),
		}

		err := handler.mongodbStore.CreateTransaction(transaction)
		if err != nil {
			// should never happen.
			resp := models.ErrorResponse{
				ErrorMessage: "error creating transaction record",
			}
			handler.responseWriter(w, resp, http.StatusInternalServerError)
			return
		}

		if err = handler.mongodbStore.UpdateAccountBalance(payload.AccountId, newBalance); err != nil {
			// should never happen.
			resp := models.ErrorResponse{
				ErrorMessage: "error updating balance",
			}
			handler.responseWriter(w, resp, http.StatusInternalServerError)
			return
		}

		resp := models.PaymentResponsePayload{
			AccountId: account.AccountID,
			Reference: payload.Reference,
			Amount:    payload.Amount,
		}

		handler.responseWriter(w, resp)

	case "credit":
		newBalance := account.Balance + payload.Amount

		transaction := &models.Transaction{
			Reference: payload.Reference,
			AccountID: payload.AccountId,
			Amount:    payload.Amount,
			Type:      models.CREDIT,
			Status:    models.SUCCESS,
			CreatedAt: time.Now().Unix(),
		}

		err := handler.mongodbStore.CreateTransaction(transaction)
		if err != nil {
			// should never happen.
			resp := models.ErrorResponse{
				ErrorMessage: "error creating transaction record",
			}
			handler.responseWriter(w, resp, http.StatusInternalServerError)
			return
		}

		if err = handler.mongodbStore.UpdateAccountBalance(payload.AccountId, newBalance); err != nil {
			// should never happen.
			resp := models.ErrorResponse{
				ErrorMessage: "error updating balance",
			}
			handler.responseWriter(w, resp, http.StatusInternalServerError)
			return
		}

		resp := models.PaymentResponsePayload{
			AccountId: account.AccountID,
			Reference: payload.Reference,
			Amount:    payload.Amount,
		}

		handler.responseWriter(w, resp)
	}
}

func (handler *HttpHandler) GetPaymentHandler(w http.ResponseWriter, r *http.Request) {
	reference := chi.URLParam(r, "reference")

	transaction, err := handler.mongodbStore.GetPaymentByReferenceId(reference)
	if err != nil {
		resp := models.ErrorResponse{
			ErrorMessage: "reference not found",
		}
		handler.responseWriter(w, resp, http.StatusNotFound)
		return
	}

	resp := models.PaymentResponsePayload{
		AccountId: transaction.AccountID,
		Reference: transaction.Reference,
		Amount:    transaction.Amount,
	}

	handler.responseWriter(w, resp)
}
