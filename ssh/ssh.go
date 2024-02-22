package ssh

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"io"
	"strings"
	"time"

	"github.com/containerd/console"
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
	session.Signal(ssh.SIGINT)
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
		return nil, fmt.Errorf("Parsing plain private key failed %v", err)
	}
	return cfg(user, ssh.PublicKeys(signer)), nil
}

func cfg(user string, auth ...ssh.AuthMethod) *ssh.ClientConfig {
	c := &ssh.ClientConfig{
		User:            user,
		Auth:            append(auth, ssh.KeyboardInteractive(keyboardInteractiveChallenge)),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	return c
}

func keyboardInteractiveChallenge(name, instruction string, questions []string, echos []bool) (answers []string, err error) {

	for _, v := range questions {
		var res = ""
		if len(v) > 0 {
			prompt := promptui.Prompt{
				Label: strings.Join(questions, ""),
				Templates: &promptui.PromptTemplates{
					Prompt:  "{{ . }} ",
					Valid:   "{{ . | green }} ",
					Success: "{{ . | green }} ",
				},
			}
			res, err = prompt.Run()
		}
		if err != nil {
			return nil, err
		}
		answers = append(answers, res)
	}
	return
}
