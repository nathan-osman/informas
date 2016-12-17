package db

import (
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
func migrateConfigTable() error {
	_, err := db.Exec(
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
func NewConfig() (*Config, error) {
	c := &Config{
		values: make(map[string]string),
	}
	r, err := db.Query(`SELECT Key, Value FROM Config`)
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

// Get retrieves the value for the configuration entry with the specified key.
// An empty string is returned if the value does not exist.
func (c *Config) Get(key string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	v, _ := c.values[key]
	return v
}

// Set stores a new value for the specified key.
func (c *Config) Set(key, value string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, err := db.Exec(
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
