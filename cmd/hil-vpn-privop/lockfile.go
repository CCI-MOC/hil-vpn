package main

import (
	"golang.org/x/sys/unix"
	"os"
)

// This file implements a simple file lock, allowing us to assume that only one instance of
// hil-vpn-privop is running at a time.

const lockFilePath = "/tmp/hil-vpn.lock"

var lockFileRef *os.File

func lockFile() {
	file, err := os.Create(lockFilePath)
	chkfatal("Opening lock file", err)
	lockFileRef = file
	chkfatal("Locking file", unix.Flock(int(file.Fd()), unix.LOCK_EX))
}

func unlockFile() {
	chkfatal("Unlocking file", unix.Flock(int(lockFileRef.Fd()), unix.LOCK_UN))
	lockFileRef.Close()
}
