package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestFindAllFiles(t *testing.T) {
	parentDir, tempFile1, tempFile2 := setupTempDir()
	defer os.RemoveAll(parentDir)

	files, err := findAllFiles(parentDir, ".go")
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		expected string
		actual   string
	}{
		{tempFile1, files[1]},
		{tempFile2, files[0]},
	}

	for _, test := range tests {
		if test.expected != test.actual {
			t.Errorf("Found files don't correspond with present files: expected=%q, got=%q", test.expected, test.actual)
		}
	}

}

func TestFindTodosInFile(t *testing.T) {
	_, tempFile1, _ := setupTempDir()
	todos, err := findTodosInFile(tempFile1)
	if err != nil {
		fmt.Println(err)
	}

	tests := []struct {
		expected todo
		actual   todo
	}{
		{todo{tempFile1, 1, "a lot of things"}, todos[0]},
		{todo{tempFile1, 4, "find this please"}, todos[1]},
	}

	for _, test := range tests {
		if test.expected != test.actual {
			t.Errorf("Found todo does not correspond to the actual todo: expected=%q, got=%q", test.expected, test.actual)
		}
	}
}

func TestFindTodoInString(t *testing.T) {
	todoTxt1 := findTodoInString(`
// todo make tests
`)

	todoTxt2 := findTodoInString(`
// Some random text todo: begins here
//                   and continues here
//                   and ends here
`)

	tests := []struct {
		expected string
		actual   string
	}{
		{"make tests", todoTxt1},
		{"begins here and continues here and ends here", todoTxt2},
	}

	for _, test := range tests {
		if test.expected != test.actual {
			t.Errorf("Found todo text does not correspond to the actual todo text: expected=%q, got=%q", test.expected, test.actual)
		}
	}
}

func BenchmarkMain(b *testing.B) {
	for n := 0; n < b.N; n++ {
		processTodos(".", ".go", "./found_todos.txt")
	}
}

// setupTempDir is a helper func to setup a temporary dir with temporary files
// it returns the paths of the directory and the two files inside this dir
func setupTempDir() (string, string, string) {
	parentDir, err := ioutil.TempDir("", "findtodos_test")
	if err != nil {
		log.Fatal(err)
	}

	childDir, err := ioutil.TempDir(parentDir, "child_test")
	if err != nil {
		log.Fatal(err)
	}

	tempFile1, err := ioutil.TempFile(parentDir, "tempFile1.*.go")
	if err != nil {
		log.Fatal(err)
	}
	tempFile2, err := ioutil.TempFile(childDir, "tempFile2.*.go")
	if err != nil {
		log.Fatal(err)
	}

	tempFile1.WriteString(`//todo: a lot of things
package main
//this should not be found.
// todo find this please`)

	return parentDir, tempFile1.Name(), tempFile2.Name()
}
