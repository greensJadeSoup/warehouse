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
type SellerDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *SellerDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewSeller())
}

func (this *SellerDAV) DBGetTopID() (uint64, error) {
	var id uint64

	sql := fmt.Sprintf(`SELECT id FROM %s order by id desc limit 1`, this.GetModel().TableName())

	cp_log.Debug(sql)
	_, err := this.SQL(sql).Get(&id)
	if err != nil {
		return 0, cp_error.NewSysError("[SellerDAV][DBGetTopID]:" + err.Error())
	}

	return id, nil
}

func (this *SellerDAV) DBGetModelByAccount(account string) (*model.SellerMD, error) {
	md := model.NewSeller()

	sql := fmt.Sprintf(`SELECT * FROM %s WHERE account='%s'`,
		md.TableName(), account)

	cp_log.Debug(sql)
	hasRow, err := this.SQL(sql).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[SellerDAV][DBGetModelByAccount]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *SellerDAV) DBListSeller(in *cbd.ListSellerReqCBD) (*cp_orm.ModelList, error) {
	searchSQL := fmt.Sprintf(`SELECT s.id,account,real_name,phone,email,s.note,
			allow_login,vs.balance,ds.discount_id,d.name discount_name,d.enable discount_enable
		FROM %[1]s s
		JOIN t_vendor_seller vs
		on s.id = vs.seller_id and vs.enable = 1
		LEFT JOIN t_discount_seller ds
		on vs.vendor_id = ds.vendor_id and s.id = ds.seller_id
		LEFT JOIN t_discount d
		on ds.discount_id = d.id
		WHERE vs.vendor_id = %[2]d`,
		this.GetModel().TableName(), in.VendorID)

	if in.Account != "" {
		keyword := "%" + in.Account + "%"
		searchSQL += fmt.Sprintf(` AND s.account like '%s'`, keyword)
	}

	if in.RealName != "" {
		searchSQL += ` AND s.real_name like '%` + in.RealName + `%'`
	}

	searchSQL += ` order by s.id`

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListSellerRespCBD{})
}

func (this *SellerDAV) DBInsertAccount(md interface{}) error {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[SellerDAV][DBInsertAccount]注册失败,系统繁忙")
	}

	return nil
}

func (this *SellerDAV) DBUpdateSeller(md *model.SellerMD) (int64, error) {
	return this.Session.ID(md.ID).Cols(`real_name,phone,email,allow_login,note`).Update(md)
}

func (this *SellerDAV) DBGetPassword(in *cbd.CheckPasswordReqCBD) (string, error) {
	var storedPassword string

	sql := fmt.Sprintf(`SELECT password FROM %s WHERE account='%s'`,
		this.GetModel().TableName(), in.Account)

	cp_log.Debug(sql)
	hasRow, err := this.SQL(sql).Get(&storedPassword)
	if err != nil {
		return "", cp_error.NewSysError("[SellerDAV][DBGetPassword]:" + err.Error())
	} else if !hasRow {
		return "", nil
	}

	return storedPassword, nil
}

func (this *SellerDAV) DBModifyPassword(in *cbd.ModifyPasswordReqCBD) (int64, error) {
	execSQL := fmt.Sprintf(`UPDATE %s SET password='%s' WHERE account='%s'`,
			this.GetModel().TableName(), in.NewPassword, in.Account)

	cp_log.Debug(execSQL)
	res, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[SellerDAV][DBModifyPassword]:" + err.Error())
	}

	return res.RowsAffected()
}

//strAccount 账号(账号、邮箱、手机)，账号类型（account,email,phone）
func (this *SellerDAV) DBGetModelByID(id uint64) (*model.SellerMD, error) {
	md := model.NewSeller()

	searchSQL := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[SellerDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *SellerDAV) DBDelSeller(in *cbd.DelSellerReqCBD) (int64, error) {
	md := model.NewSeller()
	md.ID = in.SellerID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *SellerDAV) UpdateProfile(md *model.SellerMD) (int64, error) {
	return this.Session.ID(md.ID).Cols("real_name","phone","company_name","email","wechat_num").Update(md)
}
