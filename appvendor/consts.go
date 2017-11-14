package appvendor

import (
	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss/confirm"
)

const (
	ErrorGeneral uint = 102017000 + iota
	ErrorAccountNotConfirmed
	ErrorAccountBeingLocked
	ErrorAccountCannotConfirm
	ErrorAccountAlreadyInUsed
	ErrorAccountAuthorizedFailed

	ErrorTokenExpired
	ErrorTokenIssuedAt
	ErrorTokenNotValidYet

	PropPrimaryID    = "primaryID"
	PropEmail        = authboss.StoreEmail
	PropPassword     = authboss.StorePassword
	PropFullName     = "fullname"
	PropConfirmToken = confirm.StoreConfirmToken
	PropConfirmed    = confirm.StoreConfirmed
	PropJwtError     = "error"

	JFieldToken = "token"
)
