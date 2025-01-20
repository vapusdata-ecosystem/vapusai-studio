package pkgs

import "time"

const (
	ACCESS_TOKEN                = "access_token"
	REFRESH_TOKEN               = "refresh_token"
	USER_PROFILE                = "user_profile"
	ROUND_TRIP_HEADER           = "www-authentication"
	ROUND_TRIP_HEADER_VAL_TMPLT = "Bearer realm='%s'"
	CURRENT_DATAMARKETPLACE     = "current_datamarketplace"
	ID_TOKEN                    = "id_token"
)

var (
	COOKIE_LIKE_LEVEL_1 = time.Now().Add(2 * time.Hour)
)

const (
	IDEN = "Identifier"
	SVCS = "services"
)
