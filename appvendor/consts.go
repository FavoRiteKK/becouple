package appvendor

import (
	"gopkg.in/authboss.v1"
	"gopkg.in/authboss.v1/confirm"
)

const (
	ErrorGeneral             = 102017000
	ErrorAccountNotConfirmed = 102017001
	ErrorAccountBeingLocked  = 102017002

	PropPrimaryID    = "primaryID"
	PropEmail        = authboss.StoreEmail
	PropPassword     = authboss.StorePassword
	PropFullName     = "fullname"
	PropConfirmToken = confirm.StoreConfirmToken
)
