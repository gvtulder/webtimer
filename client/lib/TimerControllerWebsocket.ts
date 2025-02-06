// Copyright (C) 2025 Gijs van Tulder
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import { Encoder, Decoder, decode } from "@msgpack/msgpack";

import { TimerController, TimerEvent, TimerEventListener, TimerEventType } from "./TimerController.js";

const encoder = new Encoder();
const decoder = new Decoder();

export class TimerControllerWebsocket implements TimerController {
    private url : string;
    private ws : WebSocket;
    private timeout : number;
    private lastRemainingSeconds : number;
    private listeners : TimerEventListener[];
    private lastEvent : TimerEvent;

    constructor() {
        this.listeners = [];
    }

    connect(url : string) {
        this.url = url;
        this.ws = new WebSocket(this.url);
        let reconnect = () => {
            this.ws.close();
            this.connect(url);
        }
        this.ws.addEventListener('message', async (ev : MessageEvent) => {
            let msg = null;
            if (ev.data instanceof Blob) {
                msg = decoder.decode(await (ev.data as Blob).arrayBuffer());
            } else if (ev.data instanceof ArrayBuffer) {
                msg = await decoder.decodeAsync((ev.data as any).stream());
            }
            if (msg !== null && msg.version) {
                console.log(`Connected to webtimer.cc server ${url} (version ${msg.version})`)
            } else if (msg !== null) {
                this.triggerEvent({
                    type: TimerEventType.UpdateRemainingSeconds,
                    connected: true,
                    active: msg.a,
                    black: msg.b,
                    countdown: msg.c,
                    running: msg.r,
                    remainingSeconds: msg.s,
                    clients: msg.C,
                });
            }
            if (this.timeout) {
                window.clearTimeout(this.timeout);
            }
            this.timeout = window.setTimeout(reconnect, 5000);
        });
        if (this.timeout) {
            window.clearTimeout(this.timeout);
        }
        this.timeout = window.setTimeout(reconnect, 5000);
        this.triggerEvent({
            type: TimerEventType.Connecting,
            connected: false,
        });
    }

    disconnect(): void {
        this.url = null;
        if (this.timeout) {
            window.clearTimeout(this.timeout);
        }
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    private send(msg : any) {
        this.ws.send(encoder.encode(msg));
    }

    setRemainingSeconds(seconds : number) {
        this.send({"cmd": "set", "sec": seconds});
    }

    addRemainingSeconds(seconds : number) {
        if (this.lastEvent && !this.lastEvent.countdown) {
            // if we are not counting down, swap meaning of buttons
            seconds = -seconds;
        }
        this.send({"cmd": "add", "sec": seconds});
    }

    startTimer() {
        this.send({"cmd": "start"});
    }

    pauseTimer() {
        this.send({"cmd": "pause"});
    }

    resetTimer() {
        this.send({"cmd": "reset"});
    }

    toggleBlack() {
        this.send({"cmd": "toggleblack"});
    }

    addListener(listener : TimerEventListener) {
        this.listeners.push(listener);
    }

    triggerEvent(event : TimerEvent) {
        this.lastEvent = event;
        for (const listener of this.listeners) {
            listener(event);
        }
    }
}
