package cp_api

import (
	"net/http"
)

type SvrApi struct {
	SvrName 	string
	Method		string
	Port		int
	UriQuery	string
	BodyFormat	string
}

const (
	SVRAPI_SSO_SESSION_CHECK = "sso_check"
	SVRAPI_SHOPEE_FILE_MILE_SHIP_ORDER = "shopee_ship_order"
	SVRAPI_SHOPEE_GET_TRACK_INFO = "get_track_info"
)

var SvrApiList = map[string]SvrApi{
	SVRAPI_SSO_SESSION_CHECK: {
		"cangboss",
		http.MethodPost,
		25020,
		"/api/v2/common/cangboss/sso/check",
		"",
	},
	SVRAPI_SHOPEE_FILE_MILE_SHIP_ORDER: {
		"shopee",
		http.MethodPost,
		25021,
		"/api/v2/common/shopee/order/first_mile_ship_order",
		"",
	},
	SVRAPI_SHOPEE_GET_TRACK_INFO: {
		"shopee",
		http.MethodGet,
		25021,
		"/api/v2/admin/shopee/order/get_track_info",
		"",
	},
}

