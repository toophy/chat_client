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
		name_fix = fmt.Sprintf("Conn[%d]", m.SessionId)
	}

	switch m.Msg {
	case "listen failed":
		this.LogFatal("%s : Listen failed[%s]", name_fix, m.Info)

	case "listen ok":
		this.LogInfo("%s : Listen(%s) ok.", name_fix, toogo.GetSessionById(m.SessionId).GetIPAddress())

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
		msgLogin := new(proto.C2G_login)
		msgLogin.Account = "liusl"
		msgLogin.Time = 123
		msgLogin.Sign = "wokao"
		msgLogin.Write(p)
		this.LogInfo("send C2G_login")

		p.PacketWriteOver()
		session := toogo.GetSessionById(m.SessionId)
		x := new(toogo.Tmsg_packet)
		x.Data = p.GetData()
		x.Len = uint32(p.GetPos())
		x.Count = uint16(p.Count)

		toogo.PostThreadMsg(session.MailId, x)

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
func (this *MasterThread) On_packetError(sessionId uint32) {
	toogo.CloseSession(this.Get_thread_id(), sessionId)
}

// 注册消息
func (this *MasterThread) On_registNetMsg() {
	this.RegistNetMsg(proto.G2C_login_ret_Id, this.on_g2c_login_ret)
}

func (this *MasterThread) on_g2c_login_ret(pack *toogo.PacketReader, sessionId uint32) bool {
	msg := new(proto.G2C_login_ret)
	msg.Read(pack)

	this.LogInfo("on_g2c_login_ret")

	p := new(toogo.PacketWriter)
	d := make([]byte, 64)
	p.InitWriter(d)
	msgSend := new(proto.C2S_chat)
	msgSend.Channel = 1
	msgSend.Data = "你好,世界!"
	msgSend.Write(p)

	p.PacketWriteOver()
	session := toogo.GetSessionById(sessionId)
	x := new(toogo.Tmsg_packet)
	x.Data = p.GetData()
	x.Len = uint32(p.GetPos())
	x.Count = uint16(p.Count)

	toogo.PostThreadMsg(session.MailId, x)

	return true
}

func main() {
	main_thread := new(MasterThread)
	main_thread.Init_thread(main_thread, toogo.Tid_master, "master", 1000, 100, 10000)
	toogo.Run(main_thread)
}
