//go:build windows

package ssh

import (
	"time"

	"github.com/containerd/console"
	"golang.org/x/crypto/ssh"
)

// listenWindowResize Windows 平台通过轮询方式检测窗口大小变化
func listenWindowResize(session *ssh.Session, current console.Console) func() {
	done := make(chan struct{})

	go func() {
		var lastWidth, lastHeight uint16

		// 获取初始窗口大小
		if size, err := current.Size(); err == nil {
			lastWidth = size.Width
			lastHeight = size.Height
		}

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if size, err := current.Size(); err == nil {
					if size.Width != lastWidth || size.Height != lastHeight {
						lastWidth = size.Width
						lastHeight = size.Height
						_ = session.WindowChange(int(size.Height), int(size.Width))
					}
				}
			}
		}
	}()

	return func() {
		close(done)
	}
}
