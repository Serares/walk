### cli file crawler

- tests are written using table-driven testing
- because this tool also creates logfiles you can run it with a cronjob and check the logs later on

Example archive some file:

`./walk -root ./ -archive /tmp/archive_tmp -ext ".go .log"`

The above command will archive all the files on path `./` and store them in `/tmp/archive_tmp` if some of the files are
stored in a directory the tool will keep the structure creating the directories

The above command was run on the structure:
```
./
├── Makefile
├── README.md
├── actions.go
├── actions_test.go
├── archived
├── bin
│   └── walk
├── go.mod
├── helper.go
├── main.go
├── main_test.go
├── testdata
│   ├── dir.log
│   └── dir2
│       └── script.sh
```
And it generated the following archives:
```
/tmp/archive_tmp
├── actions.go.gz
├── actions_test.go.gz
├── helper.go.gz
├── main.go.gz
├── main_test.go.gz
└── testdata
    └── dir.log.gz
```