package events

import (
	"encoding/json"
	"sync"
)

type JsonEncoder struct {
	data any
	once sync.Once
	bs   []byte
	err  error
}

func NewJsonEncoder(data any) *JsonEncoder {
	return &JsonEncoder{data: data}
}

func (j *JsonEncoder) Bytes() {
	j.once.Do(func() {
		j.bs, j.err = json.Marshal(j.data)
	})
}

func (j *JsonEncoder) Encode() ([]byte, error) {
	j.Bytes()
	return j.bs, j.err
}

func (j *JsonEncoder) Length() int {
	j.Bytes()
	if j.err != nil {
		return 0
	}
	return len(j.bs)
}
