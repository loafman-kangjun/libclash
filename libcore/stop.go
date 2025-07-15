package libcore

import "github.com/metacubex/mihomo/hub/executor"

// Stop 优雅关闭 mihomo 内核。
func Stop() {
	stopOnce.Do(func() {
		if cancel != nil {
			cancel()
		} else {
			// 即便 Run 失败也兜底关停
			executor.Shutdown()
		}
	})
}