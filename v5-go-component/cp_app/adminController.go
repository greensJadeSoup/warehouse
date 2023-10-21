package cp_app

type AdminController struct {
	BaseController
}

func (this *AdminController) IsAdmin() bool {
	return true
}

