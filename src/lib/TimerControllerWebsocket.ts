import { TimerController, TimerEvent, TimerEventListener, TimerEventType } from "./TimerController";

export class TimerControllerWebsocket implements TimerController {
    private url : string;
    private ws : WebSocket;
    private timeout : number;
    private lastRemainingSeconds : number;
    private listeners : TimerEventListener[];
    private lastEvent : TimerEvent;

    constructor(url : string) {
        this.url = url;
        this.listeners = [];
    }

    run() {
        this.ws = new WebSocket(this.url);
        let reconnect = () => {
            this.ws.close();
            this.run();
        }
        this.ws.addEventListener('message', (ev) => {
            if (ev.data == "") {
                // ping
            } else {
                const msg = JSON.parse(ev.data)
                this.triggerEvent({
                    type: TimerEventType.UpdateRemainingSeconds,
                    connected: true,
                    active: msg.active,
                    black: msg.black,
                    countdown: msg.countdown,
                    running: msg.running,
                    remainingSeconds: msg.remaining,
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

    setRemainingSeconds(seconds : number) {
        this.ws.send(JSON.stringify({"cmd": "set", "seconds": seconds}));
    }

    addRemainingSeconds(seconds : number) {
        this.ws.send(JSON.stringify({"cmd": "add", "seconds": seconds}));
    }

    startTimer() {
        this.ws.send(JSON.stringify({"cmd": "start"}));
    }

    pauseTimer() {
        this.ws.send(JSON.stringify({"cmd": "pause"}));
    }

    resetTimer() {
        this.ws.send(JSON.stringify({"cmd": "reset"}));
    }

    toggleBlack() {
        this.ws.send(JSON.stringify({"cmd": "toggleblack"}));
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
