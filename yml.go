package hammer

import (
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type yml struct {
	filename string
}

func YmlFacilityLoader(filename string) facilityLoader {
	return &yml{
		filename: filename,
	}
}

func (y *yml) load() (facility, error) {
	var filename string
	if y.filename != "" {
		filename = y.filename
	} else {
		filename = "./hammer.yml"
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		if y.filename == "" && strings.Contains(err.Error(), "no such file or directory") {
			return nil, nil
		}
		return nil, err
	}

	fac := make(facility)
	if err := yaml.Unmarshal(b, fac); err != nil {
		return nil, err
	}

	/** active配置 */
	var (
		activeFilename string
		strs           []string
		suffix         string
	)

	// 通过active配置找出活跃文件。
	if strings.Contains(filename, ".yml") {
		suffix = ".yml"
	} else {
		return nil, errors.New("yml文件后缀应该是.yml, filename: " + filename)
	}

	strs = strings.Split(filename, suffix)
	if len(strs) == 2 && (strs[1] == "") {
		activeFilename = strs[0] + "_" + fac[_facilityType_active].(string) + suffix
	} else {
		return nil, errors.New("活跃的yml文件名违反active约定, filename: " + filename)
	}

	// 读取活跃文件。
	activebody, err := os.ReadFile(activeFilename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(activebody, fac); err != nil {
		return nil, err
	}

	return fac, nil
}
