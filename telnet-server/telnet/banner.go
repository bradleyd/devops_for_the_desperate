package telnet

import "fmt"

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func banner() string {
	b := "____________ ___________\r\n|  _  \\  ___|_   _|  _  \\\r\n| | | | |_    | | | | | |\r\n| | | |  _|   | | | | | |\r\n| |/ /| |     | | | |/ /\r\n|___/ \\_|     \\_/ |___/\r\n"
	return fmt.Sprintf("%s%s%s", colorYellow, b, colorReset)
}
