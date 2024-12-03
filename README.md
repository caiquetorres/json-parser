# JSON Parser
This repository provides a solution for [Coding Challenge #02](https://codingchallenges.fyi/challenges/challenge-json-parser), which involves creating a JSON parser.

## Build and Run

To build and run the tool, use the following commands:

```sh
# Build the project
make build
```

```sh
# Parsing files
./bin/json_parser file1.json file2.json file3.json

# Parsing from stdin
cat file.json | ./bin/json_parser
```

## Test

To run tests for the tool, use the following command:

```sh
# Run tests
make test
```
