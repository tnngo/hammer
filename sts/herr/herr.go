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
		Code:    400,
		Success: false,
	}
}

// 401
func Unauthenticated() *sts.Status {
	return &sts.Status{
		Code:    401,
		Success: false,
	}
}

// 500
func Internal() *sts.Status {
	return &sts.Status{
		Code:    500,
		Success: false,
	}
}
