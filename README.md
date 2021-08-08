# Dreamcatcher

Dreamcatcher combines a tcp connection listener (eg. reverse shell listener) and an HTTP file server (ie. similar to python SimpleHTTPServer).

This means you can listen for http requests and reverse shells on the same TCP port.

Also reduces the amount of tab, tmux panes, windows, etc.. that you need since everything is served on the same port.

This tool is provided as is. YMMV.

> It's like running an infinite number of netcat listeners and a python HTTP server at once. This is brilliant!
> - Emon Lusk 


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

```
$ ./dreamcatcher -d /tmp
:: Catcher listening on 127.0.0.1:4444
:: Serving directory /tmp
:: HTTP request for file: /tmp/test.txt
[127.0.0.1:35776] :: $ whoami
[127.0.0.1:35776] :: root
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

- fix bug when no data is sent (ie. add timeout to consider it a reverse shell if nothing happens after x seconds)
- add menu with prefix
    - prefix ":"
    - action: help, list/ls, switch/s \<id based on list\>
- add a way to switch between reverse shells (ie. via menu action)
