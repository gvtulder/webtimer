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
