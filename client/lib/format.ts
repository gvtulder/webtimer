
export function formatTime(seconds : number) {
    const negative = seconds < 0;
    seconds = Math.round(Math.abs(seconds));
    const components = [];
    while (seconds > 0) {
        components.push(seconds % 60);
        seconds = Math.floor(seconds / 60);
    }
    while (components.length < 2) {
        components.push(0);
    }
    components.reverse();
    for (let i = 1; i<components.length; i++) {
        components[i] = components[i] < 10 ? `0${components[i]}` : `${components[i]}`;
    }
    return components.join(':');
    return (negative ? '-' : '') + components.join(':');
}