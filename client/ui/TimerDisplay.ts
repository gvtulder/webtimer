import textFit from "textfit";

import { formatTime } from "../lib/format.js";
import { TimerEvent, TimerEventType, TimerController } from "../lib/TimerController.js";

export class TimerDisplay {
    element : HTMLDivElement;

    private controller : TimerController;
    private timeDisplay : TimeDisplay;
    private connecting : HTMLDivElement;

    constructor(controller : TimerController) {
        this.controller = controller;
        this.build();
        controller.addListener((event: TimerEvent) => { this.handleTimerEvent(event) });
    }

    build() {
        const div = document.createElement('div');
        div.className = 'timer-display';
        this.element = div;

        const card = document.createElement('div');
        card.className = 'card';
        div.appendChild(card);

        // countdown showing remaining time
        this.timeDisplay = new TimeDisplay();
        card.appendChild(this.timeDisplay.element);

        // connecting...
        const connecting = document.createElement('div');
        connecting.innerHTML = 'Connecting...';
        connecting.className = 'connecting';
        div.appendChild(connecting);
        this.connecting = connecting;
    }

    showTime(seconds : number) {
        this.timeDisplay.showTime(seconds);
    }

    handleTimerEvent(event : TimerEvent) {
        if (event.type == TimerEventType.UpdateRemainingSeconds) {
            this.showTime(event.active ? event.remainingSeconds : 0);

            this.element.classList.toggle('black', event.black);
            this.element.classList.toggle('countdown', event.countdown);
            this.element.classList.toggle('running', event.active && event.running);
            this.element.classList.toggle('timeout', event.active && event.remainingSeconds !== null && event.remainingSeconds <= 0);
            this.element.classList.toggle('warning', event.active && event.remainingSeconds !== null && event.remainingSeconds <= 60 && event.remainingSeconds > 0);
            this.connecting.classList.remove('disconnected');
        } else if (event.type == TimerEventType.Connecting) {
            this.connecting.classList.add('disconnected');
        }
    }
}

class TimeDisplay {
    element : HTMLDivElement;

    private currentText : string;

    constructor() {
        this.build();

        window.addEventListener('resize', () => { this.updateTextFit(); });
    }

    build() {
        const div = document.createElement('div');
        div.className = 'time';
        this.element = div;
    }

    showTime(seconds : number) {
        const newString = seconds === null ? '' : formatTime(seconds);
        if (newString != this.currentText) {
            this.element.innerHTML = this.currentText = newString;
            this.updateTextFit();
        }
    }

    updateTextFit() {
        textFit(this.element, { alignVert: true, alignHoriz: true, detectMultiLine: false, maxFontSize: 10000 });
    }
}
