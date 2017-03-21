package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"path/filepath"
)
var staticMap map[string]string

type Mux struct {
}
func AddstaticMap(webdir, localdir string) {
	staticMap[webdir] = localdir
}
func (mux *Mux)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.URL.Path",r.URL)
	sli := strings.Split(r.URL.Path, "/")
	fmt.Println("sli",sli)
	prefix := "/" + sli[1]                              //find webdir prefix such as "/asset"
	fmt.Println("prefix",prefix)
	if localdir, ok := staticMap[prefix]; ok != false { //assertion to find localdir
		fmt.Printf("\n******the prefix is:%s  the localdir is:%s\n\n", prefix, localdir)
		file := localdir + r.URL.Path[len(prefix):]
		http.ServeFile(w, r, file) //return local resource as static resource
		return
	}

	http.NotFound(w, r)
	return
}

func main() {
	workdir, _ := os.Getwd()
	file, _ := exec.LookPath(os.Args[0])
	ApplicationPath, _ := filepath.Abs(file)
	ApplicationDir, _ := filepath.Split(ApplicationPath)
	confPath := fmt.Sprintf("%shitball", ApplicationDir)
	f, err := os.Open(confPath)
	if err != nil {
		//如果执行文件目录中找不到的话就用工作目录试试
		workDirconfPath := fmt.Sprintf("%s/hitball", workdir)
		f, err = os.Open(workDirconfPath)
		if err != nil {
			panic(err)
		}
	}
	//staticMap = make(map[string]string)
	//AddstaticMap("/static", f.Name()) //add dir to resource in local server
	//mux := &Mux{}
	listener, _ := net.Listen("tcp", "0.0.0.0:6060")
	fmt.Println("web client start :0.0.0.0:6060")
	go func() {
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(f.Name()))))
		http.Serve(listener, nil)
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	fmt.Println("web client Shutting down...")
	listener.Close()
}
