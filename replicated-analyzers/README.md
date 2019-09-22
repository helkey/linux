The Replicated Troubleshoot product enables you to collect information from Kubernetes, Docker and the host OS to help analyze and diagnose issues in your cluster. Replicated offers a command that will generate an archive of files and command output collected from your server. Extract from bundle:
  - Host OS
  - Host OS version
  - Number of cores
  - Load average in seconds over the past 15 minutes
  - Disk usage in bytes on the root device
  - Docker version
  - Docker storage driver

The simplest solution processes the data using Bash shell commands, by extracting the data in the shell using `tar`, 
then finding the relevant items in the extracted files using `grep` and `awk`.
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
to produce machine-readable output (which needs more formatting work):
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
