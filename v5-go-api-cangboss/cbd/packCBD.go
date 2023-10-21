package cbd

import (
	"strconv"
	"time"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_obj"
)

func NewOrder(t int64) *model.OrderMD {
	ym := strconv.Itoa(time.Unix(t, 0).Year()) + "_" + strconv.Itoa(int(time.Unix(t, 0).Month()))
	return &model.OrderMD{Yearmonth: ym}
}

// ------------------------ req ------------------------
type PackModelDetailCBD struct {
	Platform        string `json:"platform"`
	ShopID          uint64 `json:"shop_id"`
	PlatformShopID  string `json:"platform_shop_id"`
	Region          string `json:"region"`
	ShopName        string `json:"shop_name"`
	ItemID          uint64 `json:"item_id,string"`
	PlatformItemID  string `json:"platform_item_id"`
	ItemName        string `json:"item_name"`
	ItemSKU         string `json:"item_sku"`
	ModelID         uint64 `json:"model_id,string" binding:"required,gte=1"`
	PlatformModelID string `json:"platform_model_id"`
	ModelSku        string `json:"model_sku"`
	Remark          string `json:"remark"`
	Image           string `json:"image"`
	Count           int    `json:"count" binding:"required,gte=1"`
	StoreCount      int    `json:"store_count" binding:"omitempty,gte=0"`
	EnterCount      int    `json:"enter_count" binding:"omitempty,gte=0"`
	DependID        string `json:"depend_id"`
	Note            string `json:"note" binding:"omitempty,lte=2000"`
}

type PackDetailCBD struct {
	ID              uint64 `json:"id,string" binding:"omitempty,gte=1"`
	Type            string `json:"type" binding:"required,lte=16"`
	TrackNum        string `json:"track_num" binding:"omitempty,lte=32"`
	StockID         uint64 `json:"stock_id,string"`
	ExpressCodeType int    `json:"express_code_type"`

	PackID         uint64
	Status         string
	SourceRecvTime int64
	ToRecvTime     int64
	DeliveryCount  int
	DeliveryTime   int64
	ReturnCount    int
	ReturnTime     int64

	PackModelDetailCBD
}

type AddReportReqCBD struct {
	VendorID       uint64                  `json:"vendor_id" binding:"omitempty,gte=1"`
	SellerID       uint64                  `json:"seller_id"  binding:"required,gte=1"`
	WarehouseID    uint64                  `json:"warehouse_id"  binding:"omitempty,gte=1"`
	LineID         uint64                  `json:"line_id"  binding:"omitempty,gte=1"`
	SendWayID      uint64                  `json:"sendway_id" binding:"omitempty,gte=1"`
	ReportType     string                  `json:"report_type"  binding:"required,eq=order|eq=stock_up"`
	OrderID        uint64                  `json:"order_id,string"  binding:"omitempty,gte=1"`
	OrderTime      int64                   `json:"order_time"  binding:"omitempty,gte=1"`
	Note           string                  `json:"note"  binding:"omitempty,lte=512"`
	ShipOrder      bool                    `json:"ship_order" binding:"omitempty"`
	ConsumableList []ConsumablePriceDetail `json:"consumable_list" binding:"omitempty"`
	Detail         []PackDetailCBD         `json:"detail" binding:"required,dive,required"`

	WarehouseName   string
	WarehouseRole   string
	MdOrder         *model.OrderMD
	MdOrderSimple   *model.OrderSimpleMD
	MdSourceWh      model.WarehouseMD
	MdToWh          model.WarehouseMD
	MdSw            model.SendWayMD
	StockUpAddrInfo OrderRecvAddrCBD

	OnlyStock bool //新子项是否纯库存
	SkuDetail SkuDetail
}

