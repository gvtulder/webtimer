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
        let c = this.input.value.replace(' ', '-').toLowerCase().replace(/[^-a-z0-9]/, '').substring(0, 20);
        if (v != c) this.input.value = c;
        const ok = c.length >= 5;
        this.button.disabled = !ok;
        return ok;
    }
}
