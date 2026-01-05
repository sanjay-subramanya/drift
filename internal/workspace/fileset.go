package workspace

import (
	"path"
	"strings"
)

func IsEnvFile(filename string) bool {
	if filename == ".env" {
		return true
	}
	return len(filename) >= 4 && filename[:4] == ".env"
}

func IsDeploymentFile(filename string) bool {
	if strings.ToLower(filename) == "dockerfile" {
		return true
	}
	if len(filename) >= 10 && filename[:10] == "Dockerfile" {
		return true
	}
	if filename == "docker-compose.yml" || filename == "docker-compose.yaml" {
		return true
	}
	if match, _ := path.Match("k8s/*.yaml", filename); match {
		return true
	}
	return false
}
