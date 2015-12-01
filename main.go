package main

import (
	"github.com/toophy/toogo"
)

// 主线程
type MasterThread struct {
	toogo.Thread
}

// 首次运行
func (this *MasterThread) On_first_run() {
}

// 响应线程最先运行
func (this *MasterThread) On_pre_run() {
	// 处理各种最先处理的问题
}

// 响应线程运行
func (this *MasterThread) On_run() {
}

// 响应线程退出
func (this *MasterThread) On_end() {
}

func main() {
	main_thread := new(MasterThread)
	main_thread.Init_thread(main_thread, toogo.Tid_master, "master", 100, 10000)
	toogo.Run(main_thread)
}
