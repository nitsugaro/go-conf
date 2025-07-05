## Basic Config Handler

Install

```bash
go get github.com/nitsugaro/go-conf@latest
```

#### LoadConfig

```go
err := goconf.LoadConfig() //default file .config.json

//or

err := goconf.LoadConfig("/path/to/my/config.json")
```

#### Get Config

```go
valStr, err := goconf.GetField[string]("my-key-str")

valStr := goconf.GetOpField[string]("my-key-str", "default-value")

valInt := goconf.GetOpField[int]("my.sub.key-int", 2001) // { "my": { "sub": { "key-int": 2002 } } }

type MyConfig struct {
	X int    `json:"x"`
	Y string `json:"y"`
	Z bool   `json:"z"`
}

val, err := goconf.GetField[MyConfig]("my-custom-obj")
```
