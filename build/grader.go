package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	partId := os.Args[2]
	solutionFilesPath := "/grader/solutionFiles"
	switch partId {
	case "HW1_game":
		runTestHW("/grader/HW1_game", solutionFilesPath)
	default:
		log.Fatalln("No valid partId.")
	}
}

func runTestHW(path string, filesPath string) {
	dir := path
	files, err := os.ReadDir(filesPath)
	if err != nil {
		log.Fatalf("FAIL\ncompilation error\n\n%v", err)
	}

	for _, file := range files {
		if err = copyFile(filesPath, path, file.Name()); err != nil {
			log.Fatalf("FAIL\ncompilation error\n\n%v", err)
		}
	}

	cmd := exec.Command("go", "test", "-v")
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	errorText := fmt.Sprintf("STDOUT:\n%sFAIL\n\nSTDERR:\n%s", stdout.String(), stderr.String())
	if err != nil {
		errorText = fmt.Sprintf("FAIL\ncompilation error\n\n%s", errorText)
		log.Fatalf(errorText)
	}

	os.Exit(0)
}

func copyFile(srcDir, dstDir, fileName string) error {
	srcFile := filepath.Join(srcDir, fileName)
	dstFile := filepath.Join(dstDir, fileName)

	original, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer original.Close()

	file, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, original); err != nil {
		return err
	}

	return nil
}
