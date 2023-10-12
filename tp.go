package chLogger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func Getime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func FckText(text string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR HEAD: %s\n ERROT TEXT: %s", text, err)
		ExitWithSecTimeout(1)
	}
}

func ExitWithSecTimeout(status int) {
	// Го любит завершать свою работу раньше чем сделать все завершающие операции, но все же он остается очень быстрым одной секунды ему достатотчно
	// 0 - norm
	// 1 - error
	fmt.Println("Завершение работы программы через 2 сек.")
	time.Sleep(500 * time.Millisecond)
	os.Exit(status)
}
func BinDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func WorkDir() (string, error) {
	ex, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

// fmt.Sprintf("Ошибка, директория не может быть создана: %v", err)
func CheckMkdir(workFolder string) error {
	if _, err := os.Stat(workFolder); os.IsNotExist(err) {
		err = os.Mkdir(workFolder, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateOpenFile(filePath string) (file *os.File, err error) {
	file, err = os.Create(filePath)
	return
}

func OpenWriteFile(filePath string) (file *os.File, err error) {
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	return
}

func OpenReadFile(filePath string) (dat []byte, err error) {
	dat, err = ioutil.ReadFile(filePath)
	return
}
