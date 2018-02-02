package cadvisor

import (
	"fmt"
	"github.com/google/cadvisor/client"
	"github.com/google/cadvisor/info/v1"
	"github.com/spf13/viper"
	"os"
	// "reflect"
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
func Containers() {
	c := connect()

	request := v1.ContainerInfoRequest{NumStats: 1}
	sInfo, err := c.SubcontainersInfo("/docker", &request)
	if err != nil {
		fmt.Println("Error getting containers info %v", err)
		return
	}
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
