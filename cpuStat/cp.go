package cpuStat

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
)

func cpuInfo() {
	cpuStat, err := cpu.Info()
	if err != nil {
		fmt.Println("error", err)
		return
	}

	percentage, errP := cpu.Percent(0, true)
	if errP != nil {
		fmt.Println("error", errP)
		return
	}

	fmt.Println("CPU Percent ", percentage)

	for _, val := range cpuStat {
		fmt.Println("CPU Info \n", val)
	}
}
