package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/tegal1337/telegram-cli/internal/app"
	"github.com/tegal1337/telegram-cli/internal/config"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/telegram"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Telegram.APIID == 0 || cfg.Telegram.APIHash == "" {
		if err := setupWizard(cfg); err != nil {
			log.Fatalf("Setup failed: %v", err)
		}
	}

	s := store.NewStore()
	authorizer := telegram.NewTUIAuthorizer(cfg)

	// Start TDLib client in background — it blocks on auth.
	tgClient := telegram.NewClientAsync(cfg, authorizer)

	// Create root model.
	model := app.New(cfg, tgClient, s, authorizer)

	// Create bubbletea program.
	p := tea.NewProgram(model)

	// Wire auth state changes into bubbletea via p.Send().
	authorizer.SetStateCallback(func(state telegram.AuthState, hint string) {
		p.Send(app.AuthStateChangedMsg{State: int(state), Hint: hint})
	})

	// Once client is ready, start the update listener.
	go func() {
		tgClient.WaitReady()
		td := tgClient.TD()
		if td != nil {
			listener := telegram.NewListener(td, p)
			listener.Start()
		}
		// Notify TUI that we're authenticated.
		me, err := tgClient.GetMe()
		if err == nil && me != nil {
			p.Send(app.AuthenticatedMsg{
				UserId:    me.Id,
				FirstName: me.FirstName,
				LastName:  me.LastName,
			})
		}
	}()

	// Graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		tgClient.Close()
		os.Exit(0)
	}()

	if _, err := p.Run(); err != nil {
		tgClient.Close()
		log.Fatalf("Error running TUI: %v", err)
	}

	tgClient.Close()
}

func setupWizard(cfg *config.Config) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════════╗")
	fmt.Println("  ║         Telegram CLI - First Run         ║")
	fmt.Println("  ╚══════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  Get your API credentials from:")
	fmt.Println("  https://my.telegram.org/apps")
	fmt.Println()

	for {
		fmt.Print("  Enter API ID: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		id, err := strconv.Atoi(input)
		if err != nil || id <= 0 {
			fmt.Println("  Invalid API ID. Must be a number.")
			continue
		}
		cfg.Telegram.APIID = int32(id)
		break
	}

	for {
		fmt.Print("  Enter API Hash: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if len(input) < 10 {
			fmt.Println("  Invalid API Hash. Too short.")
			continue
		}
		cfg.Telegram.APIHash = input
		break
	}

	fmt.Print("  Enter phone number (optional, press Enter to skip): ")
	phone, _ := reader.ReadString('\n')
	phone = strings.TrimSpace(phone)
	if phone != "" {
		cfg.Telegram.Phone = phone
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Println()
	fmt.Println("  Config saved! Starting Telegram CLI...")
	fmt.Println()

	return nil
}
