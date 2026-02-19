package metrics

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type DBStatsProvider func() DBStats

type DBStats struct {
	AcquireCount         int64
	AcquiredConns        int32
	IdleConns            int32
	TotalConns           int32
	ConstructingConns    int32
	MaxConns             int32
	AcquireDurationNanos int64
}

type Collector struct {
	mu                 sync.RWMutex
	requestCount       map[string]uint64
	errorCount         map[string]uint64
	latencyBuckets     map[string][]uint64
	dbQueryCount       map[string]uint64
	dbQueryErrorCount  map[string]uint64
	dbLatencyBuckets   map[string][]uint64
	latencyBucketEdges []float64
	dbStats            DBStatsProvider
}

func NewCollector(dbStats DBStatsProvider) *Collector {
	return &Collector{
		requestCount:       make(map[string]uint64),
		errorCount:         make(map[string]uint64),
		latencyBuckets:     make(map[string][]uint64),
		dbQueryCount:       make(map[string]uint64),
		dbQueryErrorCount:  make(map[string]uint64),
		dbLatencyBuckets:   make(map[string][]uint64),
		latencyBucketEdges: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
		dbStats:            dbStats,
	}
}

func (c *Collector) ObserveRequest(method, route string, status int, duration time.Duration) {
	key := sanitizeLabel(method + "_" + route)
	seconds := duration.Seconds()

	c.mu.Lock()
	defer c.mu.Unlock()
	c.requestCount[key]++
	if status >= 400 {
		c.errorCount[key]++
	}
	c.observeDuration(c.latencyBuckets, key, seconds)
}

func (c *Collector) ObserveDBQuery(name string, duration time.Duration, failed bool) {
	key := sanitizeLabel(name)
	seconds := duration.Seconds()

	c.mu.Lock()
	defer c.mu.Unlock()
	c.dbQueryCount[key]++
	if failed {
		c.dbQueryErrorCount[key]++
	}
	c.observeDuration(c.dbLatencyBuckets, key, seconds)
}

func (c *Collector) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		c.mu.RLock()
		requestCount := cloneU64Map(c.requestCount)
		errorCount := cloneU64Map(c.errorCount)
		latencyBuckets := cloneBucketMap(c.latencyBuckets)
		dbQueryCount := cloneU64Map(c.dbQueryCount)
		dbQueryError := cloneU64Map(c.dbQueryErrorCount)
		dbLatency := cloneBucketMap(c.dbLatencyBuckets)
		edges := append([]float64(nil), c.latencyBucketEdges...)
		c.mu.RUnlock()

		writeCounterMetrics(w, "http_requests_total", requestCount)
		writeCounterMetrics(w, "http_request_errors_total", errorCount)
		writeHistogramMetrics(w, "http_request_duration_seconds_bucket", latencyBuckets, edges)

		writeCounterMetrics(w, "db_query_total", dbQueryCount)
		writeCounterMetrics(w, "db_query_errors_total", dbQueryError)
		writeHistogramMetrics(w, "db_query_duration_seconds_bucket", dbLatency, edges)

		if c.dbStats != nil {
			stats := c.dbStats()
			fmt.Fprintf(w, "db_pool_acquire_count %d\n", stats.AcquireCount)
			fmt.Fprintf(w, "db_pool_acquired_conns %d\n", stats.AcquiredConns)
			fmt.Fprintf(w, "db_pool_idle_conns %d\n", stats.IdleConns)
			fmt.Fprintf(w, "db_pool_total_conns %d\n", stats.TotalConns)
			fmt.Fprintf(w, "db_pool_constructing_conns %d\n", stats.ConstructingConns)
			fmt.Fprintf(w, "db_pool_max_conns %d\n", stats.MaxConns)
			fmt.Fprintf(w, "db_pool_acquire_duration_seconds %.6f\n", float64(stats.AcquireDurationNanos)/1e9)
		}
	}
}

func (c *Collector) observeDuration(target map[string][]uint64, key string, seconds float64) {
	buckets, ok := target[key]
	if !ok {
		buckets = make([]uint64, len(c.latencyBucketEdges)+1)
	}
	placed := false
	for i, edge := range c.latencyBucketEdges {
		if seconds <= edge {
			buckets[i]++
			placed = true
			break
		}
	}
	if !placed {
		buckets[len(buckets)-1]++
	}
	target[key] = buckets
}

func writeCounterMetrics(w http.ResponseWriter, name string, values map[string]uint64) {
	keys := sortedKeys(values)
	for _, key := range keys {
		fmt.Fprintf(w, "%s{route=\"%s\"} %d\n", name, key, values[key])
	}
}

func writeHistogramMetrics(w http.ResponseWriter, name string, values map[string][]uint64, edges []float64) {
	keys := sortedBucketKeys(values)
	for _, key := range keys {
		for i, count := range values[key] {
			if i < len(edges) {
				fmt.Fprintf(w, "%s{route=\"%s\",le=\"%.2f\"} %d\n", name, key, edges[i], count)
			} else {
				fmt.Fprintf(w, "%s{route=\"%s\",le=\"+Inf\"} %d\n", name, key, count)
			}
		}
	}
}

func sanitizeLabel(value string) string {
	v := strings.ReplaceAll(value, " ", "_")
	v = strings.ReplaceAll(v, "/", "_")
	v = strings.ReplaceAll(v, "-", "_")
	return strings.Trim(v, "_")
}

func cloneU64Map(source map[string]uint64) map[string]uint64 {
	result := make(map[string]uint64, len(source))
	for k, v := range source {
		result[k] = v
	}
	return result
}

func cloneBucketMap(source map[string][]uint64) map[string][]uint64 {
	result := make(map[string][]uint64, len(source))
	for k, v := range source {
		result[k] = append([]uint64(nil), v...)
	}
	return result
}

func sortedKeys(m map[string]uint64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedBucketKeys(m map[string][]uint64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
