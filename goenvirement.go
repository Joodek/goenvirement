package goenvirement

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// load the envirements variable from files
//
// if no file provided, we fall back to .env file in the current working directory
// The variables will be loaded and set to the current process env so you can get them via os.Getenv
//
// The existing variables won't be overriden
func Load(files ...string) error {
	envs, err := parseOrDefault(files...)

	if err != nil {
		return err
	}

	expand(envs, &err)

	if err != nil {
		return err
	}

	for key, value := range envs {
		if _, ok := os.LookupEnv(key); ok {
			continue
		}

		err := os.Setenv(key, value)

		if err != nil {
			return err
		}
	}

	return nil
}

// load the envirements variable from files
//
// if no file provided, we fall back to .env file in the current working directory
// The variables will be loaded and set to the current process env so you can get them via os.Getenv
//
// unlike Load , Overload will override the existing variables
func Overload(files ...string) error {
	envs, err := parseOrDefault(files...)

	if err != nil {
		return err
	}

	expand(envs, &err)

	if err != nil {
		return err
	}

	for key, value := range envs {
		err := os.Setenv(key, value)

		if err != nil {
			return err
		}
	}

	return nil
}

// read envirement variables from files and return them as map
// if files is empty we fall back to .env file in the current directory
func Read(files ...string) (map[string]string, error) {
	envs, err := parseOrDefault(files...)

	if err != nil {
		return nil, err
	}

	expand(envs, &err)

	if err != nil {
		return nil, err
	}

	return envs, nil
}

// read envirement variables from s and return them as map
func Unmarshal(s string) (map[string]string, error) {
	envs, err := buildKeyValuePairs(s)

	if err != nil {
		return nil, err
	}

	expand(envs, &err)

	if err != nil {
		return nil, err
	}

	return envs, nil
}

// stringify the envirement variables from m and return them as string
func Marshal(m map[string]string) (string, error) {
	var lines []string

	for k, v := range m {
		if err := validateKeyValuePair(k, v); err != nil {
			return "", err
		}

		if isComment(v) {
			return "", fmt.Errorf("invalid value [%s] for key %s", v, k)
		}

		lines = append(lines, fmt.Sprintf("%s=%s\n", k, v))
	}

	sort.Strings(lines)

	return strings.Join(lines, ""), nil
}

// write  variables from m to f file
func Write(m map[string]string, f string) error {
	content, err := Marshal(m)

	if err != nil {
		return err
	}

	err = os.WriteFile(f, []byte(content), 0644)

	if err != nil {
		return err
	}

	return nil
}

func Append(key string, value string, file string) (err error) {

	envs, err := parseOrDefault(file)

	if err != nil {
		return
	}

	envs[key] = value
	return
}
