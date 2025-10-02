# Haar

[![license](https://img.shields.io/badge/license-MIT-red.svg)](LICENSE)
[![go version](https://img.shields.io/github/go-mod/go-version/TuanKiri/haar)](go.mod)
[![go report](https://goreportcard.com/badge/github.com/TuanKiri/haar)](https://goreportcard.com/report/github.com/TuanKiri/haar)

This project provides a Golang implementation of the Haar 2D transform, following an approach similar to [imgSeek](https://sourceforge.net/projects/imgseek). The library can generate image hashes for indexing and similarity search in [iqdb](https://github.com/danbooru/iqdb).

> [!WARNING]
> Output from the Go package shows small deviations compared to the [C reference](https://github.com/danbooru/iqdb/blob/master/src/haar.cpp).
> This does not affect the final result, and all images are successfully found. Refer to [tests](haar_test.go) for specifics.

## Installation

```shell
go get github.com/TuanKiri/haar
```

## Quick Start

Create your `.go` file. For example: `main.go`.

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TuanKiri/haar"
)

func main() {
	image, err := os.ReadFile("your_image_file.jpg")
	if err != nil {
		log.Fatal(err)
	}

	signature, err := haar.SignatureFromBlob(image)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(signature)
}
```

Run your program:

```shell
go run main.go
```

## Motivation

Most booru sites need an image file to search. [Danbooru](https://github.com/danbooru) can use an image's hash instead, which is much faster.
