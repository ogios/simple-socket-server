# Simple Socket Server

**made with Golang**


# Config
**Default Path:**  
- Unix: `$XDG_CONFIG_HOME`, if not then `$HOME/.config/`
- Darwin: `$HOME/Library/Application`
- Windows: `%AppData%`

with path: `/transfer-go/base_server.yml`



**Custom Path:**  
Provide args: `-c <your_custom_path>`


# Usage

## Normal Server

### Create Server
```go
import "github.com/ogios/simple-socket-server/server/normal"

func main() {
	server, err := normal.NewSocketServer()
	if err != nil {
		panic(err)
	}
	fmt.Println("server created")
	if err := server.Serv(); err != nil {
		panic(err)
	}
}
```

### Data Transfer
First we add a custom type callback for type `push` and start the server  
In the callback function, we read and print bytes one by one:
```go
server.AddTypeCallback("push", func(conn *normal.Conn) error {
    // read and print every single byte till the end
    fmt.Printf("Type: %s\n", conn.Type)
    for {
        b, readerr := conn.Reader.ReadByte()
        if readerr != nil {
            if readerr.Error() == "EOF" {
                fmt.Print("\n")
                conn.Close()
                return nil
            } else {
                return readerr
            }
        }
        fmt.Printf("%d ", b)
    }
})
server.Serv()
```
And we send a UTF-8 string bytes using python or whatever  
The request body should be `<type>\n<body>`, and there shouldn't be any `\n` in type:
```python
import socket

HOST = "127.0.0.1"
PORT = 15001

a = "push\ndata"
b = a.encode("utf-8")
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

s.connect((HOST, PORT))
s.send(b)
s.close()
```
The outputs are shown in stdout:
```
Type: push
100 97 116 97 
```

## With Proxy Server
Require a server with **public IP address** for transmitting ips and ports of each other to both clients.

In progress...
