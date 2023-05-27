package ssh

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/containerd/console"
	"github.com/trzsz/trzsz-go/trzsz"
	"golang.org/x/crypto/ssh"
)

func Connect(ip string, port int, cfg *ssh.ClientConfig) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	sshClient, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return nil, err
	}
	return sshClient, nil
}

func RunTerminal(c *ssh.Client, in io.Reader, stdOut, stdErr io.WriteCloser) error {
	session, err := c.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Signal(ssh.SIGINT)
	session.Stderr = stdErr
	var (
		current = console.Current()
		ws      console.WinSize
	)

	if err = current.SetRaw(); err != nil {
		return err
	}

	if ws, err = current.Size(); err != nil {
		return err
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     //打开回显
		ssh.TTY_OP_ISPEED: 14400, //输入速率 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, //输出速率 14.4kbaud
		ssh.VSTATUS:       1,
	}

	//Request pseudo terminal
	if err = session.RequestPty("xterm-256color", int(ws.Height), int(ws.Width), modes); err != nil {
		return err
	}

	// Support trzsz ( trz / tsz )
	serverIn, err := session.StdinPipe()
	if err != nil {
		return err
	}
	serverOut, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	trzszFilter := trzsz.NewTrzszFilter(in, stdOut, serverIn, serverOut,
		trzsz.TrzszOptions{TerminalColumns: int32(ws.Width)})

	if err = session.Shell(); err != nil {
		return err
	}
	go consoleMonitor(current, session, trzszFilter)
	return session.Wait()
}

func consoleMonitor(c console.Console, session *ssh.Session, trzszFilter *trzsz.TrzszFilter) {
	var (
		t     = time.NewTicker(time.Second / 10)
		ws, _ = c.Size()
	)
	for {
		select {
		case <-t.C:
			cws, err := c.Size()
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			if cws.Height != ws.Height || cws.Width != ws.Width {
				_ = session.WindowChange(int(cws.Height), int(cws.Width))
				trzszFilter.SetTerminalColumns(int32(cws.Width))
				ws = cws
			}
		}
	}
}

func PwdCfg(user, pwd string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
}

func PkCfg(user string, pemBytes []byte) (*ssh.ClientConfig, error) {
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("Parsing plain private key failed %v", err)
	}
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}, nil
}
