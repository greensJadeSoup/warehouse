package cbd

import (
	"mime/multipart"
	"warehouse/v5-go-api-cangboss/model"
)

// ------------------------ req ------------------------
type GetSingleOrderReqCBD struct {
	VendorID  uint64 `json:"vendor_id" form:"vendor_id"`
	SellerID  uint64 `json:"seller_id" form:"seller_id"`
	OrderID   uint64 `json:"order_id,string" form:"order_id"`
	OrderTime int64  `json:"order_time" form:"order_time"`
}

type GetOrderBySNReqCBD struct {
	VendorID uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	SN       string `json:"sn"  form:"sn" binding:"required,lte=32"`
}

type UploadOrderDocumentReqCBD struct {
	SellerID  uint64 `json:"seller_id" form:"seller_id"`
	OrderID   uint64 `json:"order_id,string" form:"order_id"`
	OrderTime int64  `json:"order_time" form:"order_time"`
	Pdf       *multipart.FileHeader
	TmpPath   string
	Url       string
}

type OrderRecvAddr struct {
	Name        string `json:"name" binding:""`
	Phone       string `json:"phone" binding:""`
	City        string `json:"city" binding:""`
	State       string `json:"state" binding:""`   //区
	Zipcode     string `json:"zipcode" binding:""` //区
	FullAddress string `json:"full_address" binding:""`
}

type AddManualOrderReqCBD struct {
	SellerID           uint64  `json:"seller_id" binding:"required,gte=1"`
	SN                 string  `json:"sn"  xorm:"sn" binding:"required,lte=32"`
	Region             string  `json:"region" binding:"omitempty,lte=8"`
	ShippingCarrier    string  `json:"shipping_carrier" binding:"required,lte=16"`
	TotalAmount        float64 `json:"total_amount" binding:"omitempty,gte=0"`
	CashOnDelivery     uint8   `json:"cash_on_delivery" binding:"eq=0|eq=1"`
	NoteBuyer          string  `json:"note_buyer" binding:"omitempty,lte=255"`
	IsCb               *uint8  `json:"is_cb" binding:"required,eq=0|eq=1"`
	PlatformCreateTime int64

	RecvAddr OrderRecvAddr `json:"recv_addr" binding:""`

	ItemDetail []struct {
		ModelID         uint64 `json:"model_id,string" binding:"required"`
		Count           int    `json:"count" binding:"required"`
		ItemName        string `json:"item_name"`
		ItemSku         string `json:"item_sku"`
		ItemID          uint64 `json:"item_id,string"`
		PlatformItemID  string `json:"platform_item_id"`
		ModelName       string `json:"model_name"`
		ModelSku        string `json:"model_sku"`
		PlatformModelID string `json:"platform_model_id"`
		Image           string `json:"image"`
		Remark          string `json:"remark"`
	} `json:"item_detail" binding:"required,dive,required"`
}

type AddOrderReqCBD struct {
	ID                 uint64  `json:"id"  xorm:"id"  binding:"required,gte=1"`
	SellerID           uint64  `json:"seller_id"  xorm:"seller_id"  binding:"required,gte=1"`
	Platform           string  `json:"platform"  xorm:"platform"  binding:"required,gte=1"`
	ShopID             uint64  `json:"shop_id"  xorm:"shop_id"  binding:"required,gte=1"`
	SN                 string  `json:"sn"  xorm:"sn"  binding:"required,lte=32"`
	PickNum            string  `json:"pick_num"  xorm:"pick_num"`
	Status             string  `json:"status"  xorm:"status"  binding:"required,lte=16"`
	ItemDetail         string  `json:"item_detail"  xorm:"item_detail"  binding:"required,lte=255"`
	Region             string  `json:"region"  xorm:"region"  binding:"required,lte=8"`
	ShippingCarrier    string  `json:"shipping_carrier"  xorm:"shipping_carrier"  binding:"required,lte=16"`
	TotalAmount        float64 `json:"total_amount"  xorm:"total_amount"  binding:""`
	PayTime            int64   `json:"pay_time"  xorm:"pay_time"  binding:""`
	ReportTime         int64   `json:"report_time"  xorm:"report_time"`
	PaymentMethod      string  `json:"payment_method"  xorm:"payment_method"  binding:""`
	CashOnDelivery     uint8   `json:"cash_on_delivery"  xorm:"cash_on_delivery"  binding:""`
	RecvAddr           string  `json:"recv_addr"  xorm:"recv_addr"  binding:""`
	BuyerUserID        uint64  `json:"buyer_user_id"  xorm:"buyer_user_id"  binding:""`
	BuyerUsername      string  `json:"buyer_username"  xorm:"buyer_username"  binding:""`
	PlatformCreateTime int64   `json:"platform_create_time"  xorm:"platform_create_time"  binding:""`
	PlatformUpdateTime int64   `json:"platform_update_time"  xorm:"platform_update_time"  binding:""`
	NoteBuyer          string  `json:"note_buyer"  xorm:"note_buyer"  binding:""`
	NoteSeller         string  `json:"note_seller"  xorm:"note_seller"  binding:""`
	PickupTime         int64   `json:"pickup_time"  xorm:"pickup_time"  binding:""`
	ReportVendorTo     uint64  `json:"report_vendor_to"  xorm:"report_vendor_to"`

	FeeStatus   string `json:"fee_status"  xorm:"fee_status"`
	PriceDetail string `json:"price_detail"  xorm:"price_detail"`
	IsCb        uint8  `json:"is_cb"  xorm:"is_cb"`
}

