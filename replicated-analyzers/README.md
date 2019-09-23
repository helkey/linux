# **
## Linux / Kubernetes Information

This project extracts a summary of Kubernetes system status, by parsing key data from Linux process file
`/proc/cpuinfo`, Linux commands `df` amd `loadavg`. and custom Replicated json files 
`docker_info.json` and `docker_version.json` summarizing Docker information.

The Replicated Troubleshoot product enables you to collect information from Kubernetes, Docker and the host OS 
to help analyze and diagnose issues in your cluster. Replicated offers a command that will generate an archive of files 
and command output collected from your server. This project extracts from this archieve bundle:
  - Host OS
  - Host OS version
  - Number of cores
  - Load average in seconds over the past 15 minutes
  - Disk usage in bytes on the root device
  - Docker version
  - Docker storage driver

## Bash Data Extraction
The location of this key summary information within this bundle is:
```
default
  commands
    df:
	  stdout: disk free
	loadavg
	  loadavg
        $ 0.26 0.14 0.05 5/233 5186
	docker
    docker_info.json
	  "Driver": "overlay2"
	  "OSType": "linux"
	  "OperatingSystem": "Ubuntu 18.04.2 LTS"
	docker_version.json
	  "Version": "18.09.6"
  proc/cpuinfo: processor:0 to processor:3
      grep processor /proc/cpuinfo | wc -l
```

The simplest solution processes the data using Bash shell commands, by extracting the data in the shell using `tar`, 
then finding the relevant items in the extracted files using `grep` and `awk`, to produce machine-readable output:
```sh
# tar -xzvf ../**
echo -n "Host OS Type: "; grep '"OSType"' default/docker/docker_info.json | awk -F'"' '{print $4}'
OS=$(grep '"OperatingSystem"' default/docker/docker_info.json | awk -F'"' '{print $4}')
echo -n "  Host OS: "; echo $OS | awk -F' ' '{print $1}'
echo -n "  OS Version: "; echo $OS | awk -F' ' '{print $2}'
echo -n "cores="; grep processor default/proc/cpuinfo | wc -l
echo -n "load average (sec): 15 * 60 * "; cat default/commands/loadavg/loadavg | awk -F' ' '{print $3}'
echo -n "Disk Usage - 1k blocks: "; grep overlay default/commands/df/stdout | awk -F' ' '{print $3}'
echo -n "Driver: "; grep '"Driver"' default/docker/docker_info.json | awk -F'"' '{print $4}'
echo -n "Docker OS Verson: "; grep '"Version"' default/docker/docker_version.json | awk -F'"' '{print $4}'
```
```sh
Host OS Type: linux
  Host OS: Ubuntu
  OS Version: 18.04.2
cores=4
load average (sec): 15 * 60 * 0.05
Disk Usage - 1k blocks: 2050392
Driver: overlay2
Docker OS Verson: 18.09.6
```
It is difficult to reliabily distribute shell files for data processing, as different systems provide different shells.
This Bash approach to data extraction also is brittle - it relies on the particular file structure that the data is written in.

## Go(lang) Data Extraction
A [more robust solution](https://github.com/helkey/linux/blob/master/replicated-analyzers/README.md) would be to first parse 
the JSON files, then extract parameters as key/value pairs. Such an approach can be written in a modern programming language such as Go, 
which can be more easily distributed, and which also provides (at least in principal) better error handling, as well as allowing
better integration with other custom analyis tools.
```sh
go run extract.go supportbundle.tar.gz
```

```go
// extract.go
// 
package main

import...

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
...
```
which produces output:
```sh
OS Type: linux
Host OS: Ubuntu
Host OS Version: 18.04.2
Num proc: 4; Num cores: 8
Load average (over 15min): 45 sec
Disk Usage: 2099601408 bytes
Docker Version: 18.09.6
Docker Driver: overlay2
```

## Go Program Distribution
One goal for this project is to build free-standing executables. The standard Go executables are dynamically linked to the Go runtime.
```sh
go build extract.go
```
Go can also produce statically linked applications, e.g.:
```sh

```
and cross-compiled for other platforms.
```sh
```

### Docker
This analysis application can be packaged up as a Docker executable:
```
```

### AWS / Packer
This application can also be run on cloud computing hardware such as Amazon Web Services (AWS),
which can run the previous Docker image, or a custom AMI can be built using [Hashicorp Packer]().



