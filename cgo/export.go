package main

/*
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

enum {
	LOG_TRACE, LOG_DEBUG, LOG_INFO,
	LOG_WARN,  LOG_ERROR, LOG_FATAL,
};

typedef void (*mihomo_log_cb)(int level, const char* msg);
static inline void call_log_cb(mihomo_log_cb cb, int lvl, const char* msg) {
	if (cb) cb(lvl, msg);
}
*/
import "C"

import (
	"errors"
	"sync/atomic"
	"unsafe"

	"github.com/metacubex/mihomo/log"
	"github.com/loafman-kangjun/libclash/libcore"
)

/* --------------- 内部工具 --------------- */

// 最后一条错误信息（线程安全）
var lastErr atomic.Pointer[C.char]

func recordErr(err error) int32 {
	// 清理之前的
	if p := lastErr.Swap(nil); p != nil {
		C.free(unsafe.Pointer(p))
	}
	if err == nil {
		return 0
	}
	cs := C.CString(err.Error()) // malloc(C)
	lastErr.Store(cs)
	return 1 // 统一返回 1 表示失败，可拓展更多错误码
}

//export Mihomo_GetLastError
func Mihomo_GetLastError() *C.char {
	return lastErr.Load() // 由调用方 free()
}

//export Mihomo_FreeCString
func Mihomo_FreeCString(p *C.char) { C.free(unsafe.Pointer(p)) }

/* --------------- 日志回调 --------------- */

var cLogCb atomic.Pointer[C.mihomo_log_cb]

//export Mihomo_SetLogCallback
func Mihomo_SetLogCallback(cb C.mihomo_log_cb) {
	cLogCb.Store(&cb) // 保存的是 C 世界的指针，合法
}

// 将 mihomo/log 的事件向 C 转发
func init() {
	sub := log.Subscribe()
	go func() {
		for ev := range sub {
			if pcb := cLogCb.Load(); pcb != nil && *pcb != nil {
				cmsg := C.CString(ev.Payload)           // malloc(C)
				C.call_log_cb(*pcb, C.int(ev.LogLevel), cmsg)
				C.free(unsafe.Pointer(cmsg))           // 立即释放
			}
		}
	}()
}

/* --------------- 核心导出 --------------- */

//export Mihomo_Run
func Mihomo_Run(cPath, cB64 *C.char, geo C.bool) C.int {
	err := libcore.Run(C.GoString(cPath), C.GoString(cB64), bool(geo))
	return C.int(recordErr(err))
}

//export Mihomo_Reload
func Mihomo_Reload(cPath, cB64 *C.char) C.int {
	err := libcore.Reload(C.GoString(cPath), C.GoString(cB64))
	return C.int(recordErr(err))
}

//export Mihomo_Stop
func Mihomo_Stop() {
	libcore.Stop()
	recordErr(nil) // 清掉旧错误
}

/* --------------- main --------------- */

// 空 main 防止 “no non-test Go files” 报错
func main() {}

/* --------------- 额外：确保静态检查不移除关键引用 --------------- */

var _ = errors.New // import keep-alive