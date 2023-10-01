package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AntonPashechko/gophkeeper/internal/client/config"
	"github.com/AntonPashechko/gophkeeper/internal/client/sender"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"

	reader *bufio.Reader
)

//go build -ldflags="-X 'main.Version=v1.0.0' -X 'app/build.Time=$(date)'"

func readLine(title string) string {
	fmt.Print(fmt.Sprintf("Enter %s:", title))
	line, _ := reader.ReadString('\n')
	line = strings.TrimSuffix(line, "\r\n")
	return line
}

func main() {

	reader = bufio.NewReader(os.Stdin)

	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("cannot load config: %s\n", err)
	}

	sender := sender.NewSender(cfg)
	if sender.Init() != nil {
		log.Fatalf("cannot initialize sender: %s\n", err)
	}

	for {
		cmd := readLine(`command`)

		switch cmd {
		case `register`:
			login := readLine(`login`)
			password := readLine(`password`)

			fmt.Print(login, password)
		}
	}

}
