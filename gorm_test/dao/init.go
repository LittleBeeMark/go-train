package dao

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	// postgres 驱动
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strings"
	"time"
)

var DebugHook bool

// IsNotFindErr 数据不存在错误
func IsNotFindErr(err error) bool {
	return err != nil && (err == gorm.ErrRecordNotFound || err == sql.ErrNoRows)
}

// IsDuplicateKeyErr 主键冲突
func IsDuplicateKeyErr(err error) bool {
	return strings.HasPrefix(err.Error(), "pq: duplicate key value violates unique constraint")
}

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

const (
	uniqueShortIDFuncSQL = `CREATE EXTENSION IF NOT EXISTS "pgcrypto";
	CREATE OR REPLACE FUNCTION unique_short_id()
	RETURNS TRIGGER AS $$
	DECLARE
	  key TEXT;
	  qry TEXT;
	  found TEXT;
	BEGIN
	  qry := 'SELECT id FROM ' || quote_ident(TG_TABLE_NAME) || ' WHERE id=';
	  LOOP
		key := encode(gen_random_bytes(6), 'base64');
		key := replace(key, '/', 'x'); -- url safe replacement
		key := replace(key, '+', 'v'); -- url safe replacement
		EXECUTE qry || quote_literal(key) INTO found;
		IF found IS NULL THEN
		  EXIT;
		END IF;
	  END LOOP;
	  NEW.id = key;
	  RETURN NEW;
	END;
	$$ language 'plpgsql';`

	CreateTimeDecs = "create_at desc"
)

func genShortIDTriggerSQL(tableName string) string {
	return fmt.Sprintf(`DO $do$ BEGIN IF EXISTS (SELECT 1 FROM pg_trigger WHERE  NOT tgisinternal AND tgname = 'trigger_%s_shortid') THEN ELSE
	CREATE TRIGGER trigger_%s_shortid BEFORE INSERT ON %s FOR EACH ROW EXECUTE PROCEDURE unique_short_id();
	END IF; END $do$`, tableName, tableName, tableName)
}

// InitDBConnection 初始数据库连接
func InitDBConnection() {
	var err error
	for i := 3; i > 0; i-- {
		db, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=mark password=123 dbname=test sslmode=disable")
		if err != nil {
			logrus.Warningf("DB init failed %v", err.Error())
			logrus.Infof("Try to reconnected after 5s[%d]", 4-i)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		logrus.Panicf("DB init failed %v", err.Error())
	}

	if DebugHook {
		db = db.Debug()
	}

	if err != nil {
		logrus.Panicf("DB init failed %v", err.Error())
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(80)
}

// Init 初始化数据库
func Init() {
	InitDBConnection()

	// 添加ShortID 生成函数
	db.Exec(uniqueShortIDFuncSQL)
	db.Exec(genShortIDTriggerSQL("test_jsons"))
	db.Exec(genShortIDTriggerSQL("test_jsons2"))

}
