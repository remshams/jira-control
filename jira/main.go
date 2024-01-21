package main

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/logger"
)

func main() {
	logger.PrepareLogger()
	log.Debug("Hello, World!")
	fmt.Println("Hello, World!")
}
