package dav

import (
	"fmt"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type VendorSellerDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *VendorSellerDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewVendorSeller())
}

func (this *VendorSellerDAV) DBGetModelByID(id uint64) (*model.VendorSellerMD, error) {
	md := model.NewVendorSeller()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,seller_id,balance FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[VendorSellerDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *VendorSellerDAV) DBGetModelByVendorIDSellerID(vendorID, sellerID uint64) (*model.VendorSellerMD, error) {
	md := model.NewVendorSeller()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,seller_id,balance FROM %s WHERE vendor_id=%d and seller_id=%d and enable = 1`,
		md.TableName(), vendorID, sellerID)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[VendorSellerDAV][DBGetModelByVendorIDSellerID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *VendorSellerDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[VendorSellerDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *VendorSellerDAV) DBListBySellerID(in *cbd.ListVendorSellerReqCBD) (*[]cbd.ListVendorSellerRespCBD, error) {
	list := &[]cbd.ListVendorSellerRespCBD{}

	searchSQL := fmt.Sprintf(`SELECT vendor_id,seller_id,balance FROM %s WHERE seller_id=%d and enable = 1`,
		this.GetModel().TableName(), in.SellerID)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *VendorSellerDAV) DBListByVendorID(in *cbd.ListVendorSellerReqCBD) (*[]cbd.ListVendorSellerRespCBD, error) {
	var condSQL string

	list := &[]cbd.ListVendorSellerRespCBD{}

	if in.SellerKey != "" {
		condSQL += ` AND (s.id = '` + in.SellerKey + `' or s.real_name like '%` + in.SellerKey + `%')`
	}

	searchSQL := fmt.Sprintf(`SELECT vs.vendor_id,vs.seller_id,vs.balance,s.real_name
			FROM %[1]s vs
			INNER JOIN t_seller s
			on vs.seller_id = s.id
			WHERE vs.vendor_id=%[2]d and enable = 1 %[3]s`,
		this.GetModel().TableName(), in.VendorID, condSQL)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func (this *VendorSellerDAV) DBUpdateVendorSeller(md *model.VendorSellerMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *VendorSellerDAV) DBDelVendorSeller(in *cbd.DelVendorSellerReqCBD) (int64, error) {
	md := model.NewVendorSeller()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *VendorSellerDAV) DBListBalance(in *cbd.ListBalanceReqCBD) (*[]cbd.ListBalanceRespCBD, error) {
	list := &[]cbd.ListBalanceRespCBD{}

	searchSQL := fmt.Sprintf(`SELECT vs.vendor_id,vs.seller_id,vs.balance
			FROM %[1]s vs
			WHERE vs.seller_id=%[2]d and enable = 1`,
		this.GetModel().TableName(), in.SellerID)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return list, nil
}

func DBUpdateSellerBalance(da *cp_orm.DA, md *model.VendorSellerMD) (int64, error) {
	execSQL := fmt.Sprintf(`update db_base.t_vendor_seller set balance=%0.2f where vendor_id = %d and seller_id = %d`,
		md.Balance, md.VendorID, md.SellerID)

	cp_log.Debug(execSQL)
	execRow, err := da.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

