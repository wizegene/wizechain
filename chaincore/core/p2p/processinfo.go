package p2p

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/davecgh/go-spew/spew"
	"github.com/shirou/gopsutil/load"
	"log"
)

func GetCPUStats() {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}

	for _, s := range stat.CPUStats {
		// s.User
		// s.Nice
		// s.System
		// s.Idle
		// s.IOWait
		spew.Dump(s)
	}
}

func GetLoadAverage() float64 {
	load, _ := load.Avg()
	return load.Load1

}
