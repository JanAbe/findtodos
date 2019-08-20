package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/*
	input = directory
	args -> option = welke extensie (.go | .java | .js)
		findtodos -e .go
	of geef config bestand mee
		findtodos -c [path.to.config]

	eerst normaal, daarna concurrent maken
	concurrent:
	ga langs alle .go bestanden (bijv.)
	lees alle bestanden
	ga opzoek naar alle comments die beginnen met todo: ...
		todo: kunnen vervangen met een door de gebruiker gespecificeerd patroon
			bijv. todo 1:
	als .go -> comment = // of /*
	als .py -> comment = #

	extraheer de tekst die daar achter staat
	// todo: refactor code -> [filepath] - [lineNumber] - [todo text]

	misschien nog iets van een flag meegeven dat je de uitvoeringstijd wilt meten
*/
func main() {
	// todo: make the code concurrent
	args := setupFlags()

	outputFile, err := createOutputFile(*args.output)
	if err != nil {
		fmt.Println(err)
	}

	foundTodos, err := findTodosInDir(*args.directory, *args.extension)
	if err != nil {
		fmt.Println(err)
	}

	err = writeTodos(foundTodos, outputFile)
	if err != nil {
		fmt.Println(err)
	}

}

type todo struct {
	FileName   string
	LineNumber int
	Text       string
}

func (t todo) toString() string {
	return fmt.Sprintf("%s - %d - %s\n", t.FileName, t.LineNumber, t.Text)
}

// CreateOutputFile create a file at the specified location to which
// all found todo's are written
func createOutputFile(outputLocation string) (string, error) {
	outputFile, err := os.Create(outputLocation)
	if err != nil {
		// todo: check if this is ok, else -> rewrite error handling
		return "", err
	}
	defer outputFile.Close()

	return outputLocation, err
}

// FindTodosInDir find all todo's that are present in the directory
func findTodosInDir(dir, extension string) ([]todo, error) {
	var todos []todo

	files, err := findAllFiles(dir, extension)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		foundTodos, err := findTodosInFile(file)
		if err != nil {
			return nil, err
		}

		for _, foundTodo := range foundTodos {
			todos = append(todos, foundTodo)
		}
	}

	return todos, nil
}

// FindTodosInFile find all todo's in the provided file
func findTodosInFile(path string) ([]todo, error) {
	var todos []todo

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("error opening file %q with error: %v\n", path, err)
		return todos, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNumber := 0
	var line string
	for scanner.Scan() {
		lineNumber++
		line = scanner.Text()
		todoText := findTodoInString(line)
		if todoText == "" {
			continue
		}
		todos = append(todos, todo{file.Name(), lineNumber, todoText})
	}

	return todos, nil
}

// todo: make constants of "//" and startPos+len("todo"), what to trim, etc.
// todo: improve todo findings
// FindTodoInString finds and returns the body of the todo from a string,
// but only if the string contains a todo
func findTodoInString(line string) string {
	var todoText string
	var startPos int
	line = strings.TrimSpace(line)
	line = strings.ToLower(line)
	if strings.HasPrefix(line, "//") {
		startPos = strings.Index(line, "todo")
		if startPos == -1 {
			return ""
		}
		todoText = line[startPos+len("todo"):] // get todo's text (everything from todo -> till the end of line)
	}
	todoText = strings.TrimLeft(todoText, ": ")
	return todoText
}

// WriteTodos writes all provided todos to the specified file
func writeTodos(todos []todo, outputFile string) error {
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, todo := range todos {
		_, err := file.WriteString(todo.toString())
		if err != nil {
			return err
		}
	}
	return nil
}

// findAllFiles returns a slice of filePaths of files that live inside the
// provided directory and have the provided extension
func findAllFiles(dir, extension string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("failed accessing path: %q, with error: %v\n", path, err)
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == extension {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func welcome() string {
	msg := "Find todo's in your project.\nType -h for help."
	return msg
}

type args struct {
	directory *string
	extension *string
	output    *string
}

func setupFlags() args {
	directory := flag.String("directory", ".", "This flag is used to specify which directory should be scanned.")
	extension := flag.String("extension", ".go", "This flag is used to specify the extension type the program should look for.")
	output := flag.String("output", "./found_todos.txt", "This flag is used to specify the output location of the found todo's.")
	flag.Parse()

	flags := args{directory, extension, output}

	return flags
}
