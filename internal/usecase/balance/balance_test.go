package balance

import (
	"context"
	"errors"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository/mocks"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type BalanceTestSuite struct {
	suite.Suite
	ctx            context.Context
	balance        Balance
	repositoryMock *mocks.BalanceMock
}

func Test_Balance_TestSuite(t *testing.T) {
	suite.Run(t, &BalanceTestSuite{})
}

func (suite *BalanceTestSuite) SetupTest() {
	suite.repositoryMock = &mocks.BalanceMock{}
	suite.ctx = context.Background()
	suite.balance = New(suite.repositoryMock)
}

func (suite *BalanceTestSuite) TearDownTest() {
	suite.repositoryMock.AssertExpectations(suite.T())
}

func (suite *BalanceTestSuite) Test_Balance_ReturnsResult_InCaseOfNoErrors() {
	// Arrange
	userID := uuid.New()
	suite.repositoryMock.
		On("Balance", suite.ctx, userID).
		Return(model.Balance{Current: 12.3, Withdrawn: 45.6}, nil).
		Once()

	// Act
	res, err := suite.balance.Balance(suite.ctx, userID)

	// Assert
	suite.Nil(err)
	suite.Equal(12.3, res.Current)
	suite.Equal(45.6, res.Withdrawn)
}

func (suite *BalanceTestSuite) Test_Balance_ReturnsError_InCaseOfRepositoryError() {
	// Arrange
	userID := uuid.New()
	expectedError := errors.New("repository error")
	suite.repositoryMock.
		On("Balance", suite.ctx, userID).
		Return(model.Balance{}, expectedError).
		Once()

	// Act
	_, err := suite.balance.Balance(suite.ctx, userID)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, expectedError)
}

func (suite *BalanceTestSuite) Test_Withdraw_ReturnsError_WhenOrderNumberIsEmpty() {
	// Arrange
	userID := uuid.New()
	orderNumber := ""
	amount := 10.0

	// Act
	err := suite.balance.Withdraw(suite.ctx, userID, orderNumber, amount)

	// Assert
	suite.Error(err)
	suite.Equal(usecase.ErrOrderNumberIsInvalid, err)
}

func (suite *BalanceTestSuite) Test_Withdraw_ReturnsError_WhenOrderNumberIsInvalid() {
	// Arrange
	userID := uuid.New()
	orderNumber := "123456789" // Invalid Luhn number
	amount := 10.0

	// Act
	err := suite.balance.Withdraw(suite.ctx, userID, orderNumber, amount)

	// Assert
	suite.Error(err)
	suite.Equal(usecase.ErrOrderNumberIsInvalid, err)
}

func (suite *BalanceTestSuite) Test_Withdraw_ReturnsResult_InCaseOfNoErrors() {
	// Arrange
	userID := uuid.New()
	orderNumber := "4561261212345467" // Valid Luhn number
	amount := 10.0
	suite.repositoryMock.
		On("Withdraw", suite.ctx, userID, orderNumber, amount).
		Return(nil).
		Once()

	// Act
	err := suite.balance.Withdraw(suite.ctx, userID, orderNumber, amount)

	// Assert
	suite.Nil(err)
}

func (suite *BalanceTestSuite) Test_Withdraw_ReturnsError_InCaseOfRepositoryError() {
	// Arrange
	userID := uuid.New()
	orderNumber := "4561261212345467" // Valid Luhn number
	amount := 10.0
	expectedError := errors.New("repository error")
	suite.repositoryMock.
		On("Withdraw", suite.ctx, userID, orderNumber, amount).
		Return(expectedError).
		Once()

	// Act
	err := suite.balance.Withdraw(suite.ctx, userID, orderNumber, amount)

	// Assert
	suite.Error(err)
	suite.ErrorIs(err, expectedError)
}

func (suite *BalanceTestSuite) Test_AllWithdrawals_ReturnsResult_InCaseOfNoErrors() {
	// Arrange
	userID := uuid.New()
	expectedWithdrawals := []model.Withdrawal{
		{
			Order:       "1234567890",
			Sum:         100.0,
			ProcessedAt: time.Now(),
		},
		{
			Order:       "0987654321",
			Sum:         50.0,
			ProcessedAt: time.Now(),
		},
	}
	suite.repositoryMock.
		On("WithdrawalHistory", suite.ctx, userID).
		Return(expectedWithdrawals, nil).
		Once()

	// Act
	result, err := suite.balance.AllWithdrawals(suite.ctx, userID)

	// Assert
	suite.Nil(err)
	suite.Equal(expectedWithdrawals, result)
}

func (suite *BalanceTestSuite) Test_AllWithdrawals_ReturnsError_InCaseOfRepositoryError() {
	// Arrange
	userID := uuid.New()
	expectedError := errors.New("repository error")
	suite.repositoryMock.
		On("WithdrawalHistory", suite.ctx, userID).
		Return([]model.Withdrawal{}, expectedError).
		Once()

	// Act
	_, err := suite.balance.AllWithdrawals(suite.ctx, userID)

	// Assert
	suite.ErrorIs(err, expectedError)
}