type BatchAddReportReqCBD struct {
	VendorID    uint64 `json:"vendor_id" binding:"omitempty,gte=1"`
	SellerID    uint64 `json:"seller_id"  binding:"required,gte=1"`
	WarehouseID uint64 `json:"warehouse_id"  binding:"omitempty,gte=1"`
	LineID      uint64 `json:"line_id"  binding:"omitempty,gte=1"`
	SendWayID   uint64 `json:"sendway_id" binding:"omitempty,gte=1"`
	ReportList  []struct {
		ReportType     string                  `json:"report_type"  binding:"required,eq=order|eq=stock_up"`
		OrderID        uint64                  `json:"order_id,string"  binding:"omitempty,gte=1"`
		OrderTime      int64                   `json:"order_time"  binding:"omitempty,gte=1"`
		Note           string                  `json:"note"  binding:"omitempty,lte=512"`
		ShipOrder      bool                    `json:"ship_order" binding:"omitempty"`
		Detail         []PackDetailCBD         `json:"detail" binding:"required,dive,required"`
		ConsumableList []ConsumablePriceDetail `json:"consumable_list" binding:"omitempty"`
	} `json:"report_list" binding:"required,dive,required"`
}

type EditReportReqCBD struct {
	VendorID       uint64 `json:"vendor_id" binding:"omitempty,gte=1"`
	SellerID       uint64 `json:"seller_id"  binding:"omitempty,gte=1"`
	OrderID        uint64 `json:"order_id,string" binding:"required,gte=1"`
	OrderTime      int64  `json:"order_time"  binding:"omitempty,gte=1"`
	WarehouseID    uint64 `json:"warehouse_id"  binding:"required,gte=1"`
	WarehouseName  string
	WarehouseRole  string
	LineID         uint64                  `json:"line_id"  binding:"omitempty,gte=1"`
	SendWayID      uint64                  `json:"sendway_id"  binding:"omitempty,gte=1"`
	Note           string                  `json:"note"  binding:"omitempty,lte=512"`
	ConsumableList []ConsumablePriceDetail `json:"consumable_list" binding:"omitempty"`

	MdOrder         *model.OrderMD
	MdOrderSimple   *model.OrderSimpleMD
	MdSourceWh      model.WarehouseMD
	MdToWh          model.WarehouseMD
	MdSw            model.SendWayMD
	ReportType      string
	StockUpAddrInfo OrderRecvAddrCBD

	OnlyStock bool //新子项是否纯库存
	SkuDetail SkuDetail

	Detail []PackDetailCBD `json:"detail" binding:"required,dive,required"`
}

type DelReportReqCBD struct {
	MdOrder *model.OrderMD
}

type DelPackReqCBD struct {
	SellerID uint64 `json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
}

type GetReportReqCBD struct {
	VendorID    uint64 `json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID    uint64 `json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	OrderID     uint64 `json:"order_id" form:"order_id" binding:"required,gte=1"`
	OrderTime   int64  `json:"order_time" form:"order_time" binding:"required,gte=1"`
	WarehouseID uint64 `json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
}

type GetTrackInfoReqCBD struct {
	SellerID uint64 `json:"seller_id" form:"seller_id" binding:"required,gte=1"`
	TrackNum string `json:"track_num" form:"track_num" binding:"required,lte=32"`
}

type GetPackDetailReqCBD struct {
	VendorID uint64 `json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID uint64 `json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	PackID   uint64 `json:"pack_id" form:"pack_id" binding:"required,gte=1"`
}

type EnterPackDetailReqCBD struct {
	VendorID  uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	SearchKey string `json:"search_key" form:"search_key" binding:"required,lte=32"`
	IsReturn  bool   `json:"is_return" form:"is_return"`
}

type CheckNumReqCBD struct {
	VendorID  uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	SearchKey string `json:"search_key" form:"search_key" binding:"required,lte=32"`
}

type GetReadyOrderReqCBD struct {
	VendorID  uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	SearchKey string `json:"search_key" form:"search_key" binding:"required,lte=32"`
}

type BatchPrintOrderReqCBD struct {
	VendorID    uint64 `json:"vendor_id" binding:"omitempty,gte=1"`
	WarehouseID uint64 `json:"warehouse_id" binding:"omitempty,gte=1"`
	SellerID    uint64 `json:"seller_id" binding:"omitempty,gte=1"`

	Detail []struct {
		OrderID   uint64 `json:"order_id,string"`
		OrderTime int64  `json:"order_time"`
	} `json:"detail"`
}

