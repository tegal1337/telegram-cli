package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/tegal1337/telegram-cli/internal/app"
	"github.com/tegal1337/telegram-cli/internal/config"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/telegram"
)

func main() {
	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate required config.
	if cfg.Telegram.APIID == 0 || cfg.Telegram.APIHash == "" {
		fmt.Println("Telegram API credentials not configured.")
		fmt.Println("")
		fmt.Println("1. Get credentials from https://my.telegram.org/apps")
		fmt.Println("2. Copy config.example.toml to ~/.config/tele-tui/config.toml")
		fmt.Println("3. Fill in api_id and api_hash")
		fmt.Println("")
		fmt.Println("Or set TELETUI_CONFIG to your config file path.")
		os.Exit(1)
	}

	// Create in-memory store.
	s := store.NewStore()

	// Create TUI authorizer.
	authorizer := telegram.NewTUIAuthorizer(cfg)

	// Create TDLib client.
	tgClient, err := telegram.NewClient(cfg, authorizer)
	if err != nil {
		log.Fatalf("Failed to create Telegram client: %v", err)
	}

	// Create root model.
	model := app.New(cfg, tgClient, s, authorizer)

	// Create bubbletea program.
	p := tea.NewProgram(model)

	// Start TDLib update listener — bridges TDLib updates to bubbletea.
	listener := telegram.NewListener(tgClient.TD(), p)
	listener.Start()

	// Handle graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		tgClient.Close()
		os.Exit(0)
	}()

	// Run the TUI.
	if _, err := p.Run(); err != nil {
		tgClient.Close()
		log.Fatalf("Error running TUI: %v", err)
	}

	tgClient.Close()
}
