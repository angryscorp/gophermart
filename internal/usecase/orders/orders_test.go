package orders

import (
	"context"
	"errors"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository/mocks"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type OrdersTestSuite struct {
	suite.Suite
	ctx            context.Context
	orders         Orders
	repositoryMock *mocks.OrdersMock
	requestChan    chan string
	responseChan   chan *model.Accrual
}

func Test_Orders_TestSuite(t *testing.T) {
	suite.Run(t, &OrdersTestSuite{})
}

func (suite *OrdersTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.repositoryMock = &mocks.OrdersMock{}
	suite.requestChan = make(chan string, 1)
	suite.responseChan = make(chan *model.Accrual, 1)
	suite.orders = New(suite.repositoryMock, suite.requestChan, suite.responseChan, zerolog.Nop())
}

func (suite *OrdersTestSuite) TearDownTest() {
	suite.repositoryMock.AssertExpectations(suite.T())
	suite.requestChan = nil
	suite.responseChan = nil
}

func (suite *OrdersTestSuite) Test_AllOrders_ReturnsResult_InCaseOfNoErrors() {
	// Arrange
	userID := uuid.New()
	order := model.Order{
		Number:     "1234567890",
		UserID:     userID,
		Status:     model.OrderStatusProcessing,
		UploadedAt: time.Now(),
	}
	suite.repositoryMock.
		On("AllOrders", suite.ctx, userID).
		Return([]model.Order{order}, nil).
		Once()

	// Act
	res, err := suite.orders.AllOrders(suite.ctx, userID)

	// Assert
	suite.Nil(err)
	suite.Equal([]model.Order{order}, res)
}

func (suite *OrdersTestSuite) Test_Balance_ReturnsError_InCaseOfRepositoryError() {
	// Arrange
	userID := uuid.New()
	expectedError := errors.New("repository error")
	suite.repositoryMock.
		On("AllOrders", suite.ctx, userID).
		Return([]model.Order{}, expectedError).
		Once()

	// Act
	_, err := suite.orders.AllOrders(suite.ctx, userID)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, expectedError)
}

func (suite *OrdersTestSuite) Test_UploadOrder_ReturnsError_InCaseOfOrderNumberIsEmpty() {
	// Arrange
	userID := uuid.New()
	orderNumber := ""

	// Act
	err := suite.orders.UploadOrder(suite.ctx, orderNumber, userID)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, usecase.ErrOrderNumberIsInvalid)
}

func (suite *OrdersTestSuite) Test_UploadOrder_ReturnsError_InCaseOfOrderNumberIsInvalid() {
	// Arrange
	userID := uuid.New()
	orderNumber := "1234567891" // Invalid Luhn number

	// Act
	err := suite.orders.UploadOrder(suite.ctx, orderNumber, userID)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, usecase.ErrOrderNumberIsInvalid)
}

func (suite *OrdersTestSuite) Test_UploadOrder_ReturnsError_InCaseOfRepositoryError() {
	// Arrange
	userID := uuid.New()
	orderNumber := "79927398713" // Valid Luhn number
	expectedError := errors.New("database error")

	suite.repositoryMock.
		On("OrderInfoForUpdate", suite.ctx, orderNumber).
		Return(nil, expectedError).
		Once()

	// Act
	err := suite.orders.UploadOrder(suite.ctx, orderNumber, userID)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, model.ErrUnknownInternalError)
}

func (suite *OrdersTestSuite) Test_UploadOrder_ReturnsError_InCaseOfOrderAlreadyUploadedBySameUser() {
	// Arrange
	userID := uuid.New()
	orderNumber := "79927398713" // Valid Luhn number
	existingOrder := &model.Order{
		Number: orderNumber,
		UserID: userID,
		Status: model.OrderStatusProcessing,
	}

	suite.repositoryMock.
		On("OrderInfoForUpdate", suite.ctx, orderNumber).
		Return(existingOrder, nil).
		Once()

	// Act
	err := suite.orders.UploadOrder(suite.ctx, orderNumber, userID)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, usecase.ErrOrderIsAlreadyUploaded)
}

func (suite *OrdersTestSuite) Test_UploadOrder_ReturnsError_InCaseOfOrderUploadedByAnotherUser() {
	// Arrange
	userID := uuid.New()
	anotherUserID := uuid.New()
	orderNumber := "79927398713" // Valid Luhn number
	existingOrder := &model.Order{
		Number: orderNumber,
		UserID: anotherUserID,
		Status: model.OrderStatusProcessing,
	}

	suite.repositoryMock.
		On("OrderInfoForUpdate", suite.ctx, orderNumber).
		Return(existingOrder, nil).
		Once()

	// Act
	err := suite.orders.UploadOrder(suite.ctx, orderNumber, userID)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, usecase.ErrOrderWasUploadedAnotherUser)
}

