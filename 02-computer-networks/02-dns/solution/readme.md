Three solutions to the problem:

- Using Go's "net" module: `go run dns.go using_net.go wikipedia.org A`
- Using Go's "syscall" module: `go run dns.go using_syscalls.go wikipedia.org A`
- Using Python: `python3 simple_dns.py wikipedia.org A`

If you'd like to extract `dns.go` into its own module, see the instructions in [this guide](https://go.dev/doc/tutorial/call-module-code).
