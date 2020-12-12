package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/omecodes/nc2i/nc2i"
	_ "github.com/omecodes/nc2i/statik"
)

func main() {
	if err := nc2i.Cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
