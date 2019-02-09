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
	ErrorRefreshTokenInvalid
	ErrorRequestBodyTooLarge

	PropPrimaryID    = "primaryID"
	PropUserID       = "userID"
	PropEmail        = authboss.StoreEmail
	PropPassword     = authboss.StorePassword
	PropFullName     = "fullname"
	PropNickName     = "nickname"
	PropDateOfBirth  = "date_of_birth"
	PropGender       = "gender"
	PropJob          = "job"
	PropDeviceName   = "device_name"
	PropConfirmToken = confirm.StoreConfirmToken
	PropConfirmed    = confirm.StoreConfirmed
	PropJwtError     = "error"
	PropShortAbout   = "short_about"
	PropLivingAt     = "living_at"
	PropWorkingAt    = "working_at"
	PropHomeTown     = "home_town"
	PropStatus       = "status"
	PropWeight       = "weight"
	PropHeight       = "height"

	JFieldToken        = "access_token"
	JFieldRefreshToken = "refresh_token"
	JFieldExpireIn     = "expire_in"
	JFieldUserProfile  = "user_profile"
)
