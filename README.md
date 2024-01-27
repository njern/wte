# wte

This library provides an implementation of the Whittaker-Eilers smoothing algorithm in Go. It is designed to smooth time series data or other sequences, offering support for both equally and non-equally spaced data. The library is particularly useful for reducing noise in data while preserving important signals.

## Features

- Whittaker-Eilers smoothing for sequences of `float64`.
- Support for both equally and non-equally spaced data.
- Customizable smoothing parameters.

## Installation

To use this library, first ensure you have Go installed on your system. Then, you can install the library using `go get`:

```bash
go get -u github.com/njern/wte
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/njern/wte"
)

func main() {
    data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    lambda := 10.0
    order := 2

    // For equally spaced data, pass nil for spacing
    smoothed, err := wte.Smooth(data, lambda, order, nil)
    if err != nil {
        panic(err)
    }

    fmt.Println("Smoothed data:", smoothed)
}
```

## Contributing

Contributions to improve this library are welcome. Feel free to fork the repository, make your changes, and submit a pull request.

## License

This library is licensed under the [MIT License](LICENSE).

## Acknowledgements

- I was inspired to develop this library after reading [The perfect way to smooth your noisy data](https://towardsdatascience.com/the-perfect-way-to-smooth-your-noisy-data-4f3fe6b44440) by [Andrew Bowell](https://github.com/AnBowell) and used his (excellent) [Rust library](https://github.com/AnBowell/whittaker-eilers) to generate the test cases.
- This library was developed using the [gonum](https://github.com/gonum/gonum) package for Go matrix operations. Special thanks to the Gonum contributors.