package libcore

import (
	"encoding/base64"

	"github.com/metacubex/mihomo/hub"
)

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

	// 传入 nil 时 hub.Parse 会使用文件配置
	return hub.Parse(cfgBytes)
}