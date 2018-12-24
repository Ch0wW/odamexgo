# odamexgo
A simple Golang class to query an Odamex server. It currently supports Odamex 0.7.0 servers at the moment (and probably more in the future)

## Installation
```bash
go get github.com/ch0ww/odamexgo
```

## Usage
- First of all, import this class to your project.
```golang
import "github.com/ch0ww/odamexgo"
```

Then, create an ServerQuery class by parsing an Odamex URI :
```golang
odasv, err := odamexgo.NewOdaURI("odamex://<ip>[:<port>]")
if err != nil {
    fmt.Println(err)
    return
}

// Receive and parse all Odamex data.
sv, err := odasv.GetServerInfo()
if err != nil {
    fmt.Println(err)  // In case of a server unreachable
    return
}
```
