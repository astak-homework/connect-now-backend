package errors

import (
	"errors"
	"net/http"

	"github.com/astak-homework/connect-now-backend/resources"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func HttpResponseErrorHandler(c *gin.Context) {
	c.Next()

	httpResponseError := &HttpResponseError{}
	for _, err := range c.Errors {
		if ok := errors.As(err, httpResponseError); ok {
			log.Error().Err(err).Msg(httpResponseError.internalErrorMessage)
			switch httpResponseError.status {
			case http.StatusInternalServerError, http.StatusServiceUnavailable:
				c.Header("Retry-After", "120")
				c.IndentedJSON(httpResponseError.status, gin.H{
					"message":    resources.GetMessage(c, httpResponseError.resourceKey),
					"request_id": requestid.Get(c),
					"code":       httpResponseError.status,
				})
			case http.StatusBadRequest:
				validationErrors := validator.ValidationErrors{}
				trans, _ := resources.Translator.GetTranslator(c.GetHeader("Accept-Language"))
				if ok := errors.As(httpResponseError.originalError, &validationErrors); ok {
					fieldErrors := make(map[string]string)
					for _, fieldError := range validationErrors {
						fieldErrors[fieldError.Field()] = fieldError.Translate(trans)
					}
					c.JSON(httpResponseError.status, gin.H{
						"message":          resources.GetMessage(c, httpResponseError.resourceKey),
						"validationErrors": fieldErrors,
					})
				} else {
					c.IndentedJSON(httpResponseError.status, gin.H{"message": resources.GetMessage(c, httpResponseError.resourceKey)})
				}
			default:
				c.IndentedJSON(httpResponseError.status, gin.H{"message": resources.GetMessage(c, httpResponseError.resourceKey)})
			}
			break
		}
	}
}
