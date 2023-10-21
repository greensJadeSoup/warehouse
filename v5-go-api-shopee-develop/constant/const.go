package constant

const (
	USER_TYPE_SUPER_MANAGER = "super_manager"
	USER_TYPE_MANAGER       = "manager"
	USER_TYPE_SELLER        = "seller"
	USER_TYPE_SERVICE       = "service"

	LOGIN_TYPE_ACCOUNT = "account"
	LOGIN_TYPE_EMAIL   = "email"
	LOGIN_TYPE_PHONE   = "phone"
	LOGIN_TYPE_THIRD   = "third"
	LOGIN_TYPE_FACE    = "face"

	SHOPEE_URI_AUTH_PARTNER                       = "/api/v2/shop/auth_partner"
	SHOPEE_URI_GET_ACCESSTOKEN                    = "/api/v2/auth/token/get"
	SHOPEE_URI_REFRESH_ACCESSTOKEN                = "/api/v2/auth/access_token/get"
	SHOPEE_URI_GET_SHOP_INFO                      = "/api/v2/shop/get_shop_info"
	SHOPEE_URI_GET_SHOP_PROFILE                   = "/api/v2/shop/get_profile"
	SHOPEE_URI_GET_ITEM_LIST                      = "/api/v2/product/get_item_list"
	SHOPEE_URI_GET_ITEM_BASE_INFO                 = "/api/v2/product/get_item_base_info"
	SHOPEE_URI_GET_MODEL_LIST                     = "/api/v2/product/get_model_list"
	SHOPEE_URI_GET_ORDER_LIST                     = "/api/v2/order/get_order_list"
	SHOPEE_URI_GET_ORDER_DETAIL                   = "/api/v2/order/get_order_detail"
	SHOPEE_URI_GET_SHIPPING_PARAM                 = "/api/v2/logistics/get_shipping_parameter"
	SHOPEE_URI_GET_TRACKING_NUM                   = "/api/v2/logistics/get_tracking_number"
	SHOPEE_URI_GET_TRACKING_INFO                  = "/api/v2/logistics/get_tracking_info"
	SHOPEE_URI_GET_ADDRESS_LIST                   = "/api/v2/logistics/get_address_list"
	SHOPEE_URI_SHIP_ORDER                         = "/api/v2/logistics/ship_order"
	SHOPEE_URI_CREATE_SHIPPING_DOCUMENT           = "/api/v2/logistics/create_shipping_document"
	SHOPEE_URI_GET_RESULT_SHIPPING_DOCUMENT       = "/api/v2/logistics/get_shipping_document_result"
	SHOPEE_URI_DOWNLOAD_SHIPPING_DOCUMENT         = "/api/v2/logistics/download_shipping_document"
	SHOPEE_URI_GET_SHIPPING_DOCUMENT_INFO         = "/api/v2/logistics/get_shipping_document_data_info"
	SHOPEE_URI_GET_CHANNEL_LIST                   = "/api/v2/first_mile/get_channel_list"
	SHOPEE_URI_GENERATE_FIRST_MILE_TRACKING_NUM   = "/api/v2/first_mile/generate_first_mile_tracking_number"
	SHOPEE_URI_BIND_FIRST_MILE_TRACKING_NUM       = "/api/v2/first_mile/bind_first_mile_tracking_number"
	SHOPEE_URI_GET_FIRST_MILE_TRACKING_NUM_DETAIL = "/api/v2/first_mile/get_detail"
	SHOPEE_URI_GET_RETURN_DETAIL                  = "/api/v2/returns/get_return_detail"
	SHOPEE_URI_GET_RETURN_LIST                    = "/api/v2/returns/get_return_list"

	OSS_REGION_SZ            = "oss-cn-shenzhen"
	BUCKET_NAME_PUBLICE_PDF  = "publice-pdf"
	OSS_PATH_SHOPEE_DOCUMENT = "shopee-document"
)

// SHOPEE PUSH CODE
const (
	SHOPEE_PUSH_CODE_SHOP_AUTH           = 1 //店铺授权
	SHOPEE_PUSH_CODE_SHOP_UNAUTH         = 2 //店铺解绑
	SHOPEE_PUSH_CODE_ORDER_STATUS_UPDATE = 3 //订单状态更新
	SHOPEE_PUSH_CODE_ORDER_TRACKNUM_PUSH = 4 //订单trackNum更新
)

