// Package email provides functionalities for configuring and sending emails.
package mailer

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	ht "html/template"
	"log/slog" // Import the slog package
	tt "text/template"

	"go-boilerplate/internal/config"

	"github.com/wneessen/go-mail"
)

// ContentConfig holds the actual content and sender/subject for an email.
type ContentConfig struct {
	sender       string
	subject      string
	textBody     string
	htmlBody     string
	textTemplate *tt.Template
	htmlTemplate *ht.Template
	templateData interface{}
}

// Option defines a function type that modifies a ContentConfig.
type Option func(*ContentConfig) error

// --- Content Configuration Options ---

func WithSender(sender string) Option {
	return func(c *ContentConfig) error {
		if sender == "" {
			return errors.New("email sender cannot be empty")
		}
		c.sender = sender
		return nil
	}
}

func WithSubject(subject string) Option {
	return func(c *ContentConfig) error {
		c.subject = subject
		return nil
	}
}

func WithTextBody(body string) Option {
	return func(c *ContentConfig) error {
		if c.textTemplate != nil {
			return errors.New("cannot use WithTextBody and WithTextTemplate simultaneously")
		}
		c.textBody = body
		return nil
	}
}

func WithHTMLBody(body string) Option {
	return func(c *ContentConfig) error {
		if c.htmlTemplate != nil {
			return errors.New("cannot use WithHTMLBody and WithHTMLTemplate simultaneously")
		}
		c.htmlBody = body
		return nil
	}
}

func WithTextTemplate(tmpl *tt.Template, data interface{}) Option {
	return func(c *ContentConfig) error {
		if tmpl == nil {
			return errors.New("text template cannot be nil")
		}
		if c.textBody != "" {
			return errors.New("cannot use WithTextTemplate and WithTextBody simultaneously")
		}
		c.textTemplate = tmpl
		c.templateData = data
		return nil
	}
}

func WithHTMLTemplate(tmpl *ht.Template, data interface{}) Option {
	return func(c *ContentConfig) error {
		if tmpl == nil {
			return errors.New("html template cannot be nil")
		}
		if c.htmlBody != "" {
			return errors.New("cannot use WithHTMLTemplate and WithHTMLBody simultaneously")
		}
		c.htmlTemplate = tmpl
		c.templateData = data
		return nil
	}
}

// --- Service Definition ---
// EmailService provides methods for sending emails.
type EmailService struct {
	appConfig  config.AppConfig
	authConfig config.AuthConfig
	smtpConfig config.SMTPConfig
	Logger     *slog.Logger // Injectable logger
}

// NewEmailService creates a new email service with the given SMTP configuration and logger.
// If logger is nil, a default slog logger (discarding output) will be used to prevent nil panics.
func NewEmailService(cfg config.Settings, logger *slog.Logger) (*EmailService, error) {
	smtpConfig := cfg.Email
	if smtpConfig.Host == "" {
		return nil, errors.New("SMTP host cannot be empty")
	}
	if smtpConfig.Port <= 0 {
		return nil, errors.New("SMTP port must be a positive integer")
	}

	// Handle nil logger case: use a discard logger to avoid panics
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), &slog.HandlerOptions{Level: slog.LevelError + 100})) // Effectively discards
		// Or, more simply if you have a common "discard" logger in your app:
		// logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	return &EmailService{
		appConfig:  cfg.App,
		authConfig: cfg.Auth,
		smtpConfig: smtpConfig,
		Logger:     logger.With("service", "email"), // Add a common attribute for this service
	}, nil
}