type ListPackManagerReqCBD struct {
	VendorID        uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	Problem         bool   `json:"problem" form:"problem"`
	SellerKey       string `json:"seller_key" form:"seller_key" binding:"omitempty,lte=32"`
	SN              string `json:"sn" form:"sn" binding:"omitempty,lte=32"`
	TrackNum        string `json:"track_num" form:"track_num" binding:"omitempty,lte=32"`
	WarehouseID     uint64 `json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
	OnlyCount       bool   `json:"only_count" form:"only_count" binding:"omitempty"`
	WarehouseIDList []string
	LineIDList      []string
	Type            string `json:"type" form:"type" binding:"omitempty,eq=stock_up|eq=order"`
	Status          string `json:"status" form:"status" binding:"omitempty,eq=init|eq=enter_source|eq=enter_to"`
	Reason          string `json:"reason" form:"reason" binding:"omitempty,eq=destroy|eq=lose|eq=lose_destroy|eq=no_report|eq=no_report_destroy"`
	Source          int64  `json:"source" form:"source" binding:"omitempty,gte=1"`
	To              int64  `json:"to" form:"to" binding:"omitempty,gte=1"`

	IsPaging  bool `json:"is_paging" form:"is_paging"`
	PageIndex int  `json:"page_index" form:"page_index" binding:"required"`
	PageSize  int  `json:"page_size" form:"page_size" binding:"required"`

	ExcelOutput bool
}

type ListPackSellerReqCBD struct {
	SellerID     uint64 `json:"seller_id" form:"seller_id" binding:"required,gte=1"`
	Problem      bool   `json:"problem" form:"problem"`
	SellerKey    string `json:"seller_key" form:"seller_key" binding:"omitempty,lte=32"`
	SN           string `json:"sn" form:"sn" binding:"omitempty,lte=32"`
	TrackNum     string `json:"track_num" form:"track_num" binding:"omitempty,lte=32"`
	WarehouseID  uint64 `json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
	Type         string `json:"type" form:"type" binding:"omitempty,eq=stock_up|eq=order"`
	Status       string `json:"status" form:"status" binding:"omitempty,eq=init|eq=enter_source|eq=enter_to"`
	Reason       string `json:"reason" form:"reason" binding:"omitempty,eq=destroy|eq=lose|eq=lose_destroy|eq=no_report|eq=no_report_destroy"`
	OnlyCount    bool   `json:"only_count" form:"only_count" binding:"omitempty"`
	Source       int64  `json:"source" form:"source" binding:"omitempty,gte=1"`
	To           int64  `json:"to" form:"to" binding:"omitempty,gte=1"`
	VendorIDList []string

	IsPaging  bool `json:"is_paging" form:"is_paging"`
	PageIndex int  `json:"page_index" form:"page_index" binding:"required"`
	PageSize  int  `json:"page_size" form:"page_size" binding:"required"`
}

type EditPackWeightReqCBD struct {
	VendorID uint64  `json:"vendor_id" binding:"required,gte=1"`
	PackID   uint64  `json:"pack_id,string" binding:"required,gte=1"`
	Weight   float64 `json:"weight"  binding:"omitempty,gte=0"`
}

type OrderWeightCBD struct {
	OrderID   uint64  `json:"order_id,string" binding:"required,gte=0"`
	OrderTime int64   `json:"order_time" binding:"required,gte=0"`
	Weight    float64 `json:"weight"  binding:"omitempty,gte=0"`
	Length    float64 `json:"length"  binding:"omitempty,gte=0"`
	Width     float64 `json:"width"  binding:"omitempty,gte=0"`
	Height    float64 `json:"height"  binding:"omitempty,gte=0"`
}

type EditPackOrderWeightReqCBD struct {
	VendorID uint64           `json:"vendor_id" binding:"required,gte=1"`
	PackID   uint64           `json:"pack_id,string" binding:"required,gte=1"`
	Detail   []OrderWeightCBD `json:"detail" binding:"required,dive,required"`
}

type EditPackTrackNumReqCBD struct {
	VendorID uint64 `json:"vendor_id" binding:"required,gte=1"`
	PackID   uint64 `json:"pack_id,string" binding:"required,gte=1"`
	TrackNum string `json:"track_num"  binding:"required,lte=255"`
}

type EditPackManagerNoteReqCBD struct {
	VendorID    uint64 `json:"vendor_id" binding:"required,gte=1"`
	PackID      uint64 `json:"pack_id,string" binding:"required,gte=1"`
	ManagerNote string `json:"manager_note"  binding:"required,lte=255"`
}

