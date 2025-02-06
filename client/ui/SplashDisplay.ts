
export class SplashDisplay {
    input : HTMLInputElement;

    constructor(key : string) {
        const input = document.querySelector('input.key') as HTMLInputElement;
        const form = input.closest('form') as HTMLFormElement;
        const button = form.querySelector('input[type="submit"]') as HTMLInputElement;

        if (key) {
            input.value = key;
        }

        function cleanString() {
            const v = input.value;
            let c = input.value.replace(' ', '-').toLowerCase().replace(/[^-a-z0-9]/, '').substring(0, 20);
            if (v != c) input.value = c;
            button.disabled = c.length < 5;
        }
        cleanString();

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
        this.input.select();
    }
}
