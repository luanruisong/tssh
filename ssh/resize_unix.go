//go:build darwin || freebsd || linux || netbsd || openbsd || solaris

package ssh

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/containerd/console"
	"golang.org/x/crypto/ssh"
)

// listenWindowResize 监听终端窗口大小变化并通知远程会话
func listenWindowResize(session *ssh.Session, current console.Console) func() {
	sigwinchCh := make(chan os.Signal, 1)
	signal.Notify(sigwinchCh, syscall.SIGWINCH)

	go func() {
		for range sigwinchCh {
			if newSize, err := current.Size(); err == nil {
				_ = session.WindowChange(int(newSize.Height), int(newSize.Width))
			}
		}
	}()

	return func() {
		signal.Stop(sigwinchCh)
		close(sigwinchCh)
	}
}
