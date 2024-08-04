package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"thirdparty-service/environment"
	"thirdparty-service/mocks"
	"thirdparty-service/models"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_HttpHandler_PostPaymentHandlerDebit(t *testing.T) {
	const (
		success = iota
		errorReadingRequestBody
		errorBadRequest
		errorGettingAccount
		errorInsufficientBalance
		errorMakingWithdrawal
		errorCreatingTransaction
		errorUpdatingAccountBalance
	)

	testCases := []struct {
		name     string
		testType int
	}{
		{
			name:     "Test success",
			testType: success,
		},

		{
			name:     "Test error reading request body",
			testType: errorReadingRequestBody,
		},

		{
			name:     "Test error with bad request",
			testType: errorBadRequest,
		},

		{
			name:     "Test error fetching account",
			testType: errorGettingAccount,
		},

		{
			name:     "Test error insufficient balance",
			testType: errorInsufficientBalance,
		},

		{
			name:     "Test error making withdrawal",
			testType: errorMakingWithdrawal,
		},

		{
			name:     "Test error creating transaction record",
			testType: errorCreatingTransaction,
		},

		{
			name:     "Test error while updating account balance",
			testType: errorUpdatingAccountBalance,
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	cfg := &environment.Config{}

	mockDataStore := mocks.NewMockMongoDBStore(controller)

	handler := NewHTTPHandler(cfg, mockDataStore)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			switch testCase.testType {
			case success:
				mockRequest := models.PostPaymentRequestPayload{
					AccountId: "acc_001",
					Reference: "ref-001",
					Amount:    1.50,
				}

				mockPayload, err := json.Marshal(mockRequest)
				assert.NoError(t, err)

				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   10,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/third-party/payments", bytes.NewBuffer(mockPayload))

				query := r.URL.Query()
				query.Set("type", "debit")
				r.URL.RawQuery = query.Encode()

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(&mockAccount, nil)

				mockDataStore.
					EXPECT().
					CreateTransaction(gomock.Any()).
					Return(nil)

				newBalance := mockAccount.Balance - mockRequest.Amount

				mockDataStore.
					EXPECT().
					UpdateAccountBalance(mockRequest.AccountId, newBalance).
					Return(nil)

				handler.PostPaymentHandler(w, r)
				assert.Equal(t, http.StatusOK, w.Code)

			case errorReadingRequestBody:
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/third-party/payments", errReader(0))

				handler.PostPaymentHandler(w, r)
				assert.Equal(t, http.StatusBadRequest, w.Code)

			case errorBadRequest:

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/third-party/payments", bytes.NewBuffer(nil))

				handler.PostPaymentHandler(w, r)
				assert.Equal(t, http.StatusBadRequest, w.Code)

			case errorGettingAccount:
				mockRequest := models.PostPaymentRequestPayload{
					AccountId: "acc_001",
					Reference: "ref-001",
					Amount:    1.50,
				}

				mockPayload, err := json.Marshal(mockRequest)
				assert.NoError(t, err)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/third-party/payments", bytes.NewBuffer(mockPayload))

				query := r.URL.Query()
				query.Set("type", "debit")
				r.URL.RawQuery = query.Encode()

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(nil, errors.New(""))

				handler.PostPaymentHandler(w, r)
				assert.Equal(t, http.StatusNotFound, w.Code)

			case errorInsufficientBalance:
				mockRequest := models.PostPaymentRequestPayload{
					AccountId: "acc_001",
					Reference: "ref-001",
					Amount:    10,
				}

				mockPayload, err := json.Marshal(mockRequest)
				assert.NoError(t, err)
				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   1,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/third-party/payments", bytes.NewBuffer(mockPayload))

				query := r.URL.Query()
				query.Set("type", "debit")
				r.URL.RawQuery = query.Encode()

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(&mockAccount, nil)

				handler.PostPaymentHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorCreatingTransaction:
				mockRequest := models.PostPaymentRequestPayload{
					AccountId: "acc_001",
					Reference: "ref-001",
					Amount:    1.50,
				}

				mockPayload, err := json.Marshal(mockRequest)
				assert.NoError(t, err)

				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   10,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/third-party/payments", bytes.NewBuffer(mockPayload))

				query := r.URL.Query()
				query.Set("type", "debit")
				r.URL.RawQuery = query.Encode()

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(&mockAccount, nil)

				mockDataStore.
					EXPECT().
					CreateTransaction(gomock.Any()).
					Return(errors.New(""))

				handler.PostPaymentHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorUpdatingAccountBalance:
				mockRequest := models.PostPaymentRequestPayload{
					AccountId: "acc_001",
					Reference: "ref-001",
					Amount:    1.50,
				}

				mockPayload, err := json.Marshal(mockRequest)
				assert.NoError(t, err)

				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   10,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/third-party/payments", bytes.NewBuffer(mockPayload))

				query := r.URL.Query()
				query.Set("type", "debit")
				r.URL.RawQuery = query.Encode()

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(&mockAccount, nil)

				mockDataStore.
					EXPECT().
					CreateTransaction(gomock.Any()).
					Return(nil)

				newBalance := mockAccount.Balance - mockRequest.Amount

				mockDataStore.
					EXPECT().
					UpdateAccountBalance(mockRequest.AccountId, newBalance).
					Return(errors.New(""))

				handler.PostPaymentHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			}
		})
	}
}

