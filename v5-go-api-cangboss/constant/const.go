package constant

const (
	USER_TYPE_SUPER_MANAGER = "super_manager"
	USER_TYPE_MANAGER = "manager"
	USER_TYPE_SELLER = "seller"
	USER_TYPE_SERVICE = "service" //客服

	LOGIN_TYPE_ACCOUNT = "account"
	LOGIN_TYPE_EMAIL = "email"
	LOGIN_TYPE_PHONE = "phone"
	LOGIN_TYPE_THIRD = "third"
	LOGIN_TYPE_FACE = "face"

	OSS_REGION_SZ = "oss-cn-shenzhen"
	BUCKET_NAME_PUBLICE_IMAGE = "publice-images"
	BUCKET_NAME_PUBLICE_PDF = "publice-pdf"
	OSS_PATH_ITEM_PICTURE = "item_images"
	OSS_PATH_ORDER_PICTURE = "order_images"
	OSS_PATH_SHOPEE_DOCUMENT = "shopee-document"
)


const (
	COMMON_PASSWORD_MD5 = "dead98e7833aa3604dba74fefe107d59" //前端原密码5bd492b39b3db1630a21942910f7e2bf
)

//calculate type
const (
	CALCULATE_TYPE_ADD = "add"
	CALCULATE_TYPE_SUB = "sub"
)

//photo type
const (
	PHOTO_SUFFIX_JPG = ".jpg"
	PHOTO_SUFFIX_JPEG = ".jpeg"
	PHOTO_SUFFIX_PNG = ".png"

	PHOTO_CONTENT_TYPE_JPG = "image/jpg"
	PHOTO_CONTENT_TYPE_JPEG = "image/jpeg"
	PHOTO_CONTENT_TYPE_PNG = "image/png"
)

//time const
const (
	REDIS_EXPIRE_TIME_AUTH_SHOP = 10
)

//order status
const (
	ORDER_STATUS_UNPAID = "unpaid" 		//未付款
	ORDER_STATUS_PAID = "paid" 		//已付款
	ORDER_STATUS_PRE_REPORT = "pre_report" 	//已预报
	ORDER_STATUS_READY = "ready" 		//已到齐
	ORDER_STATUS_PACKAGED = "packaged" 	//已打包
	ORDER_STATUS_STOCK_OUT = "stock_out" 	//已出库
	ORDER_STATUS_CUSTOMS = "customs" 	//清关中（改成了已通关）
	ORDER_STATUS_ARRIVE = "arrive" 		//已达目的仓库
	ORDER_STATUS_DELIVERY = "delivery" 	//已派送
	ORDER_STATUS_TO_CHANGE = "to_change" 	//改单中
	ORDER_STATUS_CHANGED = "changed" 	//已改单
	ORDER_STATUS_TO_RETURN = "to_return" 	//转囤货（之前是退货中）
	ORDER_STATUS_RETURNED = "returned" 	//已上架（之前是已退货）
	ORDER_STATUS_OTHER = "other" 		//其他
)

//shopee order status
const (
	SHOPEE_ORDER_STATUS_UNPAID = "UNPAID" 				//未付款
	SHOPEE_ORDER_STATUS_READY_TO_SHIP = "READY_TO_SHIP" 		//可以发货
	SHOPEE_ORDER_STATUS_PROCESSED = "PROCESSED" 			//卖家已安排发货
	SHOPEE_ORDER_STATUS_RETRY_SHIP = "RETRY_SHIP" 			//需要重新发货
	SHOPEE_ORDER_STATUS_SHIPPED = "SHIPPED" 			//面单已打印
	SHOPEE_ORDER_STATUS_TO_CONFIRM_RECEIVE = "TO_CONFIRM_RECEIVE" 	//卖家已收货
	SHOPEE_ORDER_STATUS_IN_CANCEL = "IN_CANCEL" 			//用户申请取消订单
	SHOPEE_ORDER_STATUS_CANCELLED = "CANCELLED" 			//订单已取消
	SHOPEE_ORDER_STATUS_TO_RETURN = "TO_RETURN" 			//退货中
	SHOPEE_ORDER_STATUS_COMPLETED = "COMPLETED" 			//已完成
)

//platform const
const (
	PLATFORM_SHOPEE = "shopee"
	PLATFORM_STOCK_UP = "stock_up" 		//囤货
	PLATFORM_MANUAL = "manual"		//手动录单
)

//Num type
const (
	NUM_TYPE_ORDER = "order" 		//订单
	NUM_TYPE_EXPRESS = "express" 		//快递
	NUM_TYPE_PLATFORM_TRACKNUM = "platform_tracknum" 	//平台物流追踪号
)

