package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

// LoginModel manages Twitch Device Code login screen state.
type TwitchLoginState int

const (
	Idle TwitchLoginState = iota
	Requesting
	WaitingForAuthorization
	Polling
	Success
	Error
)

// Message types for Bubbletea update loop
// These are strictly returned from Update, never called directly.
type (
	RequestDeviceCodeMsg struct{}
	ReceiveDeviceCodeMsg struct {
		DeviceCode      string
		UserCode        string
		VerificationURI string
		ExpiresIn       int
		Interval        int
		Err             error
	}
	PollTokenMsg    struct{ Attempt int }
	ReceiveTokenMsg struct {
		AccessToken    string
		RefreshToken   string
		TokenExpiresIn int
		Err            error
	}
	TickMsg struct{}
)

// LoginSuccessMsg allows screen transition on Twitch login success.
type LoginSuccessMsg struct {}

// TwitchLoginModel manages Twitch Device Code Flow login screen state.
type LoginModel struct {
	state TwitchLoginState

	// Twitch Device Code Flow fields
	deviceCode      string
	userCode        string
	verificationURI string
	expiresIn       int // seconds left for user to authorize
	interval        int // seconds between polls
	pollCount       int

	// Tokens
	accessToken    string
	refreshToken   string
	tokenExpiresIn int

	// UI
	Width   int
	Height  int
	ErrMsg  string
	Success bool
}

// LoginSuccessMsg allows screen transition on Twitch login success.
// Already defined above, do not redeclare.

func NewLoginModel(width int, height int) *LoginModel {
	// Start in Idle state with just dimensions and blank fields.
	return &LoginModel{
		state:   Idle,
		Width:   width,
		Height:  height,
		ErrMsg:  "",
		Success: false,
	}
}

func (m *LoginModel) Init() tea.Cmd {
	// Start in Idle, no initial commands
	return nil
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == Idle && msg.String() == "enter" {
			m.state = Requesting
			clientID := "8pbsu0inj1huddl1inp1800p4vtmwy" // replace with your actual client ID
			return m, requestDeviceCodeCmd(clientID)
		}
		if m.state == Error && msg.String() == "enter" {
			// Allow retry: reset to idle for user
			m.state = Idle
			m.ErrMsg = ""
			m.pollCount = 0
			m.deviceCode = ""
			m.userCode = ""
			m.verificationURI = ""
			m.expiresIn = 0
			m.interval = 0
			m.accessToken = ""
			m.refreshToken = ""
			m.tokenExpiresIn = 0
			return m, nil
		}
	case ReceiveDeviceCodeMsg:
		if msg.Err != nil {
			m.state = Error
			m.ErrMsg = msg.Err.Error()
			return m, nil
		}
		m.state = WaitingForAuthorization
		m.deviceCode = msg.DeviceCode
		m.userCode = msg.UserCode
		m.verificationURI = msg.VerificationURI
		m.expiresIn = msg.ExpiresIn
		m.interval = msg.Interval
		m.pollCount = 0
		// Start polling after instruction is displayed
		return m, pollTokenCmd(m.deviceCode, m.pollCount, m.expiresIn)
	case TickMsg:
		if m.state == WaitingForAuthorization {
			m.pollCount++
			m.expiresIn = m.expiresIn - m.interval
			return m, pollTokenCmd(m.deviceCode, m.pollCount, m.expiresIn)
		}
	case ReceiveTokenMsg:
		if msg.Err != nil {
			if msg.Err.Error() == "poll_pending" && m.pollCount+1 < 20 && m.expiresIn > 0 {
				// Schedule timer for the next poll
				return m, waitIntervalCmd(time.Duration(m.interval) * time.Second)
			}
			m.state = Error
			m.ErrMsg = msg.Err.Error()
			return m, nil
		}
		// Success!
		m.state = Success
		m.accessToken = msg.AccessToken
		m.refreshToken = msg.RefreshToken
		m.tokenExpiresIn = msg.TokenExpiresIn
		// Transition to main app
		return m, func() tea.Msg { return LoginSuccessMsg{} }
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}
	return m, nil
}

// tea.Cmd to request device code from Twitch
func requestDeviceCodeCmd(clientID string) tea.Cmd {
	return func() tea.Msg {
		endpoint := "https://id.twitch.tv/oauth2/device"
		data := "client_id=" + clientID + "&scope=user:read:email"
		resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(data))
		if err != nil {
			return ReceiveDeviceCodeMsg{Err: err}
		}
		defer resp.Body.Close()
		var out struct {
			DeviceCode      string `json:"device_code"`
			UserCode        string `json:"user_code"`
			VerificationURI string `json:"verification_uri"`
			ExpiresIn       int    `json:"expires_in"`
			Interval        int    `json:"interval"`
		}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&out); err != nil {
			return ReceiveDeviceCodeMsg{Err: err}
		}
		return ReceiveDeviceCodeMsg{
			DeviceCode:      out.DeviceCode,
			UserCode:        out.UserCode,
			VerificationURI: out.VerificationURI,
			ExpiresIn:       out.ExpiresIn,
			Interval:        out.Interval,
			Err:             nil,
		}
	}
}

// tea.Cmd to wait for the polling interval and emit TickMsg
func waitIntervalCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(_ time.Time) tea.Msg { return TickMsg{} })
}

