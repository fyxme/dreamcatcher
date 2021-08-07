# Dreamcatcher

Dreamcatcher combines a tcp connection listener (eg. reverse shell listener) and an HTTP file server (ie. similar to python SimpleHTTPServer).

This means you can listen for http requests and reverse shells on the same TCP port.

Also reduces the amount of tab, tmux panes, windows, etc.. that you need since everything is served on the same port.

This tool is provided as is. YMMV.

> It's like running an infinite number of netcat listeners and a python HTTP server at once. This is brilliant! - Emon Lusk 


## Compiling

```
go build dreamcatcher.go
```

## Usage

```
$ ./dreamcatcher -h
flag needs an argument: -h
Usage of ./dreamcatcher:
  -d string
        Directory you want to serve files from (default ".")
  -h string
        Host to listen on (default "127.0.0.1")
  -p int
        Port to listen on (default 4444)
```

## Known bugs and future improvements

It can currently catch multiple reverse shells at once but it doesnt handle them properly. Instead it will switch between them:
```
$ go run dreamcatcher.go
:: Catcher listening on 127.0.0.1:4444
:: Serving directory .
[127.0.0.1:35686] :: $ [127.0.0.1:35688] :: $
[127.0.0.1:35686] :: $ [127.0.0.1:35690] :: $
[127.0.0.1:35688] :: $
[127.0.0.1:35686] :: $
[127.0.0.1:35690] :: $
[127.0.0.1:35688] :: $
[127.0.0.1:35686] :: $
```

### Todo

- add menu with prefix
    - prefix ":"
    - action: help, list/ls, switch/s \<id based on list\>
- add a way to switch between reverse shells (ie. via menu action)
