# Goenv

this is an easy to use go package to intract with `.env` files and manage envirement variables

## Usage

install the package via :

```bash
    go get github.com/Joodek/goenv
```

set your variables in a `.env` file

```bash
APP_URL=http://localhost
APP_PORT=8080

```

then call the `Load` function

```go
package main

import (
    "fmt"
    "os"

    "github.com/Joodek/goenv"
)

func main(){
    goenv.Load()

    fmt.Println(os.Getenv("APP_URL")) // http://localhost
    fmt.Println(os.Getenv("APP_PORT")) // 8080
}

```

by default , the load function will try to load a file named `.env` from the current working directory , you can for sure change this behavior by providing a file name or even multiple files to load

```go

func main(){
    // one file
    goenv.Load("path/to/file")

    // multiple files
    goenv.Load("path/to/file1","path/to/file2")

  // ...
}

```

the load function will never override the existing variables, even if they're empty
but if you want to them to be overriden , use the `Overload` function instead

```go

func main(){

    goenv.Overload()
   // or
    goenv.Overload("path/to/file1","path/to/file2")

  // ...
}

```

sometimes you may wish to get those variables back, and take controle of what to do with them
this is possible by using `Read` function, it will read the variables and return them as map

```go

func main(){

    envs := goenv.Read()
    // or
    envs := goenv.Read("path/to/file1","path/to/file2")

    fmt.println(envs["SOME_VARIABLE"])

  // ...
}

```

sometimes you may have the variables as a string, not in a file, you can parse them using the `Unmarshal` function like so :

```go

func main(){

   envs := goenv.Unmarshal(
		`APP_URL="http://localhost"
		APP_PORT =8080 `,
	)

	fmt.Println("APP_URL    :   ", envs["APP_URL"])
	fmt.Println("APP_PORT   :   ", envs["APP_PORT"])

  // ...
}

```

#### Comments

in all cases , whichever a function you use , you can write comments

```bash
# this is valid
APP_URL=http://localhost

APP_PORT=8080 # this is also valid

```

#### Variables

you can use variables to represent other keys in your envirement, like you do in a bash script , look at the following

```bash
# you can use local evirements ,
APP_URL="http://localhost${APP_PORT}" # http://localhost:8080
APP_PORT=8080

# or you may want to use the system envirements as well
USER_CACHE="$HOME/programs/cache" # /home/yassinebenaid/programs/cache

```
