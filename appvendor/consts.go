package appvendor

import (
	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss/confirm"
)

const (
	ErrorGeneral              = 102017000
	ErrorAccountNotConfirmed  = 102017001
	ErrorAccountBeingLocked   = 102017002
	ErrorAccountCannotConfirm = 102017003

	PropPrimaryID    = "primaryID"
	PropEmail        = authboss.StoreEmail
	PropPassword     = authboss.StorePassword
	PropFullName     = "fullname"
	PropConfirmToken = confirm.StoreConfirmToken
	PropConfirmed    = confirm.StoreConfirmed
	PropJwtError     = "error"
)
