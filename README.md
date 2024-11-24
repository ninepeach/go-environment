# go-environment
Global Variable Management

The `go-environment` package provides a singleton pattern for managing global variables in Go applications. It allows for centralized registration, storage, and retrieval of environment configurations, application-wide variables, and constants. This package ensures that only one instance of the global variable registry exists, simplifying management and access across your Go project.

## Installation

You can install the `go-environment` package using `go get`:

```bash
go get github.com/ninepeach/go-environment

```

##Usage

The package uses a singleton instance to manage global variables. Here is an example of how to use it:

```go
package main

import "github.com/ninepeach/go-environment"

func main() {
    // Register global variables
    env := goenv.GetInstance()
    env.Set("AppName", "MyGoApp")
    env.Set("MaxRetries", 3)

    // Access global variables
    appName := env.Get("AppName").(string)
    maxRetries := env.Get("MaxRetries").(int)

    // Print values
    fmt.Println("AppName:", appName)
    fmt.Println("MaxRetries:", maxRetries)
}

```

