package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/eaciit/toolkit"

	_ "github.com/go-sql-driver/mysql"

	"github.com/eaciit/sqlh"
)

const (
	sqlconn = "root:Database.1024@tcp(localhost:3306)/ectestdb"
)

var (
	db  *sql.DB
	err error
)

type Employee struct {
	ID       string
	Name     string
	Level    int
	JoinDate time.Time
}

func TestConnect(t *testing.T) {
	db, err = sqlh.Connect("mysql", sqlconn)
	if err != nil {
		fmt.Println("error connecting database", err)
		os.Exit(100)
	}
}

func TestCreateTable(t *testing.T) {
	sql := "CREATE TABLE test_table_model2(\n" +
		"	id      	VARCHAR(32)     NOT NULL,\n" +
		"	name  		VARCHAR(200)    NOT NULL,\n" +
		"	level  		INT    			NOT NULL,\n" +
		"	datejoin   	DATE   			NOT NULL,\n" +
		"	PRIMARY KEY (id)\n" +
		");\n"

	qr := sqlh.Exec(db, sqlh.ExecNonScalar, sql)
	if qr.Error() != nil {
		t.Error(qr.Error())
	}
}

func TestInsert(t *testing.T) {
	t.Skip()
	sql := "insert into test_table_model2 (id, name, level, datejoin) values(?,?,?,?)"
	id := toolkit.RandomString(32)
	name := "Name " + id
	qr := sqlh.Exec(db, sqlh.ExecNonScalar, sql, id, name, 100, "2017-01-01")
	if qr.Error() != nil {
		t.Error(qr.Error())
	} else {
		affected, _ := qr.CUDAResult().RowsAffected()
		fmt.Println("Data inserted: ", affected)
	}
}

func TestSelect(t *testing.T) {
	sql := "select * from test_table_model2 order by id desc limit 2"
	es := []Employee{}

	qr := sqlh.Exec(db, sqlh.ExecQuery, sql)
	if qr.Error() != nil {
		t.Error(qr.Error())
	}

	err = qr.Fetch(&es, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Returned record:%d\n%s", len(es), toolkit.JsonStringIndent(es, "\n"))
}

func TestSelectM(t *testing.T) {
	sql := "select * from test_table_model2 where id like 'test%' order by id desc limit 2"
	es := []toolkit.M{}

	qr := sqlh.Exec(db, sqlh.ExecQuery, sql)
	if qr.Error() != nil {
		t.Error(qr.Error())
	}

	err = qr.Fetch(&es, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Returned record:%d\n%s", len(es), toolkit.JsonStringIndent(es, "\n"))
}

func TestClose(t *testing.T) {
	if db != nil {
		db.Close()
	}
}
