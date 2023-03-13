package collector

type OutputData struct {
	TaskName string
	RuleName string
	Url      string
	Time     string
	Data     interface{}
}

type Store interface {
	Save(datas ...OutputData) error
}
