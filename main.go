package main

import (
	"github.com/x893675/certadmin/app"
	"github.com/x893675/certadmin/app/cmd/util"
)

func main() {
	util.CheckErr(app.Run())
}
