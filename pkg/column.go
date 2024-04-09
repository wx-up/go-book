package pkg

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JsonColumn 理论上来说一切可以被 json 库所处理的类型都能被用作 T
// 不建议使用指针作为 T 的类型
// 如果 T 是指针，那么在 Val 为 nil 的情况下，一定要把 Valid 设置为 false
type JsonColumn[T any] struct {
	Val   T
	Valid bool
}

func (j JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	res, err := json.Marshal(j.Val)
	return res, err
}

func (j *JsonColumn[T]) Scan(src any) error {
	var bs []byte
	switch val := src.(type) {
	case nil:
		return nil
	case []byte:
		bs = val
	case string:
		bs = []byte(val)
	default:
		return fmt.Errorf("JsonColumn.Scan 不支持 src 类型 %v", src)
	}
	if err := json.Unmarshal(bs, &j.Val); err != nil {
		return err
	}
	j.Valid = true
	return nil
}
