import { TimerDisplay } from "./ui/TimerDisplay";
import { TimerControllerWebsocket } from "./lib/TimerControllerWebsocket";
import { ControlsDisplay } from "./ui/ControlsDIsplay";

export function startApp(wsUrl : string) {
    const controller = new TimerControllerWebsocket(wsUrl);

    const container = document.getElementById('container');

    const timerDisplay = new TimerDisplay(controller);
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

    window.timerController = controller;
}
