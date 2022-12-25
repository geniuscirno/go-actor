package attributes

type Attributes struct {
	m map[string]string
}

func New(m map[string]string) *Attributes {
	a := &Attributes{m: make(map[string]string)}
	for k, v := range m {
		a.m[k] = v
	}
	return a
}

func (a *Attributes) WithValues(values map[string]string) *Attributes {
	n := &Attributes{m: make(map[string]string, len(a.m)+len(values))}
	for k, v := range a.m {
		n.m[k] = v
	}
	for k, v := range values {
		n.m[k] = v
	}
	return n
}

func (a *Attributes) WithValue(key string, value string) *Attributes {
	return a.WithValues(map[string]string{key: value})
}

func (a *Attributes) Match(o *Attributes) bool {
	if o == nil {
		return true
	}

	for k, v := range o.m {
		val, ok := a.m[k]
		if !ok {
			return false
		}

		if val != v {
			return false
		}
	}
	return true
}
