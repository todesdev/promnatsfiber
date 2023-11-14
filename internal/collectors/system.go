package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"runtime"
)

type SystemMetricsCollector interface {
	Collect(ch chan<- prometheus.Metric)
	Describe(ch chan<- *prometheus.Desc)
}

const (
	SystemSubsystem            = "system"
	SystemCPUUsagePercent      = "system_cpu_usage_percent"
	SystemCpuUsagePercentHelp  = "CPU usage as a percentage."
	SystemMemoryUsageBytes     = "system_memory_usage_bytes"
	SystemMemoryUsageBytesHelp = "Memory usage in bytes."
	SystemMemoryTotalBytes     = "system_memory_total_bytes"
	SystemMemoryTotalBytesHelp = "Total memory in bytes."
	SystemGCStats              = "system_gc_stats"
	SystemGCStatsHelp          = "GC stats."
	SystemGoRoutineCount       = "system_go_routine_count"
	SystemGoRoutineCountHelp   = "Number of go routines."
)

type ODSystemMetricsCollector struct {
	proc               *process.Process
	cpuUsageDesc       *prometheus.Desc
	memoryUsageDesc    *prometheus.Desc
	memoryTotalDesc    *prometheus.Desc
	gcStatsDesc        *prometheus.Desc
	goRoutineCountDesc *prometheus.Desc
}

func NewODSystemMetricsCollector(reg *prometheus.Registry, serviceName string) SystemMetricsCollector {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic(err)
	}

	collector := &ODSystemMetricsCollector{
		proc: proc,
		cpuUsageDesc: prometheus.NewDesc(
			prometheus.BuildFQName(serviceName, SystemSubsystem, SystemCPUUsagePercent),
			SystemCpuUsagePercentHelp,
			nil, nil,
		),
		memoryUsageDesc: prometheus.NewDesc(
			prometheus.BuildFQName(serviceName, SystemSubsystem, SystemMemoryUsageBytes),
			SystemMemoryUsageBytesHelp,
			nil, nil,
		),
		memoryTotalDesc: prometheus.NewDesc(
			prometheus.BuildFQName(serviceName, SystemSubsystem, SystemMemoryTotalBytes),
			SystemMemoryTotalBytesHelp,
			nil, nil,
		),
		gcStatsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(serviceName, SystemSubsystem, SystemGCStats),
			SystemGCStatsHelp,
			nil, nil,
		),
		goRoutineCountDesc: prometheus.NewDesc(
			prometheus.BuildFQName(serviceName, SystemSubsystem, SystemGoRoutineCount),
			SystemGoRoutineCountHelp,
			nil, nil,
		),
	}

	reg.MustRegister(collector)
	return collector

}

func (c *ODSystemMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	// CPU Usage
	cpuPercent, _ := c.proc.CPUPercent()
	ch <- prometheus.MustNewConstMetric(c.cpuUsageDesc, prometheus.GaugeValue, cpuPercent)

	// Memory Usage
	memInfo, _ := c.proc.MemoryInfo()
	ch <- prometheus.MustNewConstMetric(c.memoryUsageDesc, prometheus.GaugeValue, float64(memInfo.RSS))

	// System Total Memory
	vmStat, _ := mem.VirtualMemory()
	ch <- prometheus.MustNewConstMetric(c.memoryTotalDesc, prometheus.GaugeValue, float64(vmStat.Total))

	// GC Stats
	var gcStats runtime.MemStats
	runtime.ReadMemStats(&gcStats)
	ch <- prometheus.MustNewConstMetric(c.gcStatsDesc, prometheus.GaugeValue, float64(gcStats.PauseTotalNs)/1e9)

	// Goroutines Count
	ch <- prometheus.MustNewConstMetric(c.goRoutineCountDesc, prometheus.GaugeValue, float64(runtime.NumGoroutine()))
}

func (c *ODSystemMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.cpuUsageDesc
	ch <- c.memoryUsageDesc
	ch <- c.memoryTotalDesc
	ch <- c.gcStatsDesc
	ch <- c.goRoutineCountDesc
}
