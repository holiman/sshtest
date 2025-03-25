package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"bufio"
	ssh2 "github.com/holiman/sshtest/ssh"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:   "boom",
		Usage:  "make an explosive entrance",
		Action: testSsh,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Value:    "127.0.0.1",
				Usage:    "Hostname or IP for the ssh server",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "port",
				Value: "22",
				Usage: "Ssh port",
			},
			&cli.StringFlag{
				Name:     "keyfile",
				Value:    "",
				Usage:    "File containing public keys to attempt",
				Required: true,
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Info("Error occurred", "err", err)
		os.Exit(1)
	}
}

func testSsh(ctx context.Context, cmd *cli.Command) error {
	var (
		user    = "root"
		host    = cmd.String("host")
		port    = cmd.String("port")
		keyFile *os.File
		err     error
		addr    = net.JoinHostPort(host, port)
	)
	// Open the file containing pubkeys
	if keyFile, err = os.Open(cmd.String("keyfile")); err != nil {
		return err
	}
	defer keyFile.Close()

	scanner := bufio.NewScanner(keyFile)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		pubstr := scanner.Text()
		if len(pubstr) == 0 {
			continue
		}
		// Create the public key to test with
		pubkey, _, _, _, err := ssh2.ParseAuthorizedKey([]byte(pubstr))
		if err != nil {
			return err
		}

		doAttempt(addr, user, pubkey)
	}

	return nil
}

func doAttempt(addr string, user string, pubkey ssh2.PublicKey) error {
	// Connect to the host
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		slog.Error("Failed to connect", "addr", addr, "err", err)
		return err
	}
	slog.Info("TCP connected", "addr", addr)
	slog.Info("Testing", "user", user, "pubkey", fmt.Sprintf("%v %v", pubkey.Type(), base64.RawStdEncoding.EncodeToString(pubkey.Marshal())))
	// Trigger handshake
	ssh2.NewClientConn(conn, addr, &ssh2.ClientConfig{
		User: user,
		Auth: []ssh2.AuthMethod{
			ssh2.PublicKeys(&publicOnlySigner{
				key: pubkey,
			}),
		},
		HostKeyCallback: ssh2.InsecureIgnoreHostKey(),
	})
	return nil
}

type publicOnlySigner struct {
	key ssh2.PublicKey
}

func (p *publicOnlySigner) PublicKey() ssh2.PublicKey {
	return p.key
}

func (p *publicOnlySigner) Sign(rand io.Reader, data []byte) (*ssh2.Signature, error) {
	panic("not supported")
}
