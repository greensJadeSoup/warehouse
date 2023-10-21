package bll

import (
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
)

//接口业务逻辑层
type NoticeBL struct {
	Si *cp_api.CheckSessionInfo
}

func NewNoticeBL(si *cp_api.CheckSessionInfo) *NoticeBL {
	return &NoticeBL{Si: si}
}

func (this *NoticeBL) AddNotice(in *cbd.AddNoticeReqCBD) error {
	//todo 查验是否重名

	err := dal.NewNoticeDAL(this.Si).AddNotice(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *NoticeBL) ListNotice(in *cbd.ListNoticeReqCBD) (*cp_orm.ModelList, error) {
	if this.Si.IsManager {
		in.VendorIDList = append(in.VendorIDList, strconv.FormatUint(in.VendorID, 10))
	} else {
		vsList, err := dal.NewVendorSellerDAL(this.Si).ListBySellerID(&cbd.ListVendorSellerReqCBD{
			SellerID: in.SellerID,
			IsPaging: false,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range *vsList {
			in.VendorIDList = append(in.VendorIDList, strconv.FormatUint(v.VendorID, 10))
		}
	}

	ml, err := dal.NewNoticeDAL(this.Si).ListNotice(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *NoticeBL) EditNotice(in *cbd.EditNoticeReqCBD) error {
	md, err := dal.NewNoticeDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("NoticeID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	_, err = dal.NewNoticeDAL(this.Si).EditNotice(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *NoticeBL) DelNotice(in *cbd.DelNoticeReqCBD) error {
	md, err := dal.NewNoticeDAL(this.Si).GetModelByID(in.ID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("NoticeID不存在:" + strconv.FormatUint(in.ID, 10))
	}

	_, err = dal.NewNoticeDAL(this.Si).DelNotice(in)
	if err != nil {
		return err
	}

	return nil
}

