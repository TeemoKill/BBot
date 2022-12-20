package utils

import (
	"bufio"
	"os"
	"strings"
	"time"
)

// FileExist 判断文件是否存在
func FileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func readLine() (str string, err error) {
	reader := bufio.NewReader(os.Stdin)
	str, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	str = strings.TrimSpace(str)
	return str, nil
}

func ReadLine() (str string, err error) {
	return readLine()
}

func ReadLineTimeout(timeout time.Duration, defaultValue string) (str string, err error) {
	strChan := make(chan string)
	defer close(strChan)
	errChan := make(chan error)
	defer close(errChan)

	go func() {
		line, err := readLine()
		if err != nil {
			errChan <- err
		} else {
			strChan <- line
		}
	}()

	select {
	case str = <-strChan:
		return str, nil
	case err = <-errChan:
		return "", err
	case <-time.After(timeout):
		return defaultValue, nil
	}
}
