package telegram

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tegal1337/telegram-cli/internal/config"
	"github.com/zelenin/go-tdlib/client"
)

// AuthState represents the current state of the authorization flow.
type AuthState int

const (
	AuthStateWaitPhone AuthState = iota
	AuthStateWaitCode
	AuthStateWaitPassword
	AuthStateWaitQR
	AuthStateReady
	AuthStateClosed
)

// TUIAuthorizer implements client.AuthorizationStateHandler for the TUI.
// It sends auth state changes as tea.Msg through a channel that the
// auth UI component reads.
type TUIAuthorizer struct {
	tdlibParams *client.SetTdlibParametersRequest
	phoneCh     chan string
	codeCh      chan string
	passwordCh  chan string
	phone       string
}

// NewTUIAuthorizer creates a new authorizer for the TUI auth flow.
func NewTUIAuthorizer(cfg *config.Config) *TUIAuthorizer {
	return &TUIAuthorizer{
		tdlibParams: TdlibParameters(cfg),
		phoneCh:     make(chan string, 1),
		codeCh:      make(chan string, 1),
		passwordCh:  make(chan string, 1),
		phone:       cfg.Telegram.Phone,
	}
}

// SubmitPhone sends the phone number to the auth flow.
func (a *TUIAuthorizer) SubmitPhone(phone string) {
	a.phoneCh <- phone
}

// SubmitCode sends the verification code to the auth flow.
func (a *TUIAuthorizer) SubmitCode(code string) {
	a.codeCh <- code
}

// SubmitPassword sends the 2FA password to the auth flow.
func (a *TUIAuthorizer) SubmitPassword(password string) {
	a.passwordCh <- password
}

// Handle implements client.AuthorizationStateHandler.
func (a *TUIAuthorizer) Handle(c *client.Client, state client.AuthorizationState) error {
	switch s := state.(type) {
	case *client.AuthorizationStateWaitTdlibParameters:
		_, err := c.SetTdlibParameters(a.tdlibParams)
		return err

	case *client.AuthorizationStateWaitPhoneNumber:
		phone := a.phone
		if phone == "" {
			phone = <-a.phoneCh
		}
		_, err := c.SetAuthenticationPhoneNumber(&client.SetAuthenticationPhoneNumberRequest{
			PhoneNumber: phone,
			Settings: &client.PhoneNumberAuthenticationSettings{
				AllowFlashCall:       false,
				AllowMissedCall:      false,
				IsCurrentPhoneNumber: false,
			},
		})
		return err

	case *client.AuthorizationStateWaitCode:
		code := <-a.codeCh
		_, err := c.CheckAuthenticationCode(&client.CheckAuthenticationCodeRequest{
			Code: code,
		})
		return err

	case *client.AuthorizationStateWaitPassword:
		_ = s // contains password hint
		password := <-a.passwordCh
		_, err := c.CheckAuthenticationPassword(&client.CheckAuthenticationPasswordRequest{
			Password: password,
		})
		return err

	case *client.AuthorizationStateReady:
		return nil

	case *client.AuthorizationStateClosed:
		return nil

	default:
		return fmt.Errorf("unexpected auth state: %T", state)
	}
}

// CLIAuthorizer creates a simple CLI-based authorizer for testing.
func CLIAuthorizer(cfg *config.Config) client.AuthorizationStateHandler {
	authorizer := client.ClientAuthorizer(TdlibParameters(cfg))
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-authorizer.PhoneNumber:
				fmt.Print("Enter phone number: ")
				phone, _ := reader.ReadString('\n')
				authorizer.PhoneNumber <- strings.TrimSpace(phone)
			case <-authorizer.Code:
				fmt.Print("Enter code: ")
				code, _ := reader.ReadString('\n')
				authorizer.Code <- strings.TrimSpace(code)
			case <-authorizer.Password:
				fmt.Print("Enter 2FA password: ")
				pw, _ := reader.ReadString('\n')
				authorizer.Password <- strings.TrimSpace(pw)
			}
		}
	}()
	return authorizer
}
