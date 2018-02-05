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

# Configuration

```
mkdir ~/.zabbix-docker
cp config.example.yml ~/.zabbix-docker/config.yml
```
Edit ~/.zabbix-docker/config.yml

# Running

Add to crontab

```
5/* * * * * /home/myuser/zabbix-docker -z >/dev/null 2>&1
```


