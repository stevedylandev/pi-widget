# Pi Widget

<p align="center">
  <img src="https://dweb.mypinata.cloud/ipfs/QmXL9vahr78uxmRQ4LNEFB7k2rQRP8wgYg4jcvZgdVfBPz?img-format=webp" alt="Alt text" width="200" height="auto" style="border-radius:10px;">
</p>

This little Golang server can be used as a starting point to display your Raspberry Pi stats or any other information you want to display via the web. For my case it also shares information about my [IPFS](https://ipfs.io) node.

It uses some simple functions to get the data and then sends it through a server-sides events endpoint called  `/events`. There’s a separate `index.html` file that handles the data and renders it on the page. For ease or editing I have made this file embedded so it will be included during the build process.

## Usage and Deployment

If you plan to use this for yourself then you will likely want to remove the IPFS stats code from the repo (unless you too happen to be a nerd like me and enjoy messing with that stuff).

Building is pretty straight forward, however you will want to make sure you have the right target build using the following command:

```
GOOS=linux GOARCH=arm GOARM=6 go build -o ~/pi-widget .
```

You can either build this on device or move it over via SMTP. Once on the Pi it can be spun up simply by running `./pi-widget`. To make this more persistent to you can set up a service like so:

1. Make a new file under `/etc/systemd/user/pi-widget.service` and put the following contents inside.
```
[Unit]
Description=A little Pi Widget
Documentation=https://stevedylan.dev
After=network.target

[Service]
Type=simple
ExecStart=/home/steve/pi-widget
Restart=on-failure
KillSignal=SIGINT

[Install]
WantedBy=multi-user.target
```

2. Run the following commands to enable, run, and persist the daemon

```
systemctl --user enable pi-widget
systemctl --user start pi-widget
loginctl enable-linger $USER
```

With the daemon running out of port `4321` you can setup a [Cloudflare Tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/) pointing to `http://localhost:4321` and have a working site with your stats!
