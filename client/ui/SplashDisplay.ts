
export class SplashDisplay {
    element : HTMLDivElement;
    input : HTMLInputElement;

    constructor() {
        const div = document.createElement('div');
        div.className = 'splash-display';
        this.element = div;

        const form = document.createElement('form');
        const p1 = document.createElement('p');
        const p2 = document.createElement('p');
        const input = document.createElement('input');
        const button = document.createElement('input');

        p1.innerHTML = 'Enter the ID to join a running timer,<br/>or choose a new ID to start one.';

        form.action = '';
        input.placeholder = 'Session ID';
        input.type = 'text';
        button.value = 'Go';
        button.type = 'submit';

        form.appendChild(p1);
        p2.appendChild(input);
        p2.appendChild(button);
        form.appendChild(p2);
        div.appendChild(form);

        function cleanString() {
            const v = input.value;
            let c = input.value.replace(' ', '-').toLowerCase().replace(/[^-a-z0-9]/, '').substring(0, 20);
            if (v != c) input.value = c;
            button.disabled = c.length < 5;
        }

        input.addEventListener('keydown', cleanString);
        input.addEventListener('keyup', cleanString);
        input.addEventListener('change', cleanString);

        form.addEventListener('submit', function(ev : SubmitEvent) {
            ev.stopPropagation();
            ev.preventDefault();
            cleanString();
            if (input.value.length >= 5) {
                window.location.href = input.value;
            }
            return false;
        });

        this.input = input;
    }

    focus() {
        this.input.focus();
    }
}
