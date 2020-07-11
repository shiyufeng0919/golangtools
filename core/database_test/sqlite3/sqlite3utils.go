package sqlite3

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

/*
   示例：sqlite3 数据库
        SQLite，是一款轻型的数据库，是遵守ACID的关系型数据库管理系统。它的设计目标是嵌入式的。
        主要的通信协议是在编程语言内的直接API调用
*/
// People - database fields
type People struct {
	id   int
	name string
	age  int
}

type appContext struct {
	db *sql.DB
}

//连接sqlite3
func connectDB(driverName string, dbName string) (*appContext, string) {
	db, err := sql.Open(driverName, dbName)
	if err != nil {
		return nil, err.Error()
	}
	if err = db.Ping(); err != nil {
		return nil, err.Error()
	}
	sql_table := `
		CREATE TABLE IF NOT EXISTS "user_info" (
		   "user_name" VARCHAR(64) NULL,
		   "user_age" VARCHAR(64) NULL
		);
   `
	db.Exec(sql_table)
	return &appContext{db}, ""
}

// Create
func (c *appContext) Create() {
	stmt, err := c.db.Prepare("INSERT INTO user_info(user_name,user_age) values(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	result, err := stmt.Exec("Jack", 1)
	if err != nil {
		fmt.Printf("add error: %v", err)
		return
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted id is ", lastID)
}

// Read
func (c *appContext) Read() {
	rows, err := c.db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := new(People)
		err := rows.Scan(&p.id, &p.name, &p.age)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(p.id, p.name, p.age)
	}
}

// UPDATE
func (c *appContext) Update() {
	stmt, err := c.db.Prepare("UPDATE users SET age = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	result, err := stmt.Exec(10, 1)
	if err != nil {
		log.Fatal(err)
	}
	affectNum, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("update affect rows is ", affectNum)
}

// DELETE
func (c *appContext) Delete() {
	stmt, err := c.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	result, err := stmt.Exec(1)
	if err != nil {
		log.Fatal(err)
	}
	affectNum, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("delete affect rows is ", affectNum)
}