type EnterReqCBD struct {
	VendorID    uint64  `json:"vendor_id" binding:"required,gte=1"`
	WarehouseID uint64  `json:"warehouse_id" binding:"required,gte=1"`
	SearchKey   string  `json:"search_key" binding:"required,lte=32"`
	Weight      float64 `json:"weight" binding:"omitempty,gte=0"`
	RackID      uint64  `json:"rack_id" binding:"omitempty,gte=1"`
	IsReturn    bool    `json:"is_return"`

	ShopID        uint64
	SellerID      uint64
	WarehouseRole string
	WarehouseName string
	OrderStatus   string

	Detail []struct {
		ID          uint64 `json:"id,string"`
		ModelID     uint64 `json:"model_id,string"`
		StockID     uint64 `json:"stock_id,string"`
		CheckCount  int    `json:"check_count"`
		Type        string
		DeliverTime int64
		RackDetail  []RackDetailCBD `json:"rack_detail" binding:"required,dive,required"`
	} `json:"detail" binding:"required,dive,required"`
}

type ProblemPackManagerReqCBD struct {
	VendorID    uint64 `json:"vendor_id" binding:"required,gte=1"`
	WarehouseID uint64 `json:"warehouse_id" binding:"required,gte=1"`

	TrackNum string  `json:"track_num" binding:"required,lte=32"`
	SellerID uint64  `json:"seller_id" binding:"omitempty,gte=1"`
	Reason   string  `json:"reason" binding:"required,eq=lose|eq=lose_destroy|eq=no_report|eq=no_report_destroy|eq=destroy"`
	Weight   float64 `json:"weight" binding:"omitempty,gte=0"`
	RackID   uint64  `json:"rack_id" binding:"omitempty,gte=1"`

	ManagerNote   string `json:"manager_note"  binding:"required,lte=255"`
	WarehouseName string
	WarehouseRole string
}

type EditTmpRackReqCBD struct {
	VendorID  uint64 `json:"vendor_id" binding:"required,gte=1"`
	TmpType   string `json:"tmp_type" binding:"required,eq=pack|eq=order"`
	ObjectID  uint64 `json:"object_id,string" binding:"omitempty"`
	NewRackID uint64 `json:"new_rack_id" binding:"required,gte=1"`
}

type DownPackReqCBD struct {
	VendorID     uint64   `json:"vendor_id" binding:"required,gte=1"`
	TrackNumList []string `json:"track_num" binding:"required"`
}

type CheckDownPackReqCBD struct {
	VendorID     uint64   `json:"vendor_id" binding:"required,gte=1"`
	TrackNumList []string `json:"track_num" binding:"required"`
}

type LogisticsInfoCBD struct {
	OrderID  uint64 `json:"order_id" xorm:"order_id"`
	SellerID uint64 `json:"seller_id" xorm:"seller_id"`

	WarehouseID    uint64 `json:"warehouse_id" xorm:"warehouse_id"`
	WarehouseName  string `json:"warehouse_name" xorm:"warehouse_name"`
	LineID         uint64 `json:"line_id" xorm:"line_id"`
	SourceID       uint64 `json:"source_id" xorm:"source_id"`
	SourceName     string `json:"source_name" xorm:"source_name"`
	SourceAddress  string `json:"source_address" xorm:"source_address"`
	SourceReceiver string `json:"source_receiver" xorm:"source_receiver"`
	SourcePhone    string `json:"source_phone" xorm:"source_phone"`
	ToID           uint64 `json:"to_id" xorm:"to_id"`
	ToName         string `json:"to_name" xorm:"to_name"`
	ToAddress      string `json:"to_address" xorm:"to_address"`
	ToReceiver     string `json:"to_receiver" xorm:"to_receiver"`
	ToPhone        string `json:"to_phone" xorm:"to_phone"`
	ToNote         string `json:"to_note" xorm:"to_note"`
	SendWayID      uint64 `json:"sendway_id"  xorm:"sendway_id"`
	SendWayType    string `json:"sendway_type"  xorm:"sendway_type"`
	SendWayName    string `json:"sendway_name" xorm:"sendway_name"`
	TmpRackCBD     `xorm:"extends"`
}