// order type
const (
	ORDER_TYPE_SHOPEE   = "shopee"   //shopee平台的订单
	ORDER_TYPE_STOCK_UP = "stock_up" //囤货
	ORDER_TYPE_MANUAL   = "manual"   //手动录单
)

// Shipping Carrier type
const (
	SHIPPING_CARRIER_7_11                  = "7-ELEVEN"
	SHIPPING_CARRIER_LAIERFU               = "萊爾富"
	SHIPPING_CARRIER_SHOPEE_SHOP_TO_SHOP   = "蝦皮店到店"
	SHIPPING_CARRIER_FULL_HOUSE            = "全家"
	SHIPPING_CARRIER_OK_MART               = "OK Mart"
	SHIPPING_CARRIER_BLACK_CAT             = "黑貓宅急便"
	SHIPPING_CARRIER_SHOPEE_DELIVERY       = "蝦皮宅配"
	SHIPPING_CARRIER_SELLER_DELIVERY       = "賣家宅配"
	SHIPPING_CARRIER_SELLER_DELIVERY_BIG   = "賣家宅配：大型/超重物品運送"
	SHIPPING_CARRIER_OFFLINE_SHOP_TO_SHOP  = "线下店到店"
	SHIPPING_CARRIER_OFFLINE_DELIVERY      = "线下宅配"
	SHIPPING_CARRIER_HOUSE_COMMON_DELIVERY = "宅配通"
	//蝦皮海外 - 宅配（海運）
	//蝦皮海外 - 蝦皮店到店
	//蝦皮海外 - 7-11
	//蝦皮海外 - 萊爾富（空運）
	//蝦皮海外 - 全家
	//蝦皮海外 - 萊爾富（海運）
	//萊爾富-經濟包
	//蝦皮海外 - 宅配（空運）
)

// common redis key
const (
	REDIS_KEY_AUTH_SHOP = "[shopee]auth_shop:"
)

// time const
const (
	REDIS_EXPIRE_TIME_AUTH_SHOP = 10
)

// platform const
const (
	PLATFORM_SHOPEE = "shopee"
)

// shopee order status
const (
	SHOPEE_ORDER_STATUS_UNPAID             = "UNPAID"             //未付款
	SHOPEE_ORDER_STATUS_READY_TO_SHIP      = "READY_TO_SHIP"      //可以发货
	SHOPEE_ORDER_STATUS_PROCESSED          = "PROCESSED"          //卖家已安排发货
	SHOPEE_ORDER_STATUS_RETRY_SHIP         = "RETRY_SHIP"         //需要重新发货
	SHOPEE_ORDER_STATUS_SHIPPED            = "SHIPPED"            //面单已打印
	SHOPEE_ORDER_STATUS_TO_CONFIRM_RECEIVE = "TO_CONFIRM_RECEIVE" //卖家已收货
	SHOPEE_ORDER_STATUS_IN_CANCEL          = "IN_CANCEL"          //用户申请取消订单
	SHOPEE_ORDER_STATUS_CANCELLED          = "CANCELLED"          //订单已取消
	SHOPEE_ORDER_STATUS_TO_RETURN          = "TO_RETURN"          //退货中
	SHOPEE_ORDER_STATUS_COMPLETED          = "COMPLETED"          //已完成
)

// order status
const (
	ORDER_STATUS_UNPAID     = "unpaid"     //未付款
	ORDER_STATUS_PAID       = "paid"       //已付款
	ORDER_STATUS_PRE_REPORT = "pre_report" //已预报
	ORDER_STATUS_READY      = "ready"      //已到齐
	ORDER_STATUS_PACKAGED   = "packaged"   //已打包
	ORDER_STATUS_STOCK_OUT  = "stock_out"  //已出库
	ORDER_STATUS_CUSTOMS    = "customs"    //清关中
	ORDER_STATUS_ARRIVE     = "arrive"     //已达目的仓库
	ORDER_STATUS_DELIVERY   = "delivery"   //已派送
	ORDER_STATUS_TO_CHANGE  = "to_change"  //改单中
	ORDER_STATUS_CHANGED    = "changed"    //已改单
	ORDER_STATUS_TO_RETURN  = "to_return"  //退货中
	ORDER_STATUS_RETURNED   = "returned"   //已退货
	ORDER_STATUS_OTHER      = "other"      //其他
)
