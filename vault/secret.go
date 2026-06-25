package vault

import "fmt"

// SecretSnapshot holds a resolved secret with its version and source metadata.
type SecretSnapshot struct {
	Env       string
	Mount     string
	Path      string
	Version   int
	Data      map[string]interface{}
}

// String returns a human-readable identifier for the snapshot.
func (s *SecretSnapshot) String() string {
	return fmt.Sprintf("%s::%s/%s@v%d", s.Env, s.Mount, s.Path, s.Version)
}

// Keys returns all top-level keys present in the secret data.
func (s *SecretSnapshot) Keys() []string {
	keys := make([]string, 0, len(s.Data))
	for k := range s.Data {
		keys = append(keys, k)
	}
	return keys
}

// GetString returns the string representation of a secret key's value.
// Returns an empty string and false if the key does not exist.
func (s *SecretSnapshot) GetString(key string) (string, bool) {
	v, ok := s.Data[key]
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%v", v), true
}

// FetchSnapshot retrieves a secret snapshot using the provided client.
func FetchSnapshot(client *Client, mount, path string, version int) (*SecretSnapshot, error) {
	data, resolvedVersion, err := client.ReadSecretVersion(mount, path, version)
	if err != nil {
		return nil, err
	}

	return &SecretSnapshot{
		Env:     client.Env,
		Mount:   mount,
		Path:    path,
		Version: resolvedVersion,
		Data:    data,
	}, nil
}
