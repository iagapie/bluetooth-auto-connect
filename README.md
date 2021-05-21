# bluetooth-auto-connect

A linux command line tool try to automatically connect to all paired and trusted bluetooth devices.

## Build

```sh
make build
```

## Installation

### Manual

* Install the `bluetooth-auto-connect` app in `/usr/bin/`
* If you are using systemd, consider installing the `init/bluetooth-auto-connect.service` file to `/etc/systemd/system/`
    - Enable the service with `sudo systemctl enable bluetooth-auto-connect`

#### PulseAudio

If you are using a bluetooth headset with pulseaudio running as your user, then the above service will fail to connect to your headset on boot because pulseaudio won't have been started when to bluetooth service comes up.

* If you are using systemd, consider installing the `init/pulseaudio-bluetooth-auto-connect.service` file to `/etc/systemd/user/` or `~/.config/systemd/user/`
    - Enable the service with `systemctl --user enable pulseaudio-bluetooth-auto-connect`

## License

MIT License

Copyright (c) 2021 Igor Agapie