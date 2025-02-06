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

import { TimerDisplay } from "./ui/TimerDisplay.js";
import { TimerControllerWebsocket } from "./lib/TimerControllerWebsocket.js";
import { ControlsDisplay } from "./ui/ControlsDisplay.js";
import { SplashDisplay } from "./ui/SplashDisplay.js";
import { Router, WsUrlFunction } from "./ui/Router.js";

declare var VERSION : string;

export function startApp(basePath : string, wsUrl : WsUrlFunction) {
    console.log(`Running webtimer.cc client version ${VERSION}`);

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
}
