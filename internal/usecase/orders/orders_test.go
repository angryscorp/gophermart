package orders

import (
	"context"
	"errors"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository/mocks"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
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
	suite.requestChan = make(chan string)
	suite.responseChan = make(chan *model.Accrual)
	suite.orders = New(suite.repositoryMock, suite.requestChan, suite.responseChan, zerolog.Nop())
}

func (suite *OrdersTestSuite) TearDownTest() {
	suite.repositoryMock.AssertExpectations(suite.T())
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
