package main

import (
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Server struct {
	Url       string `json:"url" yaml:"url"`
	Variables struct {
		Hostname struct {
			Default string `json:"default" yaml:"default"`
		} `json:"hostname" yaml:"hostname"`
	} `json:"variables" yaml:"variables"`
}

type Openapi struct {
	Openapi string   `json:"openapi" yaml:"openapi"`
	Servers []Server `json:"servers" yaml:"servers"`
	Info    struct {
		Version string `json:"version" yaml:"version"`
		Title   string `json:"title" yaml:"title"`
	} `json:"info" yaml:"info"`
	Tags []struct {
		Name string `json:"name" yaml:"name"`
	} `json:"tags" yaml:"tags"`
	Paths      map[string]interface{} `json:"paths" yaml:"paths"`
	Components struct {
		RequestBodies   map[string]interface{} `json:"requestBodies" yaml:"requestBodies"`
		SecuritySchemes map[string]interface{} `json:"securitySchemes" yaml:"securitySchemes"`
		Responses       map[string]interface{} `json:"responses" yaml:"responses"`
		Schemas         map[string]interface{} `json:"schemas" yaml:"schemas"`
	} `json:"components" yaml:"components"`
}

func main() {

	destinationFile := os.Args[1]
	urls := os.Args[2:]

	openapis := make([]Openapi, 0, len(urls))

	for _, url := range urls {
		openapiYaml, err := os.ReadFile(url)
		if err != nil {
			log.Fatal(err)
		}
		var structOpenapi Openapi

		err = yaml.Unmarshal(openapiYaml, &structOpenapi)
		if err != nil {
			log.Fatal(err)
		}

		openapis = append(openapis, structOpenapi)
	}

	finalOpenapi := new(Openapi)

	finalOpenapi.Info.Title = "api"
	finalOpenapi.Info.Version = "1.0.0"
	finalOpenapi.Openapi = "3.0.0"
	finalOpenapi.Servers = append(
		finalOpenapi.Servers,
		Server{
			Url: "http://{hostname}/api",
			Variables: struct {
				Hostname struct {
					Default string `json:"default" yaml:"default"`
				} `json:"hostname" yaml:"hostname"`
			}{
				Hostname: struct {
					Default string `json:"default" yaml:"default"`
				}{
					Default: "localhost:8080",
				},
			},
		},
		Server{
			Url: "https://{hostname}/api",
			Variables: struct {
				Hostname struct {
					Default string `json:"default" yaml:"default"`
				} `json:"hostname" yaml:"hostname"`
			}{
				Hostname: struct {
					Default string `json:"default" yaml:"default"`
				}{
					Default: "localhost",
				},
			},
		},
	)

	for _, o := range openapis {
		finalOpenapi.Paths = utils.SafeMerge(finalOpenapi.Paths, o.Paths)
		finalOpenapi.Components.RequestBodies = utils.SafeMerge(finalOpenapi.Components.RequestBodies, o.Components.RequestBodies)
		finalOpenapi.Components.SecuritySchemes = utils.SafeMerge(finalOpenapi.Components.SecuritySchemes, o.Components.SecuritySchemes)
		finalOpenapi.Components.Responses = utils.SafeMerge(finalOpenapi.Components.Responses, o.Components.Responses)
		finalOpenapi.Components.Schemas = utils.SafeMerge(finalOpenapi.Components.Schemas, o.Components.Schemas)
		finalOpenapi.Tags = append(finalOpenapi.Tags, o.Tags...)
	}

	d, err := yaml.Marshal(finalOpenapi)

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(destinationFile, d, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("created file %s", destinationFile)
}
