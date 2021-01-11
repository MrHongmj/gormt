package gtools

import (
	"flag"
	"os/exec"

	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/gormt/data/view/model"
	"github.com/xxjwxc/gormt/data/config"
	"github.com/xxjwxc/gormt/data/view/model/genmysql"
	"github.com/xxjwxc/public/tools"
)

var mysqlInfo config.MysqlDbInfo
var outDir string
var singularTable bool
var foreignKey bool
var funcKey bool
var ui bool
var urlTag string


func init() {
	var project string
	var table string
	//rootCmd.PersistentFlags().StringVarP(&project, "project", "P", "gadget", "请输入项目名称")
	//rootCmd.MarkFlagRequired("project")
	flag.StringVar(&project, "P", "", "please input project")
	flag.StringVar(&table, "T", "", "please input table")
	flag.Parse()
	if project == "" {
		panic("please input project 。eg：gravidity、pregnancy、ybb 、gadget")
	}

	sqlConf := map[string]config.MysqlDbInfo{
		
	}
	var ok bool
	mysqlInfo,ok = sqlConf[project]
	if !ok {
		panic("error project")
	}
	outDir = "./orm/"+ project
	singularTable = true
	foreignKey = false
	funcKey = false
	ui = false
	urlTag = ""
	//outFileName = ""
	config.SetMysqlDbInfo(&mysqlInfo)
	config.SetOutDir(outDir)
	config.SetTableName(table)
}

func Execute() {

	modeldb := genmysql.GetMysqlModel()
	pkg := modeldb.GenModel()

	list, _ := model.Generate(pkg)

	for _, v := range list {
		path := config.GetOutDir() + "/" + v.FileName
		tools.WriteFile(path, []string{v.FileCtx}, true)

		mylog.Info("formatting differs from goimport's:")
		cmd, _ := exec.Command("goimports", "-l", "-w", path).Output()
		mylog.Info(string(cmd))

		mylog.Info("formatting differs from gofmt's:")
		cmd, _ = exec.Command("gofmt", "-l", "-w", path).Output()
		mylog.Info(string(cmd))
	}
}
