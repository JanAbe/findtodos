package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	args := setupFlags()
	processTodos(*args.directory, *args.extension, *args.output)
}

// ProcessTodos processes all todos, it searches for them and writes them to the output file
func processTodos(dir, ext, out string) {
	outputFile, err := createOutputFile(out)
	if err != nil {
		fmt.Println(err)
	}

	files, err := findAllFiles(dir, ext)
	if err != nil {
		fmt.Println(err)
	}

	var wg sync.WaitGroup
	c := make(chan []todo)
	for _, file := range files {
		wg.Add(1)
		go fetchTodos(file, c, &wg)
	}

	// Close the channel when all goroutines have finished
	go func() {
		wg.Wait()
		close(c)
	}()

	for todos := range c {
		writeTodos(todos, outputFile)
	}

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

// fetchTodos fetches all todo's from the file and sends them to the channel
func fetchTodos(file string, c chan []todo, wg *sync.WaitGroup) {
	// todos, err := findTodosInFile(file)
	todos, err := newFindInFile(file)
	if err != nil {
		fmt.Println(err)
	}
	c <- todos
	wg.Done()
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

func newFindInFile(path string) ([]todo, error) {
	var todos []todo

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("error opening file %q with error: %v\n", path, err)
		return todos, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	var line string
	var text string

	// todo: make it recursive, atm it only checks the next comment
	// todo: rewrit code
	// if line contains a todo
	// check if next line starts with text indented to the todo of the line above and does not contain a todo

	// maybe do something like:
	// while nextLine isContinuation ...

	for scanner.Scan() {
		lineNumber++
		line = scanner.Text()
		todoText, todoPos := findTodoInLine(line)
		if todoPos == -1 {
			continue
		}

		text += todoText
		scanner.Scan()
		line = scanner.Text()
		if isContinuationOfTodo(line, todoPos) {
			text += " " + line[todoPos+1:]
		}
		todos = append(todos, todo{file.Name(), lineNumber, text})
		lineNumber++
		text = ""
		continue
	}

	return todos, nil
}

// isContinuationOfTodo check if the line is a continuation of the previous line by looking at the position of the todo
// e.g:
// todo: some task
//  that needs to be fulfilled
// has a correct indentation and is thus, a continuation
func isContinuationOfTodo(line string, todoPos int) bool {
	commentIdentifier := "//"
	space := " "
	key := "todo"
	prefix := commentIdentifier + strings.Repeat(space, todoPos-1)
	containsTodo := strings.Contains(line, key)
	hasCorrectPrefix := strings.HasPrefix(line, prefix)

	return hasCorrectPrefix && !containsTodo
}

// is the same as findTodoInString, is here for experimenting purposes
func findTodoInLine(line string) (string, int) {
	line = strings.ToLower(line)
	todoPos := -1
	todoText := ""
	if strings.HasPrefix(line, "//") {
		todoPos = strings.Index(line, "todo")
		if todoPos == -1 {
			return todoText, -1
		}
		todoText = line[todoPos:]
	}

	return todoText, todoPos
}

// todo: rename to FindTodoInLine()
// FindTodoInString finds and returns the body of the todo from a string,
// but only if the string contains a todo
func findTodoInString(line string) string {
	var todoText string
	var startPos int
	commentIndicator := "//"
	key := "todo"
	n := len(key)
	charsToRemove := ": "

	line = strings.TrimSpace(line)
	line = strings.ToLower(line)
	if strings.HasPrefix(line, commentIndicator) {
		startPos = strings.Index(line, key)
		if startPos == -1 {
			return ""
		}
		todoText = line[startPos+n:] // get todo's text (everything from todo -> till the end of line)
	}
	todoText = strings.TrimLeft(todoText, charsToRemove)
	return todoText
}

// todo struct to capture todo's
type todo struct {
	FileName   string
	LineNumber int
	Text       string
}

// toString returns a string representation of the Todo
func (t todo) toString() string {
	return fmt.Sprintf("%s - %d - %s\n", t.FileName, t.LineNumber, t.Text)
}

// args acts as a wrapper for all provided arguments
type args struct {
	directory *string
	extension *string
	output    *string
}

// setupFlags initializes all supported flags and parses the provided values to these flags
// returning them at the end
func setupFlags() args {
	directory := flag.String("directory", ".", "This flag is used to specify which directory should be scanned.")
	extension := flag.String("extension", ".go", "This flag is used to specify the extension type the program should look for.")
	output := flag.String("output", "./found_todos.txt", "This flag is used to specify the output location of the found todo's.")
	flag.Parse()

	flags := args{directory, extension, output}

	return flags
}
