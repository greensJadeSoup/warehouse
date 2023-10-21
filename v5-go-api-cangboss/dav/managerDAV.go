package dav

import (
	"fmt"
	"strings"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type ManagerDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ManagerDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewManager())
}

func (this *ManagerDAV) DBGetModelByAccount(account string) (*model.ManagerMD, error) {
	md := model.NewManager()

	sql := fmt.Sprintf(`SELECT * FROM %s WHERE account='%s'`,
		md.TableName(), account)

	cp_log.Debug(sql)
	hasRow, err := this.SQL(sql).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ManagerDAV][DBGetModelByAccount]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ManagerDAV) DBInsertAccount(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[ManagerDAV][DBInsertAccount]注册失败,系统繁忙")
	}

	return nil
}

func (this *ManagerDAV) DBGetPassword(in *cbd.CheckPasswordReqCBD) (string, error) {
	var storedPassword string

	sql := fmt.Sprintf(`SELECT password FROM %s WHERE account='%s'`,
		this.GetModel().TableName(), in.Account)

	cp_log.Debug(sql)
	hasRow, err := this.SQL(sql).Get(&storedPassword)
	if err != nil {
		return "", cp_error.NewSysError("[ManagerDAV][DBGetPassword]:" + err.Error())
	} else if !hasRow {
		return "", nil
	}

	return storedPassword, nil
}

func (this *ManagerDAV) DBModifyPassword(in *cbd.ModifyPasswordReqCBD) (int64, error) {
	execSQL := fmt.Sprintf(`UPDATE %s SET password='%s' WHERE account='%s'`,
			this.GetModel().TableName(), in.NewPassword, in.Account)

	cp_log.Debug(execSQL)
	res, err := this.Exec(execSQL)
	if err != nil {
		return 0, cp_error.NewSysError("[ManagerDAV][DBModifyPassword]:" + err.Error())
	}

	return res.RowsAffected()
}

//strAccount 账号(账号、邮箱、手机)，账号类型（account,email,phone）
func (this *ManagerDAV) DBGetModelByID(id uint64) (*model.ManagerMD, error) {
	md := model.NewManager()

	sql := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(sql)
	hasRow, err := this.SQL(sql).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ManagerDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ManagerDAV) DBUpdateManager(md *model.ManagerMD) (int64, error) {
	return this.Session.ID(md.ID).Cols("warehouse_id","real_name","phone","email","allow_login","warehouse_role","note").Update(md)
}

func (this *ManagerDAV) DBListManager(in *cbd.ListManagerReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string

	//如果是管理员的话，需要根据条件过滤属于哪些仓库
	if in.WarehouseID != "" && in.Type != cp_constant.USER_TYPE_SUPER_MANAGER {
		condSQL = " AND ("
		for _, v := range strings.Split(in.WarehouseID, ",") {
			condSQL += fmt.Sprintf(` FIND_IN_SET("%s",warehouse_id) or`, v)
		}
		condSQL = strings.TrimRight(condSQL, "or") + ")"
	}

	if in.Type != "" {
		condSQL += fmt.Sprintf(` AND m.type = '%s'`, in.Type)
	}

	if in.Account != "" {
		condSQL += ` AND account like '%` + in.Account + `%'`
	}

	if in.RealName != "" {
		condSQL += ` AND real_name like '%` + in.RealName + `%'`
	}

	searchSQL := fmt.Sprintf(`SELECT m.id,m.vendor_id,warehouse_id,w.name warehouse_name,account,
			m.type,m.warehouse_role,real_name,phone,email,allow_login,m.note
			FROM %s m
			LEFT JOIN db_warehouse.t_warehouse w
			on m.warehouse_id = w.id
			WHERE m.vendor_id=%d %s`,
		this.GetModel().TableName(), in.VendorID, condSQL)

	searchSQL += ` order by warehouse_id,id`

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListManagerRespCBD{})
}

func (this *ManagerDAV) DBDelManager(in *cbd.DelManagerReqCBD) (int64, error) {
	md := model.NewManager()
	md.ID = in.ManagerID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

func (this *ManagerDAV) UpdateProfile(md *model.ManagerMD) (int64, error) {
	return this.Session.ID(md.ID).Cols("real_name","phone","company_name","email","wechat_num").Update(md)
}


