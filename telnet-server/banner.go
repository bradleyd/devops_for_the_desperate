package main

import "fmt"

const (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func banner() string {
	b :=
		`
____________ ___________
|  _  \  ___|_   _|  _  \
| | | | |_    | | | | | |
| | | |  _|   | | | | | |
| |/ /| |     | | | |/ /
|___/ \_|     \_/ |___/
`
	return fmt.Sprintf("%s%s%s", colorGreen, b, colorReset)
}
