package libcore

import (
	"github.com/metacubex/mihomo/component/updater"
	"github.com/metacubex/mihomo/hub"
	"github.com/metacubex/mihomo/hub/executor"

	"context"
)

// Run 启动 mihomo 内核。
// cfgPath:  配置文件路径，为空时默认 $HOME/.config/mihomo/config.yaml
// cfgB64 :  base64 编码的配置，为 "" 时优先使用 cfgPath
// geoMode:  开启 geodata 模式
func Run(cfgPath, cfgB64 string, geoMode bool) error {
	var retErr error

	startOnce.Do(func() {
		startLogObserver()

		// ---------- 解析配置 ----------
		cfgBytes, err := loadConfigBytes(cfgPath, cfgB64, geoMode)
		if err != nil {
			retErr = err
			return
		}

		// ---------- 启动核心 ----------
		if err := hub.Parse(cfgBytes); err != nil {
			retErr = err
			return
		}

		if updater.GeoAutoUpdate() {
			updater.RegisterGeoUpdater()
		}

		// ---------- 生命周期 ----------
		ctx, cancelFn := context.WithCancel(context.Background())
		cancel = cancelFn

		// 收到 ctx.Done() 时优雅关停
		go func() {
			<-ctx.Done()
			executor.Shutdown()
		}()
	})

	return retErr
}