type EditOrderReportReqCBD struct {
	OrderID   uint64 `json:"order_id"  xorm:"order_id"  binding:"required,gte=1"`
	OrderTime int64  `json:"order_time"  binding:"omitempty,gte=1"`

	NoteSeller     string `json:"note_seller"  xorm:"note_seller"`
	ItemDetail     string `json:"item_detail"  xorm:"item_detail"`
	RecvAddr       string `json:"recv_addr"  xorm:"recv_addr"  binding:""`
	ReportTime     int64  `json:"report_time" xorm:"report_time"`
	ReportVendorTo uint64 `json:"report_vendor_to"  xorm:"report_vendor_to"`
	PickupTime     int64  `json:"pickup_time" xorm:"pickup_time"`
	OnlyStock      bool   `json:"only_stock"  xorm:"only_stock"`

	Status      string  `json:"status"  xorm:"status"`
	Price       float64 `json:"price"  xorm:"price"`
	PriceReal   float64 `json:"price_real"  xorm:"price_real"`
	PriceDetail string  `json:"price_detail"  xorm:"price_detail"`
}

type EditOrderReqCBD struct {
	VendorID    uint64  `json:"vendor_id" binding:"omitempty,gte=1"`
	SellerID    uint64  `json:"seller_id" binding:"omitempty,gte=1"`
	OrderID     uint64  `json:"order_id,string" binding:"required,gte=1"`
	Weight      float64 `json:"weight" binding:"omitempty,gte=0"`
	WarehouseID uint64  `json:"warehouse_id" binding:"omitempty,gte=1"`
	LineID      uint64  `json:"line_id" binding:"omitempty,gte=1"`
	SendWayID   uint64  `json:"sendway_id" binding:"omitempty,gte=1"`
	Status      string  `json:"status" binding:"omitempty,lte=16"`

	ShippingCarrier string        `json:"shipping_carrier" binding:"omitempty"`
	RecvAddr        OrderRecvAddr `json:"recv_addr" binding:""`
}

type EditManualOrderReqCBD struct {
	VendorID  uint64 `json:"vendor_id" binding:"omitempty,gte=1"`
	SellerID  uint64 `json:"seller_id" binding:"omitempty,gte=1"`
	OrderID   uint64 `json:"order_id,string" binding:"required,gte=1"`
	OrderTime int64  `json:"order_time"  binding:"required,gte=1"`

	ShippingCarrier string        `json:"shipping_carrier" binding:"omitempty"`
	RecvAddr        OrderRecvAddr `json:"recv_addr"`
	Region          string        `json:"region" binding:"omitempty,lte=8"`
	TotalAmount     float64       `json:"total_amount" binding:"omitempty,gte=0"`
	CashOnDelivery  uint8         `json:"cash_on_delivery" binding:"eq=0|eq=1"`
	IsCb            *uint8        `json:"is_cb" binding:"required,eq=0|eq=1"`

	MdOrder *model.OrderMD
}

type EditOrderPriceRealReqCBD struct {
	VendorID  uint64  `json:"vendor_id" binding:"required,gte=1"`
	OrderID   uint64  `json:"order_id,string"  binding:"required,gte=1"`
	OrderTime int64   `json:"order_time"  binding:"required,gte=1"`
	PriceReal float64 `json:"price_real"  binding:"omitempty,gte=0"`
}

