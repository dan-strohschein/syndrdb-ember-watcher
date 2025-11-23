package shared

import "time"

type Metrics struct {
	Magic     uint32 // Validation
	Version   uint32
	Timestamp int64

	// Connection metrics
	ActiveConns uint32
	TotalConns  uint64

	// Query metrics
	QueriesPerSec uint32
	SlowQueries   uint32

	// Resource metrics
	MemoryUsedMB uint32
	CacheHitRate float32
}

const MetricsSize = 64          // bytes, aligned
const MetricsMagic = 0x53594E44 // "SYND"

type GlobalScanMetrics struct {
	// Cross-scanner metrics
	TotalScanners      int           `json:"total_scanners"`
	TotalScans         int64         `json:"total_scans"`
	AverageLatency     time.Duration `json:"average_latency"`
	GlobalCacheHitRate float64       `json:"global_cache_hit_rate"`

	// Hot key insights
	GlobalHotKeys   []string            `json:"global_hot_keys"`
	HotKeysByBundle map[string][]string `json:"hot_keys_by_bundle"`

	// Performance insights
	SlowestBundles     []BundlePerformance    `json:"slowest_bundles"`
	MostQueriedBundles []BundleQueryFrequency `json:"most_queried_bundles"`

	// System health
	TotalMemoryPressureGCs int64 `json:"total_memory_pressure_gcs"`
	TotalErrors            int64 `json:"total_errors"`

	// Query planner recommendations
	RecommendedOptimizations []OptimizationRecommendation `json:"recommended_optimizations"`

	// Last updated timestamp
	LastUpdated time.Time `json:"last_updated"`
}

// BundlePerformance tracks performance metrics for a specific bundle
type BundlePerformance struct {
	BundleName     string        `json:"bundle_name"`
	AverageLatency time.Duration `json:"average_latency"`
	ScanCount      int64         `json:"scan_count"`
	CacheHitRate   float64       `json:"cache_hit_rate"`
	ErrorRate      float64       `json:"error_rate"`
}

// BundleQueryFrequency tracks query frequency for a specific bundle
type BundleQueryFrequency struct {
	BundleName     string  `json:"bundle_name"`
	QueryCount     int64   `json:"query_count"`
	QueriesPerHour float64 `json:"queries_per_hour"`
	UniqueHotKeys  int     `json:"unique_hot_keys"`
}

// OptimizationRecommendation suggests performance improvements
type OptimizationRecommendation struct {
	Type           string    `json:"type"` // "index", "cache", "partition", etc.
	BundleName     string    `json:"bundle_name"`
	KeyName        string    `json:"key_name,omitempty"`
	Priority       string    `json:"priority"` // "high", "medium", "low"
	Reason         string    `json:"reason"`
	EstimatedGain  string    `json:"estimated_gain"`
	Implementation string    `json:"implementation"`
	CreatedAt      time.Time `json:"created_at"`
}

type GlobalMemoryMetrics struct {
	// Counters
	totalQueriesChecked int64 // Total number of queries that had memory tracking enabled
	totalLimitExceeded  int64 // Number of queries that exceeded the memory limit

	// Memory statistics
	totalProjectedMemory int64 // Sum of all projected memory values (for calculating average)
	maxProjectedMemory   int64 // Maximum projected memory across all queries
}
