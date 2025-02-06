import { TimerController, TimerEvent, TimerEventType } from "../lib/TimerController.js";

export class ControlsDisplay {
    element : HTMLDivElement;

    private controller : TimerController;

    private btnStart : Button;
    private btnPause : Button;
    private btnReset : Button;
    private btnBlack : Button;

    constructor(controller : TimerController) {
        this.controller = controller;
        this.build();
        controller.addListener((event: TimerEvent) => { this.handleTimerEvent(event) });
    }

    build() {
        const div = document.createElement('div');
        div.className = 'controls-display';
        this.element = div;

        const ternary = document.createElement('div');
        ternary.className = 'controls ternary';
        div.appendChild(ternary);

        const secondary = document.createElement('div');
        secondary.className = 'controls secondary';
        div.appendChild(secondary);

        const main = document.createElement('div');
        main.className = 'controls main';
        div.appendChild(main);

        const addButton = (tgt : HTMLElement, name : string, text : string, icon : string, fn : ButtonHandler) => {
            const btn = new Button(name, text, icon, fn);
            tgt.appendChild(btn.element);
            return btn;
        };

        this.btnReset = addButton(
            main,
            'reset', 'Reset', 'template-icon-reset',
            () => { this.controller.resetTimer() }
        );

        this.btnPause = addButton(
            main,
            'pause', 'Pause', 'template-icon-pause',
            () => { this.controller.pauseTimer() }
        );
        this.btnPause.disable();
        
        this.btnStart = addButton(
            main,
            'start', 'Start', 'template-icon-start',
            () => { this.controller.startTimer() }
        );

        this.btnBlack = addButton(
            main,
            'black', 'Black', 'template-icon-black',
            () => { this.controller.toggleBlack() }
        );

        const secondarySet = document.createElement('div');
        secondarySet.className = 'controls secondarySet';
        secondary.appendChild(secondarySet);

        for (const i of [5, 10, 15, 30, 45, 60]) {
            addButton(
                secondarySet,
                `set${i}00`, `${i}:00`, null,
                () => { this.controller.setRemainingSeconds(i * 60) }
            );
        }

        const secondaryAdd = document.createElement('div');
        secondaryAdd.className = 'controls secondaryAdd';
        secondary.appendChild(secondaryAdd);

        addButton(
            secondaryAdd,
            'sub1000', '&ndash;10:00', null,
            () => { this.controller.addRemainingSeconds(-600) }
        );
        addButton(
            secondaryAdd,
            'sub500', '&ndash;5:00', null,
            () => { this.controller.addRemainingSeconds(-300) }
        );
        addButton(
            secondaryAdd,
            'sub100', '&ndash;1:00', null,
            () => { this.controller.addRemainingSeconds(-60) }
        );
        addButton(
            secondaryAdd,
            'add1000', '+10:00', null,
            () => { this.controller.addRemainingSeconds(600) }
        );
        addButton(
            secondaryAdd,
            'add500', '+5:00', null,
            () => { this.controller.addRemainingSeconds(300) }
        );
        addButton(
            secondaryAdd,
            'add100', '+1:00', null,
            () => { this.controller.addRemainingSeconds(60) }
        );
    }

    handleTimerEvent(event : TimerEvent) {
        if (event.type == TimerEventType.UpdateRemainingSeconds) {
            this.btnStart.toggleEnabled(!event.running);
            this.btnPause.toggleEnabled(event.running);
            this.btnBlack.toggleOnOff(event.black);
        }
    }
}


type ButtonHandler = (target : Button) => void;

class Button {
    element : HTMLAnchorElement;
    handler : ButtonHandler;

    constructor(name : string, text : string, icon? : string, handler? : ButtonHandler) {
        this.build(name, text, icon);
        this.handler = handler;
    }

    build(name : string, text : string, icon : string) {
        const a = document.createElement('a');
        a.className = `button button-${name}`;
        if (!icon) {
            a.innerHTML = text;
        } else {
            a.title = text;
            a.appendChild((document.getElementById(icon) as HTMLTemplateElement).content.cloneNode(true));
            a.classList.add('icon');
        }
        a.addEventListener('click', (event : MouseEvent) => {
            event.stopPropagation();
            if (this.handler) {
                this.handler(this);
            }
            return false;
        });
        this.element = a;
    }

    toggleEnabled(enabled : boolean) {
        this.element.classList.toggle('disabled', !enabled);
    }

    enable() {
        this.element.classList.remove('disabled');
    }

    disable() {
        this.element.classList.add('disabled');
    }

    toggleOnOff(on : boolean) {
        this.element.classList.toggle('on', on);
    }

    on() {
        this.element.classList.add('on');
    }

    off() {
        this.element.classList.remove('on');
    }
}
