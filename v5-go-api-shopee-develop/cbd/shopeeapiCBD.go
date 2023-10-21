package cbd

//------------------------ push ------------------------
type OrderStatusPush struct {
	Data		struct{
		OrderSN		string		`json:"ordersn"`
		Status		string		`json:"status"`
		TrackNum	string		`json:"tracking_no"`
		UpdateTime	int64		`json:"update_time"`
	}	`json:"data"`
	ShopID 		uint64		`json:"shop_id"`
	Code		int		`json:"code"`
}
//------------------------ req ------------------------
type CommonReqCBD struct {
	ApiUrl		string
	PartnerID 	uint64
	Key 		string
}

type GetAccessTokenReqCBD struct {
	PartnerID 		uint64		`json:"partner_id"`
	ShopID	 		uint64		`json:"shop_id,omitempty"`
	MainAccount		uint64		`json:"main_account_id,omitempty"`
	Code			string		`json:"code"`
}

type RefreshAccessTokenReqCBD struct {
	RefreshToken		string		`json:"refresh_token"`
	PartnerID 		uint64		`json:"partner_id"`
	ShopID	 		uint64		`json:"shop_id"`
}

type GetShopInfoReqCBD struct {
	AccessToken 		string		`json:"access_token"`
	ShopID	 		uint64		`json:"shop_id"`
}

type OrderItemCBD struct {
	SN 			string		`json:"order_sn"`
	TrackingNumber		string		`json:"tracking_number,omitempty"`
}

type GetDocumentDataInfoItem struct {
	Key		string		`json:"key"`
}

type GetDocumentDataInfoCBD struct {
	SN 			string		`json:"order_sn"`
	RecAddressInfo		[]GetDocumentDataInfoItem	`json:"recipient_address_info"`
}

type CreateShippingDocumentReqCBD struct {
	OrderList 		[]OrderItemCBD		`json:"order_list"`
}

type DownloadShippingDocumentReqCBD struct {
	OrderList 		[]OrderItemCBD		`json:"order_list"`
}

type SellerInfoCBD struct {
	Address 		string		`json:"address"`
	Name			string		`json:"name"`
	Zipcode			string		`json:"zipcode"`
	Region			string		`json:"region"`
	Phone			string		`json:"phone"`
}

type GenerateFirstMileTrackingNumReqCBD struct {
	DeclareDate		string			`json:"declare_date"`
	SellerInfo 		SellerInfoCBD		`json:"seller_info"`
}

type BindFirstMileTrackingNumReqCBD struct {
	FirstMileTrackingNumber		string			`json:"first_mile_tracking_number"`
	ShipmentMethod			string			`json:"shipment_method"`
	Region				string			`json:"region"`
	LogisticsChannelID		int			`json:"logistics_channel_id"`
	OrderList 			[]OrderItemCBD		`json:"order_list"`
}

//------------------------ resp ------------------------
type GetAccessTokenRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`

	RefreshToken		string		`json:"refresh_token"`
	AccessToken		string		`json:"access_token"`
	ExpireIn		int64		`json:"expire_in"`

	ShopIDList		[]uint64	`json:"shop_id_list"`
}

type RefreshAccessTokenRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`

	RefreshToken		string		`json:"refresh_token"`
	AccessToken		string		`json:"access_token"`
	ExpireIn		int64		`json:"expire_in"`
}

type GetShopInfoRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`

	ShopName		string		`json:"shop_name"`
	ShopID	 		uint64		`json:"shop_id"`
	Status			string		`json:"status"`
	IsCB			bool		`json:"is_cb"`
	IsSIP			bool		`json:"is_sip"`
	IsCNSC			bool		`json:"is_cnsc"`
	ExpireTime	 	int64		`json:"expire_time"`
	Region			string		`json:"region"`
}

type GetShopProfileRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`

	Response struct {
		ShopName		string		`json:"shop_name"`
		ShopLogo		string		`json:"shop_logo"`
		Description		string		`json:"description"`
	}`json:"response"`
}


type GetItemListRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`

	Response struct {
		TotalCount		uint64		`json:"total_count"`
		HasNext			bool		`json:"has_next_page"`
		NextOffset		int		`json:"next_offset"`
		Item	[]struct{
			ItemID			uint64		`json:"item_id"`
			ItemStatus		string		`json:"item_status"`
			UpdateTime		int64		`json:"update_time"`

		}	`json:"item"`
	}		`json:"response"`
}

type ShopeeItemBaseInfoCBD struct {
	ID			int64		`json:"-"`
	ItemID			uint64		`json:"item_id"`
	CategoryID		uint64		`json:"category_id"`
	ItemName		string		`json:"item_name"`
	ItemStatus		string		`json:"item_status"`
	Description		string		`json:"description"`
	ItemSku			string		`json:"item_sku"`
	Weight			string		`json:"weight"`
	HasModel		bool		`json:"has_model"`
	UpdateTime		int64		`json:"update_time"`

	Image			struct{
		ImageUrlList		[]string		`json:"image_url_list"`
	}		`json:"image"`
}

type GetItemBaseInfoRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`

	Response struct {
		ItemList	[]ShopeeItemBaseInfoCBD	`json:"item_list"`
	}		`json:"response"`
}

type GetModelListRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`

	Response struct {
		ItemID		uint64	`json:"-"`

		TierVariation	[]struct{
			Name			string		`json:"name"`
			OptionList		[]struct{
				Option	string		`json:"option"`
				Image	struct{
					ImageUrl	string		`json:"image_url"`
				}	`json:"image"`
			}		`json:"option_list"`
		}	`json:"tier_variation"`

		Model	[]struct{
			ModelID		uint64		`json:"model_id"`
			ModelSku	string		`json:"model_sku"`
			TierIndex	[]int		`json:"tier_index"`
			Images		string		`json:"-"`
		}	`json:"model"`
	}		`json:"response"`
}


type GetOrderListRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`

	Response struct {
		More			bool		`json:"more"`
		NextCursor		string		`json:"next_cursor"`
		OrderList	[]struct{
			OrderSN			string		`json:"order_sn"`
			OrderStatus		string		`json:"order_status"`
		}	`json:"order_list"`
	}		`json:"response"`
}


type GetOrderDetailRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`

	Response struct {
		OrderList	[]struct{
			SN			string		`json:"order_sn"`
			Region			string		`json:"region"`
			Currency		string		`json:"currency"`
			CashOnDelivery		bool		`json:"cod"`
			CashOnDeliveryInt	int8
			TotalAmount		float64		`json:"total_amount"`
			Status			string
			PlatformStatus		string		`json:"order_status"`
			ShippingCarrier		string		`json:"shipping_carrier"`
			ShipByDate		int64		`json:"ship_by_date"`
			PaymentMethod		string		`json:"payment_method"`
			NoteBuyer		string		`json:"message_to_seller"`
			PlatformCreateTime	int64		`json:"create_time"`
			PlatformUpdateTime	int64		`json:"update_time"`
			BuyerUserID		uint64		`json:"buyer_user_id"`
			BuyerUsername		string		`json:"buyer_username"`
			NoteSeller		string		`json:"-"`
			PayTime			int64		`json:"pay_time"`
			PickupTime		int64		`json:"-"`
			CancelBy		string		`json:"cancel_by"`
			CancelReason		string		`json:"cancel_reason"`
			BuyerCancelReason	string		`json:"buyer_cancel_reason"`
			IsCb			int8		`json:"is_cb"`

			RecvAddrStr		string
			RecvAddr		struct{
				Name		string		`json:"name"`
				Phone		string		`json:"phone"`
				Town		string		`json:"town"`
				District	string		`json:"district"`
				City		string		`json:"city"`
				State		string		`json:"state"`
				Region		string		`json:"region"`
				Zipcode		string		`json:"zipcode"`
				FullAddress	string		`json:"full_address"`
			} `json:"recipient_address"`

			ItemListStr		string
			ItemCount		int
			ItemList		[]struct{
				ItemID			int64		`json:"item_id"`
				ItemName		string		`json:"item_name"`
				ItemSKU			string		`json:"item_sku"`
				ModelID			int64		`json:"model_id"`
				ModelName		string		`json:"model_name"`
				ModelSKU		string		`json:"model_sku"`
				Weight			float64		`json:"weight"`
				Count			int		`json:"model_quantity_purchased"`
				OriPri			float64		`json:"model_original_price"`
				DiscPri			float64		`json:"model_discounted_price"`
				ImageInfo		struct{
					ImageUrl	string		`json:"image_url"`
				} `json:"image_info"`
			} `json:"item_list"`

			PackageListStr		string
			PackageList []struct{
				PackageNumber		string		`json:"package_number"`
				LogisticsStatus		string		`json:"logistics_status"`
				ShippingCarrier		string		`json:"shipping_carrier"`
				ItemList		[]struct{
					ItemID			int64		`json:"item_id"`
					ModelID			int64		`json:"model_id"`
				} `json:"item_list"`
			} `json:"package_list"`

		}	`json:"order_list"`
	}		`json:"response"`
}

type GetShipParamRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
	Response struct {
		InfoNeeded	struct {
			DropOff			[]string		`json:"dropoff"` //线下（可能这么翻译）
			Pickup			[]string		`json:"pickup"`  //
			NonIntegrated		[]string		`json:"non_integrated"`  //非集成
		}	`json:"info_needed"`
		DropOff	struct {
			SlugList	[]struct {
				Slug		string		`json:"slug"`
				SlugName	string		`json:"slug_name"`
			}	`json:"slug_list"`
		}	`json:"dropoff"`
		PickUp	struct {
			AddressList	[]struct {
				AddressID		uint64		`json:"address_id"`
				Region			string		`json:"region"`
				State			string		`json:"state"`
				City			string		`json:"city"`
				Address			string		`json:"address"`
				Zipcode			string		`json:"zipcode"`
				District		string		`json:"district"`
				Town			string		`json:"town"`
				AddressFlag		[]string	`json:"address_flag"`
				TimeSlotList		[]struct {
					Date			int64	`json:"date"`
					PickupTimeID		string	`json:"pickup_time_id"`
				}  `json:"time_slot_list"`
			}	`json:"address_list"`
		}	`json:"pickup"`
	}		`json:"response"`
}

type CreateFaceDocumentRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
	Response struct {
		ResultList	[]struct{
			OrderSN			string		`json:"order_sn"`
			FailError		string		`json:"fail_error"`
			FailMessage		string		`json:"fail_message"`
		}	`json:"result_list"`
	}		`json:"response"`
}

type GetDocumentResultRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
	Response struct {
		ResultList	[]struct{
			OrderSN			string		`json:"order_sn"`
			Status			string		`json:"status"`
			FailError		string		`json:"fail_error"`
			FailMessage		string		`json:"fail_message"`
		}	`json:"result_list"`
	}		`json:"response"`
}

type GetTrackNumRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
	Response struct {
		TrackingNumber			string		`json:"tracking_number"`
		FirstMileTrackingNumber		string		`json:"first_mile_tracking_number"`
		LastMileTrackingNumber		string		`json:"last_mile_tracking_number"`
	}		`json:"response"`
}

type GetTrackInfoItem struct {
	UpdateTime			int64		`json:"update_time"`
	Description			string		`json:"description"`
	LogisticsStatus			string		`json:"logistics_status"`
}

type GetTrackInfoRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
	Response struct {
		TrackingInfo			[]GetTrackInfoItem		`json:"tracking_info"`
	}		`json:"response"`
}

type GetAddressListRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
	Response struct {
		ShowPickupAddress		bool		`json:"show_pickup_address"`
		AddressList []struct {
			AddressID		uint64		`json:"address_id"`
			Region			string		`json:"region"`
			State			string		`json:"state"`
			City			string		`json:"city"`
			Address			string		`json:"address"`
			Zipcode			string		`json:"zipcode"`
			District		string		`json:"district"`
			Town			string		`json:"town"`
			AddressType		[]string	`json:"address_type"`
		}	`json:"address_list"`
	}		`json:"response"`
}

type DownloadFaceDocumentRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
}

type GetChannelListRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
	Response struct {
		AddressList []struct {
			LogisticsChannelID		int		`json:"logistics_channel_id"`
			LogisticsChannelName		string		`json:"logistics_channel_name"`
			ShipmentMethod			string		`json:"shipment_method"`
		}	`json:"logistics_channel_list"`
	}		`json:"response"`
}

type GetFirstMileTrackingNumDetailRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
	Response struct {
		LogisticsChannelID		int		`json:"logistics_channel_id"`
		FirstMileTrackingNumber		string		`json:"first_mile_tracking_number"`
		Status				string		`json:"status"`
		ShipmentMethod			string		`json:"shipment_method"`
		DeclareDate			string		`json:"declare_date"`
		OrderList []struct {
			OrderSN				string		`json:"order_sn"`
			PackageNumber			string		`json:"package_number"`
			SlsTrackingNumber		string		`json:"sls_tracking_number"`
			PickUpDone			bool		`json:"pick_up_done"`
			ArrivedTransitWarehouse		bool		`json:"arrived_transit_warehouse"`
		}		`json:"order_list"`
	}		`json:"response"`
}

type GenerateFirstMileTrackNumRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
	Response struct {
		FirstMileTrackingNumberList		[]string		`json:"first_mile_tracking_number_list"`
	}
}

type OrderResultCBD struct {
	OrderSN			string		`json:"order_sn"`
	FailError		string		`json:"fail_error"`
	FailMessage		string		`json:"fail_message"`
}

type BindFirstMileTrackNumRespCBD struct {
	RequestID		string		`json:"request_id"`
	Error			string		`json:"error"`
	Message			string		`json:"message"`
	Warning			[]struct {
		OrderSN		string		`json:"order_sn"`
	}
	Response struct {
		FirstMileTrackingNumber		string			`json:"first_mile_tracking_number"`
		ResultList			[]OrderResultCBD	`json:"order_list"`
	}	`json:"response"`
}
