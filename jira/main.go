package main

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/logger"
)

func main() {
	logger.PrepareLogger()
	log.Debug("Hello, World!")
	fmt.Println("Hello, World!")
}
