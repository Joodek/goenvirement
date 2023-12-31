# Goenv

<div align="center">
	
![goenvirement](https://github.com/Joodek/goenvirement/assets/101285507/9e67c7fa-a361-4670-a440-1b7ac8fa0067)
</div>

this is an easy to use go package to intract with `.env` files and manage envirement variables

## Usage

install the package via :

```bash
go get github.com/Joodek/goenvirement
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

    "github.com/Joodek/goenvirement"
)

func main(){
    goenvirement.Load()

    fmt.Println(os.Getenv("APP_URL")) // http://localhost
    fmt.Println(os.Getenv("APP_PORT")) // 8080
}

```

by default , the load function will try to load a file named `.env` from the current working directory , alternatively you can provide a file path or even multiple files to load

```go

func main(){
    // one file
    goenvirement.Load("path/to/file")

    // multiple files
    goenvirement.Load("path/to/file1","path/to/file2")

  // ...
}

```

the load function will never override the existing variables, even if they're empty
but if you want to them to be overriden , use the `Overload` function instead

```go

func main(){

    goenvirement.Overload()
   // or
    goenvirement.Overload("path/to/file1","path/to/file2")

  // ...
}

```

sometimes you may wish to get those variables back, and take controle of what to do with them
this is possible by using `Read` function, it will read the variables and return them as map

```go

func main(){

    envs := goenvirement.Read()
    // or
    envs := goenvirement.Read("path/to/file1","path/to/file2")

    fmt.println(envs["SOME_VARIABLE"])

  // ...
}

```

In case you have the variables as a string, you can parse them using the `Unmarshal` function like so :

```go

func main(){

   envs := goenvirement.Unmarshal(
		`APP_URL="http://localhost"
		APP_PORT =8080 `,
	)

	fmt.Println("APP_URL    :   ", envs["APP_URL"])
	fmt.Println("APP_PORT   :   ", envs["APP_PORT"])

  // ...
}

```

#### Comments

in all cases , whichever function you use , you can write comments

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

### Hint

Unlike other libraries, you don't have to stick with a specific order, this means that you can have somthing similar to this :

```bash

KEY_1="value1-$KEY_2"
KEY_2="value2-$KEY_3"
KEY_3="value3-somthing"
```

and they will all be expanded correctly

```go

func main(){

   goenvirement.Load()

	fmt.Println(os.Getenv("KEY_1")) // value1-value2-value3-somthing
	fmt.Println(os.Getenv("KEY_2")) // value2-value3-somthing

  // ...
}

```

so far everything looks amazing, but, look at this example :

```bash

KEY_1="value1-$KEY_2"
KEY_2="value2-$KEY_3"
KEY_3="value3-$KEY_1"
```

as you notice here, we are tring to read a key that will never be reached, this example will return the folowing error :

```bash
recursion detected : trying to read KEY_2 by KEY_1
```

## Author

[yassinebenaid](https://github.com/yassinebenaid)
