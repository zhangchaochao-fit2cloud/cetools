package docker

import (
	"context"
	"fmt"
	"inspect/pkg/dto"
	"inspect/pkg/global"
	"inspect/pkg/utils"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"

	//"github.com/1Panel-dev/1Panel/backend/app/model"
	//"github.com/1Panel-dev/1Panel/backend/global"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func NewClient() (Client, error) {
	//sock := global.CONF.Docker.Sock
	//if sock == "" || len(sock) == 0 {
	sock := "unix:///var/run/docker.sock"
	//}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithHost(sock), client.WithAPIVersionNegotiation())
	if err != nil {
		return Client{}, err
	}

	return Client{
		cli: cli,
	}, nil
}

func (c Client) Close() {
	_ = c.cli.Close()
}

func NewDockerClient() (*client.Client, error) {
	sock := global.CONF.Docker.Sock
	if len(sock) == 0 {
		sock = "unix:///var/run/docker.sock"
	}
	info, err := os.Stat(sock)
	if err != nil {
		return nil, fmt.Errorf("please check whether the '%s' file exists", sock)
	}
	// 检查是否为 socket 文件
	if !info.Mode().IsRegular() && info.Mode()&os.ModeSocket == os.ModeSocket {
		return nil, fmt.Errorf("'%s' is not a socket file", sock)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithHost(sock), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (c Client) ListContainersStats(names []string) ([]*dto.ServiceInfo, error) {
	var (
		options  container.ListOptions
		namesMap = make(map[string]bool)
		services []*dto.ServiceInfo
	)
	options.All = true
	if len(names) > 0 {
		var array []filters.KeyValuePair
		for _, n := range names {
			namesMap["/"+n] = true
			array = append(array, filters.Arg("name", n))
		}
		options.Filters = filters.NewArgs(array...)
	}
	containers, err := c.cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	wg.Add(len(containers))
	for i := 0; i < len(containers); i++ {
		go func(timeoutCh <-chan time.Time, con types.Container) {
			defer wg.Done()

			serviceInfo := &dto.ServiceInfo{}

			serviceInfo.Name = getName(&con)
			serviceInfo.Status = con.State
			serviceInfo.Ports = getUniquePublicPorts(&con)
			serviceInfo.Runtime = con.Status

			select {
			case <-timeoutCh:
				mu.Lock()
				services = append(services, serviceInfo)
				mu.Unlock()
				global.LOG.Errorf("load container info from %s failed, err: timeout", con.Names)
			default:
				stats, err := c.cli.ContainerStats(context.Background(), con.ID, false)
				if err != nil {
					mu.Lock()
					services = append(services, serviceInfo)
					mu.Unlock()
					global.LOG.Errorf("load container info from %s failed, err: %v", con.Names, err)
					return
				}

				defer stats.Body.Close()
				statMap := utils.Stat(stats)

				memStatInfo := statMap["memory_stats"].(map[string]interface{})
				cpuStatInfo := statMap["cpu_stats"].(map[string]interface{})
				systemCpu := cpuStatInfo["system_cpu_usage"].(float64)
				cpuUsage := cpuStatInfo["cpu_usage"].(map[string]interface{})
				//totalUsage := cpuUsage["total_usage"].(float64)
				usageInKernelmode := cpuUsage["usage_in_kernelmode"].(float64)
				usageInUsermode := cpuUsage["usage_in_usermode"].(float64)

				memLimit := utils.SpaceFloat(memStatInfo["limit"].(float64))
				memUsed := utils.SpaceFloat(memStatInfo["usage"].(float64))

				cpuPercent := utils.CalculatePercent(systemCpu, usageInKernelmode+usageInUsermode)

				serviceInfo.CPUPercent = cpuPercent
				serviceInfo.MemUsed = memUsed
				serviceInfo.MemLimit = memLimit
				mu.Lock()
				services = append(services, serviceInfo)
				mu.Unlock()
			}
		}(time.After(5*time.Second), containers[i])
	}
	wg.Wait()
	//for _, con := range containers {
	//
	//	stats, err := c.cli.ContainerStats(context.Background(), con.ID, false)
	//	if err != nil {
	//		return nil, err
	//	}
	//	defer stats.Body.Close()
	//
	//	var nameStr string
	//	for i, name := range con.Names {
	//		name = strings.ReplaceAll(name, "/", "")
	//		nameStr += name
	//		if len(con.Names)-1 != i {
	//			nameStr += ","
	//		}
	//	}
	//	statMap := utils.Stat(stats)
	//
	//	var portStr string
	//	for i, p := range con.Ports {
	//		if p.PublicPort == 0 {
	//			continue
	//		}
	//		portStr += strconv.Itoa(int(p.PublicPort))
	//		if len(con.Ports)-1 != i {
	//			portStr += ","
	//		}
	//	}
	//
	//	memStatInfo := statMap["memory_stats"].(map[string]interface{})
	//	cpuStatInfo := statMap["cpu_stats"].(map[string]interface{})
	//	systemCpu := cpuStatInfo["system_cpu_usage"].(float64)
	//	cpuUsage := cpuStatInfo["cpu_usage"].(map[string]interface{})
	//	//totalUsage := cpuUsage["total_usage"].(float64)
	//	usageInKernelmode := cpuUsage["usage_in_kernelmode"].(float64)
	//	usageInUsermode := cpuUsage["usage_in_usermode"].(float64)
	//
	//	memLimit := utils.SpaceFloat(memStatInfo["limit"].(float64))
	//	memUsed := utils.SpaceFloat(memStatInfo["usage"].(float64))
	//
	//	cpuPercent := utils.CalculatePercent(systemCpu, usageInKernelmode+usageInUsermode)
	//
	//	serviceInfo := &dto.ServiceInfo{
	//		Name:       nameStr,
	//		Status:     con.State,
	//		Ports:      portStr,
	//		CPUPercent: cpuPercent,
	//		MemUsed:    memUsed,
	//		MemLimit:   memLimit,
	//		Runtime:    con.Status,
	//	}
	//	services = append(services, serviceInfo)
	//}
	return services, nil
}

func getName(con *types.Container) string {
	var nameStr string
	for i, name := range con.Names {
		name = strings.ReplaceAll(name, "/", "")
		nameStr += name
		if len(con.Names)-1 != i {
			nameStr += ","
		}
	}
	return nameStr
}

func getUniquePublicPorts(con *types.Container) string {
	uniquePorts := make(map[uint16]bool) // 用于记录已出现的端口
	var portStrBuilder strings.Builder   // 使用 strings.Builder 来构建字符串

	for _, p := range con.Ports {
		if p.PublicPort == 0 {
			continue
		}
		if _, exists := uniquePorts[p.PublicPort]; !exists {
			uniquePorts[p.PublicPort] = true
			portStrBuilder.WriteString(strconv.Itoa(int(p.PublicPort)))
			portStrBuilder.WriteString(",")
		}
	}

	// 移除最后一个逗号
	portStr := portStrBuilder.String()
	if len(portStr) > 0 {
		portStr = portStr[:len(portStr)-1]
	}

	return string(portStr)
}

func (c Client) ListContainersByName(names []string) ([]types.Container, error) {
	var (
		options  container.ListOptions
		namesMap = make(map[string]bool)
		res      []types.Container
	)
	options.All = true
	if len(names) > 0 {
		var array []filters.KeyValuePair
		for _, n := range names {
			namesMap["/"+n] = true
			array = append(array, filters.Arg("name", n))
		}
		options.Filters = filters.NewArgs(array...)
	}
	containers, err := c.cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	for _, con := range containers {
		if _, ok := namesMap[con.Names[0]]; ok {
			res = append(res, con)
		}
	}
	return res, nil
}
func (c Client) ListAllContainers() ([]types.Container, error) {
	var (
		options container.ListOptions
	)
	options.All = true
	containers, err := c.cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (c Client) CreateNetwork(name string) error {
	_, err := c.cli.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Driver: "bridge",
	})
	return err
}

func (c Client) DeleteImage(imageID string) error {
	if _, err := c.cli.ImageRemove(context.Background(), imageID, image.RemoveOptions{Force: true}); err != nil {
		return err
	}
	return nil
}

func (c Client) InspectContainer(containerID string) (types.ContainerJSON, error) {
	return c.cli.ContainerInspect(context.Background(), containerID)
}

func (c Client) PullImage(imageName string, force bool) error {
	if !force {
		exist, err := c.CheckImageExist(imageName)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	if _, err := c.cli.ImagePull(context.Background(), imageName, image.PullOptions{}); err != nil {
		return err
	}
	return nil
}

func (c Client) GetImageIDByName(imageName string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("reference", imageName)
	list, err := c.cli.ImageList(context.Background(), image.ListOptions{
		Filters: filter,
	})
	if err != nil {
		return "", err
	}
	if len(list) > 0 {
		return list[0].ID, nil
	}
	return "", nil
}

func (c Client) CheckImageExist(imageName string) (bool, error) {
	filter := filters.NewArgs()
	filter.Add("reference", imageName)
	list, err := c.cli.ImageList(context.Background(), image.ListOptions{
		Filters: filter,
	})
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}

func (c Client) NetworkExist(name string) bool {
	var options types.NetworkListOptions
	options.Filters = filters.NewArgs(filters.Arg("name", name))
	networks, err := c.cli.NetworkList(context.Background(), options)
	if err != nil {
		return false
	}
	return len(networks) > 0
}

func CreateDefaultDockerNetwork() error {
	cli, err := NewClient()
	if err != nil {
		global.LOG.Errorf("init docker client error %s", err.Error())
		return err
	}
	defer cli.Close()
	if !cli.NetworkExist("1panel-network") {
		if err := cli.CreateNetwork("1panel-network"); err != nil {
			global.LOG.Errorf("create default docker network  error %s", err.Error())
			return err
		}
	}
	return nil
}
