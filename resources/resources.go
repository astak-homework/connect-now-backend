package resources

import (
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

type LocalizableMessageKey string

const (
	MsgInvalidData         LocalizableMessageKey = "invalid_data"
	MsgUserNotFound        LocalizableMessageKey = "user_not_found"
	MsgProfileNotFound     LocalizableMessageKey = "profile_not_found"
	MsgIncorrectPassword   LocalizableMessageKey = "incorrect_password"
	MsgInternalServerError LocalizableMessageKey = "internal_server_error"
)

func GetMessage(context *gin.Context, key LocalizableMessageKey) string {
	return i18n.MustGetMessage(context, string(key))
}

func GetMessageInvalidData(context *gin.Context) string {
	return GetMessage(context, MsgInvalidData)
}

func GetMessageUserNotFound(context *gin.Context) string {
	return GetMessage(context, MsgUserNotFound)
}

func GetMessageProfileNotFound(context *gin.Context) string {
	return GetMessage(context, MsgProfileNotFound)
}

func GetMessageIncorrectPassword(context *gin.Context) string {
	return GetMessage(context, MsgIncorrectPassword)
}

func GetMessageInternalServerError(context *gin.Context) string {
	return GetMessage(context, MsgInternalServerError)
}
