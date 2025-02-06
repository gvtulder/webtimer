import { TimerDisplay } from "./ui/TimerDisplay.js";
import { TimerControllerWebsocket } from "./lib/TimerControllerWebsocket.js";
import { ControlsDisplay } from "./ui/ControlsDisplay.js";
import { SplashDisplay } from "./ui/SplashDisplay.js";
import { Router, WsUrlFunction } from "./ui/Router.js";

export function startApp(basePath : string, wsUrl : WsUrlFunction) {
    const controller = new TimerControllerWebsocket();

    const router = new Router(basePath, wsUrl);

    const container = document.getElementById('container');

    const timerDisplay = new TimerDisplay(controller, router);
    container.appendChild(timerDisplay.element);

    const controlsDisplay = new ControlsDisplay(controller);
    container.appendChild(controlsDisplay.element);

    const splashDisplay = new SplashDisplay(router);

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


    router.run(basePath, controller, timerDisplay, splashDisplay);

    globalThis.timerController = controller;
    globalThis.router = router;
}