package herr

import "github.com/tnngo/hammer/sts"

// 200
func OK() *sts.Status {
	return &sts.Status{
		Code:    0,
		Success: true,
	}
}

// 400
func InvalidArgument() *sts.Status {
	return &sts.Status{
		Code:    3,
		Success: false,
		Msg:     "参数错误",
	}
}

// 401
func Unauthenticated() *sts.Status {
	return &sts.Status{
		Code:    16,
		Success: false,
		Msg:     "签权失败",
	}
}

// 403
func PermissionDenied() *sts.Status {
	return &sts.Status{
		Code:    7,
		Success: false,
		Msg:     "没有操作权限",
	}
}

// 500
func Internal() *sts.Status {
	return &sts.Status{
		Code:    13,
		Success: false,
		Msg:     "服务器内部发生异常",
	}
}
