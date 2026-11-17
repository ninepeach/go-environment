# go-environment

A lightweight, thread-safe key-value registry for Go applications, with an optional process-wide singleton.

## Installation

```bash
go get github.com/ninepeach/go-environment
```

## Usage

```go
package main

import (
	"fmt"

	goenv "github.com/ninepeach/go-environment"
)

func main() {
	env := goenv.GetInstance()
	env.Set("AppName", "MyGoApp")
	env.Set("MaxRetries", 3)

	appName := env.Get("AppName").(string)
	maxRetries := env.Get("MaxRetries").(int)
	fmt.Println("AppName:", appName)
	fmt.Println("MaxRetries:", maxRetries)

	if value, ok := env.Lookup("AppName"); ok {
		fmt.Println("Found:", value)
	}
	env.Delete("MaxRetries")
}
```

Use `New()` for an independent registry. The zero value `var env goenv.Environment` is also usable.

- `Get(key)` returns `nil` for a missing key. `Lookup(key)` distinguishes a missing key from a stored `nil`.
- `Delete(key)` removes a key and is safe for missing keys.
- `Clone()` makes a **shallow copy**: the registries have independent key maps, but stored pointers, slices, maps, and other reference-type values remain shared.
- Registry operations are synchronized; **mutable stored values are not automatically thread-safe**.
- `GetFieldValue` reads exported struct fields and returns errors for invalid input or inaccessible fields.
