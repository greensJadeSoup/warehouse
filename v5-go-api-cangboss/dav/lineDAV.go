package dav

import (
	"fmt"
	"strings"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type LineDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *LineDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewLine())
}

func (this *LineDAV) DBGetModelByID(id uint64) (*model.LineMD, error) {
	md := model.NewLine()

	searchSQL := fmt.Sprintf(`select * from %[1]s WHERE id=%[2]d`,
		md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[LineDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *LineDAV) DBGetModelDetailByID(id uint64) (*cbd.GetLineCBD, error) {
	md := &cbd.GetLineCBD{}

	searchSQL := fmt.Sprintf(`select l.*,w1.name source_name,w2.name to_name 
		from %[1]s l
		LEFT JOIN t_warehouse w1
		on l.source = w1.id
		LEFT JOIN t_warehouse w2
		on l.to = w2.id
		WHERE l.id=%[2]d`,
		this.GetModel().TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[LineDAV][DBGetModelDetailByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *LineDAV) DBGetModelDetailByIDList(idList []string) (*[]cbd.GetLineCBD, error) {
	list := &[]cbd.GetLineCBD{}

	searchSQL := fmt.Sprintf(`select l.*,w1.name source_name,w2.name to_name 
		from %[1]s l
		LEFT JOIN t_warehouse w1
		on l.source = w1.id
		LEFT JOIN t_warehouse w2
		on l.to = w2.id
		WHERE l.id in (%[2]s)`,
		this.GetModel().TableName(), strings.Join(idList, ","))

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[LineDAV][DBGetModelDetailByIDList]:" + err.Error())
	}

	return list, nil
}

func (this *LineDAV) DBInsertAccount(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if execRow == 0 {
		return cp_error.NewSysError("[LineDAV][DBInsertAccount]注册失败,系统繁忙")
	}

	return nil
}

func (this *LineDAV) DBListLine(in *cbd.ListLineReqCBD) (*cp_orm.ModelList, error) {
	var whCondSQL string

	if in.VendorID > 0 {
		whCondSQL += fmt.Sprintf(` AND l.vendor_id = %[1]d`, in.VendorID)
	}

	if in.Source > 0 {
		whCondSQL += fmt.Sprintf(` AND l.source = %[1]d`, in.Source)
	}

	if in.To > 0 {
		whCondSQL += fmt.Sprintf(` AND l.to = %[1]d`, in.To)
	}

	if len(in.WarehouseIDList) > 0 {
		whCondSQL += fmt.Sprintf(` AND (l.source in (%[1]s) or l.to in (%[1]s))`, strings.Join(in.WarehouseIDList, ","))
	}

	if len(in.LineIDList) > 0 {
		whCondSQL += fmt.Sprintf(` AND (l.id in (%[1]s))`, strings.Join(in.LineIDList, ","))
	}

	searchSQL := fmt.Sprintf(`SELECT l.id,l.vendor_id,l.source,l.to,l.sort,l.note,w1.name source_whr,w2.name to_whr
				FROM %[1]s l
				LEFT JOIN t_warehouse w1
				on l.source=w1.id and l.vendor_id = w1.vendor_id
				LEFT JOIN t_warehouse w2
				on l.to=w2.id and l.vendor_id = w2.vendor_id
				WHERE 1=1 %[2]s order by l.sort,l.id ASC`,
		this.GetModel().TableName(), whCondSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListLineRespCBD{})
}

func (this *LineDAV) DBListLineInternal(in *cbd.ListLineReqCBD) (*[]cbd.ListLineRespCBD, error) {
	var whCondSQL string

	if in.Source > 0 {
		whCondSQL += fmt.Sprintf(` AND l.source = %[1]d`, in.Source)
	}

	if in.To > 0 {
		whCondSQL += fmt.Sprintf(` AND l.to = %[1]d`, in.To)
	}

	list := &[]cbd.ListLineRespCBD{}
	searchSQL := fmt.Sprintf(`SELECT * FROM %[1]s l WHERE l.vendor_id = %[2]d %[3]s`,
		this.GetModel().TableName(), in.VendorID, whCondSQL)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(list)
	if err != nil {
		return nil, cp_error.NewSysError("[LineDAV][DBListLineInternal]:" + err.Error())
	}

	return list, nil
}

func (this *LineDAV) DBUpdateLine(md *model.LineMD) (int64, error) {
	return this.Session.ID(md.ID).Cols("source","to","sort","note").Update(md)
}

func (this *LineDAV) DBDelLine(in *cbd.DelLineReqCBD) (int64, error) {
	md := model.NewLine()
	md.ID = in.LineID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}

	return execRow, nil
}

