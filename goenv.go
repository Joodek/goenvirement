package goenv

import (
	"os"
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

	expand(envs)

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

	expand(envs)

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

	expand(envs)

	return envs, nil
}

// read envirement variables from s and return them as map
func Unmarshal(s string) (map[string]string, error) {
	envs, err := buildKeyValuePairs(s)

	if err != nil {
		return nil, err
	}

	expand(envs)

	return envs, nil
}
