package main

import "fmt"
import (
	"flag"
	"lm"
	"os"
	"ImServer/service"
	"crypto/tls"
)

var host = flag.String("host", "", "please input -host")
var port = flag.String("port", "3333", "please input -port")
var tcpMap = make(map[string]chan *lm.Data_Message) ///////创建chan map

func main() {

	flag.Parse()
	//////tsl1.2 创建

	cer, err := tls.LoadX509KeyPair("/Users/shanksYao/Documents/work/im/src/ImServer/cert.pem","/Users/shanksYao/Documents/work/im/src/ImServer/key.pem")
	if err!=nil{
		println("-------",err.Error())
		return
	}
	config:=&tls.Config{Certificates:[]tls.Certificate{cer}}
	/// tsl 建立完毕
	listener, err := tls.Listen("tcp", *host+":"+*port,config)

	if err != nil {
		fmt.Println("error listening :", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("listening on " + *host + "port :" + *port)

	for {
		fmt.Println("-----start listen")
		conn, _ := listener.Accept()
		fmt.Println("=======handle")
		//////并发支持
		userManager := new(service.UserManager)

		userManager.Conn=conn
		go userManager.Login(tcpMap)

		fmt.Println(">>>>>>>>end")
	}
}