type TmpRackCBD struct {
	RackID            uint64 `json:"rack_id"  xorm:"rack_id"`
	RackWarehouseID   uint64 `json:"rack_warehouse_id"  xorm:"rack_warehouse_id"`
	RackWarehouseRole string `json:"rack_warehouse_role"  xorm:"rack_warehouse_role"`
	RackNum           string `json:"rack_num"  xorm:"rack_num"`
	AreaNum           string `json:"area_num"  xorm:"area_num"`
}

type FreezeStockCBD struct {
	StockID uint64 `json:"stock_id" xorm:"stock_id"`
	ModelID uint64 `json:"model_id" xorm:"model_id"`
	Count   int    `json:"count" xorm:"count"`
}

// --------------------resp-------------------------------
type PackSubCBD struct {
	ID          uint64 `json:"id,string"  xorm:"id"`
	PackID      uint64 `json:"pack_id,string" xorm:"pack_id"`
	TrackNum    string `json:"track_num"  xorm:"track_num"`
	Problem     uint8  `json:"problem"  xorm:"problem"`
	Reason      string `json:"reason"  xorm:"reason"`
	ManagerNote string `json:"manager_note"  xorm:"manager_note"`

	SellerID uint64 `json:"seller_id"  xorm:"seller_id"`
	ShopID   uint64 `json:"shop_id"  xorm:"shop_id"`
	ShopName string `json:"shop_name"  xorm:"shop_name"`
	Platform string `json:"platform"  xorm:"platform"`
	SN       string `json:"sn"  xorm:"sn"`
	PickNum  string `json:"pick_num"  xorm:"pick_num"`

	OrderID   uint64 `json:"order_id,string"  xorm:"order_id"`
	OrderTime int64  `json:"order_time"  xorm:"order_time"`

	Type         string `json:"type"  xorm:"type"`
	Count        int    `json:"count"  xorm:"count"`
	StoreCount   int    `json:"store_count"  xorm:"store_count"`
	EnterCount   int    `json:"enter_count"  xorm:"enter_count"`
	CheckCount   int    `json:"check_count"  xorm:"check_count"`
	DeliverCount int    `json:"deliver_count"  xorm:"deliver_count"`
	ReturnCount  int    `json:"return_count"  xorm:"return_count"`

	Region         string `json:"region"  xorm:"region"`
	PlatformShopID string `json:"platform_shop_id"  xorm:"platform_shop_id"`

	ItemID         uint64 `json:"item_id,string"  xorm:"item_id"`
	PlatformItemID string `json:"platform_item_id"  xorm:"platform_item_id"`
	ItemName       string `json:"item_name"  xorm:"item_name"`
	ItemStatus     string `json:"item_status"  xorm:"item_status"`

	ModelID         uint64 `json:"model_id,string"  xorm:"model_id"`
	PlatformModelID string `json:"platform_model_id"  xorm:"platform_model_id"`
	ModelSku        string `json:"model_sku"  xorm:"model_sku"`
	Remark          string `json:"remark" xorm:"remark"`
	Images          string `json:"images"  xorm:"images"`
	ModelIsDelete   uint8  `json:"model_is_delete"  xorm:"model_is_delete"`
	DependID        string `json:"depend_id"  xorm:"depend_id"`
	Status          string `json:"status"  xorm:"status"`
	Note            string `json:"note"  xorm:"note"`
	ExpressCodeType int    `json:"express_code_type"  xorm:"express_code_type"`

	SourceRecvTime int64 `json:"source_recv_time"  xorm:"source_recv_time"`
	ToRecvTime     int64 `json:"to_recv_time"  xorm:"to_recv_time"`
	DeliverTime    int64 `json:"deliver_time"  xorm:"deliver_time"`
	ReturnTime     int64 `json:"return_time"  xorm:"return_time"`

	CreateTime cp_obj.Datetime `json:"create_time"  xorm:"create_time"`

	Total   int    `json:"total"  xorm:"total"`
	Freeze  int    `json:"freeze"  xorm:"freeze"`
	StockID uint64 `json:"stock_id,string" xorm:"stock_id"`
	AreaNum string `json:"area_num" xorm:"area_num"` //临时包裹的
	RackID  uint64 `json:"rack_id" xorm:"rack_id"`   //临时包裹的

	RackWarehouseID   uint64 `json:"rack_warehouse_id" xorm:"rack_warehouse_id"`     //临时包裹的
	RackWarehouseRole string `json:"rack_warehouse_role" xorm:"rack_warehouse_role"` //临时包裹的
	RackNum           string `json:"rack_num" xorm:"rack_num"`                       //临时包裹的

	HasGift    uint8 `json:"has_gift"  xorm:"has_gift"`
	AutoImport uint8 `json:"auto_import"  xorm:"auto_import"`

	RackDetail []RackDetailCBD `json:"rack_detail" xorm:"rack_detail"`
}