//order type
const (
	ORDER_TYPE_SHOPEE = "shopee" 		//shopee平台的订单
	ORDER_TYPE_STOCK_UP = "stock_up" 	//囤货
	ORDER_TYPE_MANUAL = "manual" 		//手动录单
)

//sub pack type
const (
	PACK_SUB_TYPE_EXPRESS = "express" 	//快递
	PACK_SUB_TYPE_STOCK = "stock" 		//库存
	PACK_SUB_TYPE_STOCK_UP = "stock_up" 	//囤货
)

//report type
const (
	REPORT_TYPE_ORDER = "order" 		//快递
	REPORT_TYPE_STOCK_UP = "stock_up" 	//囤货
)

//warehouse type
const (
	WAREHOUSE_ROLE_SOURCE = "source" 	//始发仓库
	WAREHOUSE_ROLE_TO = "to" 		//目的仓库
)

//pack problem reason
const (
	PACK_PROBLEM_DESTROY = "destroy" 			//包裹破损
	PACK_PROBLEM_LOSE = "lose" 				//无人认领
	PACK_PROBLEM_LOSE_DESTROY = "lose_destroy" 		//无人认领+破损
	PACK_PROBLEM_NO_REPORT = "no_report" 			//无预报（可以入库）
	PACK_PROBLEM_NO_REPORT_DESTROY = "no_report_destroy" 	//无预报+破损
)

//fee status
const (
	FEE_STATUS_UNHANDLE = "un_handle" 		//未扣款
	FEE_STATUS_SUCCESS = "success" 			//扣款成功
	FEE_STATUS_FAIL = "fail" 			//扣款失败
	FEE_STATUS_RETURN = "return" 			//已退款
)

//fee fail reason
const (
	FEE_FAIL_REASON_BALANCE = "balance" 		//余额不足
)

//pack reserved
const (
	PACK_TRACK_NUM_RESERVED = "reserved" 			//快递保留
)

//pack status
const (
	PACK_STATUS_INIT = "init" 			//未到达
	PACK_STATUS_ENTER_SOURCE = "enter_source" 	//已达始发仓
	PACK_STATUS_ENTER_TO = "enter_to" 		//已达目的仓
	PACK_STATUS_RETURN_SOURCE = "return_source" 	//已退货到始发仓
	PACK_STATUS_RETURN_TO = "return_to" 		//已退货到目的仓
)

//sku type
const (
	SKU_TYPE_EXPRESS = "express" 		//快递
	SKU_TYPE_STOCK = "stock" 		//库存
	SKU_TYPE_MIX = "mix" 			//混合
	SKU_TYPE_EXPRESS_RETURN = "express_return" 	//退到目的仓的快递
)

//sku unit type
const (
	SKU_UNIT_TYPE_COUNT = "count" 		//快递
	SKU_UNIT_TYPE_ROW = "row" 		//库存
)

//connection status
const (
	CONNECTION_STATUS_INIT = "init" 		//未出库
	CONNECTION_STATUS_STOCK_OUT = "stock_out" 	//已出库
	CONNECTION_STATUS_CUSTOMS = "customs" 		//清关中
	CONNECTION_STATUS_ARRIVE = "arrive" 		//已达目的仓库
)

//rack action type
const (
	RACK_ACTION_ADD = "add"
	RACK_ACTION_SUB = "sub"
)

//mid connection type
const (
	MID_CONNECTION_NORMAL = "normal"
	MID_CONNECTION_SPECIAL = "special"
	MID_CONNECTION_SPECIAL_A = "special_a"
	MID_CONNECTION_SPECIAL_B = "special_b"
)

