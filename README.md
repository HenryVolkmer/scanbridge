scanbridge
==========

## Motivation

some low-budget Multifunction-Devices (Printer+Scanner) dont support the "scan-to-mail"-Feature and having no PDF-Support at all. Scanbridge is the bridge between the Scanner-Device and E-Mail-Gateway, it also generates PDFs.

## Build

make sure you have [go](https://go.dev/) and [npm](https://www.npmjs.com/) installed on your Buildsystem.

`go run tools/build.go` will build the Frontend and Backend and create a binary `./bin/scanbridge`. The Frontend and all Assets are embeded in that binary.

## Run

You have to specify the Scan-Device with the `-device` Flag.

In order to show all available devices run `scanimage -L`

Something like this should be outputed:

```
device escl:https://192.168.0.157:443' is a HP 0138B60D6ED (HP Color Laser MFP 178 179) platen,adf scanner
``` 

now you can run scanbridge with `scanbridge -device=escl:https://192.168.0.157:443 -bind=127.0.0.1:8080` and open http://127.0.0.1:8080 in your Webbrowser.

## API

`/api/scan` will start a Scan and return the UUID of Scanresult

`/api/download/{uuid}` will download a Scanresult (PDF) by given UUID.


### optional configuration file

one can start scanbridge by passing a configuration file: `scanbridge -config=/path/to/scanbridge.conf`.

A Sample-Configuration can be found [here](./scanbridge.conf.dist).

## systemd unit

move the scanbridge binary to `/usr/local/bin/scanbridge`
Scanbridge needs permissions to write into `/var/tmp` and `/tmp` which is also a POSIX requirement.

```
[Unit]
Description=scanbridge
Wants=network-online.target
After=syslog.target network.target nss-lookup.target network-online.target

[Service]
ExecStart=/usr/local/bin/scanbridge
User=root
Group=root

[Install]
WantedBy=multi-user.target
```