type OrderDeductReqCBD struct {
	VendorID  uint64 `json:"vendor_id" binding:"required,gte=1"`
	OrderID   uint64 `json:"order_id,string"  binding:"required,gte=1"`
	OrderTime int64  `json:"order_time"  binding:"required,gte=1"`

	MdOrder  *model.OrderMD
	MdSeller *model.SellerMD
}

type OrderRefundReqCBD struct {
	VendorID    uint64  `json:"vendor_id" binding:"required,gte=1"`
	OrderID     uint64  `json:"order_id,string"  binding:"required,gte=1"`
	OrderTime   int64   `json:"order_time"  binding:"required,gte=1"`
	PriceRefund float64 `json:"price_refund"  binding:"required,gte=0"`
}

type EditNoteManagerReqCBD struct {
	VendorID  uint64 `json:"vendor_id" binding:"required,gte=1"`
	OrderID   uint64 `json:"order_id,string"  binding:"required,gte=1"`
	OrderTime int64  `json:"order_time"  binding:"required,gte=1"`
	Note      string `json:"note"  binding:"omitempty,lte=255"`
}

type EditManagerImagesReqCBD struct {
	VendorID  uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	OrderID   uint64 `json:"order_id,string" form:"order_id" binding:"required,gte=1"`
	OrderTime int64  `json:"order_time" form:"order_time" binding:"required,gte=1"`
	OriImages string `json:"ori_images" form:"ori_images" binding:"omitempty"`
	ImageList string `json:"image_list" form:"image_list"`
	Detail    []OrderImageDetailCBD
}

type OrderImageDetailCBD struct {
	Name    string `json:"name" form:"name"`
	Url     string `json:"url" form:"url"`
	Image   *multipart.FileHeader
	TmpPath string
}

type EditNoteSellerReqCBD struct {
	SellerID  uint64 `json:"seller_id" binding:"required,gte=1"`
	OrderID   uint64 `json:"order_id,string"  binding:"required,gte=1"`
	OrderTime int64  `json:"order_time"  binding:"required,gte=1"`
	Note      string `json:"note"  binding:"omitempty,lte=255"`
}

type OrderDeliveryReqCBD struct {
	VendorID          uint64 `json:"vendor_id" binding:"required,gte=1"`
	OrderID           uint64 `json:"order_id,string"  binding:"required,gte=1"`
	OrderTime         int64  `json:"order_time"  binding:"required,gte=1"`
	DeliveryNum       string `json:"delivery_num"  binding:"omitempty,lte=32"`
	DeliveryLogistics string `json:"delivery_logistics"  binding:"omitempty,lte=32"`

	MdOrder       *model.OrderMD
	MdOrderSimple *model.OrderSimpleMD
	PackList      *[]model.PackMD
	//TiList			[]cp_api.GetTrackInfoItemResp //改单使用的，A改B，A订单要去判断物流是否到达商超
}

type GetPriceDetailReqCBD struct {
	VendorID  uint64 `json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID  uint64 `json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	OrderID   uint64 `json:"order_id" form:"order_id" binding:"required,gte=1"`
	OrderTime int64  `json:"order_time" form:"order_time" binding:"required,gte=1"`
}

