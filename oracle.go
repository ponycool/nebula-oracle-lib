package oradb

import (
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

var (
	oracleLock        sync.RWMutex
	oracleInitialized bool
	oracle            *sql.DB
)

type Ora struct {
}

// OracleInit 初始化Oracle
func OracleInit(conf *Config, logger *zap.Logger) {
	oracleLock.Lock()
	defer oracleLock.Unlock()

	var err error
	if oracleInitialized {
		err = fmt.Errorf("[db] oracle already initialized")
	}

	connStr := FormatOracleConnUri(conf)

	if len(connStr) == 0 {
		err = fmt.Errorf("[db] oracle connection uri invalid")
	}

	oracle, err = sql.Open("godror", connStr)

	err = oracle.Ping()

	if err != nil {
		defer logger.Error(err.Error())
		panic(err)
	}

	logger.Info(fmt.Sprintf("[db] oracle connection successful"))
}

// FormatOracleConnUri 格式化连接字符串
func FormatOracleConnUri(conf *Config) string {
	port, err := strconv.ParseInt(conf.Port, 10, 32)
	if err != nil || port == 0 {
		port = 1521
	}
	url := fmt.Sprintf("%s/%s@%s:%d/%s?connect_timeout=30",
		conf.User,
		conf.Password,
		conf.Host,
		port,
		conf.SID,
	)
	return url
}
