package dal

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)


//数据逻辑层
type ConsumableDAL struct {
	dav.ConsumableDAV
	Si *cp_api.CheckSessionInfo
}

func NewConsumableDAL(si *cp_api.CheckSessionInfo) *ConsumableDAL {
	return &ConsumableDAL{Si: si}
}

func (this *ConsumableDAL) GetModelByID(id uint64) (*model.ConsumableMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(id)
}

func (this *ConsumableDAL) GetModelByName(vendorID uint64, name string) (*model.ConsumableMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByName(vendorID, name)
}

func (this *ConsumableDAL) AddConsumable(in *cbd.AddConsumableReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	md := &model.ConsumableMD {
		VendorID: in.VendorID,
		Name: in.Name,
		Note: in.Note,
	}
	err = this.DBInsert(md)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	//err = DiscountAddConsumable(&this.DA, in.VendorID, md.ID, in.Name)
	//if err != nil {
	//	return err
	//}

	return this.Commit()
}

func (this *ConsumableDAL) EditConsumable(in *cbd.EditConsumableReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	md, err := this.DBGetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("ID不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.VendorID != in.VendorID {
		return cp_error.NewNormalError("无该耗材访问权:" + strconv.FormatUint(in.ID, 10))
	}

	if md.Name != in.Name {
		mdNew, err := this.DBGetModelByName(in.VendorID, in.Name)
		if err != nil {
			return err
		} else if mdNew != nil {
			return cp_error.NewNormalError("相同耗材名已存在:" + in.Name)
		}
	}

	mdNew := &model.ConsumableMD {
		ID: in.ID,
		Name: in.Name,
		Note: in.Note,
	}

	_, err = this.DBUpdateConsumable(mdNew)
	if err != nil {
		return err
	}

	if md.Name != in.Name { //把计价组中json的仓库名也改了
		err = DiscountEditConsumable(&this.DA, in.VendorID, in.ID, in.Name)
		if err != nil {
			return err
		}
	}

	return this.Commit()
}

func (this *ConsumableDAL) ListConsumable(in *cbd.ListConsumableReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListConsumable(in)
}

func (this *ConsumableDAL) DelConsumable(in *cbd.DelConsumableReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	_, err = this.DBDelConsumable(in)
	if err != nil {
		return err
	}

	err = DiscountDelConsumable(&this.DA, in.VendorID, in.ID)
	if err != nil {
		return err
	}

	return this.Commit()
}