// tea.Cmd to poll for token from Twitch
func pollTokenCmd(deviceCode string, attempt int, expiresIn int) tea.Cmd {
	return func() tea.Msg {
		const maxAttempts = 20
		clientID := "8pbsu0inj1huddl1inp1800p4vtmwy" // replace with your Twitch client ID
		clientSecret := ""                           // Only if your app needs it
		endpoint := "https://id.twitch.tv/oauth2/token"
		grantType := "urn:ietf:params:oauth:grant-type:device_code"
		data := "client_id=" + clientID + "&device_code=" + deviceCode + "&grant_type=" + grantType
		if clientSecret != "" {
			data += "&client_secret=" + clientSecret
		}
		resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(data))
		if err != nil {
			return ReceiveTokenMsg{Err: err}
		}
		defer resp.Body.Close()
		var out struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int    `json:"expires_in"`
			Error        string `json:"error"`
			ErrorDesc    string `json:"error_description"`
			Status       int    `json:"status"`
			Message      string `json:"message"`
		}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&out); err != nil {
			return ReceiveTokenMsg{Err: err}
		}
		if out.AccessToken != "" {
			return ReceiveTokenMsg{
				AccessToken:    out.AccessToken,
				RefreshToken:   out.RefreshToken,
				TokenExpiresIn: out.ExpiresIn,
				Err:            nil,
			}
		}
		// Support both Twitch's response styles: "error" or "message"
		if out.Error == "authorization_pending" || out.Message == "authorization_pending" || out.Error == "slow_down" || out.Message == "slow_down" {
			if attempt+1 >= maxAttempts || expiresIn <= 0 {
				return ReceiveTokenMsg{Err: errTimeout()}
			}
			return ReceiveTokenMsg{Err: errors.New("poll_pending")}
		}
		if out.Error == "expired_token" || out.Message == "expired_token" || out.Error == "access_denied" || out.Message == "access_denied" {
			return ReceiveTokenMsg{Err: errDenied(out.ErrorDesc)}
		}
		return ReceiveTokenMsg{Err: errUnknown(out.Error+"/"+out.Message, out.ErrorDesc)}
	}
}

func errTimeout() error {
	return errors.New("Authorization timed out, please try again.")
}

func errDenied(desc string) error {
	return errors.New("Access denied: " + desc)
}

func errUnknown(err, desc string) error {
	return errors.New("Unknown Twitch error: " + err + ": " + desc)
}

func (m *LoginModel) View() tea.View {
	var content string
	lines := 5

	switch m.state {
	case Idle:
		// Initial login screen
		title := AppTitle("Login with Twitch", m.Width)
		button := "[ Login ]"
		buttonPad := (m.Width - len(button)) / 2
		if buttonPad < 0 {
			buttonPad = 0
		}
		buttonLine := strings.Repeat(" ", buttonPad) + button
		footer := FooterStyle.Render("Press Enter to log in with Twitch.  Ctrl+C: Quit")
		content = title + "\n" + buttonLine + "\n\n" + footer
		lines = 5
	case Requesting:
		// Show loading
		title := AppTitle("Logging in…", m.Width)
		loading := "Requesting device code from Twitch..."
		loadingPad := (m.Width - len(loading)) / 2
		if loadingPad < 0 {
			loadingPad = 0
		}
		loadingLine := strings.Repeat(" ", loadingPad) + loading
		content = title + "\n" + loadingLine
		lines = 4
	case WaitingForAuthorization:
		// Show instructions with polling progress
		title := AppTitle("Authorize Twitch Device", m.Width)
		codeLine := "Code: " + m.userCode
		urlLine := "URL: " + m.verificationURI
		attemptLine := "Polling attempt: " + fmt.Sprintf("%d/20", m.pollCount+1)
		codePad := (m.Width - len(codeLine)) / 2
		urlPad := (m.Width - len(urlLine)) / 2
		attemptPad := (m.Width - len(attemptLine)) / 2
		if codePad < 0 {
			codePad = 0
		}
		if urlPad < 0 {
			urlPad = 0
		}
		if attemptPad < 0 {
			attemptPad = 0
		}
		instructions := "Open the URL above in your browser and enter the code to authorize."
		instPad := (m.Width - len(instructions)) / 2
		if instPad < 0 {
			instPad = 0
		}
		content = title + "\n" +
			strings.Repeat(" ", codePad) + codeLine + "\n" +
			strings.Repeat(" ", urlPad) + urlLine + "\n" +
			strings.Repeat(" ", attemptPad) + attemptLine + "\n\n" +
			strings.Repeat(" ", instPad) + instructions
		lines = 8
	case Error:
		title := AppTitle("Twitch Login Error", m.Width)
		errLine := RenderError(m.ErrMsg)
		errPad := (m.Width - len(m.ErrMsg)) / 2
		if errPad < 0 {
			errPad = 0
		}
		footer := FooterStyle.Render("Press Enter to retry.  Ctrl+C: Quit")
		content = title + "\n" + strings.Repeat(" ", errPad) + errLine + "\n\n" + footer
		lines = 6
	case Success:
		title := AppTitle("Twitch Login Successful!", m.Width)
		content = title + "\n" + "Success!"
		lines = 4
	}

	pad := 0
	if m.Height > lines {
		pad = (m.Height - lines) / 2
	}

	view := tea.NewView(strings.Repeat("\n", pad) + content)
	view.AltScreen = true

	return view
}
