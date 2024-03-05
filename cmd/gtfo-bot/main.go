package main

import (
	"fmt"

	"github.com/tandreis/gtfo_telegram_bot_go/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