type CheckNumRespCBD struct {
	NumType string `json:"num_type"  xorm:"num_type"`
}

type GetReportRespCBD struct {
	SellerID uint64 `json:"seller_id"  xorm:"seller_id"`
	RealName string `json:"real_name"  xorm:"real_name"`

	ShopID         uint64 `json:"shop_id"  xorm:"shop_id"`
	IsCb           uint8  `json:"is_cb"  xorm:"is_cb"`
	ShopName       string `json:"shop_name"  xorm:"shop_name"`
	PlatformShopID string `json:"platform_shop_id" xorm:"platform_shop_id"`

	WarehouseID    uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName  string `json:"warehouse_name"  xorm:"warehouse_name"`
	LineID         uint64 `json:"line_id"  xorm:"line_id"`
	Source         uint64 `json:"source" xorm:"source"`
	SourceName     string `json:"source_name" xorm:"source_name"`
	SourceAddress  string `json:"source_address" xorm:"source_address"`
	SourceReceiver string `json:"source_receiver" xorm:"source_receiver"`
	SourcePhone    string `json:"source_phone" xorm:"source_phone"`
	To             uint64 `json:"to" xorm:"to"`
	ToName         string `json:"to_name" xorm:"to_name"`
	ToAddress      string `json:"to_address" xorm:"to_address"`
	ToReceiver     string `json:"to_receiver" xorm:"to_receiver"`
	ToPhone        string `json:"to_phone" xorm:"to_phone"`
	ToNote         string `json:"to_note" xorm:"to_note"`
	SendWayID      uint64 `json:"sendway_id"  xorm:"sendway_id"`
	SendWayType    string `json:"sendway_type"  xorm:"sendway_type"`
	SendWayName    string `json:"sendway_name" xorm:"sendway_name"`

	OrderID            uint64  `json:"order_id,string" xorm:"order_id"`
	Platform           string  `json:"platform" xorm:"platform"`
	PlatformCreateTime int64   `json:"platform_create_time"  xorm:"platform_create_time"`
	SN                 string  `json:"sn" xorm:"sn"`
	Status             string  `json:"status" xorm:"status"`
	TotalAmount        float64 `json:"total_amount"  xorm:"total_amount"`
	CashOnDelivery     uint8   `json:"cash_on_delivery"  xorm:"cash_on_delivery"`
	PlatformStatus     string  `json:"platform_status" xorm:"platform_status"`
	PickNum            string  `json:"pick_num" xorm:"pick_num"`
	DeliveryNum        string  `json:"delivery_num"  xorm:"delivery_num"`
	DeliveryLogistics  string  `json:"delivery_logistics"  xorm:"delivery_logistics"`
	ReportTime         int64   `json:"report_time" xorm:"report_time"`
	PickupTime         int64   `json:"pickup_time" xorm:"pickup_time"`
	ItemDetail         string  `json:"item_detail"  xorm:"item_detail"`
	Region             string  `json:"region"  xorm:"region"`
	ShippingCarrier    string  `json:"shipping_carrier" xorm:"shipping_carrier"`
	PackWay            string  `json:"pack_way"  xorm:"pack_way"`
	FeeStatus          string  `json:"fee_status" xorm:"fee_status"`
	PlatformTrackNum   string  `json:"platform_track_num" xorm:"platform_track_num"`
	Price              float64 `json:"price" xorm:"price"`
	Weight             float64 `json:"weight" xorm:"weight"`
	Volume             float64 `json:"volume" xorm:"volume"`
	NoteBuyer          string  `json:"note_buyer" xorm:"note_buyer"`
	NoteSeller         string  `json:"note_seller" xorm:"note_seller"`
	NoteManager        string  `json:"note_manager"  xorm:"note_manager"`
	NoteManagerID      uint64  `json:"note_manager_id"  xorm:"note_manager_id"`
	NoteManagerTime    int64   `json:"note_manager_time"  xorm:"note_manager_time"`
	NoteManagerName    string  `json:"note_manager_name"  xorm:"note_manager_name"`
	Consumable         string  `json:"consumable"  xorm:"consumable"`
	SkuType            string  `json:"sku_type"  xorm:"sku_type"`
	ChangeTo           string  `json:"change_to"  xorm:"change_to"`
	RecvAddr           string  `json:"recv_addr"  xorm:"recv_addr"`
	ChangeFrom         string  `json:"change_from"  xorm:"change_from"`
	ChangeTime         int64   `json:"change_time"  xorm:"change_time"`

	TimeNow int64 `json:"time_now" xorm:"time_now"`
	Ready   bool  `json:"ready" xorm:"ready"`
	//ReadyPack		int		`json:"ready_pack"`
	//TotalPack		int		`json:"total_pack"`

	PackSubList []PackSubCBD `json:"packsub_list" xorm:"packsub_list"`
	TmpRackCBD
}

