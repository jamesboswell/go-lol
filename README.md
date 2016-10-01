# go-lol [![GoDoc](https://godoc.org/github.com/kdy1997/go-lol?status.svg)](godoc)
A new, generated golang client for rito apis.

# Installation
```sh
go get -u github.com/kdy1997/go-lol
```

# Features
 - Clean API. See [godoc][godoc]
   - No global variable.
 - [context](https://godoc.org/golang.org/x/net/context) support.
 - Google app engine support (*http.Client from context.Context)


# FAQ
If your question is not listed here, please feel free to make an issue for it.

## Why do you generate instead of writing it by hand?
I wrote a html parser to practice html parsing and to deal with api change in more simpler way.


# Contributing
 - You need go 1.6 to run ```go generate```

```go generate github.com/kdy1997/go-lol```



# License
Apache2



[godoc]:(https://godoc.org/github.com/kdy1997/go-lol)
