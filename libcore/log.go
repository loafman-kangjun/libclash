package libcore

import (
	"fmt"

	log "github.com/metacubex/mihomo/log"
)

// startLogObserver 挂载一次日志观察者。
func startLogObserver() {
	sub := log.Subscribe()
	go func() {
		for ev := range sub {
			fmt.Printf("[OBS]%s: %s\n", ev.LogLevel, ev.Payload)
		}
	}()
}