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

		key, value, err := splitKeyValue(line)

		if err != nil {
			return nil, err
		}

		pairs[key] = value
	}

	return pairs, nil
}

func splitKeyValue(l string) (key string, value string, err error) {

	pair := strings.SplitN(l, "=", 2)

	if len(pair) != 2 {
		return "", "", fmt.Errorf("invalid syntax [%s] ", l)
	}

	key = strings.TrimSpace(pair[0])

	if !isValidKey(key) {
		return "", "", fmt.Errorf("invalid key  [%s] , the key should only contains letters, numbers and underscore", key)
	}

	// exclude comments and quotes
	value = strings.TrimSpace(strings.Split(pair[1], "#")[0])

	if !isValidValue(value) {
		return "", "", fmt.Errorf("invalid value [%s] for key [%s]", value, key)
	}

	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		value = strings.Trim(value, "\"")
	}

	return key, value, nil
}

func isValidKey(key string) bool {
	rx := regexp.MustCompile(`^[A-z_0-9]+$`)

	return rx.MatchString(key)
}

func isValidValue(value string) bool {

	var rx *regexp.Regexp

	if strings.HasPrefix(value, `"`) {
		rx = regexp.MustCompile(`^"[^"]+"$`)
	} else {
		rx = regexp.MustCompile(`^.+$`)
	}

	return rx.MatchString(value)
}

func isComment(l string) bool {
	rx := regexp.MustCompile(`^(\s)*#[^\n]*$`)

	return rx.MatchString(l)
}
