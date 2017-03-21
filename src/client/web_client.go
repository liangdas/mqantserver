package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
)

func main() {
	workdir, _ := os.Getwd()
	file, _ := exec.LookPath(os.Args[0])
	ApplicationPath, _ := filepath.Abs(file)
	ApplicationDir, _ := filepath.Split(ApplicationPath)
	confPath := fmt.Sprintf("%spublic", ApplicationDir)
	f, err := os.Open(confPath)
	if err != nil {
		//如果执行文件目录中找不到的话就用工作目录试试
		workDirconfPath := fmt.Sprintf("%s/public", workdir)
		f, err = os.Open(workDirconfPath)
		if err != nil {
			panic(err)
		}
	}

	hitballconfPath := fmt.Sprintf("%shitball", ApplicationDir)
	hitballf, err := os.Open(hitballconfPath)
	if err != nil {
		//如果执行文件目录中找不到的话就用工作目录试试
		workDirconfPath := fmt.Sprintf("%s/hitball", workdir)
		hitballf, err = os.Open(workDirconfPath)
		if err != nil {
			panic(err)
		}
	}

	listener, _ := net.Listen("tcp", "0.0.0.0:8080")
	fmt.Println("web client start :0.0.0.0:8080")
	go func() {
		http.Handle("/mqant/", http.StripPrefix("/mqant/", http.FileServer(http.Dir(f.Name()))))
		http.Handle("/hitball/", http.StripPrefix("/hitball/", http.FileServer(http.Dir(hitballf.Name()))))
		http.Serve(listener, nil)
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	fmt.Println("web client Shutting down...")
	listener.Close()
}
