package fetcher

type Temp struct {
	data map[string]interface{}
}

func (t *Temp) Get(key string) interface{} {
	return t.data[key]
}

func (t *Temp) Set(key string, value interface{}) error {
	if t.data == nil {
		t.data = make(map[string]interface{}, 8)
	}
	t.data[key] = value
	return nil
}

func (t *Temp) Copy() *Temp {
	n := &Temp{}
	n.data = make(map[string]interface{})
	for k, v := range t.data {
		n.data[k] = v
	}
	return n
}
