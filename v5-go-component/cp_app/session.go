package cp_app

import (
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
)

func SuperAdminValidityCheck(si *cp_api.CheckSessionInfo, url string, vendorID uint64) error {
	cp_log.Info(url)
	//todo 校验接口权限

	if si.AccountType != cp_constant.USER_TYPE_SUPER_MANAGER {
		return cp_error.NewNormalError("用户不是超管", cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
	}

	if vendorID != si.VendorDetail[0].VendorID {
		return cp_error.NewNormalError("用户不属于该供应商", cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
	}

	si.IsSuperManager = true
	si.ManagerID = si.UserID

	return nil
}

func AdminValidityCheck(si *cp_api.CheckSessionInfo, url string, vendorID uint64) error {
	cp_log.Info(url)
	//todo 校验接口权限

	if si.AccountType != cp_constant.USER_TYPE_MANAGER && si.AccountType != cp_constant.USER_TYPE_SUPER_MANAGER {
		return cp_error.NewNormalError("用户不是管理员", cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
	}

	if vendorID != si.VendorDetail[0].VendorID {
		return cp_error.NewNormalError("用户不属于该供应商", cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
	}

	if si.AccountType == cp_constant.USER_TYPE_SUPER_MANAGER {
		si.IsSuperManager = true
	}

	si.ManagerID = si.UserID
	si.IsManager = true

	return nil
}

func SellerValidityCheck(si *cp_api.CheckSessionInfo, vendorID uint64, sellerID uint64) error {
	if vendorID > 0 {
		found := false
		for _, v := range si.VendorDetail {
			if vendorID == v.VendorID {
				found = true
			}
		}
		if !found {
			return cp_error.NewNormalError("用户不属于该供应商", cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
		}
	}

	if si.AccountType != cp_constant.USER_TYPE_SELLER {
		return cp_error.NewNormalError("用户类型错误", cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
	}

	if sellerID != si.UserID {
		return cp_error.NewNormalError("请确认用户id是否正确", cp_constant.RESPONSE_CODE_PARAMPARSE_FAIL)
	}

	return nil
}
