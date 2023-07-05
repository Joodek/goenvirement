package goenvirement

import (
	"fmt"
	"os"
	"regexp"
)

func expand(env map[string]string, err *error) {
	defer func(err *error) {
		er := recover()

		if er != nil {
			*err = fmt.Errorf("%v", er)
		}

	}(err)

	for key, value := range env {
		env[key] = evaluate(env, key, value)
	}
}

func evaluate(env map[string]string, key string, value string) string {

	if !expectExpanding(value) {
		return value
	}

	return os.Expand(value, func(s string) string {
		// in case of key1="somthing ${key1} something"
		// we return to avoid infinit loop
		if s == key {
			return ""
		}

		throttle(key, s)

		v, exists := env[s]

		if !exists {
			// if the key is not exists in the env file, we look at the global env
			v, exists := os.LookupEnv(s)

			if exists {
				return v
			}
		}

		// in case the value needs more expanding
		return evaluate(env, s, v)
	})
}

func expectExpanding(value string) bool {
	rx := regexp.MustCompile("[$]({)?[A-z_0-9]+}?")

	return rx.MatchString(value)
}
