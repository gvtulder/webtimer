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

export interface TimerController {
    connect(url : string) : void;
    disconnect() : void;
    setRemainingSeconds(seconds : number) : void;
    addRemainingSeconds(seconds : number) : void;
    startTimer() : void;
    resetTimer() : void;
    pauseTimer() : void;
    toggleBlack() : void;
    addListener(listener : TimerEventListener) : void;
}
