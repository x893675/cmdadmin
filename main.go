package main

import (
	"github.com/x893675/cmdadmin/app"
	"github.com/x893675/cmdadmin/app/cmd/util"
)

func main() {
	util.CheckErr(app.Run())
}
