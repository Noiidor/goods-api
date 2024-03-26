package errors

import (
	"goods-api/internal/dto"

	"github.com/gin-gonic/gin"
)

var (
	internalCodes = map[int]string{
		404: "3",
		400: "2",
		500: "1",
	}
)

func JsonError(ctx *gin.Context, httpCode int, err error) {
	ctx.JSON(httpCode, dto.ErrorResponse{
		Code:    internalCodes[httpCode],
		Message: err.Error(),
		Details: nil,
	})
}
