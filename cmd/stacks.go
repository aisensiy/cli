package cmd

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"net"
	"golang.org/x/crypto/ssh/agent"
	"io"
)


var makeKeyring = func() {
}

func StackCreate() error {
	sshConfig := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
	}

	connection, err := ssh.Dial("tcp", "192.168.50.4:2222", sshConfig)
	if err != nil {
		return  fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	err = session.Run("stack-init javajersey")
	return nil
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
