import { TimerDisplay } from "./ui/TimerDisplay.js";
import { TimerControllerWebsocket } from "./lib/TimerControllerWebsocket.js";
import { ControlsDisplay } from "./ui/ControlsDisplay.js";
import { SplashDisplay } from "./ui/SplashDisplay.js";

var controller : TimerControllerWebsocket;

export function startApp(wsUrl : string, backUrl : string) {
    const controller = new TimerControllerWebsocket(wsUrl);

    const container = document.getElementById('container');

    const timerDisplay = new TimerDisplay(controller, backUrl);
    container.appendChild(timerDisplay.element);

    const controlsDisplay = new ControlsDisplay(controller);
    container.appendChild(controlsDisplay.element);

    let wakeLock = null;

    async function enableWakeLock() {
        if (!wakeLock && 'wakeLock' in navigator) {
            try {
                wakeLock = await navigator.wakeLock.request('screen');
            } catch (err) {
            }
        }
    }

    function releaseWakeLock() {
        if (wakeLock) {
            wakeLock.release().then(() => { wakeLock = null; });
        }
    }

    timerDisplay.element.addEventListener('click', () => {
        container.classList.toggle('fullscreen');
        if (container.classList.contains('fullscreen')) {
            if (!document.fullscreenElement && document.documentElement.requestFullscreen) {
                document.documentElement.requestFullscreen();
            }
            enableWakeLock();
        } else {
            if (document.exitFullscreen) {
                document.exitFullscreen();
            }
            releaseWakeLock();
        }
    });

    controller.run();

    globalThis.timerController = controller;
}

export function startSplash(key : string) {
    const container = document.getElementById('container');
    const splashDisplay = new SplashDisplay(key);
    container.appendChild(splashDisplay.element);
    splashDisplay.focus();
}
