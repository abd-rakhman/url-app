package gen

import (
	"database/sql"
	"fmt"
	"math/rand"
)

func generate(cnt int) string {
	symbols := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"
	code := ""
	for i := 0; i < (cnt/100 + 10); i++ {
		code += string(rune(symbols[rand.Intn(62)]))
	}
	return code
}

func unique(s string, Db *sql.DB) bool {
	getQuery := fmt.Sprintf("SELECT id FROM data WHERE url='%s';", s)
	rows, _ := Db.Query(getQuery)

	if rows.Next() {
		return false
	} else {
		return true
	}
}

func GetCode(db *sql.DB) string {
	x := generate(0)
	cnt := 0
	for !unique(x, db) {
		x = generate(cnt)
		cnt++
	}
	return x
}