type ListOrderReqCBD struct {
	VendorID     uint64 `json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID     uint64 `json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	SellerIDList []string

	WareHouseRole string

	Platform        string `json:"platform" form:"platform" binding:"omitempty,lte=32"`
	PlatformExclude string `json:"platform_exclude" form:"platform_exclude" binding:"omitempty,lte=32"`
	OrderType       string `json:"order_type" form:"order_type" binding:"omitempty"`
	OrderTypeList   []string
	FeeStatus       string `json:"fee_status" form:"fee_status" binding:"omitempty"`
	FeeStatusList   []string

	ShippingCarry     string `json:"shipping_carry" form:"shipping_carry" binding:"omitempty"`
	ShippingCarryList []string

	DeliveryLogistics     string `json:"delivery_logistics" form:"delivery_logistics" binding:"omitempty"`
	DeliveryLogisticsList []string

	OrderStatus          string `json:"order_status" form:"order_status" binding:"omitempty"`
	OrderStatusList      []string
	OrderStatusNotInList []string

	PlatformStatus     string `json:"platform_status" form:"platform_status" binding:"omitempty"`
	PlatformStatusList []string

	NoDisPlatformStatus     string `json:"nodis_platform_status" form:"nodis_platform_status" binding:"omitempty"`
	NoDisPlatformStatusList []string

	SellerKey        string `json:"seller_key" form:"seller_key" binding:"omitempty,lte=32"`
	ShopKey          string `json:"shop_key" form:"shop_key" binding:"omitempty,lte=64"`
	ItemKey          string `json:"item_key" form:"item_key" binding:"omitempty,lte=64"`
	SkuKey           string `json:"sku_key" form:"sku_key" binding:"omitempty,lte=64"`
	SearchKey1       string `json:"search_key_1" form:"search_key_1" binding:"omitempty,lte=512"`
	SearchKey1List   []string
	JHDList          []string
	SearchKey2       string `json:"search_key_2" form:"search_key_2" binding:"omitempty,lte=64"`
	ProblemPack      bool   `json:"problem_pack" form:"problem_pack" binding:"omitempty"`
	ConnectionFilter string `json:"connection_filter" form:"connection_filter" binding:"omitempty,eq=yes|eq=no"` //yes:已加入集包 no:未加入集包
	SkuType          string `json:"sku_type" form:"sku_type" binding:"omitempty,eq=mix|eq=stock|eq=express|eq=express_return"`

	WarehouseID     uint64 `json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
	WarehouseIDList []string
	LineIDList      []string

	RackID      string `json:"rack_id" form:"rack_id" binding:"omitempty,gte=1"`
	RackIDList  []string
	StockID     uint64 `json:"stock_id" form:"stock_id" binding:"omitempty,gte=1"`
	StockIDList []string

	CancelDays int    `json:"cancel_days" form:"cancel_days" binding:"omitempty,gte=1"`
	SkuCount   int    `json:"sku_count" form:"sku_count" binding:"omitempty,gte=1"`
	IsCb       *uint8 `json:"is_cb" form:"is_cb" binding:"omitempty"`

	From int64 `json:"from" form:"from" binding:"required,gte=1"`
	To   int64 `json:"to" form:"to" binding:"required,gte=1"`

	IsPaging    bool `json:"is_paging" form:"is_paging"`
	PageIndex   int  `json:"page_index" form:"page_index" binding:""`
	PageSize    int  `json:"page_size" form:"page_size" binding:""`
	ExcelOutput bool
}

type OrderTrendReqCBD struct {
	VendorID uint64 `json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID uint64 `json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`

	From            int64
	To              int64
	WarehouseIDList []string
	LineIDList      []string
}

type DelOrderReqCBD struct {
	SellerID  uint64 `json:"seller_id" form:"seller_id"`
	OrderID   uint64 `json:"order_id" xorm:"order_id"`
	OrderTime int64  `json:"order_time" xorm:"order_time"`
}

type OrderRecvAddrCBD struct {
	Name        string `json:"name"  xorm:"name"`
	FullAddress string `json:"full_address"  xorm:"full_address"`
	Phone       string `json:"phone"  xorm:"phone"`
}

type OrderBaseInfoCBD struct {
	SellerID  uint64 `json:"seller_id" xorm:"seller_id"`
	OrderID   uint64 `json:"order_id" xorm:"order_id"`
	OrderTime int64  `json:"order_time" xorm:"order_time"`
	SN        string `json:"sn" xorm:"sn"`
}

type OrderAppTimeInfoCBD struct {
	OrderID   uint64  `json:"order_id" xorm:"order_id"`
	Date      string  `json:"date" xorm:"date"`
	PriceReal float64 `json:"price_real" xorm:"price_real"`
}

type WeightPriceDetail struct {
	SendWayID   uint64 `json:"sendway_id"  xorm:"sendway_id"`
	SendWayName string `json:"sendway_name"  xorm:"sendway_name"`

	Price          float64 `json:"price"  xorm:"price"`
	Weight         float64 `json:"weight"  xorm:"weight"`
	PriEach        float64 `json:"pri_each"  xorm:"pri_each"`
	PriOrder       float64 `json:"pri_order"  xorm:"pri_order"`
	PriFirstWeight float64 `json:"pri_first_weight"  xorm:"pri_first_weight"`
}

type SkuPriceDetail struct {
	WarehouseID   uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"  xorm:"warehouse_name"`

	Start       int     `json:"start"  xorm:"start"`
	Price       float64 `json:"price"  xorm:"price"`
	SkuType     string  `json:"sku_type"  xorm:"sku_type"`           // 快递express or 囤货stock
	SkuUnitType string  `json:"sku_unit_type"  xorm:"sku_unit_type"` // 个count or 项row
	SkuCount    int     `json:"sku_count"  xorm:"sku_count"`         //总个数
	SkuRow      int     `json:"sku_row"  xorm:"sku_row"`             //总行数
	ExceedCount int     `json:"exceed_count"  xorm:"exceed_count"`   //超过个数
	ExceedRow   int     `json:"exceed_row"  xorm:"exceed_row"`       //超过行数
	PriEach     float64 `json:"pri_each"  xorm:"pri_each"`
	PriOrder    float64 `json:"pri_order"  xorm:"pri_order"`
}

