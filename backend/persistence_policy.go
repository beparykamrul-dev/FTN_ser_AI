package backend

import "time"

type PersistencePolicy struct {
	RAMFirst        bool          `json:"ram_first"`
	RedisHotState   bool          `json:"redis_hot_state"`
	PostgresDurable bool          `json:"postgres_durable"`
	AsyncPersist    bool          `json:"async_persist"`
	ReadPathDBFree  bool          `json:"read_path_db_free"`
	CacheTTL        time.Duration `json:"cache_ttl"`
}

func DefaultPersistencePolicy() PersistencePolicy {
	return PersistencePolicy{
		RAMFirst:        true,
		RedisHotState:   true,
		PostgresDurable: true,
		AsyncPersist:    true,
		ReadPathDBFree:  true,
		CacheTTL:        30 * time.Second,
	}
}
