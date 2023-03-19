package gorm

import (
	"encoding/json"
	"sync"

	"github.com/hedon954/go-crawler/collector"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormStore struct {
	lock       sync.Mutex
	dataDocker []collector.OutputData
	db         *gorm.DB
	options
}

// New creates a new GormStore
func New(opts ...Option) (*GormStore, error) {
	dos := defaultOption
	for _, opt := range opts {
		opt(&dos)
	}
	s := &GormStore{}
	s.options = dos

	var err error
	s.db, err = gorm.Open(mysql.Open(s.options.sqlUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Save saves datas to GormStore
func (s *GormStore) Save(datas ...collector.OutputData) error {
	s.lock.Lock()
	for _, data := range datas {
		s.dataDocker = append(s.dataDocker, data)
	}
	s.lock.Unlock()
	if len(s.dataDocker) >= s.BatchCount {
		if _, err := s.Flush(); err != nil {
			return err
		}
	}
	return nil
}

// Flush flushes datas to storage
func (s *GormStore) Flush() (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.logger.Info("flush start")
	defer s.logger.Info("flush end")

	if len(s.dataDocker) == 0 {
		return 0, nil
	}

	datas := make([]map[string]interface{}, 0)
	for i := 0; i < len(s.dataDocker); i++ {
		d := s.dataDocker[i].Struct
		bs, _ := json.Marshal(d)
		m := make(map[string]interface{})
		_ = json.Unmarshal(bs, &m)
		datas = append(datas, m)
	}
	_ = s.db.AutoMigrate(&s.dataDocker[0].Struct)
	table := s.dataDocker[0].Struct.TableName()
	s.dataDocker = nil
	return len(datas), s.db.Table(table).Create(&datas).Error
}
