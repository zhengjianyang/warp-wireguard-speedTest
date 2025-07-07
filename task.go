package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

// PingResult 存储单次ping的结果
type PingResult struct {
	IP       net.IP
	Port     int
	Duration time.Duration
}

// byDuration 实现sort.Interface，用于按Duration对PingResult切片进行排序
type byDuration []PingResult

func (a byDuration) Len() int           { return len(a) }
func (a byDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDuration) Less(i, j int) bool { return a[i].Duration < a[j].Duration }

// RunTasks 运行ping测试任务
func RunTasks(config *Config) {
	var results []PingResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	pool, _ := ants.NewPool(config.Threads)
	defer pool.Release()

	for _, cidr := range config.IPv4CIDRs {
		ips, err := ParseIPRange(cidr)
		if err != nil {
			log.Printf("无法解析CIDR %s: %v", cidr, err)
			continue
		}

		for _, ip := range ips {
			for _, port := range config.Ports {
				wg.Add(1)
				err = pool.Submit(func() {
					defer wg.Done()
					var successfulPings []time.Duration
					for i := 0; i < config.PingCount; i++ {
						ok, duration := WireGuardPing(ip, port)
						if ok {
							successfulPings = append(successfulPings, duration)
						} else {
							break
						}
					}

					if len(successfulPings) == config.PingCount {
						sort.Slice(successfulPings, func(i int, j int) bool {
							return successfulPings[i] < successfulPings[j]
						})
						fmt.Printf("IP: %s, Port: %d, Latency: %v\n", ip, port, successfulPings[0])
						mu.Lock()
						results = append(results, PingResult{IP: ip, Port: port, Duration: successfulPings[0]})
						mu.Unlock()
					}
				})
				if err != nil {
					wg.Done()
					log.Println("创建线程失败: ", err.Error())
				}
			}
		}
	}

	wg.Wait()

	log.Println("==================== 测试完成 ====================")
	if len(results) == 0 {
		return
	}
	fmt.Println("执行结果排序>>>")
	sort.Sort(byDuration(results))

	if len(results) > 10 {
		if err := saveResultsToCSV(results, config.SaveFileName); err != nil {
			log.Fatalf("无法保存结果到CSV: %v", err)
		}
		fmt.Println("结果已保存到: ", config.SaveFileName)

		results = results[:10]
	}
	fmt.Printf("前%d条延迟最低的结果>>>\n", len(results))
	for index, result := range results {
		fmt.Printf("[%d] IP: %s, Port: %d, Latency: %v\n", index+1, result.IP, result.Port, result.Duration)
	}
}

// saveResultsToCSV 将结果保存到CSV文件
func saveResultsToCSV(results []PingResult, saveFileName string) error {
	file, err := os.Create(saveFileName)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"IP", "Port", "Latency(ms)"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, result := range results {
		row := []string{
			result.IP.String(),
			fmt.Sprintf("%d", result.Port),
			fmt.Sprintf("%.2f", float64(result.Duration.Milliseconds())),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
