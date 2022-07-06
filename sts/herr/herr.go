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
func InvalidArgument(msg string) *sts.Status {
	return &sts.Status{
		Code:    3,
		Success: false,
		Msg:     msg,
	}
}

// 401
func Unauthenticated(msg string) *sts.Status {
	return &sts.Status{
		Code:    16,
		Success: false,
		Status:  msg,
	}
}

// 403
func PermissionDenied(msg string) *sts.Status {
	return &sts.Status{
		Code:    7,
		Success: false,
		Status:  msg,
	}
}

// 500
func Internal(msg string) *sts.Status {
	return &sts.Status{
		Code:    13,
		Success: false,
		Msg:     msg,
	}
}
