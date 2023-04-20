package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/incubus8/go/pkg/errors"
	"net/http"
	"time"
)

func AbortWithAPIError(c *gin.Context, err *errors.APIError) {
	c.Abort()
	err.RecordedAt = time.Now().UTC().Format(time.RFC3339)

	if err.StatusCode != 0 {
		c.JSON(err.StatusCode, err)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}

	if err.Err != nil {
		c.Error(err.Err)
	} else {
		c.Error(err)
	}
}

func AbortWithValidationError(c *gin.Context, errs []error) {
	c.Abort()

	strErrors := make([]string, len(errs))
	for i, err := range errs {
		strErrors[i] = err.Error()
	}

	ve := &errors.ValidationError{
		ValidationErrorReason: errors.ValidationErrorReason{
			Code:    "0001",
			Message: strErrors,
		},
		StatusCode: http.StatusBadRequest,
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusBadRequest, ve)
}
