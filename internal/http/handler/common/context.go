package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(ctx *gin.Context) (*uuid.UUID, error) {
	userID, exists := ctx.Get("userID")
	if !exists {
		return nil, errors.New("user is not authenticated")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return nil, errors.New("invalid user ID format")
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return &userUUID, nil
}
