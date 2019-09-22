
# tar -xzvf ../**

echo -n "Host OS Type: "; grep '"OSType"' default/docker/docker_info.json | awk -F'"' '{print $4}'

OS=$(grep '"OperatingSystem"' default/docker/docker_info.json | awk -F'"' '{print $4}')
echo -n "  Host OS: "; echo $OS | awk -F' ' '{print $1}'
echo -n "  OS Version: "; echo $OS | awk -F' ' '{print $2}'

echo -n "cores="; grep processor default/proc/cpuinfo | wc -l

echo -n "load average (sec): 15 * 60 * "; cat default/commands/loadavg/loadavg | awk -F' ' '{print $3}'

echo -n "Disk Usage - 1k blocks: "; grep overlay default/commands/df/stdout | awk -F' ' '{print $3}'

# grep '"Driver"' default/docker/docker_info.json
echo -n "Driver: "; grep '"Driver"' default/docker/docker_info.json | awk -F'"' '{print $4}'

# grep '"Version"' default/docker/docker_version.json
echo -n "Docker OS Verson: "; grep '"Version"' default/docker/docker_version.json | awk -F'"' '{print $4}'
