# Core: Troubleshoot Analyzers Coding Challenge

The Replicated Troubleshoot product enables you to collect information from Kubernetes, Docker and the host OS to help analyze and diagnose issues in your cluster. Replicated offers a command that will generate an archive of files and command output collected from your server.

You are tasked with creating a project, in the language of your choice, to analyze a Replicated Troubleshoot bundle and surface insights. You will need to read in a bundle in tar.gz format and output machine readable insights.

The following insights should be extracted from a bundle:

- Host OS
- Host OS version
- Number of cores
- Load average in seconds over the past 15 minutes
- Disk usage in bytes on the root device
- Docker version
- Docker storage driver

We have included a sample bundle in this project. The following command can be run to generate additional bundles:

```bash
docker run -it --rm \
    --name support-bundle \
    -v $PWD:/out \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --net host --pid host -w /out  \
    -e HTTP_PROXY -e HTTPS_PROXY -e NO_PROXY \
    replicated/support-bundle \
    generate --no-upload
```

Treat this project as proof-of-concept. Please don't spend more a day on this assignment even if the result is an incomplete implementation. If you find yourself short on time, focus on a few fully working features end-to-end. Feel free to use third party packages. Also please elaborate on where you would take the project next if you had more time. Don't forget to include instructions for building and running your project. Please deliver your answer as a Github repository.