type ConsumablePriceDetail struct {
	ConsumableID   uint64  `json:"consumable_id" binding:"required,gte=1"`
	ConsumableName string  `json:"consumable_name" binding:"required,lte=255"`
	PriEach        float64 `json:"pri_each" binding:"omitempty,gte=0"`
	Count          int     `json:"count" binding:"omitempty,gte=0"`
	Price          float64 `json:"price" binding:"omitempty,gte=0"`
}

type ServicePriceDetail struct {
	PricePastePick   float64 `json:"pri_paste_pick" binding:"omitempty,gte=0"`
	PricePasteFace   float64 `json:"pri_paste_face" binding:"omitempty,gte=0"`
	PriceShopToShop  float64 `json:"pri_shop_to_shop" binding:"omitempty,gte=0"`
	PriceToShopProxy float64 `json:"pri_to_shop_proxy" binding:"omitempty,gte=0"`
	PriceDelivery    float64 `json:"pri_delivery" binding:"omitempty,gte=0"`
}

type OrderPriceDetailCBD struct {
	Price                 float64                 `json:"price"  xorm:"price"`
	PriceReal             float64                 `json:"price_real"  xorm:"price_real"`
	PriceRefund           float64                 `json:"price_refund"  xorm:"price_refund"`
	Balance               float64                 `json:"balance"  xorm:"balance"`
	SN                    string                  `json:"sn"  xorm:"sn"`
	WeightPriceDetail     WeightPriceDetail       `json:"weight_price_detail"  xorm:"weight_price_detail"`
	SkuPriceDetail        []SkuPriceDetail        `json:"sku_price_detail"  xorm:"sendway_price_detail"`
	ConsumablePriceDetail []ConsumablePriceDetail `json:"consumable_price_detail"  xorm:"consumable_price_detail"`
	PlatformPriceRules    []PlatformPriceRule     `json:"platform_price_rule" binding:"required,dive,required"`
	ServicePriceDetail    ServicePriceDetail      `json:"service_price_detail"  xorm:"service_price_detail"`
}

type ListOrderAttributeByYmReqCBD struct {
	OrderID   uint64
	OrderTime int64
	YearMonth string

	ConnectionID uint64
	Weight       float64
	FeeStatus    string
	PriceReal    float64
	SellerID     uint64
	RealName     string
}

type OrderPackUpDetailCBD struct {
	OrderID   uint64  `json:"order_id,string" binding:"required,gte=1"`
	OrderTime int64   `json:"order_time" binding:"required,gte=1"`
	Weight    float64 `json:"weight" binding:"omitempty,gte=0"`
	Length    float64 `json:"length" binding:"omitempty,gte=0"`
	Width     float64 `json:"width" binding:"omitempty,gte=0"`
	Height    float64 `json:"height" binding:"omitempty,gte=0"`
	//ConnectionID		uint64		`json:"connection_id" binding:"omitempty,gte=1"`
	CustomsNum     string                  `json:"customs_num" binding:"omitempty,lte=32"`
	PackSubDetail  []PackSubCBD            `json:"pack_sub_detail" binding:"required"`
	ConsumableList []ConsumablePriceDetail `json:"consumable_list" binding:"omitempty"`
	Deduct         bool                    `json:"deduct" binding:"omitempty"`
}

type BatchOrderReqCBD struct {
	VendorID uint64 `json:"vendor_id" binding:"omitempty,gte=1"`
	SellerID uint64 `json:"seller_id" binding:"omitempty,gte=1"`

	EditStatusDetail []struct {
		OrderID   uint64 `json:"order_id,string" binding:"required,gte=1"`
		OrderTime int64  `json:"order_time" binding:"required,gte=1"`
		Status    string `json:"status" binding:"omitempty,lte=32"`
	} `json:"edit_status_detail" binding:"omitempty,dive,required"`

	DeductDetail []struct {
		OrderID   uint64 `json:"order_id,string" binding:"required,gte=1"`
		OrderTime int64  `json:"order_time" binding:"required,gte=1"`
	} `json:"deduct_detail" binding:"omitempty,dive,required"`

	DeliveryDetail []struct {
		OrderID   uint64 `json:"order_id,string" binding:"required,gte=1"`
		OrderTime int64  `json:"order_time" binding:"required,gte=1"`
	} `json:"delivery_detail" binding:"omitempty,dive,required"`

	PackUpDetail []OrderPackUpDetailCBD `json:"pack_up_detail" binding:"omitempty,dive,required"`
}

