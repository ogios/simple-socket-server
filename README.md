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
The request body should be
```
<type_content_length> (8 bytes)
+
<type>
+
<body_content_length> (8 bytes)
+
<body>
...
```
```python
import socket

HOST = "127.0.0.1"
PORT = 15001

a = """Formal letter writing: block style vs. AMS style..."""
b = b""

def get_len(content: bytes) -> list[int]:
    l: list[int] = []
    total = len(content)
    while total >= 255:
        l.append(total%255)
        total //=255
    l.append(total)
    if len(l) > 8:
        raise Exception("Length over 8")
    while len(l) < 8:
        l.append(0)
    return l

# type
t = "push".encode()
tl = get_len(t)
b = b + bytes(tl) + t

# content
c = a.encode()
cl = get_len(c)
b = b + bytes(cl) + c

print(b)

# send
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((HOST, PORT))
s.send(b)
s.close()
```

## With Proxy Server
Require a server with **public IP address** for transmitting ips and ports of each other to both clients.

In progress...
