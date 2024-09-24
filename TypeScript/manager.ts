/**
 * @Author: tangchengqin
 * @Date: 2024/9/14 10:50:35
 * @LastEditors: tangchengqin
 * @LastEditTime: 2024/9/18 16:08:37
 * Description: 
 * Copyright: Copyright (©)}) 2024 Clarify. All rights reserved.
 */

import { Event, Params } from "./event";

class EventManager {

    private static _instance: EventManager;
    private _eventMap: Map<string, Event[]>;
    private _onceEventMap: Map<string, Event[]>;
    private counter: number = 1;

    constructor() {
        this._eventMap = new Map<string, Event[]>();
        this._onceEventMap = new Map<string, Event[]>();
    }

    static getInstance(): EventManager {
        if (!this._instance) {
            this._instance = new EventManager();
        }
        return this._instance;
    }

    private listen(eventName: string, callback: Function, filter?: Function): number {
        let callbacks = this._eventMap.get(eventName);
        if (!callbacks) {
            callbacks = [];
            this._eventMap.set(eventName, callbacks);
        }
        const evtId = this.counter;
        const event = new Event(evtId, eventName, callback, filter);
        callbacks.push(event);
        this.counter ++;
        return evtId;
    }

    private listenOnce(eventName: string, callback: Function, filter?: Function): number {
        let callbacks = this._onceEventMap.get(eventName);
        if (!callbacks) {
            callbacks = [];
            this._onceEventMap.set(eventName, callbacks);
        }
        const evtId = this.counter;
        const event = new Event(evtId, eventName, callback, filter);
        callbacks.push(event);
        this.counter ++;
        return evtId;
    }

    public on(eventName: string, callback: Function, filter?: Function, once?: boolean): number {
        if (once) {
            return this.listenOnce(eventName, callback, filter);
        }
        return this.listen(eventName, callback, filter);
    }

    private cancel(eventName: string, evtId: number): void {
        const events = this._eventMap.get(eventName);
        if (!events) {
            return;
        }
        for (let i = 0; i < events.length; i++) {
            const event = events[i];
            if (event.evtId !== evtId) {
                continue;
            }
            events.splice(i, 1);
            break;
        }
    }

    private cancelOnce(eventName: string, evtId: number): void {
        const events = this._onceEventMap.get(eventName);
        if (!events) {
            return;
        }
        for (let i = 0; i < events.length; i++) {
            const event = events[i];
            if (event.evtId !== evtId) {
                continue;
            }
            events.splice(i, 1);
            break;
        }
    }

    private cancelAll(eventName: string): void {
        this._eventMap.delete(eventName);
        this._onceEventMap.delete(eventName);
    }

    public off(eventName: string, evtId?: number, once?: boolean) {
        if (!evtId) {
            this.cancelAll(eventName);
            return;
        }
        if (once) {
            this.cancelOnce(eventName, evtId);
            return;
        }
        this.cancel(eventName, evtId);
    }

    private execute(eventName: string, args: Params, evtId?: number): void {
        const events = this._eventMap.get(eventName);
        if (!events) {
            return;
        }
        if (evtId) {
            for (const event of events) {
                if (event.evtId === evtId) {
                    event.execute(args);
                    return;
                }
            }
        } else {
            for (const event of events) {
                event.execute(args);
            }
        }
    }

    private executeOnce(eventName: string, args: Params, evtId?: number): void {
        const events = this._onceEventMap.get(eventName);
        if (!events) {
            return;
        }
        if (evtId) {
            for (const event of events) {
                if (event.evtId === evtId) {
                    event.execute(args);
                    events.splice(events.indexOf(event), 1);
                    return;
                }
            }
        } else {
            for (const event of events) {
                event.execute(args);
            }
            this._onceEventMap.delete(eventName);
        }
    }

    public call(eventName: string, args: Params, evtId?: number): void {
        this.execute(eventName, args, evtId);
        this.executeOnce(eventName, args, evtId);
    }

}

const handler = {

    listen(eventName: string, callback: Function, filter?: Function, once?: boolean): number {
        return EventManager.getInstance().on(eventName, callback, filter, once);
    },

    listenMulti(eventList: Array<any>): Array<number> {
        const ret: Array<number> = [];
        for (const event of eventList) {
            const evtId = EventManager.getInstance().on(event.eventName, event.callback, event.filter, event.once);
            ret.push(evtId);
        }
        return ret;
    },

    cancel(eventName: string, evtId?: number, once?: boolean) {
        EventManager.getInstance().off(eventName, evtId, once);
    },

    onEvent(eventName: string, args: any, evtId?: number) {
        EventManager.getInstance().call(eventName, args, evtId);
    }

}

export const installEventSystem = (target: any) => {
    for (const key in handler) {
        const func = handler[key];
        target[key] = func.bind(target);
    }
}
