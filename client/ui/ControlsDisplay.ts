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

        const secondarySet = new ButtonList('secondarySet');
        secondary.appendChild(secondarySet.element);

        for (const i of [60, 45, 30, 15, 10, 5]) {
            secondarySet.add(new Button(
                `set${i}00`, `${i}:00`, null,
                () => { this.controller.setRemainingSeconds(i * 60) }
            ), i == 15);
        }

        const secondarySub = new ButtonList('secondarySub');
        secondary.appendChild(secondarySub.element);

        secondarySub.add(new Button(
            'sub1000', '&minus;10:00', null,
            () => { this.controller.addRemainingSeconds(-600) }
        ));
        secondarySub.add(new Button(
            'sub500', '&minus;5:00', null,
            () => { this.controller.addRemainingSeconds(-300) }
        ));
        secondarySub.add(new Button(
            'sub100', '&minus;1:00', null,
            () => { this.controller.addRemainingSeconds(-60) }
        ), true);
        secondarySub.add(new Button(
            'sub30', '&minus;0:30', null,
            () => { this.controller.addRemainingSeconds(-30) }
        ));

        const secondaryAdd = new ButtonList('secondaryAdd');
        secondary.appendChild(secondaryAdd.element);

        secondaryAdd.add(new Button(
            'add1000', '+10:00', null,
            () => { this.controller.addRemainingSeconds(600) }
        ));
        secondaryAdd.add(new Button(
            'add500', '+5:00', null,
            () => { this.controller.addRemainingSeconds(300) }
        ));
        secondaryAdd.add(new Button(
            'add100', '+1:00', null,
            () => { this.controller.addRemainingSeconds(60) }
        ), true);
        secondaryAdd.add(new Button(
            'add30', '+0:30', null,
            () => { this.controller.addRemainingSeconds(30) }
        ));

        this.addDropdownHandlers();
    }

    handleTimerEvent(event : TimerEvent) {
        if (event.type == TimerEventType.UpdateRemainingSeconds) {
            this.btnStart.toggleEnabled(!event.running);
            this.btnPause.toggleEnabled(event.running);
            this.btnBlack.toggleOnOff(event.black);
        }
    }

    addDropdownHandlers() {
        document.addEventListener('click', (evt) => {
            const tgtButton = evt.target && (evt.target as HTMLElement).closest('button.up');
            const tgtList = tgtButton && tgtButton.closest('.button-list');
            for (const list of document.getElementsByClassName('button-list')) {
            list.classList.toggle('expanded', (list === tgtList) ? undefined : false);
            }
        });
    }
}


type ButtonHandler = (target : Button) => void;

class ButtonList {
    element : HTMLDivElement;
    btnCurrent : HTMLButtonElement;
    btnUp : HTMLButtonElement;
    list : HTMLUListElement;
    buttons : Button[];
    lastClicked : Button;
    lastClickedHandler : ButtonHandler;

    constructor(name : string) {
        this.build(name)
        this.buttons = [];
    }

    build(name : string) {
        const div = document.createElement('div');
        div.className = `button-list ${name}`;

        const combo = document.createElement('div');
        combo.className = 'button-combo';
        div.appendChild(combo);

        const btnCurrent = document.createElement('button');
        btnCurrent.className = 'current';
        combo.appendChild(btnCurrent);
        btnCurrent.addEventListener('click', () => {
            this.lastClickedHandler(this.lastClicked);
        });
        this.btnCurrent = btnCurrent;

        const btnUp = document.createElement('button');
        btnUp.className = 'up';
        btnUp.appendChild((document.getElementById('template-icon-up') as HTMLTemplateElement).content.cloneNode(true));
        combo.appendChild(btnUp);
        this.btnUp = btnUp;

        const ul = document.createElement('ul');
        div.appendChild(ul);
        this.list = ul;

        this.element = div;
    }

    add(button : Button, setLastClicked? : boolean) {
        const li = document.createElement('li');
        li.appendChild(button.element);
        this.list.appendChild(li);
        this.buttons.push(button);

        const oldHandler = button.handler;
        button.handler = (tgt : Button) => {
            this.updateLastClicked(button, oldHandler);
            oldHandler(tgt);
        };

        if (setLastClicked) {
            this.updateLastClicked(button, oldHandler);
        }

        return button;
    }

    updateLastClicked(button : Button, handler : ButtonHandler) {
        this.lastClicked = button;
        this.lastClickedHandler = handler;
        this.btnCurrent.innerHTML = button.element.innerHTML;
    }
}

class Button {
    element : HTMLButtonElement;
    handler : ButtonHandler;

    constructor(name : string, text : string, icon? : string, handler? : ButtonHandler) {
        this.build(name, text, icon);
        this.handler = handler;
    }

    build(name : string, text : string, icon : string) {
        const button = document.createElement('button');
        button.className = `button-${name}`;
        if (!icon) {
            button.innerHTML = text;
        } else {
            button.title = text;
            button.appendChild((document.getElementById(icon) as HTMLTemplateElement).content.cloneNode(true));
            button.classList.add('icon');
        }
        button.addEventListener('click', (event : MouseEvent) => {
            if (this.handler) {
                this.handler(this);
            }
        });
        this.element = button;
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
