package lib

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

/*
	this has been adapted from https://github.com/smilingthrone13/url-builder
*/

type credentials struct {
	user     string
	password string
}

type Builder struct {
	scheme      string
	host        string
	port        int
	credentials *credentials
	path        []string
	query       map[string][]string
	anchor      string
}

func NewURLBuilder() *Builder {
	return &Builder{
		path:  make([]string, 0),
		query: make(map[string][]string),
	}
}

func (b *Builder) WithScheme(scheme string) *Builder {
	b.scheme = scheme
	return b
}

func (b *Builder) WithHost(host string) *Builder {
	b.host = strings.TrimSuffix(host, "/")
	return b
}

func (b *Builder) WithPort(port int) *Builder {
	b.port = port
	return b
}

func (b *Builder) WithCredentials(user string, password string) *Builder {
	b.credentials = &credentials{
		user:     user,
		password: password,
	}
	return b
}

func (b *Builder) WithPath(elements ...string) *Builder {
	b.path = append(b.path, elements...)
	return b
}

func (b *Builder) WithQuery(key string, values ...string) *Builder {
	b.query[key] = append(b.query[key], values...)
	return b
}

func (b *Builder) WithAnchor(anchor string) *Builder {
	b.anchor = strings.Trim(anchor, "#/")
	return b
}

func (b *Builder) Build() (string, error) {
	if b.host == "" {
		return "", fmt.Errorf("host is required")
	}

	// check given host
	// todo: can't detect if given ipv6 contains port, so result string might be broken.
	if strings.Contains(b.host, "/") || // assume host contains scheme
		strings.Count(b.host, ":") == 1 { // assume host contains port (valid ipv6 have at least 2 colons)
		return "", fmt.Errorf("host contains forbidden symbols")
	}

	rawBaseUrl := fmt.Sprintf("%s://%s", b.scheme, b.host)

	if b.port > 0 {
		if b.port > 65535 {
			return "", fmt.Errorf("port must be in range [1, 65535]")
		}
		rawBaseUrl = fmt.Sprintf("%s:%d", rawBaseUrl, b.port)
	}

	u, err := url.Parse(rawBaseUrl)
	if err != nil {
		return "", err
	}

	if b.credentials != nil {
		if b.credentials.user == "" && b.credentials.password == "" {
			return "", fmt.Errorf("empty credentials")
		}

		if b.credentials.password == "" {
			u.User = url.User(b.credentials.user)
		} else {
			u.User = url.UserPassword(b.credentials.user, b.credentials.password)
		}

	}

	if len(b.path) > 0 {
		u = u.JoinPath(b.path...)
	}

	for k, v := range b.query {
		if k == "" {
			return "", fmt.Errorf("query key is empty")
		}
		if i := slices.Index(v, ""); i != -1 {
			return "", fmt.Errorf("query query value for key %s", k)
		}
	}

	u.RawQuery = url.Values(b.query).Encode()

	if b.anchor != "" {
		u.Fragment = b.anchor
	}

	return u.String(), nil
}
