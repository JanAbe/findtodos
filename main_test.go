package main

import (
	"fmt"
	"testing"
)

// todo rewrite test
func TestCreatingFile(t *testing.T) {
	f, err := createOutputFile("dank_test.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f)
	fmt.Println()
}

// todo: make temp directory to test it
func TestFindAllFiles(t *testing.T) {
	files, err := findAllFiles(".", ".go")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(files)
	fmt.Println()
}

// todo: make test a real test lol
func TestFindTodoInString(t *testing.T) {
	txt := findTodoInString("// hello everyone :), TODO: make temp directory to test it")
	if txt == "" {
		fmt.Println("no todo found")
	}
	fmt.Println(txt)
	fmt.Println()
}

func TestFindTodosInFile(t *testing.T) {
	todos, err := findTodosInFile("./main_test.go")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(todos)
	fmt.Println()
}

func TestFindTodosInDir(t *testing.T) {
	todos, err := findTodosInDir(".", ".go")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(todos)
	fmt.Println()
}