//event type
const (
	EVENT_TYPE_PICK_UP = "pick_up"  				//打包
	EVENT_TYPE_ENTER_SOURCE = "enter_source"  			//入库起始仓
	EVENT_TYPE_ENTER_TO = "enter_to"  				//入库目的仓
	EVENT_TYPE_RETURN_SOURCE = "return_source"  			//退货入库起始仓
	EVENT_TYPE_RETURN_TO = "return_to"  				//退货入库目的仓
	EVENT_TYPE_EDIT_ORDER = "edit_order"				//编辑订单
	EVENT_TYPE_EDIT_WEIGHT = "edit_weight"				//修改订单重量
	EVENT_TYPE_DELIVER = "deliver"					//派送(库存消耗)
	EVENT_TYPE_ORDER_DEDUCT = "order_deduct"			//订单单独扣款
	EVENT_TYPE_CONNECTION_ORDER_DEDUCT = "conn_order_deduct"	//集包订单扣款
	EVENT_TYPE_ORDER_REFUND = "order_refund"			//订单退款
	EVENT_TYPE_CHARGE = "charge"					//充值(用户列表中调整余额)
	EVENT_TYPE_DEDUCT = "deduct"					//扣款(用户列表中调整余额)
	EVENT_TYPE_EDIT_PRICE_REAL = "edit_price_real"			//更改实收价格
	EVENT_TYPE_EDIT_STOCK_COUNT = "edit_stock_count" 		//编辑库存数目
	EVENT_TYPE_EDIT_STOCK_RACK = "edit_stock_rack" 	 		//调货架
	EVENT_TYPE_ADD_STOCK_RACK = "add_stock_rack" 	 		//创建库存货架
	EVENT_TYPE_EDIT_DOWN_RACK = "down_rack" 	 		//下架
	EVENT_TYPE_BIND_STOCK = "bind_stock" 		 		//绑定库存
	EVENT_TYPE_UNBIND_STOCK = "unbind_stock" 	 		//解绑库存
	EVENT_TYPE_DEL_STOCK = "del_stock" 	 			//销毁库存
	EVENT_TYPE_ORDER_TAKE_BACK = "order_take_back" 	 		//销毁库存
	EVENT_TYPE_CHANGE_ORDER = "change_order" 	 		//改单
	EVENT_TYPE_CANCEL_CHANGE_ORDER = "cancel_change_order" 	 	//撤销改单
	EVENT_TYPE_RETURN_ORDER = "return_order" 	 		//排号入库
	EVENT_TYPE_CANCEL_RETURN_ORDER = "cancel_return_order" 	 	//撤销排号入库
	EVENT_TYPE_EDIT_MANUAL_ORDER = "edit_manual_order"		//编辑自定义订单信息
)

//object type
const (
	OBJECT_TYPE_PACK = "pack"  		//包裹
	OBJECT_TYPE_PICK_NUM = "pick_num"	//拣货单
	OBJECT_TYPE_ORDER = "order"		//订单
	OBJECT_TYPE_CONNECTION = "connection"	//集包
	OBJECT_TYPE_STOCK = "stock"		//库存
	OBJECT_TYPE_RACK = "rack"		//货架
	OBJECT_TYPE_SELLER = "seller"		//卖家
)

//pack way
const (
	PACK_WAY_DIRECTLY = "directly"  	//直接贴单 一个快递对应一个订单
	PACK_WAY_MERGE = "merge"  		//包裹合并 多个快递对应一个订单
	PACK_WAY_SPLIT = "split"  		//包裹拆分 一个快递对应多个订单
	PACK_WAY_STRUCT = "struct"  		//组合订单 多个快递对应一个订单（需要拆分及合并或与台湾仓储合并的复杂件）
)

//pack way
const (
	ORDER_DOWN_RACK_TYPE_PEOPLE = "people"  	//人为下架
	ORDER_DOWN_RACK_TYPE_DELIVERY = "delivery"  	//派送下架
)

//工单状态
const (
	APPLY_STATUS_OPEN = "open"
	APPLY_STATUS_HANDLED = "handled"
	APPLY_STATUS_CLOSE = "close"
)

//excel类型
const (
	EXCEL_OUTPUT_TYPE_CUSTOMS_COMPANY = "customs_company" //清关公司
	EXCEL_OUTPUT_TYPE_AIR_COMPANY = "air_company" //航空公司
)

//Shipping Carrier type
const (
	SHIPPING_CARRIER_7_11 = "7-ELEVEN"
	SHIPPING_CARRIER_LAIERFU = "萊爾富"
	SHIPPING_CARRIER_SHOPEE_SHOP_TO_SHOP = "蝦皮店到店"
	SHIPPING_CARRIER_FULL_HOUSE = "全家"
	SHIPPING_CARRIER_OK_MART = "OK Mart"
	SHIPPING_CARRIER_BLACK_CAT = "黑貓宅急便"
	SHIPPING_CARRIER_SHOPEE_DELIVERY = "蝦皮宅配"
	SHIPPING_CARRIER_SELLER_DELIVERY = "賣家宅配"
	SHIPPING_CARRIER_SELLER_DELIVERY_BIG = "賣家宅配：大型/超重物品運送"
	SHIPPING_CARRIER_OFFLINE_SHOP_TO_SHOP = "线下店到店"
	SHIPPING_CARRIER_OFFLINE_DELIVERY = "线下宅配"
	//蝦皮海外 - 宅配（海運）
	//蝦皮海外 - 蝦皮店到店
	//蝦皮海外 - 7-11
	//蝦皮海外 - 萊爾富（空運）
	//蝦皮海外 - 全家
	//蝦皮海外 - 萊爾富（海運）
	//萊爾富-經濟包
	//蝦皮海外 - 宅配（空運）
)
