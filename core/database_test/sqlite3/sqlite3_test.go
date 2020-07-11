package sqlite3

import (
	"fmt"
	"testing"
)

func TestSqlite3(t *testing.T) {
	c, err := connectDB("sqlite3", "user_info.db")
	if err != "" {
		panic(err)
	}
	fmt.Println("connect sqlite3 success...")
	c.Create()
	fmt.Println("add action done!")

	c.Read()
	fmt.Println("get action done!")

	c.Update()
	fmt.Println("update action done!")

	c.Delete()
	fmt.Println("delete action done!")
}