func Test_HttpHandler_PostPaymentHandlerCredit(t *testing.T) {
	const (
		success = iota
		errorMakingDeposit
		errorCreatingTransaction
		errorUpdatingAccountBalance
	)

	testCases := []struct {
		name     string
		testType int
	}{
		{
			name:     "Test success",
			testType: success,
		},

		{
			name:     "Test error making deposit",
			testType: errorMakingDeposit,
		},

		{
			name:     "Test error creating transaction record",
			testType: errorCreatingTransaction,
		},

		{
			name:     "Test error while updating account balance",
			testType: errorUpdatingAccountBalance,
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	cfg := &environment.Config{}

	mockDataStore := mocks.NewMockMongoDBStore(controller)

	handler := NewHTTPHandler(cfg, mockDataStore)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			switch testCase.testType {
			case success:
				mockRequest := models.PostPaymentRequestPayload{
					AccountId: "acc_001",
					Reference: "ref-001",
					Amount:    1.50,
				}

				mockPayload, err := json.Marshal(mockRequest)
				assert.NoError(t, err)

				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   10,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/third-party/payments", bytes.NewBuffer(mockPayload))

				query := r.URL.Query()
				query.Set("type", "credit")
				r.URL.RawQuery = query.Encode()

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(&mockAccount, nil)

				mockDataStore.
					EXPECT().
					CreateTransaction(gomock.Any()).
					Return(nil)

				newBalance := mockAccount.Balance + mockRequest.Amount

				mockDataStore.
					EXPECT().
					UpdateAccountBalance(mockRequest.AccountId, newBalance).
					Return(nil)

				handler.PostPaymentHandler(w, r)
				assert.Equal(t, http.StatusOK, w.Code)

			case errorCreatingTransaction:
				mockRequest := models.PostPaymentRequestPayload{
					AccountId: "acc_001",
					Reference: "ref-001",
					Amount:    1.50,
				}

				mockPayload, err := json.Marshal(mockRequest)
				assert.NoError(t, err)

				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   10,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/third-party/payments", bytes.NewBuffer(mockPayload))

				query := r.URL.Query()
				query.Set("type", "credit")
				r.URL.RawQuery = query.Encode()

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(&mockAccount, nil)

				mockDataStore.
					EXPECT().
					CreateTransaction(gomock.Any()).
					Return(errors.New(""))

				handler.PostPaymentHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorUpdatingAccountBalance:
				mockRequest := models.PostPaymentRequestPayload{
					AccountId: "acc_001",
					Reference: "ref-001",
					Amount:    1.50,
				}

				mockPayload, err := json.Marshal(mockRequest)
				assert.NoError(t, err)

				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   10,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/third-party/payments", bytes.NewBuffer(mockPayload))

				query := r.URL.Query()
				query.Set("type", "credit")
				r.URL.RawQuery = query.Encode()

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(&mockAccount, nil)

				mockDataStore.
					EXPECT().
					CreateTransaction(gomock.Any()).
					Return(nil)

				newBalance := mockAccount.Balance + mockRequest.Amount

				mockDataStore.
					EXPECT().
					UpdateAccountBalance(mockRequest.AccountId, newBalance).
					Return(errors.New(""))

				handler.PostPaymentHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			}
		})
	}
}

func Test_HttpHandler_GetPaymentHandler(t *testing.T) {
	const (
		success = iota
		errorOccurred
	)

	testCases := []struct {
		name      string
		reference string
		testType  int
	}{
		{
			name:      "Test success",
			reference: "ref-010",
			testType:  success,
		},

		{
			name:      "Test error invalid reference",
			reference: "invalid-ref",
			testType:  errorOccurred,
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	cfg := &environment.Config{}

	mockDataStore := mocks.NewMockMongoDBStore(controller)

	handler := NewHTTPHandler(cfg, mockDataStore)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			switch testCase.testType {
			case success:
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/third-party/payments", nil)

				ctx := chi.NewRouteContext()
				ctx.URLParams.Add("reference", testCase.reference)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

				mockDataStore.
					EXPECT().
					GetPaymentByReferenceId(testCase.reference).
					Return(&models.Transaction{
						AccountID: "accountId",
						Reference: testCase.reference,
						Amount:    10,
						Type:      models.DEBIT,
					}, nil)

				handler.GetPaymentHandler(w, r)
				assert.Equal(t, http.StatusOK, w.Code)

			case errorOccurred:
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/third-party/payments", nil)

				ctx := chi.NewRouteContext()
				ctx.URLParams.Add("reference", testCase.reference)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

				mockDataStore.
					EXPECT().
					GetPaymentByReferenceId(testCase.reference).
					Return(nil, errors.New(""))

				handler.GetPaymentHandler(w, r)
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error body reader")
}
