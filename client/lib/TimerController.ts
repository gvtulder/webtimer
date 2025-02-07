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

export enum TimerEventType {
    Connecting = "Connecting",
    UpdateRemainingSeconds = "UpdateRemainingSeconds",
    StartTimer = "StartTimer",
    PauseTimer = "PauseTimer",
}

/**
 * An event following a timer state update.
 */
export type TimerEvent = {
    type : TimerEventType,
    connected? : boolean,
    active? : boolean,
    black? : boolean,
    countdown? : boolean,
    running? : boolean,
    remainingSeconds? : number,
    clients? : number,
 };

export type TimerEventListener = (event : TimerEvent) => void;

/**
 * A TimerController handles the connection with a timer, processing commands and sending TimerEvent events to subscribers.
 */
export interface TimerController {
    /**
     * Connects to the timer server (e.g., through a WebSocket).
     * @param url the URL of the timer server
     */
    connect(url : string) : void;

    /**
     * Disconnects from the timer server.
     */
    disconnect() : void;

    /**
     * Sets the remaining seconds on the timer.
     * @param seconds the new time in seconds
     */
    setRemainingSeconds(seconds : number) : void;

    /**
     * Adds (or substracts) seconds from the timer.
     * @param seconds the number of seconds to add/subtract
     */
    addRemainingSeconds(seconds : number) : void;

    /**
     * Starts the timer.
     */
    startTimer() : void;

    /**
     * Stops the timer and sets the time to zero.
     */
    resetTimer() : void;

    /**
     * Pauses the timer.
     */
    pauseTimer() : void;

    /**
     * Toggles black-screen mode.
     */
    toggleBlack() : void;

    /**
     * Adds a listener to receive TimerEvent updates.
     * @param listener an event handler
     */
    addListener(listener : TimerEventListener) : void;
}
