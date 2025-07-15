// SPDX-License-Identifier: GPL-3.0
// 这是把原 main.go 精简后的“内核可重用 API”。

package libcore

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/metacubex/mihomo/component/geodata"
	"github.com/metacubex/mihomo/component/updater"
	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub"
	"github.com/metacubex/mihomo/hub/executor"

	log "github.com/metacubex/mihomo/log"	
)

var (
	startOnce sync.Once
	stopOnce  sync.Once
	cancel    context.CancelFunc
)

// Run 启动 mihomo 内核。
// cfgPath:  文件路径，为空时使用 $HOME/.config/mihomo/config.yaml
// cfgB64 :  直接给 base64-encoded 配置，可为 "" 表示使用文件。
// geoMode:  true 表示开启 geodata 模式。
func Run(cfgPath, cfgB64 string, geoMode bool) error {
	var retErr error
	sub := log.Subscribe()
	go func() {
		for ev := range sub { // ← 具体写法看 observable 的返回结构
			// ev 是 Event，包含 Level / Payload
			fmt.Printf("[OBS]%s: %s\n", ev.LogLevel, ev.Payload)
			// 可以写文件、推 WebSocket、发 C 回调……随便
		}
	}()
	startOnce.Do(func() {
		if geoMode {
			geodata.SetGeodataMode(true)
		}

		// ----------------- 准备配置 -----------------
		var cfgBytes []byte
		if cfgB64 != "" {
			var err error
			cfgBytes, err = base64.StdEncoding.DecodeString(cfgB64)
			if err != nil {
				retErr = err
				return
			}
		} else {
			if cfgPath == "" {
				cfgPath = filepath.Join(C.Path.HomeDir(), C.Path.Config())
			}
			C.SetConfig(cfgPath)

			if err := config.Init(C.Path.HomeDir()); err != nil {
				retErr = err
				return
			}
		}

		// ----------------- 启动核心 -----------------
		if err := hub.Parse(cfgBytes); err != nil {
			retErr = err
			return
		}
		if updater.GeoAutoUpdate() {
			updater.RegisterGeoUpdater()
		}

		ctx, cancelFn := context.WithCancel(context.Background())
		cancel = cancelFn

		// 监听 ctx.Done()，退出时调用 Shutdown
		go func() {
			<-ctx.Done()
			executor.Shutdown()
		}()
	})
	return retErr
}

// Reload 重新加载配置。参数同 Run（cfgPath & cfgB64 二选一）
func Reload(cfgPath, cfgB64 string) error {
	var cfgBytes []byte
	var err error
	if cfgB64 != "" {
		cfgBytes, err = base64.StdEncoding.DecodeString(cfgB64)
		if err != nil {
			return err
		}
	}
	return hub.Parse(cfgBytes)
}

// Stop 优雅关闭 mihomo 内核。
func Stop() {
	stopOnce.Do(func() {
		if cancel != nil {
			cancel()
		} else {
			// Run 失败时仍可能被调用，直接兜底
			executor.Shutdown()
		}
	})
}
