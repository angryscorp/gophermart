package auth

import (
	"context"
	"errors"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

type AuthTestSuite struct {
	suite.Suite
	ctx            context.Context
	auth           Auth
	repositoryMock *mocks.UsersMock
}

func Test_Auth_TestSuite(t *testing.T) {
	suite.Run(t, &AuthTestSuite{})
}

func (suite *AuthTestSuite) SetupTest() {
	suite.repositoryMock = &mocks.UsersMock{}
	suite.ctx = context.Background()
	suite.auth = New(suite.repositoryMock, "jwt-secret")
}

func (suite *AuthTestSuite) TearDownTest() {
	suite.repositoryMock.AssertExpectations(suite.T())
}

func (suite *AuthTestSuite) Test_SignUp_ReturnsResult_InCaseOfNoErrors() {
	// Arrange
	username := "username"
	password := "password"
	newUserID := uuid.New()
	suite.repositoryMock.
		On("CreateUser", suite.ctx, username, mock.AnythingOfType("string")).
		Return(&newUserID, nil).
		Once()

	// Act
	res, err := suite.auth.SignUp(suite.ctx, username, password)

	// Assert
	suite.Nil(err)
	suite.NotEmpty(res)
}

func (suite *AuthTestSuite) Test_SignUp_ReturnsError_InCaseOfCreateUserFails() {
	// Arrange
	username := "username"
	password := "password"
	expectedError := errors.New("repository error")
	suite.repositoryMock.
		On("CreateUser", suite.ctx, username, mock.AnythingOfType("string")).
		Return(nil, expectedError).
		Once()

	// Act
	res, err := suite.auth.SignUp(suite.ctx, username, password)

	// Assert
	suite.Error(err)
	suite.Empty(res)
	suite.ErrorIs(err, expectedError)
}

func (suite *AuthTestSuite) Test_SignIn_ReturnsResult_InCaseOfNoErrors() {
	// Arrange
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userID := uuid.New()
	user := &model.User{
		ID:           userID,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	suite.repositoryMock.
		On("UserData", suite.ctx, username).
		Return(user, nil).
		Once()

	// Act
	res, err := suite.auth.SignIn(suite.ctx, username, password)

	// Assert
	suite.Nil(err)
	suite.NotEmpty(res)
}

func (suite *AuthTestSuite) Test_SignIn_ReturnsError_InCaseUserNotFound() {
	// Arrange
	username := "username"
	password := "password"
	expectedError := errors.New("user not found")
	suite.repositoryMock.
		On("UserData", suite.ctx, username).
		Return(nil, expectedError).
		Once()

	// Act
	res, err := suite.auth.SignIn(suite.ctx, username, password)

	// Assert
	suite.Error(err)
	suite.Empty(res)
	suite.ErrorIs(err, expectedError)
}

func (suite *AuthTestSuite) Test_SignIn_ReturnsError_InCaseOfPasswordDoesNotMatch() {
	// Arrange
	username := "username"
	password := "password"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userID := uuid.New()
	user := &model.User{
		ID:           userID,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	suite.repositoryMock.
		On("UserData", suite.ctx, username).
		Return(user, nil).
		Once()

	// Act
	res, err := suite.auth.SignIn(suite.ctx, username, wrongPassword)

	// Assert
	suite.Error(err)
	suite.Empty(res)
}

func (suite *AuthTestSuite) Test_Validate_ReturnsToken_InCaseOfTokenIsValid() {
	// Arrange
	userID := uuid.New().String()
	claims := model.Token{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(suite.auth.jwtSecret)

	// Act
	result, err := suite.auth.Validate(tokenString)

	// Assert
	suite.Nil(err)
	suite.NotNil(result)
	suite.Equal(userID, result.UserID)
}

func (suite *AuthTestSuite) Test_Validate_ReturnsError_InCaseOfTokenIsInvalid() {
	// Arrange
	invalidToken := "invalid.token.string"

	// Act
	result, err := suite.auth.Validate(invalidToken)

	// Assert
	suite.Error(err)
	suite.Nil(result)
}

func (suite *AuthTestSuite) Test_Validate_ReturnsError_InCaseOfTokenIsExpired() {
	// Arrange
	userID := uuid.New().String()
	claims := model.Token{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired token
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(suite.auth.jwtSecret)

	// Act
	result, err := suite.auth.Validate(tokenString)

	// Assert
	suite.Error(err)
	suite.Nil(result)
}

func (suite *AuthTestSuite) Test_Validate_ReturnsError_InCaseOfTokenHasInvalidSignature() {
	// Arrange
	userID := uuid.New().String()
	claims := model.Token{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign with a different secret
	tokenString, _ := token.SignedString([]byte("different-secret"))

	// Act
	result, err := suite.auth.Validate(tokenString)

	// Assert
	suite.Error(err)
	suite.Nil(result)
}

func (suite *AuthTestSuite) Test_New_CreatesAuthInstance() {
	// Arrange
	jwtSecret := "test-secret"

	// Act
	auth := New(suite.repositoryMock, jwtSecret)

	// Assert
	suite.Equal(suite.repositoryMock, auth.repository)
	suite.Equal([]byte(jwtSecret), auth.jwtSecret)
}
