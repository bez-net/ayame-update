/*
	Disk monitoring of a machine,
	written by stoney kang, sikang99@gmail.com
*/
package main

import (
	"fmt"
	"log"
	"syscall"
)

// -----------------------------------------------------------------------------------------
const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type DiskStatus struct {
	All   uint64 `json:"all"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Avail uint64 `json:"avail"`
}

// DiskUsage: disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Avail = fs.Bavail * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

func StringDiskUsage(disk DiskStatus) (report string) {
	report = fmt.Sprintf("All: %.2f GB, ", float64(disk.All)/float64(GB))
	report += fmt.Sprintf("Avail: %.2f GB, ", float64(disk.Avail)/float64(GB))
	report += fmt.Sprintf("Used: %.2f GB, ", float64(disk.Used)/float64(GB))
	report += fmt.Sprintf("Free: %.2f GB, ", float64(disk.Free)/float64(GB))
	report += fmt.Sprintf("Spare ratio: %.2f %%", (float64(disk.Avail)/float64(disk.All))*100)
	return
}

func CheckDiskWarning(disk DiskStatus, level float64) (err error) {
	ratio := (float64(disk.Avail) / float64(disk.All)) * 100
	if ratio < level {
		msg := fmt.Sprintf("WARN> Space ratio: %.2f%% by %.2f%%", ratio, level)
		log.Println(msg)
		err = SendSlackNotification(cobot_dbs, msg)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

// -----------------------------------------------------------------------------------------
const (
	cobot_dev = "https://hooks.slack.com/services/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"
	cobot_ops = "https://hooks.slack.com/services/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"
	cobot_dbs = "https://hooks.slack.com/services/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"
)

// Config info for diskmon
type CheckConfig struct {
	level  float64 `json:"level"`
	period int     `json:"period"`
	loop   bool    `json:"loop"`
}

// -----------------------------------------------------------------------------------------
// func main() {
// 	var check CheckConfig

// 	disk := DiskUsage("/")
// 	summary := StringDiskUsage(disk)
// 	log.Println(summary)
// 	err := SendSlackNotification(cobot_ops, summary)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	for conf.loop {
// 		disk = DiskUsage("/")
// 		CheckDiskWarning(disk, conf.level)
// 		time.Sleep(time.Duration(conf.period) * time.Minute)
// 	}
// }
