package utils

import (
	"os"
	"syscall"
)

// LockFile 锁定文件，返回文件锁对象和错误信息。
func LockFileWithExclusive(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}

func LockFileWithShared(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}

// UnlockFile 解锁文件。
func UnlockFile(file *os.File) error {
	err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	if err != nil {
		return err
	}

	return file.Close()
}
