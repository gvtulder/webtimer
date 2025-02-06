import { TimerController, TimerEvent, TimerEventListener, TimerEventType } from "./TimerController";

export class TimerControllerOffline implements TimerController {
    private timer : Timer;
    private lastRemainingSeconds : number;
    private listeners : TimerEventListener[];

    constructor() {
        this.timer = new Timer();
        this.listeners = [];
        window.setInterval(() => { this.checkForTimerUpdates() }, 200);
    }

    checkForTimerUpdates() {
        const s = this.timer.remainingSeconds;
        if (s != this.lastRemainingSeconds) {
            this.lastRemainingSeconds = s;
            this.triggerEvent({
                type: TimerEventType.UpdateRemainingSeconds,
                running: this.timer.running,
                remainingSeconds: s,
            });
        }
    }

    setRemainingSeconds(seconds : number) {
        this.timer.remainingSeconds = seconds;
        this.checkForTimerUpdates();
    }

    addRemainingSeconds(seconds : number) {
        this.timer.addRemainingSeconds(seconds);
        this.checkForTimerUpdates();
    }

    startTimer() {
        this.timer.startTimer();
        this.triggerEvent({ type: TimerEventType.StartTimer, running: this.timer.running });
        this.checkForTimerUpdates();
    }

    pauseTimer() {
        this.timer.pauseTimer();
        this.triggerEvent({ type: TimerEventType.PauseTimer, running: this.timer.running });
        this.checkForTimerUpdates();
    }

    resetTimer() {
        this.timer.resetTimer();
        this.triggerEvent({ type: TimerEventType.PauseTimer, running: this.timer.running });
        this.checkForTimerUpdates();
    }

    addListener(listener : TimerEventListener) {
        this.listeners.push(listener);
    }

    triggerEvent(event : TimerEvent) {
        for (const listener of this.listeners) {
            listener(event);
        }
    }
}

class Timer {
    private remainingTimeAtStart : number;
    private startTime : number;

    constructor() {
        this.remainingTimeAtStart = null;
        this.startTime = null;
    }

    get running() : boolean {
        return this.startTime !== null;
    }

    get remainingSeconds() : number {
        if (this.remainingTimeAtStart === null) {
            return null;
        }
        let remaining = this.remainingTimeAtStart;
        if (this.startTime !== null) {
            remaining -= new Date().getTime() - this.startTime;
        }
        return Math.round(remaining / 1000);
    }

    set remainingSeconds(seconds : number) {
        if (seconds === null) {
            this.remainingTimeAtStart = null;
            this.startTime = null;
        } else {
            this.remainingTimeAtStart = seconds * 1000;
            if (this.startTime !== null) {
                this.startTime = new Date().getTime();
            }
        }
    }

    addRemainingSeconds(seconds : number) {
        if (this.remainingTimeAtStart === null) {
            this.remainingTimeAtStart = 0;
        }
        this.remainingTimeAtStart += seconds * 1000;
    }

    resetTimer() {
        this.startTime = null;
        this.remainingTimeAtStart = null;
    }

    startTimer() {
        if (this.startTime === null && this.remainingSeconds !== null) {
            this.startTime = new Date().getTime();
        }
    }

    pauseTimer() {
        if (this.startTime !== null) {
            this.remainingTimeAtStart -= (new Date().getTime() - this.startTime);
            this.startTime = null;
        }
    }
}
