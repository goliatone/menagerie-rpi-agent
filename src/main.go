package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/twinj/uuid"
)

func getDefaultUUID() string {
	if dat, _ := ioutil.ReadFile("./.device_uuid"); dat != nil {
		return string(dat)
	}
	def := uuid.NewV4().String()
	return def
}

func saveUUID(uuid string) {
	ioutil.WriteFile("./.device_uuid", []byte(uuid), 0644)
}

func main() {
	url := flag.String("url", "http://localhost:8080", "The addr of the application")
	uuid := flag.String("uuid", getDefaultUUID(), "The UUID for the application")
	post := flag.Bool("post", false, "If present then we post to service")
	verbose := flag.Bool("verbose", true, "If present will log output")
	flag.Parse()

	saveUUID(*uuid)

	host, _ := os.Hostname()
	metadata, err := GetAddress()
	handleError(err, "Error GetAddress: ")

	metadata["hostname"] = host

	output := make(map[string]interface{})

	output["metadata"] = metadata
	// output["uuid"] = *uuid
	output["name"] = getNameFromHostname(host)
	output["status"] = "online"

	jsonStr, err := json.Marshal(output)

	if *verbose {
		os.Stdout.WriteString(string(jsonStr))
	}

	handleError(err, "Error json.Marshal: ")

	if *post {
		_, err = Post(*url, jsonStr)
		handleError(err, "Error Post: ")
		// body, _ := ioutil.ReadAll(resp.Body)
	}
}

func getNameFromHostname(hostname string) string {
	hostname = strings.TrimSuffix(hostname, ".local")
	return hostname
}

func handleError(err error, msg string) {
	if err != nil {
		os.Stderr.WriteString(msg + err.Error() + "\n")
		os.Exit(1)
	}
}

//Post post post
func Post(url string, jsonData []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}

//GetAddress return address
func GetAddress() (map[string]string, error) {
	output := make(map[string]string)
	inter, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, ifa := range inter {
		addrs, err := ifa.Addrs()

		if err != nil {
			return nil, err
		}

		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					output[string(ifa.Name)] = ipnet.IP.String()
				}
			}
		}
	}
	return output, nil
}