type ChangeOrderReqCBD struct {
	VendorID  uint64 `json:"vendor_id" binding:"omitempty,gte=1"`
	SellerID  uint64 `json:"seller_id" binding:"omitempty,gte=1"`
	OrderID   uint64 `json:"order_id,string" binding:"required,gte=1"`
	OrderTime int64  `json:"order_time" binding:"required,gte=1"`
	NewSn     string `json:"new_sn" binding:"omitempty,lte=64"`

	MdOrderFrom *model.OrderMD
	MdOrderTo   *model.OrderMD
	MdOsFrom    *model.OrderSimpleMD
	MdOsTo      *model.OrderSimpleMD
}

type ReturnOrderReqCBD struct {
	SellerID  uint64 `json:"seller_id" binding:"required,gte=1"`
	OrderID   uint64 `json:"order_id,string" binding:"required,gte=1"`
	OrderTime int64  `json:"order_time" binding:"required,gte=1"`
}

type TmpOrder struct {
	RackID      uint64 `json:"-"  xorm:"rack_id"`
	SellerID    uint64 `json:"seller_id"  xorm:"seller_id"`
	RealName    string `json:"real_name"  xorm:"real_name"`
	OrderID     uint64 `json:"order_id,string" xorm:"order_id"`
	OrderTime   int64  `json:"order_time" xorm:"order_time"`
	SN          string `json:"sn"  xorm:"sn"`
	ManagerNote string `json:"manager_note"  xorm:"manager_note"`
}

type DownOrderReqCBD struct {
	VendorID uint64 `json:"vendor_id" binding:"required,gte=1"`
	OrderID  uint64 `json:"order_id,string" binding:"required,gte=1"`
}

type SkuDetail struct {
	ExpressSkuCount       int //快递类sku个数数目
	ExpressSkuRow         int //快递类sku项数目
	StockSkuCount         int //库存类sku个数数目
	StockSkuRow           int //库存类sku项数目
	ExpressReturnSkuCount int //目的仓退回来的快递
	ExpressReturnSkuRow   int //目的仓退回来的快递
}

