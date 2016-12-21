package db

import (
	"strconv"
	"sync"
)

// Config provides access to configuration values for the application. In order
// to maximize efficiency, all values are loaded at startup. Any changes are
// stored in the cache and written to the database. Access is guarded by a
// mutex in order to avoid race conditions.
type Config struct {
	values map[string]string
	mutex  sync.RWMutex
}

// migrateConfigTable executes the SQL necessary to create the Config table.
func migrateConfigTable(t *Token) error {
	_, err := t.exec(
		`
        CREATE TABLE IF NOT EXISTS Config (
            Key   VARCHAR(20) PRIMARY KEY,
            Value TEXT NOT NULL
        )
        `,
	)
	return err
}

// NewConfig loads the configuration from the database.
func NewConfig(t *Token) (*Config, error) {
	c := &Config{
		values: make(map[string]string),
	}
	r, err := t.query(`SELECT Key, Value FROM Config`)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for r.Next() {
		var (
			key   string
			value string
		)
		if err := r.Scan(&key, &value); err != nil {
			return nil, err
		}
		c.values[key] = value
	}
	return c, nil
}

// GetString retrieves the string value for the configuration entry with the
// specified key.
func (c *Config) GetString(key string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	v, _ := c.values[key]
	return v
}

// GetInt retrieves the integer value for the configuration entry with the
// specified key.
func (c *Config) GetInt(key string) int {
	v, err := strconv.Atoi(c.GetString(key))
	if err != nil {
		return 0
	}
	return v
}

// SetString stores a new string value for the specified key.
func (c *Config) SetString(t *Token, key, value string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, err := t.exec(
		`
        UPDATE Config SET Value=$1
        WHERE Key = $2
        `,
		key,
		value,
	)
	if err != nil {
		return err
	}
	c.values[key] = value
	return nil
}

// SetInt stores a new integer value for the specified key.
func (c *Config) SetInt(t *Token, key string, value int) error {
	return c.SetString(t, key, strconv.Itoa(value))
}
