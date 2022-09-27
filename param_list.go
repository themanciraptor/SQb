package sqb

type ParamList struct {
	params  []interface{}
	dialect Dialect
}

func NewParamList(dialect Dialect) *ParamList {
	return &ParamList{
		params:  []interface{}{},
		dialect: dialect,
	}
}

func (p *ParamList) RecordValueAndReturnParam(v interface{}) string {
	for k := range p.params {
		if p.params[k] == v {
			return p.dialect.FormatParam(k)
		}
	}

	p.params = append(p.params, v)
	return p.dialect.FormatParam(len(p.params))
}

func (p *ParamList) GetParamList() []interface{} {
	return p.params
}
