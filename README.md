# webtimer.cc

[webtimer.cc](https://webtimer.cc/) is a simple, free, web-based presentation timer with remote control.

## Use the online version

<img alt="Screenshots of webtimer.cc on various devices." src="https://webtimer.cc/static/montage-200px.png" style="float:right;margin:1em;">

Visit [webtimer.cc](https://webtimer.cc/).

- Use the timer to **show the remaining time** to presenters at a single presentation, a workshop, or a small conference.
- **Control the timer from your phone.** Or use any other device with a web browser.
- Easy to set up. **Works anywhere with wifi.**

### Instructions

1. Start a new timer:

   - Go to [webtimer.cc](https://webtimer.cc/) in any web browser.
   - Choose a unique **session ID** to identify your timer.

2. Open the timer on another device. Connect as many devices as you like:

   - Scan the **QR code** to open the timer.
   - Or go to [webtimer.cc](https://webtimer.cc/) and enter your session ID.

3. Control the timer from any connected device:

   - Set a **countdown** timer with the allotted time.
   - **Start** or **stop** the timer at any moment.
   - **Add** or **subtract** to adjust the time.
   - Switch to a **black screen** to hide the clock.

## Run your own version

If you wish, you can run your own webtimer server locally.

### Instructions

1. Download the precompiled binary for your computer:

   | Linux                                                                                               | Windows                                                                                             | macOS                                                                                                 |
   | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
   | [Linux (AMD64)](https://github.com/gvtulder/webtimer/releases/latest/download/webtimer)             | [Windows (AMD64)](https://github.com/gvtulder/webtimer/releases/latest/download/webtimer.exe)       | [mac OS (AMD64)](https://github.com/gvtulder/webtimer/releases/latest/download/webtimer_darwin-amd64) |
   | [Linux (ARM64)](https://github.com/gvtulder/webtimer/releases/latest/download/webtimer_linux-arm64) | [Windows (ARM64)](https://github.com/gvtulder/webtimer/releases/latest/download/webtimer_arm64.exe) | [mac OS (ARM64)](https://github.com/gvtulder/webtimer/releases/latest/download/webtimer_darwin-arm64) |

2. Start the server.

3. Open the webtimer in your browser by accessing <http://localhost:8000/>.

   _Note:_ If you want to control the timer from another device, you might need to find the IP address of your computer. The server will list the addresses on which it is listening. You might also need to adjust your firewall settings to make the server accessible from other devices. By default, the app will listen to port 8000 on any IP address. Run with `webtimer --addr <IP>:<PORT>` to choose something else.

## Host with Docker

A Docker container is available at [ghcr.io/gvtulder/webtimer](https://ghcr.io/gvtulder/webtimer).

```bash
docker pull ghcr.io/gvtulder/webtimer:latest
```

With Docker Compose you might use:

```yaml
services:
  web:
    image: ghcr.io/gvtulder/webtimer:latest
    container_name: webtimer
    ports:
      - "8000:8000"
    restart: unless-stopped
```

## Command-line parameters

| Parameter            | Description                                                                                 |
| -------------------- | ------------------------------------------------------------------------------------------- |
| `--addr <IP>:<PORT>` | Run the server on a specific IP address and port. (Default: `:8000`)                        |
| `--web <DIR>`        | Load HTML, JS and CSS from this directory. (Default: use the files included in the binary.) |

## Building from scratch

webtimer is built using Go and TypeScript.

Installing requirements:

```bash
go get
npm ci
```

Building client and server:

```bash
npx webpack
go build .
```

## About webtimer.cc

[webtimer.cc](https://webtimer.cc/) is produced and hosted by [Gijs van Tulder](https://www.vantulder.net/) in the Netherlands.

## License

The webtimer sources are licensed under the [GPL-3.0 license](LICENSE.txt).

Copyright (C) 2025 Gijs van Tulder

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
