package main

/*
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

// 日志级别保持与 Go 一致
enum {
	LOG_TRACE, LOG_DEBUG, LOG_INFO,
	LOG_WARN,  LOG_ERROR, LOG_FATAL,
};

// typedef void (*mihomo_log_cb)(int level, const char* msg);
static inline void call_log_cb(void* cb, int lvl, const char* msg) {
	if (cb != NULL) ((void(*)(int,const char*))cb)(lvl, msg);
}
*/
import "C"
import (
	"unsafe"

	"github.com/metacubex/mihomo/log"
	"github.com/loafman-kangjun/libclash/libcore"   // ← 改成实际 import 路径
)

// ---------------- 公共工具 ----------------

//export free_go_str
func free_go_str(p *C.char) { C.free(unsafe.Pointer(p)) }

func goErr2C(err error) *C.char {
	if err == nil {
		return nil
	}
	return C.CString(err.Error())
}

// ---------------- 日志回调 ----------------

var cLogCb unsafe.Pointer

//export Mihomo_SetLogCallback
func Mihomo_SetLogCallback(cb unsafe.Pointer) {
	cLogCb = cb
}

// 把 libcore 的日志转发给 C
func init() {
	sub := log.Subscribe()
	go func() {
		for ev := range sub {
			if cLogCb != nil {
				cMsg := C.CString(ev.Payload)
				C.call_log_cb(cLogCb, C.int(ev.LogLevel), cMsg)
				free_go_str(cMsg)
			}
		}
	}()
}

// ---------------- 核心导出函数 ----------------

//export Mihomo_Run
func Mihomo_Run(cPath, cB64 *C.char, geo C.bool) *C.char {
	err := libcore.Run(C.GoString(cPath), C.GoString(cB64), bool(geo))
	return goErr2C(err)
}

//export Mihomo_Reload
func Mihomo_Reload(cPath, cB64 *C.char) *C.char {
	err := libcore.Reload(C.GoString(cPath), C.GoString(cB64))
	return goErr2C(err)
}

//export Mihomo_Stop
func Mihomo_Stop() {
	libcore.Stop()
}

// 空的 main 避免 “no non-test Go files” 报错
func main() {}