package service

import (
	"net"
	"fmt"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"lm"
)

type UserManager struct {
	Conn net.Conn
}



func (co *UserManager)Login(userMap map[string]chan *lm.Data_Message) (user *lm.Data_User) {
	/////前4个字节标识  数据的长度
	size := make([]byte, 4)
	index, err := co.Conn.Read(size)
	if err!=nil {
		fmt.Println("read error",err.Error())
		return
	}
	if index>0 {
		///////
		println("read login info")
		size_ := int32(binary.BigEndian.Uint32(size))
		println("size_",size_)
		data := make([]byte, size_)
		co.Conn.Read(data)
		dataObj := new(lm.Data)
		proto.Unmarshal(data, dataObj)
		if *dataObj.Type ==int32(lm.Data_USER) {
			user := dataObj.GetUser();
			/////声明当前登录用户的消息
			userMap[user.GetUser()] = make(chan *lm.Data_Message)
			//////TODO 登陆成功返回信息
			co.ResultOK()

			/////开启输出
			go co.Output(userMap, user)
			/////开启输入
			go co.Input(userMap, user)
			return user ////返回用户信息

		}else {
			fmt.Println("---- login error ---")
		}
	}

	return
}
func (co *UserManager)ResultOK() {
	result := new(lm.Data_Result)
	result.Msg=new(string)
	*result.Msg="ok"
	result.Code=new(int32)
	*result.Code=1

	data := new(lm.Data)
	data.Type=new(int32)
	*data.Type=2
	data.Time=new(uint64)
	*data.Time=11111111//TODO
	data.Version=new(int32)
	*data.Version=1
	data.Res=result
	data_, _ := proto.Marshal(data)

	co.Conn.Write(Size(data_))
	co.Conn.Write(data_)
}





func (co *UserManager)Output(userMap map[string]chan *lm.Data_Message, user *lm.Data_User) {
	defer co.Conn.Close()
	defer delete(userMap,user.GetUser())//// del
	for {
		msg := <-userMap[user.GetUser()]
		println(user.GetUser(),"======= output")
		outData:=new(lm.Data)
		outData.Version=new(int32)
		outData.Type=new(int32)
		outData.Time=new(uint64)
		*outData.Version=1
		*outData.Type=int32(lm.Data_MESSAGE)
		*outData.Time=111111
		outData.Msg=msg
		///////
		data, _ := proto.Marshal(outData)
		dataSize := make([]byte, 4)
		binary.BigEndian.PutUint32(dataSize, uint32(len(data)))////TODO data的长度要做限制,要不然会出现计数不准的情况
		co.Conn.Write(dataSize)
		co.Conn.Write(data)
	}
}

func Size(data []byte) (size []byte) {
	dataSize := make([]byte, 4)
	binary.BigEndian.PutUint32(dataSize, uint32(len(data)))////TODO data的长度要做限制,要不然会出现计数不准的情况
	return dataSize
}



func (co *UserManager)Input(userMap map[string]chan *lm.Data_Message, user *lm.Data_User) {
	defer co.Conn.Close()
	defer delete(userMap,user.GetUser())//// del
	for {
		size := make([]byte, 4)
		index, err := co.Conn.Read(size)
		if err!=nil {
			fmt.Println(err.Error())
			return
		}
		if index>0 {
			size_ := int32(binary.BigEndian.Uint32(size))
			data := make([]byte, size_)
			dataObj := new(lm.Data)
			co.Conn.Read(data)
			proto.Unmarshal(data, dataObj)
			println("========= receive start %d",*dataObj.Type)
			if *dataObj.Type == int32(lm.Data_MESSAGE) {
				println("========= receive msg")
				msg := dataObj.GetMsg()
				msg.Send=new(string)
				*msg.Send=user.GetUser()//TODO
				println(*msg.Send+"========= receive send")
				accept:=userMap[msg.GetAccept()]
				if accept!=nil{
					println("========= receive ok",msg.GetAccept())
					userMap[msg.GetAccept()] <- msg
				}
				println(*msg.Send+"========= receive ok")
				co.ResultOK()////接收信息OK
			}
		}
	}
}