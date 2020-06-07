# Base64 Site Golang edition

## What is this thing?

Have you ever wanted to encode a tiny website in base64 but needed the link to be shared in mobile or otherwise handled by devices that only know about well-formed URLs? Well, me too.

This is a rewrite of [this](//github.com/fidiego/base64-sites) identical project. This version is written in Go instead of python.

Try it out the python version [here](//base64-sites.herokuapp.com).

## Other Info

- Built with [golang](https://golang.org/). No external dependencies.
- HTML formatted with [tidy](http://www.html-tidy.org/).
- JS formatted with [prettier](https://prettier.io/).

## Dev Setup

**Make a `.env` file**

```
cp env.default .env
```

**Run**

```
go run main.go
```

**Try it out\***

```
http POST localhost:9876/api content="<marquee>hello world</marquee>" --json
```

## TODO

- [ ] figure out css for `/render` endpoint. Get rid of the excess space so iframe fills space with no scroll on body or html.
- [ ] make a js only version that can be deployed on netlify
