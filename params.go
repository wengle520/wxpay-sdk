package wxpay

import "strconv"

type Params map[string]string

// map已经是引用类型了，所以不需要 *Params
func (p Params) SetStringParam(k, v string) Params {
	p[k] = v
	return p
}

// 不需要判断是否存在，不存在的返回空值
func (p Params) GetStringParam(k string) string {
	v, _ := p[k]
	return v
}

func (p Params) SetInt64Param(k string, v int64) Params {
	p[k] = strconv.FormatInt(v, 10)
	return p
}

func (p Params) GetInt64Param(k string) int64 {
	v, _ := strconv.ParseInt(p.GetStringParam(k), 10, 64)
	return v
}

// 判断key是否存在
func (p Params) containsKey(k string) bool {
	_, ok := p[k]
	return ok
}
