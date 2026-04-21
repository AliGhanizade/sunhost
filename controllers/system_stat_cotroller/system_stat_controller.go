package system_stat_cotroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sunhost/controllers/log_controller"
	"sunhost/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemController struct{}

// استراکچر گسترش‌یافته (حدود ۲۰ آیتم)
type FullSysData struct {
	CPUPercent   float64 `json:"cpu_percent"`
	CPUCores     int     `json:"cpu_cores"`
	RAMTotal     uint64  `json:"ram_total"`
	RAMUsed      uint64  `json:"ram_used"`
	RAMPercent   float64 `json:"ram_percent"`
	SwapTotal    uint64  `json:"swap_total"`
	SwapUsed     uint64  `json:"swap_used"`
	DiskTotal    uint64  `json:"disk_total"`
	DiskUsed     uint64  `json:"disk_used"`
	DiskPercent  float64 `json:"disk_percent"`
	NetSent      uint64  `json:"net_sent"`
	NetRecv      uint64  `json:"net_recv"`
	NetConns     int     `json:"net_conns"`
	OS           string  `json:"os_name"`
	Uptime       uint64  `json:"uptime"`
	Procs        uint64  `json:"procs"`
	GoAllocMB    float64 `json:"go_alloc_mb"`
	NumGC        uint32  `json:"num_gc"`
	LiveObjects  uint64  `json:"live_objects"`
	NumGoroutine int     `json:"num_goroutine"`
}

func (sc *SystemController) LiveStream(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	saveCounter := 0

	username := c.DefaultQuery("username", "UnknownUser")
	
	c.Set("action", "Get LiveStream")
	c.Set("username", username)
	var logC log_controller.LogController
	logC.Create(c)
	for {
		select {
		case <-c.Writer.CloseNotify():
			return
		case <-ticker.C:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			
			cpuPercent, _ := cpu.Percent(0, false)
			cpuVal := 0.0
			if len(cpuPercent) > 0 {
				cpuVal = cpuPercent[0]
			}
			cores, _ := cpu.Counts(true)
			vMem, _ := mem.VirtualMemory()
			sMem, _ := mem.SwapMemory()
			dInfo, _ := disk.Usage(".") 
			netInfo, _ := net.IOCounters(false)
			netConns, _ := net.Connections("all")
			hostInfo, _ := host.Info()

			var sent, recv uint64
			if len(netInfo) > 0 {
				sent = netInfo[0].BytesSent
				recv = netInfo[0].BytesRecv
			}

			data := FullSysData{
				CPUPercent:   cpuVal,
				CPUCores:     cores,
				RAMTotal:     vMem.Total / 1024 / 1024,
				RAMUsed:      vMem.Used / 1024 / 1024,
				RAMPercent:   vMem.UsedPercent,
				SwapTotal:    sMem.Total / 1024 / 1024,
				SwapUsed:     sMem.Used / 1024 / 1024,
				DiskTotal:    dInfo.Total / 1024 / 1024,
				DiskUsed:     dInfo.Used / 1024 / 1024,
				DiskPercent:  dInfo.UsedPercent,
				NetSent:      sent,
				NetRecv:      recv,
				NetConns:     len(netConns),
				OS:           hostInfo.OS + " " + hostInfo.Platform,
				Uptime:       hostInfo.Uptime,
				Procs:        hostInfo.Procs,
				GoAllocMB:    float64(m.Alloc) / 1024 / 1024,
				NumGC:        m.NumGC,
				LiveObjects:  m.Mallocs - m.Frees,
				NumGoroutine: runtime.NumGoroutine(),
			}

			saveCounter++
			if saveCounter >= 10 {
				statModel := model.SystemStat{
					AllocRAM:    data.GoAllocMB,
					Goroutines:  data.NumGoroutine,
					LiveObjects: data.LiveObjects,
				}
				go func() {
					statModel.Create()
					statModel.CleanupOldRecords()
				}()
				saveCounter = 0
			}

			jsonData, _ := json.Marshal(data)
			fmt.Fprintf(c.Writer, "data: %s\n\n", string(jsonData))
			c.Writer.Flush()
		}
	}
}

func (sc *SystemController) GetSystemSummary(c *gin.Context) {
	username := c.DefaultQuery("username", "UnknownUser")
	
	c.Set("action", "Get Summary")
	c.Set("username", username)
	var logC log_controller.LogController
	logC.Create(c)

	statModel := model.SystemStat{}
	summary, err := statModel.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در خواندن اطلاعات"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"avg_ram_today":  fmt.Sprintf("%.1f", summary.AvgRam),
		"max_goroutines": summary.MaxThreads,
		"total_records":  summary.TotalRecord,
	})
}