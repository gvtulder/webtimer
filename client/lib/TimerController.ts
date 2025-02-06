
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
    setRemainingSeconds(seconds : number) : void;
    addRemainingSeconds(seconds : number) : void;
    startTimer() : void;
    resetTimer() : void;
    pauseTimer() : void;
    toggleBlack() : void;
    addListener(listener : TimerEventListener) : void;
}
