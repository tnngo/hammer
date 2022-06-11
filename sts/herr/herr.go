package herr

import "github.com/tnngo/hammer/sts"

// 200
func OK() *sts.Status {
	return &sts.Status{
		Code: "OK",
	}
}

// 204
func NC() *sts.Status {
	return &sts.Status{}
}

// 400
func ParamError() *sts.Status {
	return &sts.Status{
		Code: "PARAM_ERROR",
	}
}

// 401
func SignError() *sts.Status {
	return &sts.Status{
		Code: "SIGN_ERROR",
	}
}

// 500
func SystemError() *sts.Status {
	return &sts.Status{
		Code: "SYSTEM_ERROR",
	}
}
