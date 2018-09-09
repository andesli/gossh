package scp

import (
	"testing"

	"golang.org/x/crypto/ssh"
)

func newClient() (*ssh.Client, error) {
	pass := "redhat"
	user := "root"
	ip := "192.168.56.2"
	port := "22"
	psw := []ssh.AuthMethod{ssh.Password(pass)}
	Conf := ssh.ClientConfig{User: user, Auth: psw}

	ip_port := ip + ":" + port
	return ssh.Dial("tcp", ip_port, &Conf)
}

func TestPushFile(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatal(err)
	}
	scp := NewScp(client)
	srcFile := "/etc/profile"
	dest := "/tmp"
	err = scp.PushFile(srcFile, dest)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPushDir(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatal(err)
	}
	scp := NewScp(client)
	srcDir := "/project/go/src/gossh/scp"
	dest := "/tmp/"
	err = scp.PushDir(srcDir, dest)
	if err != nil {
		t.Fatal(err)
	}
}
