package main

import (
	"log"

	sshc "github.com/williamchanrico/go-sshc"
)

func main() {
	// Prepare ssh client
	client, err := sshc.NewClient(&sshc.Config{
		User:           "root",
		PrivateKeyFile: "/home/william/.ssh/test-unencrypted.pem",
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to remote host
	conn, err := client.Connect("172.21.45.20", "22")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	log.Println("Connected, running command(s)...")

	// Run command(s) on remote host
	cmds := []string{
		"echo 'Hello World!' > /root/hello.txt",
		"sleep 10",
		"echo 'Good bye!' >> /root/hello.txt",
	}
	if err = client.Run(conn, cmds); err != nil {
		log.Println(err)
	}
}
