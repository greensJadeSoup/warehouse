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
type VendorDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *VendorDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewVendor())
}

func (this *VendorDAV) DBGetModelByID(id uint64) (*model.VendorMD, error) {
	md := model.NewVendor()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *VendorDAV) DBGetModelByName(name string) (*model.VendorMD, error) {
	md := model.NewVendor()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE name='%s'`, md.TableName(), name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *VendorDAV) DBListVendor(in *cbd.ListVendorReqCBD) (*cp_orm.ModelList, error) {
	searchSQL := fmt.Sprintf(`SELECT v.id,v.name,v.balance,v.order_fee,m.account super_manager_account,
				m.real_name super_manager_real_name
				FROM %[1]s v
				LEFT JOIN t_manager m
				on v.id = m.vendor_id 
				where m.type = 'super_manager'`,
		this.GetModel().TableName())

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListVendorRespCBD{})
}

func DBVendorDeduct(da *cp_orm.DA, vendorID uint64) (int64, error) {
	execSQL := fmt.Sprintf(`update db_base.t_vendor v set balance=balance-order_fee where id = %[1]d`, vendorID)

	execRow, err := da.Session.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow.RowsAffected()
}

func (this *VendorDAV) DBInsert(md interface{}) error {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[VendorDAV][DBInsertVendor]注册失败,系统繁忙")
	}

	return nil
}
