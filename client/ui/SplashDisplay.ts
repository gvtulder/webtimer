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

import { Router } from "./Router.js";

export class SplashDisplay {
    input : HTMLInputElement;
    button : HTMLInputElement;

    constructor(router : Router) {
        const input = document.querySelector('input.key') as HTMLInputElement;
        const form = input.closest('form') as HTMLFormElement;
        const button = form.querySelector('input[type="submit"]') as HTMLInputElement;

        const h = () => { this.cleanString(); };
        input.addEventListener('keydown', h);
        input.addEventListener('keyup', h);
        input.addEventListener('change', h);

        form.addEventListener('submit', (ev : SubmitEvent) => {
            ev.stopPropagation();
            ev.preventDefault();
            if (this.cleanString()) {
                router.navigateTo(`${input.value}/`);
            }
            return false;
        });

        this.input = input;
        this.button = button;
    }

    activate(key : string) {
        if (key) {
            this.input.value = key;
            this.cleanString();
        }
        this.input.focus();
        this.input.select();
    }

    cleanString() : boolean {
        const v = this.input.value;
        const c = this.input.value.replace(' ', '-').toLowerCase().replace(/[^-a-z0-9]/, '').substring(0, 20);
        if (v != c) this.input.value = c;
        const ok = c.length >= 5;
        this.button.disabled = !ok;
        return ok;
    }
}
