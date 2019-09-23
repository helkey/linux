
# tar -xzvf ../**

echo -n "Host OS Type: "; grep '"OSType"' default/docker/docker_info.json | awk -F'"' '{print $4}'

OS=$(grep '"OperatingSystem"' default/docker/docker_info.json | awk -F'"' '{print $4}')
echo -n "Host OS: "; echo $OS | awk -F' ' '{print $1}'
echo -n "OS Version: "; echo $OS | awk -F' ' '{print $2}'

echo -n "Number of processors: "; grep processor default/proc/cpuinfo | wc -l

echo -n "Load average (sec): 15 * 60 * "; cat default/commands/loadavg/loadavg | awk -F' ' '{print $3}'

echo -n "Disk Usage - 1k blocks: "; grep overlay default/commands/df/stdout | awk -F' ' '{print $3}'

# NOTE: typical Docker versions are ~1.8.0; this (18.0..) looks like an Ubuntu ver number
echo -n "Docker Version: "; grep '"Version"' default/docker/docker_version.json | awk -F'"' '{print $4}'

echo -n "Docker Driver: "; grep '"Driver"' default/docker/docker_info.json | awk -F'"' '{print $4}'

