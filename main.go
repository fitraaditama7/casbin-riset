package main

import (
	"database/sql"
	"log"
	"runtime"
	"time"

	"github.com/casbin/casbin/v2"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	_ "github.com/go-sql-driver/mysql"
)

func finalizer(db *sql.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	}
}

func main() {
	e := Init()

	// e.AddGroupingPolicy("alice", "admmin")

	// e.AddPolicy("admin", "/data/*", "POST")
	// if err := e.SavePolicy(); err != nil {
	// 	log.Println("SavePolicy failed, err: ", err)
	// }

	has, err := e.Enforce("maintainer", "/api/category/get", "GET")
	if err != nil {
		panic(err)
	}
	if !has {
		log.Println("do not have permission")
	} else {
		log.Println("have permission")
	}
}

func Init() *casbin.Enforcer {
	db, err := sql.Open("mysql", "maman:123459@tcp(127.0.0.1:3306)/mahasiswa")
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 10)

	runtime.SetFinalizer(db, finalizer)

	a, err := sqladapter.NewAdapter(db, "mysql", "casbin_rule")
	if err != nil {
		panic(err)
	}

	e, err := casbin.NewEnforcer("example/config.conf", a)
	if err != nil {
		panic(err)
	}

	if err = e.LoadPolicy(); err != nil {
		log.Println("LoadPolicy failed, err: ", err)
	}
	return e
}
