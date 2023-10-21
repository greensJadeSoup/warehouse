package dal

import (
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/dav"
	"warehouse/v5-go-component/cp_error"
)


//数据逻辑层
type ModelDetailDAL struct {
	dav.ModelDetailDAV
}

func NewModelDetailDAL() *ModelDetailDAL {
	return &ModelDetailDAL{}
}

func (this *ModelDetailDAL) List(sellerID, shopID uint64) (*[]cbd.ModelDetailCBD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBList(sellerID, shopID)
}


func (this *ModelDetailDAL) UpdateModelDetail(list *[]cbd.ModelDetailCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBUpdateModelDetail(list)
}


