package telegram

import (
	"fmt"

	"github.com/tegal1337/telegram-cli/internal/config"
	"github.com/zelenin/go-tdlib/client"
)

type AuthState int

const (
	AuthStateWaitPhone AuthState = iota
	AuthStateWaitCode
	AuthStateWaitPassword
	AuthStateWaitQR
	AuthStateReady
	AuthStateClosed
)

// AuthStateCallback is called when the auth state changes.
// Used to notify the TUI about state transitions.
type AuthStateCallback func(AuthState, string)

type TUIAuthorizer struct {
	tdlibParams *client.SetTdlibParametersRequest
	phoneCh     chan string
	codeCh      chan string
	passwordCh  chan string
	phone       string
	onState     AuthStateCallback
}

func NewTUIAuthorizer(cfg *config.Config) *TUIAuthorizer {
	return &TUIAuthorizer{
		tdlibParams: TdlibParameters(cfg),
		phoneCh:     make(chan string, 1),
		codeCh:      make(chan string, 1),
		passwordCh:  make(chan string, 1),
		phone:       cfg.Telegram.Phone,
	}
}

// SetStateCallback sets the callback for auth state changes.
func (a *TUIAuthorizer) SetStateCallback(cb AuthStateCallback) {
	a.onState = cb
}

func (a *TUIAuthorizer) notifyState(state AuthState, hint string) {
	if a.onState != nil {
		a.onState(state, hint)
	}
}

func (a *TUIAuthorizer) SubmitPhone(phone string) {
	a.phoneCh <- phone
}

func (a *TUIAuthorizer) SubmitCode(code string) {
	a.codeCh <- code
}

func (a *TUIAuthorizer) SubmitPassword(password string) {
	a.passwordCh <- password
}

func (a *TUIAuthorizer) Handle(c *client.Client, state client.AuthorizationState) error {
	switch s := state.(type) {
	case *client.AuthorizationStateWaitTdlibParameters:
		_, err := c.SetTdlibParameters(a.tdlibParams)
		return err

	case *client.AuthorizationStateWaitPhoneNumber:
		a.notifyState(AuthStateWaitPhone, "")
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
		a.notifyState(AuthStateWaitCode, "")
		code := <-a.codeCh
		_, err := c.CheckAuthenticationCode(&client.CheckAuthenticationCodeRequest{
			Code: code,
		})
		return err

	case *client.AuthorizationStateWaitPassword:
		hint := s.PasswordHint
		a.notifyState(AuthStateWaitPassword, hint)
		password := <-a.passwordCh
		_, err := c.CheckAuthenticationPassword(&client.CheckAuthenticationPasswordRequest{
			Password: password,
		})
		return err

	case *client.AuthorizationStateReady:
		a.notifyState(AuthStateReady, "")
		return nil

	case *client.AuthorizationStateClosed:
		a.notifyState(AuthStateClosed, "")
		return nil

	default:
		return fmt.Errorf("unexpected auth state: %T", state)
	}
}

func (a *TUIAuthorizer) Close() {
	// Don't close channels — they may still be in use by the TUI.
}
