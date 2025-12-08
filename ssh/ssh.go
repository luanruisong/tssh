package ssh

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/containerd/console"
	"github.com/manifoldco/promptui"
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

func RunTerminal(c *ssh.Client, in io.Reader, stdOut, stdErr io.Writer) error {
	session, err := c.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stdout = stdOut
	session.Stderr = stdErr
	session.Stdin = in
	var (
		current = console.Current()
		ws      console.WinSize
	)

	if err = current.SetRaw(); err != nil {
		return err
	}
	// 确保在函数返回时恢复终端状态
	defer current.Reset()

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

	// 监听窗口大小变化
	stopResize := listenWindowResize(session, current)
	defer stopResize()

	if err = session.Shell(); err != nil {
		return err
	}
	return session.Wait()
}

func PwdCfg(user, pwd string) *ssh.ClientConfig {
	return cfg(user, ssh.Password(pwd))
}

func PkCfg(user string, pemBytes []byte) (*ssh.ClientConfig, error) {
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing plain private key failed: %v", err)
	}
	// 使用 RetryableAuthMethod 包装 publickey，支持 partial success 后继续认证
	pkAuth := ssh.RetryableAuthMethod(ssh.PublicKeys(signer), 3)
	return cfg(user, pkAuth), nil
}

func cfg(user string, auth ...ssh.AuthMethod) *ssh.ClientConfig {
	// keyboard-interactive 认证，用于 2FA/OTP 等场景
	kbAuth := ssh.RetryableAuthMethod(ssh.KeyboardInteractive(keyboardInteractiveChallenge), 3)
	c := &ssh.ClientConfig{
		User:            user,
		Auth:            append(auth, kbAuth),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	return c
}

func keyboardInteractiveChallenge(name, instruction string, questions []string, echos []bool) (answers []string, err error) {
	// 显示服务器发送的指令信息
	if len(instruction) > 0 {
		fmt.Println(instruction)
	}

	answers = make([]string, len(questions))
	for i, question := range questions {
		if len(question) == 0 {
			answers[i] = ""
			continue
		}

		prompt := promptui.Prompt{
			Label: strings.TrimSpace(question),
			Templates: &promptui.PromptTemplates{
				Prompt:  "{{ . }} ",
				Valid:   "{{ . | green }} ",
				Success: "{{ . | green }} ",
			},
			// 如果 echos[i] 为 false，则隐藏输入（如密码/OTP）
			Mask: func() rune {
				if i < len(echos) && !echos[i] {
					return '*'
				}
				return 0
			}(),
		}
		answers[i], err = prompt.Run()
		if err != nil {
			return nil, err
		}
	}
	return answers, nil
}
