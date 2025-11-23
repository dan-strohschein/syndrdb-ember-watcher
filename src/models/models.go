package models

import "syndrdb-ember-watcher/src/shared"

type MetricsBlock struct {
	Timestamp      int64
	ActiveConns    int32
	QueriesPerSec  int32
	AvgQueryTimeMs int32
	IndexHitRate   float32

	BaseMetrics                shared.Metrics
	GlobalScanMetrics          shared.GlobalScanMetrics
	BundlePerformance          shared.BundlePerformance
	BundleQueryFrequency       shared.BundleQueryFrequency
	OptimizationRecommendation shared.OptimizationRecommendation
	GlobalMemoryMetrics        shared.GlobalMemoryMetrics
}