// --------------------resp-------------------------------
type ListOrderRespCBD struct {
	ID                 uint64  `json:"id,string"  xorm:"id pk autoincr"`
	SellerID           uint64  `json:"seller_id"  xorm:"seller_id"`
	Platform           string  `json:"platform"  xorm:"platform"`
	ShopID             uint64  `json:"shop_id"  xorm:"shop_id"`
	IsCB               uint8   `json:"is_cb"  xorm:"is_cb"`
	IsSIP              uint8   `json:"is_sip"  xorm:"is_sip"`
	PlatformShopID     string  `json:"platform_shop_id"  xorm:"platform_shop_id"`
	SN                 string  `json:"sn"  xorm:"sn"`
	PickNum            string  `json:"pick_num"  xorm:"pick_num"`
	Status             string  `json:"status"  xorm:"status"`
	PlatformStatus     string  `json:"platform_status"  xorm:"platform_status"`
	ItemDetail         string  `json:"item_detail,omitempty"  xorm:"item_detail"`
	Region             string  `json:"region,omitempty"  xorm:"region"`
	MidNum             string  `json:"mid_num"  xorm:"mid_num"`
	CustomsNum         string  `json:"customs_num"  xorm:"customs_num"`
	DeliveryNum        string  `json:"delivery_num"  xorm:"delivery_num"`
	DeliveryLogistics  string  `json:"delivery_logistics"  xorm:"delivery_logistics"`
	PlatformTrackNum   string  `json:"platform_track_num"  xorm:"platform_track_num"`
	ShippingCarrier    string  `json:"shipping_carrier,omitempty"  xorm:"shipping_carrier"`
	TotalAmount        float64 `json:"total_amount"  xorm:"total_amount"`
	PaymentMethod      string  `json:"payment_method,omitempty"  xorm:"payment_method"`
	Currency           string  `json:"currency,omitempty"  xorm:"currency"`
	CashOnDelivery     uint8   `json:"cash_on_delivery"  xorm:"cash_on_delivery"`
	RecvAddr           string  `json:"recv_addr,omitempty"  xorm:"recv_addr"`
	BuyerUserID        uint64  `json:"buyer_user_id"  xorm:"buyer_user_id"`
	BuyerUsername      string  `json:"buyer_username"  xorm:"buyer_username"`
	PlatformCreateTime int64   `json:"platform_create_time"  xorm:"platform_create_time"`
	PlatformUpdateTime int64   `json:"platform_update_time,omitempty"  xorm:"platform_update_time"`
	NoteBuyer          string  `json:"note_buyer"  xorm:"note_buyer"`
	NoteSeller         string  `json:"note_seller"  xorm:"note_seller"`
	NoteManager        string  `json:"note_manager"  xorm:"note_manager"`
	NoteManagerID      uint64  `json:"note_manager_id"  xorm:"note_manager_id"`
	NoteManagerTime    int64   `json:"note_manager_time"  xorm:"note_manager_time"`
	NoteManagerName    string  `json:"note_manager_name"  xorm:"note_manager_name"`
	PayTime            int64   `json:"pay_time"  xorm:"pay_time"`
	PickupTime         int64   `json:"pickup_time"  xorm:"pickup_time"`
	ShipDeadlineTime   int64   `json:"ship_deadline_time"  xorm:"ship_deadline_time"`
	ReportTime         int64   `json:"report_time"  xorm:"report_time"`
	DeductTime         int64   `json:"deduct_time"  xorm:"deduct_time"`
	DeliveryTime       int64   `json:"delivery_time"  xorm:"delivery_time"`
	ToReturnTime       int64   `json:"to_return_time"  xorm:"to_return_time"`
	ChangeFrom         string  `json:"change_from,omitempty"  xorm:"change_from"`
	ChangeTo           string  `json:"change_to,omitempty"  xorm:"change_to"`
	ChangeTime         int64   `json:"change_time,omitempty"  xorm:"change_time"`
	PackageList        string  `json:"package_list,omitempty"  xorm:"package_list"`
	CancelBy           string  `json:"cancel_by,omitempty"  xorm:"cancel_by"`
	CancelReason       string  `json:"cancel_reason,omitempty"  xorm:"cancel_reason"`
	FeeStatus          string  `json:"fee_status"  xorm:"fee_status"`
	Price              float64 `json:"price"  xorm:"price"`
	PriceReal          float64 `json:"price_real"  xorm:"price_real"`
	PriceDetail        string  `json:"price_detail"  xorm:"price_detail"`
	ManagerImagesStr   string  `json:"-" xorm:"manager_images"`
	SkuType            string  `json:"sku_type"  xorm:"sku_type"`

	OnlyStock uint8   `json:"only_stock"  xorm:"only_stock"`
	Weight    float64 `json:"weight" xorm:"weight"`
	Volume    float64 `json:"-" xorm:"-"`
	Length    float64 `json:"-" xorm:"-"`
	Width     float64 `json:"-" xorm:"-"`
	Height    float64 `json:"-" xorm:"-"`

	WarehouseID   uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"  xorm:"warehouse_name"`
	LineID        uint64 `json:"line_id"  xorm:"line_id"`
	SourceID      uint64 `json:"source" xorm:"source_id"`
	SourceName    string `json:"source_name" xorm:"source_name"`
	ToID          uint64 `json:"to" xorm:"to_id"`
	ToName        string `json:"to_name" xorm:"to_name"`
	SendWayID     uint64 `json:"sendway_id"  xorm:"sendway_id"`
	SendWayType   string `json:"sendway_type"  xorm:"sendway_type"`
	SendWayName   string `json:"sendway_name"  xorm:"sendway_name"`

	ReadyPack       int               `json:"ready_pack"`
	TotalPack       int               `json:"total_pack"`
	Problem         uint8             `json:"problem"`
	ProblemTrackNum []TrackNumInfoCBD `json:"problem_track_num"`
	AllTrackNum     []TrackNumInfoCBD `json:"all_track_num"`
	ManagerImages   []ManagerImageCBD `json:"manager_images"`
	PackSubDetail   []PackSubCBD      `json:"pack_sub_detail"`

	ShopName string `json:"shop_name"  xorm:"shop_name"`
	RealName string `json:"real_name"  xorm:"real_name"`

	TmpRackCBD
}

