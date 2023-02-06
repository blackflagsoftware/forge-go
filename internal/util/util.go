package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	clearMap map[string]func()
)

func init() {
	clearMap = make(map[string]func())
	clearMap["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearMap["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearMap["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ParseInput() string {
	reader := bufio.NewReader(os.Stdout)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	return s
}

func ClearScreen() {
	clearFunc, ok := clearMap[runtime.GOOS]
	if !ok {
		fmt.Println("\n *** Your platform is not supported to clear the terminal screen ***")
		return
	}
	clearFunc()
}

func BasicPrompt(mainMessage []string, prompts []string, acceptablePrompts []string, exitString string, f func()) string {
	for {
		fmt.Println("")
		for _, msg := range mainMessage {
			fmt.Println(msg)
		}
		fmt.Println("")
		for _, prompt := range prompts {
			fmt.Println(prompt)
		}
		if exitString != "" {
			// just in case you don't want to show this line
			fmt.Printf("(%s) to exit", exitString)
		}
		fmt.Println("")
		fmt.Println("")
		fmt.Print("Selection Choice: ")
		selection := ParseInput()
		if strings.ToLower(selection) == exitString {
			return exitString
		}
		found := false
		for _, acceptablePrompt := range acceptablePrompts {
			if strings.ToLower(selection) == acceptablePrompt {
				found = true
				break
			}
		}
		if !found {
			fmt.Print("Invalid selection, try again, press 'enter' to continue:")
			ParseInput()
			f() // this can be ClearScreen or any simple function as pre-
			continue
		}
		return strings.ToLower(selection)
	}
}

func AskYesOrNo(msg string) (answer bool) {
	for {
		msg = msg + " (y/n)? "
		fmt.Print(msg)
		def := ParseInput()
		switch def {
		case "y", "Y":
			answer = true
		case "n", "N":
			answer = false
		default:
			fmt.Println("Invalid value, get it together (y or n)!")
			continue
		}
		break
	}
	return
}

func FormatSql(sqlLines []string) string {
	// join all lines
	sql := strings.Join(sqlLines, "")
	// trim ends
	sql = strings.TrimSpace(sql)
	// lowercase
	sql = strings.ToLower(sql)
	// remove (`)s
	sql = strings.ReplaceAll(sql, "`", "")
	return sql
}

func BuildRandomString(length int) string {
	randomString := ""
	for {
		randomString += "0123456789"
		if len(randomString) > length {
			break
		}
	}
	return randomString
}

func CopyDirectory(src, dest string) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		fmt.Println("Error in making project directory:", err)
		return err
	}
	dirEntry, err := os.ReadDir(src)
	if err != nil {
		fmt.Printf("Unable to read directory %s: %s", src, err)
		return err
	}
	for _, entry := range dirEntry {
		newSrc := filepath.Join(src, entry.Name())
		newDest := filepath.Join(dest, entry.Name())
		if entry.IsDir() {
			CopyDirectory(newSrc, newDest)
			continue
		}
		CopyFile(newSrc, newDest)
	}
	return nil
}

func CopyFile(src, dest string) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Printf("Unable to open source file: %s; %s\n", src, err)
		return
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		fmt.Printf("Unable to create destination file: %s; %s\n", dest, err)
		return
	}
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		fmt.Printf("Unable to copy file: %s; %s\n", src, err)
		return
	}
	err = destFile.Sync()
	if err != nil {
		fmt.Printf("Unable to close/sync destination file: %s; %s\n", dest, err)
		return
	}
}
