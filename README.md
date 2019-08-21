# Findtodos

#### Usage
findtodos -directory=... -extension=... -output=...
> -directory string
 > This flag is used to specify which directory should be scanned. (default ".")

> -extension string
> This flag is used to specify the extension type the program should look for. (default ".go")

> -output string
> This flag is used to specify the output location of the found todo's. (default "./found_todos.txt")

###### Example
`findtodos -directory=/home/user/workspace/project -extension=.go -output=/home/user/Documents/findtodos_output.txt`


#### Benchmark
Results of multiple benchmark tests.
The left part of the table contains the results of the execution of the current program (after the concurrent changes).
The right part of the table contains the results of the execution of the program before the changes, so being lineair.

|  current program (concurrent)    | | | old program (lineair) | |
|-------|--------------|-|-------|------------|
| **N** | **run time** | | **N** | **run time** | 
| 1000  | 1771085 ns   | | 500   | 2294080 ns   |
| 1000  | 1832276 ns   | | 500   | 2640026 ns   |
| 1000  | 1781122 ns   | | 1000  | 2248513 ns   |
| 1000  | 1831600 ns   | | 1000  | 2273937 ns   |
| 1000  | 2261793 ns   | | 1000  | 2281806 ns   |


As you can see, the change to make it concurrent has resulted in a performance boost 

#### Why
Just wanted to make a small program in Go, and to easily find all todo's in a project.

#### Todo:
Improve finding todo's, atm it's kinda shitty :c
Refactor code