type ListPackRespCBD struct {
	ID uint64 `json:"id,string"  xorm:"id pk autoincr"`

	SellerID uint64 `json:"seller_id"  xorm:"seller_id"`
	RealName string `json:"real_name"  xorm:"real_name"`

	TrackNum      string `json:"track_num"  xorm:"track_num"`
	WarehouseID   uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"  xorm:"warehouse_name"`
	LineID        uint64 `json:"line_id"  xorm:"line_id"`
	Source        uint64 `json:"source"  xorm:"source"`
	To            uint64 `json:"to"  xorm:"to"`
	SourceName    string `json:"source_name"  xorm:"source_name"`
	ToName        string `json:"to_name"  xorm:"to_name"`
	SendWayID     uint64 `json:"sendway_id"  xorm:"sendway_id"`
	SendWayType   string `json:"sendway_type"  xorm:"sendway_type"`
	SendWayName   string `json:"sendway_name" xorm:"sendway_name"`

	Status string  `json:"status"  xorm:"status"`
	Type   string  `json:"type"  xorm:"type"`
	Weight float64 `json:"weight"  xorm:"weight"`

	SourceRecvTime int64 `json:"source_recv_time"  xorm:"source_recv_time"`
	ToRecvTime     int64 `json:"to_recv_time"  xorm:"to_recv_time"`

	Problem     int    `json:"problem"  xorm:"problem"`
	Reason      string `json:"reason"  xorm:"reason"`
	ManagerNote string `json:"manager_note"  xorm:"manager_note"`
	RackID      uint64 `json:"rack_id"  xorm:"rack_id"`
	RackNum     string `json:"rack_num"  xorm:"rack_num"`
	AreaNum     string `json:"area_num"  xorm:"area_num"`

	Log []ListWarehouseLogRespCBD `json:"log"`
}

type TrackInfoRespCBD struct {
	TrackNum string `json:"track_num"  xorm:"track_num"  binding:"omitempty,lte=32"`

	WarehouseID   uint64 `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"  xorm:"warehouse_name"`
	LineID        uint64 `json:"line_id"  xorm:"line_id"`
	SourceID      uint64 `json:"source"  xorm:"source"`
	SourceName    string `json:"source_name" xorm:"source_name"`
	ToID          uint64 `json:"to"  xorm:"to"`
	ToName        string `json:"to_name" xorm:"to_name"`
	SendWayID     uint64 `json:"sendway_id"  xorm:"sendway_id"`
	SendWayType   string `json:"sendway_type"  xorm:"sendway_type"`
	SendWayName   string `json:"sendway_name" xorm:"sendway_name"`
}

