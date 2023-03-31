## Computer Architecture Project

This respository is the group project for CS3339

### Developing the Application


First install Go development tools, then the application can be run with "go run main.go"

When running the application there are two flags available, if they are not provided the application cannot run.

  -i (string) the input file you want to parse

  -o (string) the output file you want to write

A sample input file, addtest1_bin.txt is included in the repository.

```bash
# an example with command line flags
go run team11_project2.go -i examples/project2_test.txt
```

### Team Members

- Ethan
- Chase
- Zoe

### Potential test cases - add any you can think of

invalid binary files

### Useful git commands

check status of files: `git status`

This tracks new file and any file changes: `git add filename.txt`

commits changes locally - `git commit -m "commit message"`

push commited changes to github - `git push`

switches branch - `git checkout <branch>`

create and checkout new branch off of current - `git checkout -b <new branch name>`