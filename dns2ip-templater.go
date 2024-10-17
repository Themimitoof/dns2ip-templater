package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Services []string `yaml:"services"`
	Ranges   []string `yaml:"ranges"`
}

func readConfig(filename string) (Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("error reading configuration file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error during YAML file decoding: %w", err)
	}

	return config, nil
}

func lookupIPs(service string) ([]string, error) {
	ips, err := net.LookupIP(service)
	if err != nil {
		return nil, fmt.Errorf("error resolving DNS for %s: %w", service, err)
	}

	var ipStrings []string
	for _, ip := range ips {
		ipStrings = append(ipStrings, ip.String())
	}

	return ipStrings, nil
}

func renderTemplate(templateFile, outputFile string, data interface{}) error {
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return fmt.Errorf("error rendering template: %w", err)
	}

	return nil
}

func exec(confFile, templateFile, outputFile *string) {
	config, err := readConfig(*confFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Prepare the data for the renderer
	data := make(map[string][]string)

	for _, service := range config.Services {
		ips, err := lookupIPs(service)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		data[service] = ips
	}

	if len(config.Ranges) > 0 {
		data["Ranges"] = []string{}

		for _, r := range config.Ranges {
			data["Ranges"] = append(data["Ranges"], r)
		}
	}

	// Render the template
	if renderTemplate(*templateFile, *outputFile, data) != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	confFile := flag.String("conf", "config.yml", "Configuration file path")
	templateFile := flag.String("template", "template.tmpl", "Template file path")
	outputFile := flag.String("output", "output.txt", "Rendered file path")
	execInterval := flag.Duration("interval", 0, "Define the interval to run the routine (by default: executed once)")
	flag.Parse()

	if *execInterval > 0 {
		for {
			fmt.Printf("Executing new iteration of dns2ip-templater... ")
			exec(confFile, templateFile, outputFile)
			fmt.Println("Done.")
			time.Sleep(*execInterval)
		}
	} else {
		exec(confFile, templateFile, outputFile)
	}
}
