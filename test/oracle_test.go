package main

import (
	"github.com/joho/godotenv"
	"github.com/ponycool/nebula-lib/log"
	lib "github.com/ponycool/nebula-oracle-lib"
	"go.uber.org/zap"
	"os"

	"testing"
)

var logger *zap.Logger

// 初始化日志
func logInit() {
	if logger != nil {
		return
	}
	logger = log.Init(
		log.SetAppName("nebula-oracle-lib"),
		log.SetDevelopment(true),
		log.SetLevel(zap.DebugLevel),
		log.SetMaxSize(2),
		log.SetMaxBackups(100),
		log.SetMaxAge(30),
	)
	logger.Info("logger initial successful")
}

func TestInit(t *testing.T) {
	logInit()
}

/**
 * 尽管编译不需要 Oracle 客户端库，但在运行时需要它们。
 * 从 https://www.oracle.com/database/technologies/instant-client/downloads.html 下载免费的 Basic 或 Basic Light 软件包。
 * 安装说明 https://oracle.github.io/odpi/doc/
 */

func TestConnect(t *testing.T) {
	t.Helper()

	_ = godotenv.Load()

	config := &lib.Config{
		User:     os.Getenv("user"),
		Password: os.Getenv("password"),
		Host:     os.Getenv("host"),
		Port:     os.Getenv("port"),
		SID:      os.Getenv("sid"),
	}
	lib.OracleInit(config, logger)
}

type model struct {
	NAME1 string `json:"name_1"`
	SORTL string `json:"sortl"`
}

func TestQuery(t *testing.T) {
	t.Helper()

	sql := "select NAME1,SORTL from T_YEJI_DETAIL_GNW_KH where NAME1= 'Akram Yousif'"

	// 占位符示例
	// sql := "select NAME1 from T_YEJI_DETAIL_GNW_KH where NAME1= :1"
	// param := []string{
	// 	"泉山区顺豪服装店",
	// }
	ora := new(lib.Ora)
	rows, err := ora.Query(sql)

	var m []model
	err = lib.ScanResult(rows, &m)
	if err != nil {
		t.Error(err.Error())
		return
	}

	log.Get().Info("查询测试结果", zap.Any("m", m))
}

func TestQueryRow(t *testing.T) {
	t.Helper()

	sql := "select NAME1,SORTL from T_YEJI_DETAIL_GNW_KH where NAME1= 'Akram Yousif'"
	ora := new(lib.Ora)
	row, err := ora.QueryRow(sql)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	var m model
	err = row.Scan(&m.NAME1, &m.SORTL)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	log.Get().Info("查询一行测试结果", zap.Any("QueryRow", m))
}

func TestExec(t *testing.T) {
	t.Helper()

	sql := "update T_YEJI_DETAIL_GNW_KH set NAME1='test' where NAME1='test'"
	ora := new(lib.Ora)
	result, err := ora.Exec(sql)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	log.Get().Info("SQL执行测试", zap.Any("受影响的行数", affected))
}
