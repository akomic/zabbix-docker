Zabbix-Docker
===========================================

Zabbix Docker monitoring integration. Requires cadvisor.

# CAdvisor

```bash
docker run \
  --volume=/:/rootfs:ro \
  --volume=/var/run:/var/run:rw \
  --volume=/sys:/sys:ro \
  --volume=/var/lib/docker/:/var/lib/docker:ro \
  --volume=/dev/disk/:/dev/disk:ro \
  --volume=/dev/kmsg:/dev/kmsg \
  --publish=127.0.0.1:4560:8080 \
  --detach=true \
  --name=cadvisor \
  google/cadvisor:latest
```

For AWS ECS
```
docker run \
  --volume=/:/rootfs:ro \
  --volume=/var/run:/var/run:rw \
  --volume=/sys:/sys:ro \
  --volume=/cgroup:/sys/fs/cgroup:ro \
  --volume=/var/lib/docker/:/var/lib/docker:ro \
  --volume=/dev/disk/:/dev/disk:ro \
  --volume=/dev/kmsg:/dev/kmsg \
  --publish=127.0.0.1:4560:8080 \
  --detach=true \
  --name=cadvisor \
  google/cadvisor:latest
```

# Installation

```
go get
go build
sudo mv zabbix-docker /usr/local/bin
```

# Configuration

```
mkdir ~/.zabbix-docker
cp config.example.yml ~/.zabbix-docker/config.yml
```
Edit ~/.zabbix-docker/config.yml


## Config file

| Variable        | Required | Default                                  | Description                                        |
| --------------- | -------- | ---------------------------------------- | -------------------------------------------------- |
| addr            | Yes      | 127.0.0.1:8080                           | CAdvisor server address                            |
| zabbixAddr      | Yes      | 127.0.0.1:10051                          | Zabbix Server address                              |
| hostname        | No       | Local Hostname                           | Docker Host Hostname (for Zabbix)                  |
| hostGroup1      | No       | Docker Containers                        | Docker label name to use as host group on Zabbix   |
| hostGroup2      | No       | Docker Containers                        | Docker label name to use as host group on Zabbix   |
| hostGroup3      | No       | Docker Containers                        | Docker label name to use as host group on Zabbix   |
| hostGroup4      | No       | Docker Containers                        | Docker label name to use as host group on Zabbix   |

## On Zabbix UI

- Configuration > Action > Auto registration
Add auto-registration based on host metadata "DHost", Operations: add
  host, Link to templates = Template DHost

# Running

Add to crontab

```
5/* * * * * /usr/local/bin/zabbix-docker -z >/dev/null 2>&1
```
