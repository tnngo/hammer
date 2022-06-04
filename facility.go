package hammer

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/tnngo/hammer/orm"
)

type FacilityType string

type facilityLoader interface {
	load() (facility, error)
}

const (
	_facilityType_active FacilityType = "active"
	_facilityType_grpc   FacilityType = "grpc"
	_facilityType_http   FacilityType = "http"
	_facilityType_remote FacilityType = "remote"
	_facilityType_orm    FacilityType = "orm"
	_facilityType_logger FacilityType = "logger"
)

type facility map[FacilityType]interface{}

var (
	// 对外调用
	_facMap = make(facility)
)

func FacilityStruct(key FacilityType, v interface{}) error {
	if v1, ok := _facMap[key]; ok {
		if reflect.TypeOf(v1).Kind() == reflect.Map {
			b, err := json.Marshal(v1)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(b, v); err != nil {
				return err
			}
		} else {
			return errors.New("必须是struct")
		}

		return nil
	}

	return errors.New("未找到配置key: " + string(key))
}

func (fac facility) loadHttp() *Http {
	if v, ok := fac[_facilityType_http]; ok {
		h := &Http{}
		b, _ := json.Marshal(v)
		json.Unmarshal(b, h)
		return h
	}

	return nil
}

func (fac facility) loadOrm() *orm.Orm {
	if v, ok := fac[_facilityType_orm]; ok {
		o := &orm.Orm{}
		b, _ := json.Marshal(v)
		json.Unmarshal(b, o)
		return o
	}

	return nil
}
