package sqlstorage

import (
	"encoding/json"

	"github.com/hedon954/go-crawler/collector"
	"github.com/hedon954/go-crawler/engine"
	"github.com/hedon954/go-crawler/sqldb"
	"go.uber.org/zap"
)

type SqlStore struct {
	dataDocker  []collector.OutputData
	columnNames []sqldb.Field
	db          sqldb.DBer
	Table       map[string]struct{}
	options
}

// New creates a new SqlStore
func New(opts ...Option) (*SqlStore, error) {
	dos := defaultOption
	for _, opt := range opts {
		opt(&dos)
	}
	s := &SqlStore{}
	s.options = dos
	s.Table = make(map[string]struct{})
	var err error
	s.db, err = sqldb.New(
		sqldb.WithSqlUrl(s.sqlUrl),
		sqldb.WithLogger(s.logger),
	)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Save saves datas to SqlStore
func (s *SqlStore) Save(datas ...collector.OutputData) error {
	for _, data := range datas {
		tn := data.TaskName
		if _, ok := s.Table[tn]; !ok {
			columnNames := getFields(data)
			if err := s.db.CreateTable(sqldb.TableData{
				TableName:   tn,
				ColumnNames: columnNames,
				AutoKey:     true,
			}); err != nil {
				s.logger.Error("create table failed", zap.Error(err))
				return err
			}
		}
		if len(s.dataDocker) >= s.BatchCount {
			_ = s.Flush()
		}
		s.dataDocker = append(s.dataDocker, data)
	}
	return nil
}

// getFields parses fields according to the data struct
func getFields(data collector.OutputData) []sqldb.Field {
	taskName := data.TaskName
	ruleName := data.RuleName
	fields := engine.GetFields(taskName, ruleName)

	var columnNames []sqldb.Field
	for _, field := range fields {
		columnNames = append(columnNames, sqldb.Field{
			Title: field,
			Type:  "MEDIUMTEXT",
		})
	}

	columnNames = append(columnNames,
		sqldb.Field{Title: "Url", Type: "VARCHAR(255)"},
		sqldb.Field{Title: "Time", Type: "VARCHAR(255)"},
	)

	return columnNames
}

// Flush flushes datas to storage
func (s *SqlStore) Flush() error {
	s.logger.Info("flush start")
	defer s.logger.Info("flush end")

	if len(s.dataDocker) == 0 {
		return nil
	}

	args := make([]interface{}, 0)
	for _, data := range s.dataDocker {
		fields := engine.GetFields(data.TaskName, data.RuleName)
		d := data.Data.(map[string]interface{})
		var value []string
		for _, field := range fields {
			v := d[field]
			switch v.(type) {
			case nil:
				value = append(value, "")
			case string:
				value = append(value, v.(string))
			default:
				j, _ := json.Marshal(v)
				value = append(value, string(j))
			}
		}
		value = append(value, data.Url, data.Time)
		for _, v := range value {
			args = append(args, v)
		}
	}

	return s.db.Insert(sqldb.TableData{
		TableName:   s.dataDocker[0].TaskName,
		ColumnNames: getFields(s.dataDocker[0]),
		Args:        args,
		DataCount:   len(s.dataDocker),
	})
}
