package main

import (
	"fmt"
	"github.com/toophy/chat_client/proto"
	"github.com/toophy/toogo"
)

// 主线程
type MasterThread struct {
	toogo.Thread
}

// 首次运行
func (this *MasterThread) On_firstRun() {
}

// 响应线程最先运行
func (this *MasterThread) On_preRun() {
	// 处理各种最先处理的问题
}

// 响应线程运行
func (this *MasterThread) On_run() {
}

// 响应线程退出
func (this *MasterThread) On_end() {
}

// 响应网络事件
func (this *MasterThread) On_netEvent(m *toogo.Tmsg_net) bool {

	name_fix := m.Name
	if len(name_fix) == 0 {
		name_fix = fmt.Sprintf("Conn[%d]", m.Id)
	}

	switch m.Msg {
	case "listen failed":
		this.LogFatal("%s : Listen failed[%s]", name_fix, m.Info)

	case "listen ok":
		this.LogInfo("%s : Listen(0.0.0.0:%d) ok.", name_fix, 8001)

	case "accept failed":
		this.LogFatal(m.Info)
		return false

	case "accept ok":
		this.LogDebug("%s : Accept ok", name_fix)

	case "connect failed":
		this.LogError("%s : Connect failed[%s]", name_fix, m.Info)

	case "connect ok":
		this.LogDebug("%s : Connect ok", name_fix)
		p := new(toogo.PacketWriter)
		d := make([]byte, 64)
		p.InitWriter(d)
		msgLogin := new(proto.C2M_login)
		msgLogin.Account = "liusl"
		msgLogin.Time = 123
		msgLogin.Sign = "wokao"
		msgLogin.Write(p)

		p.PacketWriteOver()
		session := toogo.GetSessionById(m.Id)
		m := new(toogo.Tmsg_packet)
		m.Data = p.GetData()
		m.Len = uint32(p.GetPos())
		m.Count = uint32(p.Count)

		toogo.PostThreadMsg(session.MailId, m)

		////////
		p2 := new(toogo.PacketWriter)
		d2 := make([]byte, 64)
		p2.InitWriter(d2)
		msgLogin2 := new(proto.C2M_login)
		msgLogin2.Account = "wangyh"
		msgLogin2.Time = 456
		msgLogin2.Sign = "yeye"
		msgLogin2.Write(p2)

		p2.PacketWriteOver()
		m2 := new(toogo.Tmsg_packet)
		m2.Data = p2.GetData()
		m2.Len = uint32(p2.GetPos())
		m2.Count = uint32(p2.Count)

		toogo.PostThreadMsg(session.MailId, m2)
	case "read failed":
		this.LogError("%s : Connect read[%s]", name_fix, m.Info)

	case "pre close":
		this.LogDebug("%s : Connect pre close", name_fix)

	case "close failed":
		this.LogError("%s : Connect close failed[%s]", name_fix, m.Info)

	case "close ok":
		this.LogDebug("%s : Connect close ok.", name_fix)
	}

	return true
}

// -- 当网络消息包解析出现问题, 如何处理?
func (this *MasterThread) On_packetError(m *toogo.Tmsg_packet) {
	toogo.CloseSession(this.Get_thread_id(), m.SessionId)
}

// 注册消息
func (this *MasterThread) On_registNetMsg() {
	this.RegistNetMsg(proto.M2C_login_ret_Id, this.on_m2c_login_ret)
}

func (this *MasterThread) on_m2c_login_ret(pack *toogo.PacketReader, sessionId uint32) bool {
	msg := new(proto.M2C_login_ret)
	msg.Read(pack)
	fmt.Println(msg)
	return true
}

func main() {
	main_thread := new(MasterThread)
	main_thread.Init_thread(main_thread, toogo.Tid_master, "master", 1000, 100, 10000)
	toogo.Run(main_thread)
}
