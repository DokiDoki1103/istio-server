package config

type PrometheusConfig struct {
	CacheDuration   int               `yaml:"cache_duration,omitempty"`   // Cache duration per query expressed in seconds
	CacheEnabled    bool              `yaml:"cache_enabled,omitempty"`    // Enable cache for Prometheus queries
	CacheExpiration int               `yaml:"cache_expiration,omitempty"` // Global cache expiration expressed in seconds
	CustomHeaders   map[string]string `yaml:"custom_headers,omitempty"`
	HealthCheckUrl  string            `yaml:"health_check_url,omitempty"`
	IsCore          bool              `yaml:"is_core,omitempty"`
	QueryScope      map[string]string `yaml:"query_scope,omitempty"`
	URL             string            `yaml:"url,omitempty"`
}
