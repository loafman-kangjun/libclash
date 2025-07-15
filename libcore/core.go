package libcore

import (
	"context"
	"sync"
)

var (
	startOnce sync.Once
	stopOnce  sync.Once

	cancel context.CancelFunc
)