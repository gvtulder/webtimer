import textFit from "textfit";
import { QR } from "qr-svg";

import { formatTime } from "../lib/format.js";
import { TimerEvent, TimerEventType, TimerController } from "../lib/TimerController.js";

export class TimerDisplay {
    element : HTMLDivElement;

    private controller : TimerController;
    private timeDisplay : TimeDisplay;
    private connecting : HTMLDivElement;
    private statusBar : HTMLDivElement;

    constructor(controller : TimerController, backUrl : string) {
        this.controller = controller;
        this.build(backUrl);
        controller.addListener((event: TimerEvent) => { this.handleTimerEvent(event) });
    }

    build(backUrl : string) {
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

        // qr link
        const qr = document.createElement('div');
        qr.className = 'qr';
        qr.innerHTML = QR(window.location.href, 'M');
        const qrp = document.createElement('p');
        qrp.innerHTML = window.location.href;
        qr.appendChild(qrp);
        div.appendChild(qr);
        qr.addEventListener('click', (evt) => {
            qr.classList.remove('visible');
            evt.stopPropagation();
            evt.preventDefault();
            return false;
        });

        // menu
        const menu = document.createElement('div');
        menu.className = 'menu';
        div.appendChild(menu);

        menu.addEventListener('click', (evt) => {
            if (evt.target && (evt.target as HTMLElement).closest('.show-fullscreen')) {
                return;
            }
            evt.stopPropagation();
        });

        // menu buttons
        const menuLeft = document.createElement('div');
        menuLeft.className = 'left';
        menu.appendChild(menuLeft);

        const statusBar = document.createElement('div');
        statusBar.className = 'status';
        statusBar.innerHTML = '';
        menu.appendChild(statusBar);
        this.statusBar = statusBar;

        const menuRight = document.createElement('div');
        menuRight.className = 'right';
        menu.appendChild(menuRight);

        // show exit
        const exit = document.createElement('button');
        exit.className = 'exit-timer';
        exit.appendChild((document.getElementById('template-icon-exit') as HTMLTemplateElement).content.cloneNode(true));
        menuLeft.appendChild(exit);
        exit.addEventListener('click', (evt) => {
            window.location.href = backUrl;
            evt.stopPropagation();
        });

        // show qr code
        const qrShow = document.createElement('button');
        qrShow.className = 'qr-show';
        qrShow.appendChild((document.getElementById('template-icon-qr') as HTMLTemplateElement).content.cloneNode(true));
        menuRight.appendChild(qrShow);
        qrShow.addEventListener('click', (evt) => {
            qr.classList.toggle('visible');
            evt.stopPropagation();
            evt.preventDefault();
            return false;
        });

        // show full screen
        const fullScreen = document.createElement('button');
        fullScreen.className = 'show-fullscreen';
        fullScreen.appendChild((document.getElementById('template-icon-expand') as HTMLTemplateElement).content.cloneNode(true));
        menuRight.appendChild(fullScreen);
        fullScreen.addEventListener('click', (evt) => {
            qr.classList.remove('visible');
            evt.preventDefault();
        });
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

            if (event.clients == 2) {
                this.statusBar.innerHTML = '+1 connected';
            } else if (event.clients > 2) {
                this.statusBar.innerHTML = `+${event.clients - 1} connected`;
            } else {
                this.statusBar.innerHTML = '';
            }

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
        new ResizeObserver(() => { this.updateTextFit(); }).observe(this.element);
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
