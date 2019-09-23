// extract.go
// go run extract.go supportbundle.tar.gz
package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide sourceFile name")
		return
	}
	sourceFile := os.Args[1] // "supportbundle.tar.gz" 
	processGzip(sourceFile)
}

// Defines order of output parameters (not same order as tar scan)
const (
	iHostOSType = iota
	iHostOS
	iHostOSVersion
	iCores
	iLoadAve
	iDiskUsage
	iDockerVer
	iDockerDriver
)
var output = make([]string, iDockerDriver + 1)

func processGzip(fName string) (err error) {
	f, err := os.Open(fName)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gz)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// fmt.Println("reached EOF")
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg {
			switch header.Name {
			case "default/docker/docker_info.json":
				keys := []string{"OSType", "OperatingSystem", "Driver"}
				vals, _ := json_keys(keys, header, tarReader)
				output[iHostOSType] = "OS Type: " + vals[0]

				// e.g. "Ubuntu 18.04.2 LTS"
				substring := strings.Split(vals[1], " ")
				hostOs, osVer := substring[0], substring[1]

				output[iHostOS] = "Host OS: " + hostOs
				output[iHostOSVersion] = "Host OS Version: " + osVer
				output[iDockerDriver] = "Docker Driver: " + vals[2]

			case "default/proc/cpuinfo":
				nProc, nCores := count_cores(header, tarReader)
				output[iCores] = fmt.Sprintf("Num proc: %v; Num cores: %v", nProc, nCores)
			case "default/commands/loadavg/loadavg":
				loadAve := load_ave(header, tarReader)
				output[iLoadAve] = fmt.Sprintf("Load average (over 15min): %v sec", loadAve)
			case "default/commands/df/stdout":
				diskByte := disk_usage(header, tarReader)
				output[iDiskUsage] = fmt.Sprintf("Disk Usage: %v bytes", diskByte)

			case "default/docker/docker_version.json":
				val, _ := json_key_val("Version", header, tarReader)
				output[iDockerVer] = "Docker Version: " + val
			}
		}
	}
	// Print output values in correct order
	for _, val := range output {
		if len(val) > 0 {
			fmt.Println(val)
		}
	}
	return
}

// Load average in seconds over the past 15 minutes
//  Note: Load factor already accounts for number of cores
func load_ave(header *tar.Header, tarReader *tar.Reader) string {
	// example: 0.26 0.14 0.05 5/233 5186 
	const i15 = 2     // 3rd num is 15min load factor
	const numMin = 15 // number of minutes during average
	const toSec = 60  // conv minutes to seconds
	data := tar_data(header, tarReader)
	numbers := strings.Split(string(data), " ")
	if len(numbers) <= i15 {
		return ""
	}
	// Convert to CPU-seconds over 15min interval
	load15, _ := strconv.ParseFloat(string(numbers[i15]), 64)
	cpuSeconds := int(float64(numMin) * float64(toSec) * load15)
	return fmt.Sprintf("%v", cpuSeconds)
	
}

// Disk usage in bytes on the root device
func disk_usage(header *tar.Header, tarReader *tar.Reader) string {
	const block_1k = 1024 // usage given in 1k (1024byte) blocks
	data := tar_data(header, tarReader)
	overlayLines := match_substring(string(data), "overlay")
	if len(overlayLines) < 1 {
		return "" // "overlay" line not found
	}
	columns := strings.Fields(overlayLines[0])

	// Disk usage in column 3
	disk_used, _ := strconv.Atoi(columns[2])
	if len(columns) < 3 {
		return "" // Disk usage data not found
	}
	disk_used_bytes := strconv.Itoa(block_1k * disk_used)
	return disk_used_bytes
}

// Unmarshal key:value pairs from json file
func json_keys(keys []string, header *tar.Header, tarReader *tar.Reader) ([]string, error) {
	// str_json, err := extract_json("OSType", tarReader)
	data := tar_data(header, tarReader)
	m := make(map[string]string)
	_ = json.Unmarshal([]byte(data), &m)
	var vals = make([]string, len(keys))
	for i, key := range keys {
		vals[i] = m[key]
	}
	return vals, nil
}

// Unmarshal single key:value from json file
func json_key_val(key string, header *tar.Header, tarReader *tar.Reader) (string, error) {
	keys := []string{key}
	vals, err := json_keys(keys, header, tarReader)
	return vals[0], err
}

// Number of cores
func count_cores(header *tar.Header, tarReader *tar.Reader) (int, int) {
	data := tar_data(header, tarReader)
	coreLines := match_substring(string(data), "cpu cores")
	nProc := len(coreLines) // number of processors
	nCores := 0
	for _, line := range coreLines {
		substrings := strings.Split(line, ":")
		if len(substrings) == 2 {
			cores, _ := strconv.Atoi(strings.TrimSpace(substrings[1]))
			nCores += cores
		}
	}
	return nProc, nCores
}

// Return array of lines containing a substring
func match_substring(data string, substring string) []string {
	var matches []string
	// Check for Windows/Mac line endings
	lines := strings.Split(strings.Replace(data, "\r\n", "\n", -1), "\n")
	for _, line := range lines {
		if strings.Contains(line, substring) {
			matches = append(matches, line)
		}
	}
	return matches
}

// Extract data from file within tar archive
func tar_data(header *tar.Header, tarReader *tar.Reader) []byte {
	data := make([]byte, header.Size)
	tarReader.Read(data)
	return data
}