// processContentConfig creates and validates a ContentConfig from options.
func (s *EmailService) processContentConfig(options ...Option) (*ContentConfig, error) {
	contentCfg := &ContentConfig{}

	for _, opt := range options {
		if err := opt(contentCfg); err != nil {
			return nil, fmt.Errorf("failed to apply email option: %w", err)
		}
	}

	if contentCfg.textTemplate != nil {
		var tBuf bytes.Buffer
		if err := contentCfg.textTemplate.Execute(&tBuf, contentCfg.templateData); err != nil {
			return nil, fmt.Errorf("failed to execute text template: %w", err)
		}
		contentCfg.textBody = tBuf.String()
	}

	if contentCfg.htmlTemplate != nil {
		var hBuf bytes.Buffer
		if err := contentCfg.htmlTemplate.Execute(&hBuf, contentCfg.templateData); err != nil {
			return nil, fmt.Errorf("failed to execute html template: %w", err)
		}
		contentCfg.htmlBody = hBuf.String()
	}

	if contentCfg.sender == "" {
		return nil, errors.New("email sender is required (use WithSender)")
	}
	if contentCfg.subject == "" {
		s.Logger.Warn("Email subject is empty.") // Using slog
	}
	if contentCfg.textBody == "" && contentCfg.htmlBody == "" {
		return nil, errors.New("email content is empty; at least one body (text or HTML) is required")
	}

	return contentCfg, nil
}

