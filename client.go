package sshc

import (
	"bytes"
	"io/ioutil"
	"log"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// Client is the ssh client
type Client struct {
	sshCfg *ssh.ClientConfig
}

// Config contains the ssh config
type Config struct {
	User           string
	PrivateKeyFile string
}

// NewClient returns new ssh client
func NewClient(c *Config) (*Client, error) {
	key, err := ioutil.ReadFile(c.PrivateKeyFile)
	if err != nil {
		return nil, errors.Errorf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, errors.Errorf("unable to parse private key: %v", err)
	}

	sshCfg := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// I'm using InsecureIgnoreHostKey for re-used private IP(s)
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	return &Client{
		sshCfg: sshCfg,
	}, nil
}

// Connect will open connection to the remote server and perform the SSH handshake.
func (c *Client) Connect(ip, port string) (*ssh.Client, error) {
	sshClient, err := ssh.Dial("tcp", ip+":"+port, c.sshCfg)
	if err != nil {
		return nil, errors.Errorf("unable to connect: %v", err)
	}

	return sshClient, nil
}

// Run will run all commands
func (c *Client) Run(sshClient *ssh.Client, cmds []string) error {
	var b bytes.Buffer

	for i := range cmds {
		log.Printf("Running cmd[%v]: [%v]\n", i, cmds[i])

		session, err := sshClient.NewSession()
		if err != nil {
			return errors.Errorf("Failed to create session: %v", err)
		}
		defer session.Close()

		session.Stdout = &b
		if err := session.Run(cmds[i]); err != nil {
			return errors.Errorf("Failed to run: %v", err)
		}
	}

	return nil
}
