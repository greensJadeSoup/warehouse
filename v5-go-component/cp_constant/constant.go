package cp_constant

//version
const ComponentVersion = "1.1.2"

//common
const (
	FALSE = 0
	TRUE = 1
)
//salt
const (
	PASSWORD_SALT = "b93efba3"
)

//conf
const (
	BaseConf = "conf/base.conf"
)

//app id
const (
	APPID_LOCAL = "ys_local@123587dcaf0"
	APPID_SERVER = "ys_server@1239d4749c0"
)

const (
	USER_TYPE_SUPER_MANAGER = "super_manager"
	USER_TYPE_MANAGER = "manager"
	USER_TYPE_SELLER = "seller"
	USER_TYPE_SERVICE = "service"
)

//common redis key
const (
	REDIS_KEY_SESSIONKEY = "[platform]session_key:"
	REDIS_KEY_SYNC_ORDER_FLAG = "[platform]sync_order_flag:"
	REDIS_KEY_SYNC_ITEM_FLAG = "[platform]sync_item_flag:"
	REDIS_KEY_OUTPUT_ORDER_FLAG = "[platform]output_order_flag:"
)

//common redis key expire time
const (
	REDIS_EXPIRE_SESSION_KEY = 3 //20min
	REDIS_EXPIRE_SYNC_ORDER_FLAG = 3
	REDIS_EXPIRE_SYNC_ITEM_AND_MODEL_FLAG = 3
	REDIS_EXPIRE_OUTPUT_ORDER_FLAG = 1
)

//log level
const (
	LevelPanic = "Panic"
	LevelFatal = "Fatal"
	LevelError = "Error"
	LevelWarning = "Warning"
	LevelNotice = "Notice"
	LevelInformational = "Info"
	LevelDebug = "Debug"
)

//trace level
type TracingLevel string
const (
	TracingLevelCritical TracingLevel = "critical"
	TracingLevelError TracingLevel = "error"
	TracingLevelInfo TracingLevel = "info"
	TracingLevelDebug TracingLevel = "debug"
)

//expire time
const (
	EXPIRE_TIME_TIMESTAMP = 300 // 5m
	EXPIRE_TIME_ACCESSTOKEN = 7200 // 2 hours
)

//gin
const (
	SPECIAL_ID = "special_id"
	REQUEST_ID = "request_id"
)

//http header
const (
	HTTP_HEADER_TS = "x-wh-ts"
	HTTP_HEADER_NONCE = "x-wh-nonce"
	HTTP_HEADER_SIGN = "x-wh-sign"
	HTTP_HEADER_APPID = "x-wh-appid"
	HTTP_HEADER_ACCESS_TOKEN = "x-wh-at"
	HTTP_HEADER_SESSION_KEY = "x-wh-sessionkey"

	HTTP_HEADER_CHAIN_ID = "x-wh-chainid"
	HTTP_HEADER_CHAIN_LEVEL = "x-wh-chainlevel"
	HTTP_HEADER_SESSION_INFO = "x-wh-sessioninfo"
)

//response code
const (
	/*************** 通用错误码 **************/
	// 正常响应
	RESPONSE_CODE_OK = 10000
	// 常规错误
	RESPONSE_CODE_COMMON_ERROR = 90000


	/***************901系统故障类**************/
	// 系统故障
	RESPONSE_CODE_SYSTEM = 90100
	// 服务实例不存在
	RESPONSE_CODE_MS_UNEXIST = 90101
	// Redis故障
	RESPONSE_CODE_REDIS = 90102


	/***************902接口鉴权类**************/
	// AppID无效
	RESPONSE_CODE_APPID_INVALID = 90200
	// AccessToken无效
	RESPONSE_CODE_ACCESSTOKEN_INVALID = 90201
	// Sign无效
	RESPONSE_CODE_SIGN_INVALID = 90202
	// TimeStamp无效
	RESPONSE_CODE_TIMESTAMP_INVALID = 90203
	// Action无效
	RESPONSE_CODE_ACTION_INVALID = 90204
	// SvrName无效
	RESPONSE_CODE_SVRNAME_INVALID = 90205

	/***************903接口参数解析类**************/
	// 参数解析错误
	RESPONSE_CODE_PARAMPARSE_FAIL = 90306

	/***************904接口参数解析类**************/
	// 店铺授权失败
	RESPONSE_CODE_REAUTH_SHOP = 90406
	// 快递单号不存在
	RESPONSE_CODE_TRACKNUM_UNEXIST = 90407
	// shopee获取首公里追踪号失败
	RESPONSE_CODE_FIRST_MILE_SHIP_ORDER = 90408
	// shopee获取首公里追踪号失败
	RESPONSE_CODE_FIRST_MILE_BIND = 90409
	// 集包不存在
	RESPONSE_CODE_CONNECTION_UNEXIST = 90410
	// 订单不存在
	RESPONSE_CODE_ORDER_UNEXIST = 90411
	// 余额不足
	RESPONSE_CODE_BALANCE_ALARM = 90412
	// 订单未打包
	RESPONSE_CODE_ORDER_UNPICKUP = 90413
)

//MQ
type MQ_ERR_TYPE int
const (
	MQ_ERR_TYPE_OK MQ_ERR_TYPE = 0 		//成功
	MQ_ERR_TYPE_RECOVERABLE = -1 		//可恢复的错误
	MQ_ERR_TYPE_UNRECOVERABLE = -2 		//不可恢复的错误
)