// Send sends an email to one or more recipients.
func (s *EmailService) Send(recipients []string, options ...Option) error {
	if len(recipients) == 0 {
		return errors.New("at least one recipient is required")
	}
	for _, r := range recipients {
		if r == "" {
			return errors.New("recipient email address cannot be empty")
		}
	}

	contentCfg, err := s.processContentConfig(options...)
	if err != nil {
		return fmt.Errorf("failed to process email content configuration: %w", err)
	}

	m := mail.NewMsg()
	if err := m.From(contentCfg.sender); err != nil {
		return fmt.Errorf("failed to set sender <%s>: %w", contentCfg.sender, err)
	}
	if err := m.To(recipients...); err != nil {
		return fmt.Errorf("failed to set recipients %v: %w", recipients, err)
	}
	m.Subject(contentCfg.subject)

	// Body content handling (this part looks good)
	hasTextBody := contentCfg.textBody != ""
	hasHTMLBody := contentCfg.htmlBody != ""

	if hasTextBody && hasHTMLBody {
		m.SetBodyString(mail.TypeTextPlain, contentCfg.textBody)
		m.AddAlternativeString(mail.TypeTextHTML, contentCfg.htmlBody)
	} else if hasTextBody {
		m.SetBodyString(mail.TypeTextPlain, contentCfg.textBody)
	} else if hasHTMLBody {
		m.SetBodyString(mail.TypeTextHTML, contentCfg.htmlBody)
	}

	// Build client options with proper logic
	clientOpts := []mail.Option{
		mail.WithPort(s.smtpConfig.Port),
	}

	// Add authentication if credentials are provided
	if s.smtpConfig.Username != "" {
		clientOpts = append(clientOpts,
			mail.WithSMTPAuth(mail.SMTPAuthPlain), // Use Plain auth for most servers
			mail.WithUsername(s.smtpConfig.Username),
			mail.WithPassword(s.smtpConfig.Password),
		)
	}

	// Configure TLS/SSL based on port and config
	if s.smtpConfig.UseSSL || s.smtpConfig.Port == 465 {
		// Implicit SSL/TLS (port 465)
		clientOpts = append(clientOpts, mail.WithSSL())
		s.Logger.Info("Using implicit SSL/TLS (port 465 style)")
	} else if s.smtpConfig.UseTLS || s.smtpConfig.Port == 587 {
		// STARTTLS (port 587)
		clientOpts = append(clientOpts, mail.WithTLSPolicy(mail.TLSMandatory))
		s.Logger.Info("Using STARTTLS (port 587 style)")
	} else {
		// Plain connection (port 25 or custom) - use opportunistic TLS
		clientOpts = append(clientOpts, mail.WithTLSPolicy(mail.TLSOpportunistic))
		s.Logger.Info("Using opportunistic TLS (port 25 style)")
	}

	// Add timeout if configured
	if s.smtpConfig.Timeout > 0 {
		clientOpts = append(clientOpts, mail.WithTimeout(s.smtpConfig.Timeout))
	}

	s.Logger.Info("Creating SMTP client",
		slog.String("host", s.smtpConfig.Host),
		slog.Int("port", s.smtpConfig.Port),
		slog.Bool("use_ssl", s.smtpConfig.UseSSL),
		slog.Bool("use_tls", s.smtpConfig.UseTLS),
		slog.Bool("has_auth", s.smtpConfig.Username != ""),
	)

	client, err := mail.NewClient(s.smtpConfig.Host, clientOpts...)
	if err != nil {
		s.Logger.Error("Failed to create SMTP client",
			slog.String("host", s.smtpConfig.Host),
			slog.Int("port", s.smtpConfig.Port),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to create SMTP client for %s:%d: %w", s.smtpConfig.Host, s.smtpConfig.Port, err)
	}
	defer client.Close() // Always close the client

	s.Logger.Info("Attempting to send email",
		slog.String("from", contentCfg.sender),
		slog.Any("to", recipients),
		slog.String("subject", contentCfg.subject),
		slog.String("smtp_host", s.smtpConfig.Host),
		slog.Int("smtp_port", s.smtpConfig.Port),
	)

	if err := client.DialAndSend(m); err != nil {
		s.Logger.Error("Failed to dial and send email",
			slog.String("from", contentCfg.sender),
			slog.Any("to", recipients),
			slog.String("smtp_host", s.smtpConfig.Host),
			slog.Int("smtp_port", s.smtpConfig.Port),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to dial and send email: %w", err)
	}

	s.Logger.Info("Email sent successfully",
		slog.String("from", contentCfg.sender),
		slog.Any("to", recipients),
	)
	return nil
}

// SendToSingle is a convenience wrapper around Send for a single recipient.
func (s *EmailService) SendToSingle(recipient string, options ...Option) error {
	if recipient == "" {
		return errors.New("recipient email address cannot be empty")
	}
	return s.Send([]string{recipient}, options...)
}

// LoadTextTemplate loads and parses a text template from the embedded filesystem.
func LoadTextTemplate(filename string) (*tt.Template, error) {
	tmpl, err := tt.ParseFS(EmailTemplateFS, "templates/"+filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse text template %s: %w", filename, err)
	}
	return tmpl, nil
}

// LoadHTMLTemplate loads and parses an HTML template from the embedded filesystem.
func LoadHTMLTemplate(filename string) (*ht.Template, error) {
	tmpl, err := ht.ParseFS(EmailTemplateFS, "templates/"+filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML template %s: %w", filename, err)
	}
	return tmpl, nil
}

//go:embed templates/*
var EmailTemplateFS embed.FS

func (e *EmailService) SendWelcomeEmail(email string, name string, verificationLink string) error {
	welcomeTextTemplate, err := LoadTextTemplate("welcome.txt")
	if err != nil {
		e.Logger.Error("Unable to load welcome email template", "error", err)
		return err
	}

	welcomeHTMLTemplate, err := LoadHTMLTemplate("welcome.html")
	if err != nil {
		e.Logger.Error("Unable to load welcome email template", "error", err)
		return err
	}

	type TemplateData struct {
		Name             string
		VerificationLink string
		SupportEmail     string
		AppName          string
		Timeout          string
	}
	data := TemplateData{
		AppName:          e.appConfig.AppName,
		Name:             name,
		VerificationLink: verificationLink,
		SupportEmail:     e.smtpConfig.DefaultSender,
		Timeout:          e.authConfig.TokenLifetime.String(),
	}

	emailErr := e.SendToSingle(email,
		WithSender(e.smtpConfig.DefaultSender),
		WithSubject("Verify Your Account"),
		WithHTMLTemplate(welcomeHTMLTemplate, data),
		WithTextTemplate(welcomeTextTemplate, data),
	)

	if emailErr != nil {
		e.Logger.Error("Failed to send verification email post-signup", "email", email, "error", emailErr)
		return emailErr
	}

	return nil
}
