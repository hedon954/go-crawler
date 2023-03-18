package collector

type OutputData struct {
	TaskName string
	RuleName string
	Url      string
	Time     string
	Data     interface{}
	Struct   DataStruct
}

type Store interface {
	Save(datas ...OutputData) error
}

type DataStruct interface {
	TableName() string
}
