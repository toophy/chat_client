// easybuff
// 不要修改本文件, 每次消息有变动, 请手动生成本文件
// easybuff -s 描述文件目录 -o 目标文件目录 -l 语言(go,cpp)

package proto

import (
	. "github.com/toophy/toogo"
)

// 登录聊天服务器
const C2M_login_Id = 1

type C2M_login struct {
	Account string // 帐号
	Time    int32  // 登录时间戳
	Sign    string // 验证码
}

func (t *C2M_login) Read(s *Stream) {
	t.Account = s.ReadString()
	t.Time = s.ReadInt32()
	t.Sign = s.ReadString()
}

func (t *C2M_login) Write(s *Stream) {
	s.WriteString(&t.Account)
	s.WriteInt32(t.Time)
	s.WriteString(&t.Sign)
}

// 服务器响应登录
const M2C_login_ret_Id = 2

type M2C_login_ret struct {
	Ret int8   // 登录结果,0:成功,其他为失败原因
	Msg string // 登录失败描述
}

func (t *M2C_login_ret) Read(s *Stream) {
	t.Ret = s.ReadInt8()
	t.Msg = s.ReadString()
}

func (t *M2C_login_ret) Write(s *Stream) {
	s.WriteInt8(t.Ret)
	s.WriteString(&t.Msg)
}
