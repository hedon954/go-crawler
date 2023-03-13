package sqldb

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

type DBer interface {
	CreateTable(t TableData) error
	Insert(t TableData) error
}

type Field struct {
	Title string
	Type  string
}

type TableData struct {
	TableName   string
	ColumnNames []Field
	Args        []interface{} // data
	DataCount   int
	AutoKey     bool
}

type Sqldb struct {
	options
	db *sql.DB
}

func New(opts ...Option) (*Sqldb, error) {
	dos := defaultOption
	for _, opt := range opts {
		opt(&dos)
	}
	d := &Sqldb{}
	d.options = dos
	if err := d.OpenDB(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Sqldb) OpenDB() error {
	db, err := sql.Open("mysql", d.sqlUrl)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(2048)
	db.SetMaxIdleConns(2048)
	if err = db.Ping(); err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *Sqldb) CreateTable(t TableData) error {
	if len(t.ColumnNames) == 0 {
		return errors.New("Column can not be empty")
	}
	sql := `CREATE TABLE IF NOT EXISTS ` + t.TableName + " ("
	if t.AutoKey {
		sql += `id INT(12) NOT NULL PRIMARY KEY AUTO_INCREMENT,`
	}
	for _, t := range t.ColumnNames {
		sql += t.Title + ` ` + t.Type + `,`
	}
	sql = sql[:len(sql)-1] + `) ENGINE=MyISAM DEFAULT CHARSET=utf8;`

	d.logger.Debug("crate table", zap.String("sql", sql))

	_, err := d.db.Exec(sql)
	return err
}

func (d *Sqldb) Insert(t TableData) error {
	if len(t.ColumnNames) == 0 {
		return errors.New("empty column")
	}
	sql := `INSERT INTO ` + t.TableName + `(`

	for _, v := range t.ColumnNames {
		sql += v.Title + ","
	}

	sql = sql[:len(sql)-1] + `) VALUES `

	blank := ",(" + strings.Repeat(",?", len(t.ColumnNames))[1:] + ")"
	sql += strings.Repeat(blank, t.DataCount)[1:] + `;`
	d.logger.Debug("insert table", zap.String("sql", sql))
	_, err := d.db.Exec(sql, t.Args...)
	return err
}
