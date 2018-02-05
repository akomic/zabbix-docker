package cadvisor

import (
	"fmt"
	// . "github.com/blacked/go-zabbix"
	zclient "github.com/akomic/zabbix-proto/client"
	zsender "github.com/akomic/zabbix-proto/sender"
	"github.com/google/cadvisor/client"
	"github.com/google/cadvisor/info/v1"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func connect() *client.Client {
	addr := viper.GetString("addr")
	client, err := client.NewClient(addr)
	if err != nil {
		fmt.Println("Error connecting to %s : %v", addr, err)
		os.Exit(100)
	}
	return client
}

// func containers() map[string]map[string]string {
func getAllStats() []v1.ContainerInfo {
	c := connect()

	request := v1.ContainerInfoRequest{NumStats: 1}
	sInfo, err := c.SubcontainersInfo("/docker", &request)
	if err != nil {
		fmt.Println("Error getting containers info %v", err)
		os.Exit(1)
	}

	return sInfo
}

func Containers() {
	sInfo := getAllStats()
	for _, container := range sInfo {
		if container.Id != "" {
			fmt.Printf("ID: %s Name: %s\n", container.Id, container.Aliases[0])
		}
	}
}

func Container() {
	containerId := viper.GetString("containerId")
	c := connect()

	request := v1.ContainerInfoRequest{NumStats: 1}
	cInfo, err := c.ContainerInfo(fmt.Sprintf("/docker/%s", containerId), &request)
	if err != nil {
		fmt.Println("Error getting container info %s : %v", containerId, err)
		return
	}

	fmt.Println("ID:", cInfo.Id)
	fmt.Println("Name:", cInfo.Aliases[0])
	fmt.Println("Labels:", cInfo.Labels)
	fmt.Println("Image:", cInfo.Spec.Image)
	fmt.Println("CreationTime:", cInfo.Spec.CreationTime)

	fmt.Println("CPU Usage Total:", cInfo.Stats[0].Cpu.Usage.Total, "nanoseconds")
	fmt.Println("Memory Usage:", cInfo.Stats[0].Memory.Usage, "Bytes")
	fmt.Println("IO Time:", len(cInfo.Stats[0].DiskIo.IoServiceBytes))
	for _, d := range cInfo.Stats[0].DiskIo.IoServiceBytes {
		fmt.Println(d)
	}
}

type DiscoveryPayload struct {
	Data []map[string]string `json:"data"`
}

func Zabbix() {
	verbose := viper.GetBool("verbose")
	zabbixAddr := viper.GetString("zabbixAddr")

	zbxComp := strings.Split(zabbixAddr, ":")
	zabbixHost := zbxComp[0]
	zabbixPort, _ := strconv.Atoi(zbxComp[1])

	c := zclient.NewClient(zabbixHost, zabbixPort)

	var metrics []*zsender.Metric

	sInfo := getAllStats()

	for _, container := range sInfo {
		if container.Id == "" || len(container.Stats) == 0 {
			continue
		}

		containerName := fmt.Sprintf("%s - %s", container.Aliases[0], container.Id[:12])

		items, _ := c.GetActiveItems(containerName, "DContainer")
		if verbose {
			fmt.Println("Doing Active Check for:", containerName)
			fmt.Println("Got:", reflect.TypeOf(items.Data))
		}

		// fmt.Printf("%v\n", reflect.TypeOf(container.Stats[0].DiskIo.IoServiceBytes))

		// fmt.Println("CPU SPEC:", container.Spec.Cpu.Limit)
		for _, item := range items.Data {
			switch {
			// Container Specs
			case item.Key == "Spec.Image":
				metrics = append(metrics, zsender.NewMetric(containerName, item.Key, container.Spec.Image, time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", container.Spec.Image)
				}
			case item.Key == "container.Spec.Cpu.Limit":
				metrics = append(metrics, zsender.NewMetric(containerName, item.Key, strconv.FormatUint(container.Spec.Cpu.Limit, 10), time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", strconv.FormatUint(container.Spec.Cpu.Limit, 10))
				}
			case item.Key == "container.Spec.Cpu.MaxLimit":
				metrics = append(metrics, zsender.NewMetric(containerName, item.Key, strconv.FormatUint(container.Spec.Cpu.MaxLimit, 10), time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", strconv.FormatUint(container.Spec.Cpu.MaxLimit, 10))
				}
			case item.Key == "container.Spec.Cpu.Quota":
				metrics = append(metrics, zsender.NewMetric(containerName, item.Key, strconv.FormatUint(container.Spec.Cpu.Quota, 10), time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", strconv.FormatUint(container.Spec.Cpu.Quota, 10))
				}
			case item.Key == "container.Spec.Cpu.Period":
				metrics = append(metrics, zsender.NewMetric(containerName, item.Key, strconv.FormatUint(container.Spec.Cpu.Period, 10), time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", strconv.FormatUint(container.Spec.Cpu.Period, 10))
				}
			case item.Key == "container.Spec.Cpu.Mask":
				metrics = append(metrics, zsender.NewMetric(containerName, item.Key, container.Spec.Cpu.Mask, time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", container.Spec.Cpu.Mask)
				}

			// CPU Metrics
			case item.Key == "Cpu.Usage.Total":
				metrics = append(metrics, zsender.NewMetric(containerName, item.Key, strconv.FormatUint(container.Stats[0].Cpu.Usage.Total, 10), time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", strconv.FormatUint(container.Stats[0].Cpu.Usage.Total, 10))
				}
				// Memory metrics
			case item.Key == "Memory.Usage":
				metrics = append(metrics, zsender.NewMetric(containerName, item.Key, strconv.FormatUint(container.Stats[0].Memory.Usage, 10), time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", strconv.FormatUint(container.Stats[0].Memory.Usage, 10))
				}
				// Task stats
			case item.Key == "TaskStats.NrIoWait":
				metrics = append(metrics, zsender.NewMetric(containerName, item.Key, strconv.FormatUint(container.Stats[0].TaskStats.NrIoWait, 10), time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", strconv.FormatUint(container.Stats[0].TaskStats.NrIoWait, 10))
				}
			case item.Key == "TaskStats.NrRunning":
				metrics = append(metrics, zsender.NewMetric(containerName, item.Key, strconv.FormatUint(container.Stats[0].TaskStats.NrRunning, 10), time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", strconv.FormatUint(container.Stats[0].TaskStats.NrRunning, 10))
				}
				// Disk discovery
			case item.Key == "discoveryDiskIo":
				var discoveryData []map[string]string
				for _, device := range container.Stats[0].DiskIo.IoServiceBytes {
					discoveryItem := make(map[string]string)

					discoveryItem["{#DEVICE}"] = device.Device

					// for ioKey, ioVal := range device.Stats {
					// 	discoveryItem[fmt.Sprintf("{#%s}", strings.ToUpper(ioKey))] = strconv.FormatUint(ioVal, 10)
					// }
					discoveryData = append(discoveryData, discoveryItem)
				}

				metrics = append(metrics, zsender.NewDiscoveryMetric(containerName, item.Key, discoveryData, time.Now().Unix()))
				if verbose {
					fmt.Println("Sending:", containerName, item.Key, "=", discoveryData)
				}
				// DiskIo stats
			case strings.HasPrefix(item.Key, "DiskIo.IoServiceBytes.Stats."):
				for _, device := range container.Stats[0].DiskIo.IoServiceBytes {
					for ioKey, ioVal := range device.Stats {
						guessKey := fmt.Sprintf("DiskIo.IoServiceBytes.Stats.%s[%s]", ioKey, device.Device)
						if strings.ToUpper(guessKey) == strings.ToUpper(item.Key) {
							metrics = append(metrics, zsender.NewMetric(containerName, item.Key, strconv.FormatUint(ioVal, 10), time.Now().Unix()))
							if verbose {
								fmt.Println("Sending:", containerName, item.Key, "=", strconv.FormatUint(ioVal, 10))
							}
						}
					}
				}
			default:
				if verbose {
					fmt.Println("Unknown Key:", item.Key)
				}
			}
		}
	}

	if len(metrics) > 0 {
		packet := zsender.NewPacket(metrics)
		res, err := c.Send(packet)

		if err != nil || res.Response != "success" {
			fmt.Errorf("Error sending items: %s", err.Error)
			fmt.Errorf("Got response: %s", res.Response)
		} else {
			fmt.Println("Got:", res.Info)
		}
	}
}