func (suite *OrdersTestSuite) Test_UploadOrder_ReturnsError_InCaseOfCreatingOrderFailed() {
	// Arrange
	userID := uuid.New()
	orderNumber := "79927398713" // Valid Luhn number
	expectedError := errors.New("create order error")

	suite.repositoryMock.
		On("OrderInfoForUpdate", suite.ctx, orderNumber).
		Return(nil, nil).
		Once()

	suite.repositoryMock.
		On("CreateOrder", suite.ctx, model.NewOrder(orderNumber, userID)).
		Return(expectedError).
		Once()

	// Act
	err := suite.orders.UploadOrder(suite.ctx, orderNumber, userID)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, model.ErrUnknownInternalError)
}

func (suite *OrdersTestSuite) Test_UploadOrder_Success_InCaseOfNoErrors() {
	// Arrange
	userID := uuid.New()
	orderNumber := "79927398713" // Valid Luhn number

	suite.repositoryMock.
		On("OrderInfoForUpdate", suite.ctx, orderNumber).
		Return(nil, nil).
		Once()

	suite.repositoryMock.
		On("CreateOrder", suite.ctx, model.NewOrder(orderNumber, userID)).
		Return(nil).
		Once()

	// Act
	err := suite.orders.UploadOrder(suite.ctx, orderNumber, userID)

	// Assert
	suite.NoError(err)

	//Verify that order number is sent to request channel
	select {
	case sentOrderNumber := <-suite.requestChan:
		suite.Equal(orderNumber, sentOrderNumber)
	case <-time.After(100 * time.Millisecond):
		suite.Fail("Expected order number to be sent to request channel")
	}
}

func (suite *OrdersTestSuite) Test_newOrder_ReturnsProcessedOrder_InCaseOfAccrualStatusProcessed() {
	// Arrange
	accrualValue := 100.50
	accrual := &model.Accrual{
		Order:   "1234567890",
		Status:  model.AccrualStatusProcessed,
		Accrual: &accrualValue,
	}

	// Act
	order := newOrder(accrual)

	// Assert
	suite.NotNil(order)
	suite.Equal("1234567890", order.Number)
	suite.Equal(model.OrderStatusProcessed, order.Status)
	suite.Equal(&accrualValue, order.Accrual)
}

func (suite *OrdersTestSuite) Test_newOrder_ReturnsInvalidOrder_InCaseOfAccrualStatusInvalid() {
	// Arrange
	accrualValue := 100.50
	accrual := &model.Accrual{
		Order:   "1234567890",
		Status:  model.AccrualStatusInvalid,
		Accrual: &accrualValue,
	}

	// Act
	order := newOrder(accrual)

	// Assert
	suite.NotNil(order)
	suite.Equal("1234567890", order.Number)
	suite.Equal(model.OrderStatusInvalid, order.Status)
	suite.Nil(order.Accrual)
}

func (suite *OrdersTestSuite) Test_newOrder_ReturnsNil_InCaseOfAccrualStatusProcessing() {
	// Arrange
	accrualValue := 100.50
	accrual := &model.Accrual{
		Order:   "1234567890",
		Status:  "PROCESSING",
		Accrual: &accrualValue,
	}

	// Act
	order := newOrder(accrual)

	// Assert
	suite.Nil(order)
}

func (suite *OrdersTestSuite) Test_listenResponses_UpdatesOrder_InCaseOfValidAccrualReceived() {
	// Arrange
	accrualValue := 100.50
	accrual := &model.Accrual{
		Order:   "1234567890",
		Status:  model.AccrualStatusProcessed,
		Accrual: &accrualValue,
	}
	expectedOrder := model.Order{
		Number:  "1234567890",
		Status:  model.OrderStatusProcessed,
		Accrual: &accrualValue,
	}

	var receivedOrder model.Order
	suite.repositoryMock.
		On("UpdateOrder", suite.ctx, expectedOrder).
		Run(func(args mock.Arguments) {
			receivedOrder = args[1].(model.Order)
		}).
		Return(nil).
		Once()

	// Act
	go func() {
		suite.responseChan <- accrual
	}()

	// Wait a bit for the goroutine to process
	time.Sleep(100 * time.Millisecond)

	// Assert
	suite.Equal(expectedOrder, receivedOrder)
}

func (suite *OrdersTestSuite) Test_listenResponses_HandlesUpdateOrderError() {
	// Arrange
	accrualValue := 100.50
	accrual := &model.Accrual{
		Order:   "1234567890",
		Status:  model.AccrualStatusProcessed,
		Accrual: &accrualValue,
	}
	expectedOrder := model.Order{
		Number:  "1234567890",
		Status:  model.OrderStatusProcessed,
		Accrual: &accrualValue,
	}
	expectedError := errors.New("update error")

	var receivedOrder model.Order
	suite.repositoryMock.
		On("UpdateOrder", suite.ctx, expectedOrder).
		Run(func(args mock.Arguments) {
			receivedOrder = args[1].(model.Order)
		}).
		Return(expectedError).
		Once()

	// Act
	go func() {
		suite.responseChan <- accrual
	}()

	// Wait a bit for the goroutine to process
	time.Sleep(100 * time.Millisecond)

	// Assert
	suite.Equal(expectedOrder, receivedOrder)
}
