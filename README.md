# Jellyfin Server Announcer
[![platform: gitlab](https://img.shields.io/badge/gitlab-cromefire__-%23e65328?logo=gitlab)](https://gitlab.com/cromefire_/jf-server-announcer)
[![pipeline status](https://gitlab.com/cromefire_/jf-server-announcer/badges/main/pipeline.svg)](https://gitlab.com/cromefire_/jf-server-announcer/-/pipelines)

This is a small application that may be used to advertise a Jellyfin Server under a different address / name or inside another local network.

## Use cases
This is useful in 2 cases primarily:
- Your server ist behind a reverse proxy and not (securely) reachable for it local ip address that it advertises
- Your want to advertise a server inside you network that is not inside you network (e.g. a friend's server)

## Usage
The app needs 3 parameters:
- `--id`: The Jellyfin server ID (see below)
- `--address`: The address the clients should use to communicate with Jellyfin (include the protocol and (if applicable) the base url, e.g. `https://jellyfin.example.com` or `http://1.2.3.4:8006/jellyfin`)
- `--name`: The Jellyfin server name as it should be shown to the clients, useful, if you have multiple servers being announced in the local network

You can get the current parameters by navigating to your Jellyfin instance, logging in and pasting the following into your browser console, you will sometimes want to change your address though.<br>
_(only do this, if you can understand what this code does, your browser will also warn you about not pasting untrusted code)_
```javascript
JSON.parse(localStorage.getItem("jellyfin_credentials")).Servers.forEach((e, i) => console.info(`Server ${i + 1}: --id="${e.Id}" --name="${e.Name}" --address="${e.ManualAddress}"`));
```

## Installation
_It's only tested on linux so far and binaries are only available for amd64, it should word on other platforms though._

You can fetch the binary from the [packages](https://gitlab.com/cromefire_/jf-server-announcer/-/packages) and make it executable.

For the sake of security it's probably best to add a user for it, so it's not running as root:
```shell
sudo adduser --system --disabled-login --group --home /var/lib/jf-server-announcer jf-server-announcer
```

You also have to open port `7359/udp` if you have a fire wall, because this is the port it uses to communicate

If you use systemd, here's a sample unit file (with a bit of hardening applied):
<details>
<summary>jf-server-announcer.service</summary>

```unit file (systemd)
[Unit]
Description=Jellyfin Server Announcer
After=network-online.target

[Service]
Type=simple
ExecStart=/path/to/jf-server-announcer --id "..." --address "https://..." --name "..."
Restart=on-failure
RestartSec=60s

# Hardening
User=jf-server-announcer
Group=jf-server-announcer
PrivateTmp=true
ProtectSystem=strict
NoNewPrivileges=true
RestrictNamespaces=uts ipc pid user cgroup
ProtectKernelTunables=yes
ProtectKernelModules=yes
ProtectControlGroups=yes
PrivateDevices=yes
RestrictSUIDSGID=true

[Install]
WantedBy=multi-user.service
```
</details>

## Building

This is a standard go modules project, so you just have to have [go](https://golang.org) installed and run `go build`.
