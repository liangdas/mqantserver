package main

import (
	"github.com/astaxie/beego"
	"os"
	"fmt"
	"path/filepath"
	"os/exec"
)

type MainController struct {
	beego.Controller
}


func (this *MainController) Get() {
	this.Ctx.WriteString("hello world")
}



func main() {
	beego.Router("/", &MainController{})
	workdir,_:=os.Getwd()
	file, _ := exec.LookPath(os.Args[0])
	ApplicationPath, _ := filepath.Abs(file)
	ApplicationDir, _ := filepath.Split(ApplicationPath)
	confPath:= fmt.Sprintf("%s/public",ApplicationDir)
	f, err := os.Open(confPath)
	if err!=nil{
		//如果执行文件目录中找不到的话就用工作目录试试
		workDirconfPath:= fmt.Sprintf("%s/public",workdir)
		f, err = os.Open(workDirconfPath)
		if err!=nil{
			panic(err)
		}
	}
	beego.SetStaticPath("/mqant",f.Name())
	beego.Run("0.0.0.0:8080")
}
