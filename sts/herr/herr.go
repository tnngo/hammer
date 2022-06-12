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
	}
}

// 401
func Unauthenticated() *sts.Status {
	return &sts.Status{
		Code:    16,
		Success: false,
	}
}

// 500
func Internal() *sts.Status {
	return &sts.Status{
		Code:    13,
		Success: false,
	}
}
