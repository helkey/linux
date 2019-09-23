
FILE STRUCTURE

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
  
