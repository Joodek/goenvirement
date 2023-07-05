package goenvirement

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

// parse one or more files , if files not provided is empty we default to .env
func parseOrDefault(files ...string) (map[string]string, error) {

	// if no files provided we default to .env
	if len(files) < 1 {
		wd, err := os.Getwd()

		if err != nil {
			return nil, err
		}

		return parseFile(fmt.Sprintf("%v/.env", wd))
	}

	// if only one file provided no need to go through more steps
	if len(files) == 1 {
		return parseFile(files[0])
	}

	return parseMultipleFiles(files)
}

// parse the given files and merge them to a single map
func parseMultipleFiles(files []string) (map[string]string, error) {
	env := make(map[string]string)

	for _, file := range files {

		// we parse each file at once
		fileEnvs, err := parseFile(file)

		if err != nil {
			return nil, err
		}

		// we merge the envs from each file to a single map
		for key, value := range fileEnvs {
			env[key] = value
		}
	}

	return env, nil
}

// parse the given file and return it as a map
func parseFile(file string) (map[string]string, error) {
	if !fileExists(file) {
		return nil, &fs.PathError{
			Op:   "parsing the env file",
			Path: file,
			Err:  errors.New("file doesn't exists"),
		}
	}

	content, err := getFileContent(file)

	if err != nil {
		return nil, err
	}

	return buildKeyValuePairs(content)

}

// check if a file exists or not
func fileExists(file string) bool {
	_, err := os.Stat(file)

	return err == nil
}

// get the given file content , if something goes wrong
// the error will be returned with empty string
func getFileContent(file string) (string, error) {
	content, err := os.ReadFile(file)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

func buildKeyValuePairs(s string) (map[string]string, error) {
	lines := strings.Split(s, "\n")
	pairs := make(map[string]string)

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if isComment(line) || line == "" {
			continue
		}

		if !isValidPair(line) {
			return nil, errors.New("invalid key value pairs : [" + line + "]")
		}

		key, value := splitKeyValue(line)
		pairs[key] = strings.TrimSpace(strings.Trim(value, "\""))
	}

	return pairs, nil
}

func splitKeyValue(l string) (key string, value string) {

	pair := strings.SplitN(l, "=", 2)

	key = strings.TrimSpace(pair[0])

	// exclude comments
	value = strings.TrimSpace(strings.Split(pair[1], "#")[0])

	return key, value
}

func isValidPair(l string) bool {
	rx := regexp.MustCompile(`^(\s)*[A-z_0-9]+(\s)*=(\s)*[^\n]+(\s*)(#[^\n]*)*$`)

	return rx.MatchString(l)
}

func isComment(l string) bool {
	rx := regexp.MustCompile(`^(\s)*#[^\n]*$`)

	return rx.MatchString(l)
}
