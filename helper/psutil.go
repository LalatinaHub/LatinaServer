package helper

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

type ServerStat struct {
	Cpu  []float64              `json:"cpu"`
	Host *host.InfoStat         `json:"host"`
	Mem  *mem.VirtualMemoryStat `json:"mem"`
	Disk *disk.UsageStat        `json:"disk"`
	Nic  []net.IOCountersStat   `json:"nic"`
}

func GetServerStatus() ServerStat {
	var (
		host, _ = host.Info()
		cpu, _  = cpu.Percent(1*time.Second, false)
		ram, _  = mem.VirtualMemory()
		disk, _ = disk.Usage("/")
		nic, _  = net.IOCounters(false)
	)

	return ServerStat{
		Host: host,
		Cpu:  cpu,
		Mem:  ram,
		Disk: disk,
		Nic:  nic,
	}
}