type ManagerImageCBD struct {
	Url      string `json:"url"`
	Time     string `json:"time"`
	Type     string `json:"type"`
	RealName string `json:"real_name"`
}

type TrackNumInfoCBD struct {
	TrackNum    string   `json:"track_num"`
	Status      string   `json:"status"`
	Problem     uint8    `json:"problem"`
	Reason      string   `json:"reason"`
	ManagerNote string   `json:"manager_note"`
	DependID    []string `json:"depend_id"`
}

type GetPriceDetailRespCBD struct {
	OrderID     uint64  `json:"order_id,string"  binding:"required,gte=1"`
	SellerID    uint64  `json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	RealName    string  `json:"real_name"  xorm:"real_name"`
	SN          string  `json:"sn"  xorm:"sn"`
	FeeStatus   string  `json:"fee_status"  xorm:"fee_status"`
	Balance     float64 `json:"balance"  xorm:"balance"`
	Price       float64 `json:"price"  xorm:"price"`
	PriceReal   float64 `json:"price_real"  xorm:"price_real"`
	PriceRefund float64 `json:"price_refund"  xorm:"price_refund"`
	PriceDetail string  `json:"price_detail"  xorm:"price_detail"`
}

type OrderAddress struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Town        string `json:"town"`
	District    string `json:"district"`
	City        string `json:"city"`
	State       string `json:"state"`
	Region      string `json:"region"`
	Zipcode     string `json:"zipcode"`
	FullAddress string `json:"full_address"`
}

type ListOrderStatusCountRespCBD struct {
	StatusCountList   []ListOrderStatusCountCBD `json:"status_count_list"`
	ShippingCarry     []string                  `json:"shipping_carry"`
	PlatformStatus    []string                  `json:"platform_status"`
	DeliveryLogistics []string                  `json:"delivery_logistics"`
}

type ListOrderStatusCountCBD struct {
	Status string `json:"status"  xorm:"status"`
	Count  int    `json:"count" xorm:"count"`
}

type ListOrderAttributeCBD struct {
	OrderID   uint64  `json:"id"  xorm:"id"`
	SellerID  uint64  `json:"seller_id"  xorm:"seller_id"`
	RealName  string  `json:"real_name"  xorm:"real_name"`
	Status    string  `json:"status"  xorm:"status"`
	Weight    float64 `json:"weight"  xorm:"weight"`
	FeeStatus string  `json:"fee_status"  xorm:"fee_status"`
	PriceReal float64 `json:"price_real"  xorm:"price_real"`
}

type OrderAddManualRespCBD struct {
	SellerID   uint64 `json:"seller_id"`
	OrderID    uint64 `json:"order_id,string"`
	OrderTime  int64  `json:"order_time"`
	ItemDetail string `json:"item_detail"`
	SN         string `json:"sn"`
	Platform   string `json:"platform"`
	RealName   string `json:"real_name"`
	NoteBuyer  string `json:"note_buyer"`
}

type TrendDateCount struct {
	Count int    `json:"count"`
	Date  string `json:"date"`
}

type TrendDateAmount struct {
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
}

type TrendBaseAmountCBD struct {
	Today         float64           `json:"today"`
	LastSevenDay  float64           `json:"last_seven_day"`
	LastThirtyDay float64           `json:"last_thirty_day"`
	Detail        []TrendDateAmount `json:"date_count"`
}

type TrendBaseCountCBD struct {
	Today         int              `json:"today"`
	LastSevenDay  int              `json:"last_seven_day"`
	LastThirtyDay int              `json:"last_thirty_day"`
	Detail        []TrendDateCount `json:"date_count"`
}

type OrderTrendRespCBD struct {
	Balance       float64            `json:"balance"`
	ReportTrend   TrendBaseCountCBD  `json:"report_trend"`
	DeliveryTrend TrendBaseCountCBD  `json:"delivery_trend"`
	DeductTrend   TrendBaseAmountCBD `json:"deduct_trend"`
	ConsumeTrend  TrendBaseAmountCBD `json:"consume_trend"`
}

type OrderPackUpConfirmRespCBD struct {
	OrderID        uint64   `json:"order_id,string"`
	SN             string   `json:"sn"`
	Platform       string   `json:"platform"`
	Problem        uint8    `json:"problem"`
	TrackNum       []string `json:"track_num"`
	Status         string   `json:"status"`
	PlatformStatus string   `json:"platform_status"`
}

type BatchOrderRespCBD struct {
	OrderID uint64 `json:"order_id,string"`
	SN      string `json:"sn"`
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}
