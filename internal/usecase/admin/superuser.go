package admin

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/mail"
	"os"
	"strings"
	"syscall"
	"unicode"

	"golang.org/x/term"
)

type EmailOptions struct {
	RequireMX bool
}

type PasswordOptions struct {
	MinLen            int
	RequireUpper      bool
	RequireLower      bool
	RequireDigit      bool
	RequireSymbol     bool
	DisallowSpaces    bool
	DisallowEmailPart bool
}

func (s *AdminUsecase) CreateSuperuser() error {
	email, err := askEmail(EmailOptions{RequireMX: false})
	if err != nil {
		log.Fatal(err)
	}

	password, err := askPasswordHidden(email, PasswordOptions{
		MinLen:       8,
		RequireLower: true,
		RequireDigit: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	hashed, err := s.hasher.Hash(password)
	if err != nil {
		log.Fatal(err)
	}

	err = s.adminRepo.CreateSuperuser(email, hashed)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Superuser created successfully.")
	return nil
}

func askEmail(opts EmailOptions) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("📧 Enter super admin email: ")
		raw, err := readLine(reader)
		if err != nil {
			return "", err
		}

		email, err := validateEmail(raw, opts)
		if err == nil {
			fmt.Println("✅ Email accepted.")
			return email, nil
		}

		fmt.Printf("❌ %v\n\n", err)
	}
}

func askPasswordHidden(email string, opts PasswordOptions) (string, error) {
	applyPasswordDefaults(&opts)

	for {
		fmt.Print("🔒 Enter password (hidden): ")
		p1, err := readPassword()
		if err != nil {
			return "", err
		}

		fmt.Print("🔒 Confirm password: ")
		p2, err := readPassword()
		if err != nil {
			return "", err
		}

		if p1 != p2 {
			fmt.Println("❌ Passwords do not match.")
			continue
		}

		if err := validatePassword(p1, email, opts); err != nil {
			fmt.Printf("❌ %v\n", err)
			fmt.Println(passwordRulesHint(opts))
			fmt.Println()
			continue
		}

		fmt.Println("✅ Password accepted.")
		return p1, nil
	}
}

func validateEmail(input string, opts EmailOptions) (string, error) {
	email := strings.TrimSpace(input)
	if email == "" {
		return "", errors.New("email is required")
	}
	if strings.ContainsAny(email, " \t\n\r") {
		return "", errors.New("email must not contain spaces")
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", errors.New("invalid email format")
	}

	email = strings.TrimSpace(addr.Address)

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", errors.New("email must contain exactly one '@'")
	}
	local, domain := parts[0], parts[1]
	if local == "" || domain == "" {
		return "", errors.New("email local part and domain are required")
	}

	domain = strings.TrimSuffix(domain, ".")
	if !strings.Contains(domain, ".") {
		return "", errors.New("email domain must contain a dot (e.g. example.com)")
	}
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return "", errors.New("email domain is invalid")
	}

	email = local + "@" + strings.ToLower(domain)

	if opts.RequireMX {
		if err := validateMX(domain); err != nil {
			return "", err
		}
	}

	return email, nil
}

func validateMX(domain string) error {
	mx, err := net.LookupMX(domain)
	if err != nil {
		return errors.New("email domain does not resolve (MX lookup failed)")
	}
	if len(mx) == 0 {
		return errors.New("email domain has no MX records")
	}
	return nil
}

func applyPasswordDefaults(opts *PasswordOptions) {
	if opts.MinLen <= 0 {
		opts.MinLen = 10
	}
	if !opts.RequireUpper && !opts.RequireLower && !opts.RequireDigit && !opts.RequireSymbol {
		opts.RequireUpper = true
		opts.RequireLower = true
		opts.RequireDigit = true
		opts.RequireSymbol = true
	}
	opts.DisallowSpaces = true
	opts.DisallowEmailPart = true
}

func validatePassword(pw string, email string, opts PasswordOptions) error {
	if pw == "" {
		return errors.New("password is required")
	}
	if len([]rune(pw)) < opts.MinLen {
		return fmt.Errorf("password must be at least %d characters", opts.MinLen)
	}
	if opts.DisallowSpaces && strings.IndexFunc(pw, unicode.IsSpace) != -1 {
		return errors.New("password must not contain spaces")
	}

	lower := strings.ToLower(pw)
	if lower == "password" || lower == "password123" || lower == "admin" || lower == "admin123" {
		return errors.New("password is too common")
	}
	if isAllSameChar(pw) {
		return errors.New("password is too weak (repeated characters)")
	}

	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, r := range pw {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case isSymbol(r):
			hasSymbol = true
		}
	}

	if opts.RequireUpper && !hasUpper {
		return errors.New("password must include at least one uppercase letter")
	}
	if opts.RequireLower && !hasLower {
		return errors.New("password must include at least one lowercase letter")
	}
	if opts.RequireDigit && !hasDigit {
		return errors.New("password must include at least one number")
	}
	if opts.RequireSymbol && !hasSymbol {
		return errors.New("password must include at least one symbol (e.g. !@#$)")
	}

	if email != "" {
		email = strings.ToLower(strings.TrimSpace(email))
		if strings.Contains(lower, email) {
			return errors.New("password must not contain your email")
		}
		if opts.DisallowEmailPart {
			local := strings.SplitN(email, "@", 2)[0]
			if local != "" && len(local) >= 3 && strings.Contains(lower, local) {
				return errors.New("password must not contain the email username part")
			}
		}
	}

	return nil
}

func passwordRulesHint(opts PasswordOptions) string {
	var rules []string
	rules = append(rules, fmt.Sprintf("• at least %d characters", opts.MinLen))
	if opts.RequireUpper {
		rules = append(rules, "• 1 uppercase letter")
	}
	if opts.RequireLower {
		rules = append(rules, "• 1 lowercase letter")
	}
	if opts.RequireDigit {
		rules = append(rules, "• 1 number")
	}
	if opts.RequireSymbol {
		rules = append(rules, "• 1 symbol (e.g. !@#$)")
	}
	if opts.DisallowSpaces {
		rules = append(rules, "• no spaces")
	}
	return "Password rules:\n" + strings.Join(rules, "\n")
}

func isAllSameChar(s string) bool {
	runes := []rune(s)
	if len(runes) < 4 {
		return false
	}
	first := runes[0]
	for _, r := range runes[1:] {
		if r != first {
			return false
		}
	}
	return true
}

func isSymbol(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err == nil {
		return strings.TrimSpace(line), nil
	}

	if errors.Is(err, os.ErrClosed) {
		return "", err
	}
	trimmed := strings.TrimSpace(line)
	if trimmed != "" {
		return trimmed, nil
	}

	if errors.Is(err, io.EOF) || errors.Is(err, os.ErrNotExist) {
		return "", errors.New("input cancelled")
	}
	return "", errors.New("failed to read input")
}

func readPassword() (string, error) {
	b, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", errors.New("failed to read password (input cancelled?)")
	}
	pw := strings.TrimSpace(string(b))
	return pw, nil
}
