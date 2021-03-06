package scmauth

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const SSHPrivateKeyMethodName = "ssh-privatekey"

// SSHPrivateKey implements SCMAuth interface for using SSH private keys.
type SSHPrivateKey struct{}

// Setup creates a wrapper script for SSH command to be able to use the provided
// SSH key while accessing private repository.
func (_ SSHPrivateKey) Setup(baseDir string) error {
	script, err := ioutil.TempFile("", "gitssh")
	if err != nil {
		return err
	}
	defer script.Close()
	if err := script.Chmod(0711); err != nil {
		return err
	}
	if _, err := script.WriteString("#!/bin/sh\nssh -i " +
		filepath.Join(baseDir, SSHPrivateKeyMethodName) +
		" -o StrictHostKeyChecking=false \"$@\"\n"); err != nil {
		return err
	}
	// set environment variable to tell git to use the SSH wrapper
	if err := os.Setenv("GIT_SSH", script.Name()); err != nil {
		return err
	}
	return nil
}

// Name returns the name of this auth method.
func (_ SSHPrivateKey) Name() string {
	return SSHPrivateKeyMethodName
}
