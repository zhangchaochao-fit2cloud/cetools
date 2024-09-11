package cmp

import (
	"cetool/pkg/dto"
	"cetool/pkg/utils/cmd"
	"cetool/pkg/utils/docker"
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/joho/godotenv"
	"path"
	"regexp"
	"strings"
)

type CmpService struct {
	api.Service
	project *types.Project
}

func GetCmpProject(projectName, workDir string, yml []byte, env []byte, skipNormalization bool) (*types.Project, error) {
	var configFiles []types.ConfigFile
	configFiles = append(configFiles, types.ConfigFile{
		Filename: "docker-compose.yml",
		Content:  yml},
	)
	envMap, err := godotenv.UnmarshalBytes(env)
	if err != nil {
		return nil, err
	}
	details := types.ConfigDetails{
		WorkingDir:  workDir,
		ConfigFiles: configFiles,
		Environment: envMap,
	}
	projectName = strings.ToLower(projectName)
	reg, _ := regexp.Compile(`[^a-z0-9_-]+`)
	projectName = reg.ReplaceAllString(projectName, "")
	project, err := loader.LoadWithContext(context.Background(), details, func(options *loader.Options) {
		options.SetProjectName(projectName, true)
		options.ResolvePaths = true
		options.SkipNormalization = skipNormalization
	})
	if err != nil {
		return nil, err
	}
	project.ComposeFiles = []string{path.Join(workDir, "docker-compose.yml")}
	return project, nil
}

//type ComposeProject struct {
//	Version  string
//	Services map[string]Service `yaml:"services"`
//}

//func GetCmpImages(projectName string, env, yml []byte) ([]string, error) {
//	var (
//		configFiles []types.ConfigFile
//		images      []string
//	)
//	configFiles = append(configFiles, types.ConfigFile{
//		Filename: "docker-compose.yml",
//		Content:  yml},
//	)
//	envMap, err := godotenv.UnmarshalBytes(env)
//	if err != nil {
//		return nil, err
//	}
//	details := types.ConfigDetails{
//		ConfigFiles: configFiles,
//		Environment: envMap,
//	}
//
//	project, err := loader.LoadWithContext(context.Background(), details, func(options *loader.Options) {
//		options.SetProjectName(projectName, true)
//		options.ResolvePaths = true
//	})
//	if err != nil {
//		return nil, err
//	}
//	for _, service := range project.AllServices() {
//		images = append(images, service.Image)
//	}
//	return images, nil
//}

func GetCmpServices() ([]*dto.ServiceInfo, error) {
	cli, err := docker.NewClient()
	if err != nil {
		return []*dto.ServiceInfo{}, err
	}
	services, err := cli.ListContainersStats([]string{})
	return services, err

}
func GetCmpImages() ([]string, error) {
	exists := cmd.Exists("service")
	if !exists {
		return []string{}, nil
	}
	yml, err := cmd.ExecCmd("service fit2cloud raw config")
	if err != nil {
		return []string{}, nil
	}
	var (
		configFiles []types.ConfigFile
		images      []string
	)
	configFiles = append(configFiles, types.ConfigFile{
		Filename: "docker-compose.yml",
		Content:  yml},
	)
	details := types.ConfigDetails{
		ConfigFiles: configFiles,
	}

	project, err := loader.LoadWithContext(context.Background(), details, func(options *loader.Options) {
		options.SetProjectName("cmp", true)
		options.ResolvePaths = true
	})
	if err != nil {
		return nil, err
	}
	for _, service := range project.AllServices() {
		fmt.Println(service)
		images = append(images, service.Image)
	}
	return images, nil
}
