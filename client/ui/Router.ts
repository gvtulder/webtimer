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

import { TimerController } from "../lib/TimerController.js";
import { SplashDisplay } from "./SplashDisplay.js";
import { TimerDisplay } from "./TimerDisplay.js";

export type WsUrlFunction = (key : string) => string;

export class Router {
    private basePath : string;
    private basePathRegexp : RegExp;
    private wsUrl : WsUrlFunction;

    private timerController : TimerController;
    private timerDisplay : TimerDisplay;
    private splashDisplay : SplashDisplay;

    constructor(basePath : string, wsUrl : WsUrlFunction) {
        if (!basePath) basePath = '/';
        if (!basePath.match(/^\//)) basePath = `/${basePath}`;
        if (!basePath.match(/\/$/)) basePath = `${basePath}/`;
        this.basePath = basePath;
        this.basePathRegexp = new RegExp(`^${basePath}`);
        this.wsUrl = wsUrl;
    }

    run(basePath : string, timerController : TimerController, timerDisplay : TimerDisplay, splashDisplay : SplashDisplay) {
        this.timerController = timerController;
        this.timerDisplay = timerDisplay;
        this.splashDisplay = splashDisplay;

        window.addEventListener('popstate', () => {
            this.handle();
        });
        this.handle();
    }

    handle() {
        let path = window.location.pathname;
        if (!path.match(this.basePathRegexp)) {
            this.navigateTo(this.basePath);
        }
        path = path.replace(this.basePathRegexp, '');

        if (path == '') {
            // root
            let key = window.location.hash;
            if (key && key != '#') {
                key = key.replace('#', '');
            }
            document.body.classList.add('show-splash');
            document.body.classList.remove('show-timer');
            this.timerController.disconnect();
            this.splashDisplay.activate(key);

        } else if (path.match('^[-a-zA-Z0-9]{5,}/$')) {
            // subdir
            const key = path.replace('/', '');
            document.body.classList.add('show-timer');
            document.body.classList.remove('show-splash');
            this.timerController.connect(this.wsUrl(key));
            this.timerDisplay.activate();

        } else if (path.match('^[-a-zA-Z0-9]{5,}$')) {
            // add trailing slash
            this.navigateTo(`${this.basePath}${path}/`);

        } else {
            this.navigateTo(this.basePath);
        }
    }

    navigateTo(url : string) {
        history.pushState({}, '', url);
        this.handle();
    }

    navigateToSplash() {
        this.navigateTo(this.basePath);
    }
}
