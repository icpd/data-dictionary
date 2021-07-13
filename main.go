package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	MarkdownTableHeader = "\n\n\n### %s \n| 序号 | 字段名称 | 数据类型 | 是否为空 | 字段说明 |\n| :--: |----| ---- | ---- | ---- |\n"
	MarkdownTableRow    = "| %d | %s | %s | %s | %s |\n"
	SelectTableSql      = "SELECT TABLE_NAME, TABLE_COMMENT FROM information_schema.TABLES WHERE table_type='BASE TABLE' AND TABLE_SCHEMA = ?"
	SelectColumnSql     = "SELECT ORDINAL_POSITION, COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA= ? AND TABLE_NAME= ?"
)

var (
	DSN    string
	Schema string
)

type Table struct {
	TableName    string `gorm:"column:TABLE_NAME"`
	TableComment string `gorm:"column:TABLE_COMMENT"`
}

type Column struct {
	OrdinalPosition int    `gorm:"column:ORDINAL_POSITION"`
	ColumnName      string `gorm:"column:COLUMN_NAME"`
	ColumnType      string `gorm:"column:COLUMN_TYPE"`
	IsNullable      string `gorm:"column:IS_NULLABLE"`
	ColumnComment   string `gorm:"column:COLUMN_COMMENT"`
}

func NewMysql() *gorm.DB {
	db, err := gorm.Open(
		mysql.Open(DSN),
		&gorm.Config{
			PrepareStmt: true,
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: 500 * time.Millisecond,
					Colorful:      true,
					LogLevel:      logger.Info,
				})},
	)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func init() {
	flag.StringVar(&DSN, "d", "", "mysql information_schema dsn , eg. username:password@tcp(127.0.0.1:3306)/information_schema")
	flag.StringVar(&Schema, "s", "", "schema，数据库名")
	flag.Parse()
	flag.Usage()
}

func main() {
	var tables []Table
	db := NewMysql()
	err := db.Raw(SelectTableSql, Schema).
		Scan(&tables).Error
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	for _, t := range tables {
		header := t.TableName
		if t.TableComment != "" {
			header = fmt.Sprintf("%s (%s)", t.TableName, t.TableComment)
		}
		buf.WriteString(fmt.Sprintf(MarkdownTableHeader, header))

		var columns []Column
		err := db.Raw(SelectColumnSql, Schema, t.TableName).
			Scan(&columns).Error
		if err != nil {
			log.Fatal(err)
		}

		for _, c := range columns {
			buf.WriteString(fmt.Sprintf(MarkdownTableRow, c.OrdinalPosition, c.ColumnName,
				c.ColumnType, c.IsNullable, c.ColumnComment))
		}
	}

	f, err := os.Create("db.md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}