type PackOrderSimpleCBD struct {
	SellerID        uint64  `json:"seller_id" xorm:"seller_id"`
	ShopID          uint64  `json:"shop_id" xorm:"shop_id"`
	ShopName        string  `json:"shop_name" xorm:"shop_name"`
	PlatformShopID  string  `json:"platform_shop_id" xorm:"platform_shop_id"`
	OrderID         uint64  `json:"order_id,string" xorm:"order_id"`
	OrderTime       int64   `json:"order_time" xorm:"order_time"`
	Platform        string  `json:"platform" xorm:"platform"`
	Status          string  `json:"status" xorm:"status"`
	SN              string  `json:"sn" xorm:"sn"`
	PickNum         string  `json:"pick_num" xorm:"pick_num"`
	ShippingCarrier string  `json:"shipping_carrier" xorm:"shipping_carrier"`
	IsCb            uint8   `json:"is_cb" xorm:"is_cb"`
	SkuType         string  `json:"sku_type"  xorm:"sku_type"`
	NoteBuyer       string  `json:"note_buyer"  xorm:"note_buyer"`
	NoteSeller      string  `json:"note_seller"  xorm:"note_seller"`
	FeeStatus       string  `json:"fee_status" xorm:"fee_status"`
	Price           float64 `json:"price" xorm:"price"`
	Weight          float64 `json:"weight" xorm:"weight"`
	Volume          float64 `json:"volume" xorm:"volume"`
	Length          float64 `json:"length" xorm:"length"`
	Width           float64 `json:"width" xorm:"width"`
	Height          float64 `json:"height" xorm:"height"`

	PackSubDetail []PackSubCBD `json:"pack_sub_detail"`
}

type PackRespCBD struct {
	SellerID uint64 `json:"seller_id"  xorm:"seller_id"`
	RealName string `json:"real_name"  xorm:"real_name"`
	TrackNum string `json:"track_num"  xorm:"track_num"`

	WarehouseID   uint64  `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName string  `json:"warehouse_name"  xorm:"warehouse_name"`
	LineID        uint64  `json:"line_id"  xorm:"line_id"`
	SourceID      uint64  `json:"source_id"  xorm:"source_id"`
	SourceName    string  `json:"source_name" xorm:"source_name"`
	ToID          uint64  `json:"to_id"  xorm:"to_id"`
	ToName        string  `json:"to_name" xorm:"to_name"`
	SendWayID     uint64  `json:"sendway_id"  xorm:"sendway_id"`
	SendWayType   string  `json:"sendway_type"  xorm:"sendway_type"`
	SendWayName   string  `json:"sendway_name" xorm:"sendway_name"`
	Status        string  `json:"status"  xorm:"status"`
	Type          string  `json:"type"  xorm:"type"`
	Weight        float64 `json:"weight"  xorm:"weight"`
	Problem       uint8   `json:"problem"  xorm:"problem"`
	Reason        string  `json:"reason"  xorm:"reason"`
	ManagerNote   string  `json:"manager_note"  xorm:"manager_note"`
	SearchType    string  `json:"search_type"  xorm:"search_type"`
	TmpRackCBD

	PackOrderSimple []PackOrderSimpleCBD `json:"detail"`
}

type OrderPackList struct {
	OrderID        uint64 `json:"order_id,string" xorm:"order_id"`
	PackID         uint64 `json:"pack_id,string" xorm:"pack_id"`
	SourceRecvTime int64  `json:"source_recv_time" xorm:"source_recv_time"`
	Problem        uint8  `json:"problem"  xorm:"problem"`
	Reason         string `json:"reason"  xorm:"reason"`
	ManagerNote    string `json:"manager_note"  xorm:"manager_note"`
	Status         string `json:"status"  xorm:"status"`
	TrackNum       string `json:"track_num"  xorm:"track_num"`
	DependID       string `json:"depend_id"  xorm:"depend_id"`
}

type TmpPack struct {
	RackID      uint64 `json:"-"  xorm:"rack_id"`
	SellerID    uint64 `json:"seller_id"  xorm:"seller_id"`
	RealName    string `json:"real_name"  xorm:"real_name"`
	PackID      uint64 `json:"pack_id,string" xorm:"pack_id"`
	TrackNum    string `json:"track_num"  xorm:"track_num"`
	IsReturn    int    `json:"is_return"  xorm:"is_return"`
	ManagerNote string `json:"manager_note"  xorm:"manager_note"`
}
