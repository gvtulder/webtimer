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

import { Encoder, Decoder } from "@msgpack/msgpack";

import {
    TimerController,
    TimerEvent,
    TimerEventListener,
    TimerEventType,
} from "./TimerController.js";

const encoder = new Encoder();
const decoder = new Decoder();

/**
 * Manages the WebSocket connection with the timer server.
 *
 * Communication is handled through msgpack messages.
 * The controller automatically attempts to reconnect if the connection is lost.
 */
export class TimerControllerWebsocket implements TimerController {
    private url: string;
    private ws: WebSocket;
    private timeout: number;
    private listeners: TimerEventListener[];
    private lastEvent: TimerEvent;

    constructor() {
        this.listeners = [];
    }

    /**
     * Connect to the timer server at the given URL.
     * @param url a websocket URL
     */
    connect(url: string): void {
        this.url = url;
        this.ws = new WebSocket(this.url);
        const reconnect = (): void => {
            this.ws.close();
            this.connect(url);
        };
        this.ws.addEventListener("message", async (ev: MessageEvent) => {
            let msg = null;
            if (ev.data instanceof Blob) {
                msg = decoder.decode(await (ev.data as Blob).arrayBuffer());
            } else if (ev.data instanceof ArrayBuffer) {
                msg = decoder.decode(ev.data);
            }
            if (msg !== null && msg.version) {
                console.log(
                    `Connected to webtimer.cc server ${url} (version ${msg.version})`,
                );
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

    private send(msg: unknown): void {
        this.ws.send(encoder.encode(msg));
    }

    setRemainingSeconds(seconds: number): void {
        this.send({ cmd: "set", sec: seconds });
    }

    addRemainingSeconds(seconds: number): void {
        if (this.lastEvent) {
            if (seconds < 0 && this.lastEvent.remainingSeconds == 0) {
                return;
            }
            if (
                !this.lastEvent.countdown &&
                this.lastEvent.remainingSeconds < 0
            ) {
                if (this.lastEvent.remainingSeconds - seconds > 0) {
                    this.send({ cmd: "set", sec: 0 });
                    return;
                }
                seconds = -seconds;
            }
        }
        this.send({ cmd: "add", sec: seconds });
    }

    startTimer(): void {
        this.send({ cmd: "start" });
    }

    pauseTimer(): void {
        this.send({ cmd: "pause" });
    }

    resetTimer(): void {
        this.send({ cmd: "reset" });
    }

    toggleBlack(): void {
        this.send({ cmd: "toggleblack" });
    }

    addListener(listener: TimerEventListener): void {
        this.listeners.push(listener);
    }

    triggerEvent(event: TimerEvent): void {
        this.lastEvent = event;
        for (const listener of this.listeners) {
            listener(event);
        }
    }
}
