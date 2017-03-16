package main

import (
	"fmt"
	"github.com/bdudelsack/docker-registry-client/registry"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Configuration struct {
	Repositories []RepositoryMirror     `yaml:"repositories"`
	Auth         map[string]Credentials `yaml:"auth"`
}

type RepositoryMirror struct {
	Source      string   `yaml:"source"`
	Destination string   `yaml:"destination"`
	Matches     []string `yaml:"matches"`
}

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var (
	config *Configuration
)

func main() {
	if err := readConfiguration(); err != nil {
		panic(err)
	}

	for _, r := range config.Repositories {
		if err := syncRepository(r.Source, r.Destination, r.Matches); err != nil {
			panic(err)
		}
	}
}

func readConfiguration() error {
	buf, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(buf, &config); err != nil {
		return err
	}

	return nil
}

func syncRepository(source string, destination string, matches []string) error {
	fmt.Printf("Mirroring repository %s -> %s\n", source, destination)

	sourceRegistryUrl, sourceRepository := splitReposSearchTerm(source)
	sourceRegistry, err := newRegistryWithCredentials("https://"+sourceRegistryUrl, getCredentials(sourceRegistryUrl))

	if err != nil {
		return err
	}

	sourceTags, err := sourceRegistry.Tags(sourceRepository)

	if err != nil {
		return err
	}

	fmt.Printf("Found tags: %v\n", sourceTags)

	sourceTags = filterTags(sourceTags, matches)

	fmt.Printf("Found matching tags: %v\n", sourceTags)

	targetRegistryUrl, targetRepository := splitReposSearchTerm(destination)
	targetRegistry, err := newRegistryWithCredentials("https://"+targetRegistryUrl, getCredentials(targetRegistryUrl))

	if err != nil {
		return err
	}

	for _, t := range sourceTags {
		sourceDigest, err := sourceRegistry.ManifestDigest(sourceRepository, t)

		fmt.Printf("%s:%s -- ", sourceRepository, t)

		if err != nil {
			return err
		}

		targetDigest, err := targetRegistry.ManifestDigest(targetRepository, t)

		if err != nil || targetDigest != sourceDigest {
			fmt.Println("Downloading", sourceDigest.String())

			sourceManifest, err := sourceRegistry.Manifest(sourceRepository, t)

			if err != nil {
				return err
			}

			for _, layer := range sourceManifest.FSLayers {

				fmt.Printf("%s -- ", layer.BlobSum.String())

				hasLayer, err := targetRegistry.HasLayer(targetRepository, layer.BlobSum)

				if !hasLayer || err != nil {
					reader, err := sourceRegistry.DownloadLayer(sourceRepository, layer.BlobSum)

					if err != nil {
						return err
					}

					fmt.Println("Downloading")

					targetRegistry.UploadLayer(targetRepository, layer.BlobSum, reader)

					reader.Close()

				} else {
					fmt.Println("Layer up to date")
				}

				targetRegistry.PutManifest(targetRepository, t, sourceManifest)
			}
		} else {
			fmt.Println("Image up to date")
		}
	}

	return nil
}

func filterTags(tags []string, matches []string) []string {

	var filteredTags []string

	for _, pattern := range matches {
		for _, t := range tags {
			matches, _ := regexp.MatchString(pattern, t)

			if matches {
				filteredTags = append(filteredTags, t)
			}
		}
	}

	return filteredTags
}

func getCredentials(registryUrl string) *Credentials {

	credentials, ok := config.Auth[registryUrl]

	if !ok {
		return &Credentials{
			Username: "",
			Password: "",
		}
	}

	return &credentials
}

func newRegistryWithCredentials(url string, auth *Credentials) (*registry.Registry, error) {
	transport := registry.WrapTransport(http.DefaultTransport, url, auth.Username, auth.Password)

	registry := &registry.Registry{
		URL: url,
		Client: &http.Client{
			Transport: transport,
		},
		Logf: registry.Quiet,
	}

	if err := registry.Ping(); err != nil {
		return nil, err
	}

	return registry, nil
}

const IndexName = "registry.hub.docker.com"

func splitReposSearchTerm(reposName string) (string, string) {
	nameParts := strings.SplitN(reposName, "/", 2)
	var indexName, remoteName string
	if len(nameParts) == 1 || (!strings.Contains(nameParts[0], ".") &&
		!strings.Contains(nameParts[0], ":") && nameParts[0] != "localhost") {
		// This is a Docker Index repos (ex: samalba/hipache or ubuntu)
		// 'docker.io'
		indexName = IndexName
		remoteName = reposName
	} else {
		indexName = nameParts[0]
		remoteName = nameParts[1]
	}
	return indexName, remoteName
}
