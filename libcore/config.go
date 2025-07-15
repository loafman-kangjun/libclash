package libcore

import (
	"encoding/base64"
	"path/filepath"

	"github.com/metacubex/mihomo/component/geodata"
	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
)

// loadConfigBytes 返回用户传入的原始配置字节（若传入的是文件路径则返回 nil）。
// 若 cfgB64 非空，优先生效。
func loadConfigBytes(cfgPath, cfgB64 string, geoMode bool) ([]byte, error) {
	if geoMode {
		geodata.SetGeodataMode(true)
	}

	if cfgB64 != "" {
		return base64.StdEncoding.DecodeString(cfgB64)
	}

	if cfgPath == "" {
		cfgPath = filepath.Join(C.Path.HomeDir(), C.Path.Config())
	}
	C.SetConfig(cfgPath)

	return nil, config.Init(C.Path.HomeDir())
}