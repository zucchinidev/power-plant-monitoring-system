package store

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Cfg struct {
	User   string
	Passwd string
	Addr   string
	DBName string
}

type Storage struct {
	DB  *sql.DB
	cfg Cfg
}

func NewSqlManager(c Cfg) *Storage {
	return &Storage{cfg: c}
}

func (s *Storage) Connect() error {
	var err error
	// s.DB is a pool of connections
	if s.DB, err = sql.Open("mysql", s.config().FormatDSN()); err != nil {
		return err
	}
	if err = s.DB.Ping(); err != nil {
		_ = s.DB.Close()
		return err
	}
	return nil
}

func (s *Storage) config() *mysql.Config {
	return &mysql.Config{
		User:                 s.cfg.User,
		Net:                  "tcp",
		Addr:                 s.cfg.Addr,
		DBName:               s.cfg.DBName,
		Passwd:               s.cfg.Passwd,
		AllowNativePasswords: true,
		ParseTime:            true,
	}
}

// Calling Connect does not block for server discovery.
// If you wish to know if a server has been found and connected to, use the Ping method
func (s *Storage) Ping() error {
	err := s.DB.Ping()
	if err != nil {
		return fmt.Errorf("sql service error: %s", err)
	}
	return nil
}

// Close calling close method you will close the sql pool connections
func (s *Storage) Close() error {
	return s.DB.Close()
}
