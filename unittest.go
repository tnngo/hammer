package hammer

import "github.com/tnngo/hammer/logger"

// 单元测试，且仅能用于单元测试。
func UnitTest(builds ...Builder) {
	(&logger.Console{}).Build()
	for _, v := range builds {
		v.Build()
	}
}
