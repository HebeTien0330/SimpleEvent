/**
 * @Author: tangchengqin
 * @Date: 2024/9/18 11:47:38
 * @LastEditors: tangchengqin
 * @LastEditTime: 2024/9/18 16:27:20
 * Description:
 * Copyright: Copyright (©)}) 2024 Clarify. All rights reserved.
 */
package event

import (
	"sync"
)

var (
	instance *EventManager
	once  	 sync.Once
)

type EventList []Event

func GetEventManager() *EventManager {
	once.Do(func() {
		instance = &EventManager{
			eventMap: make(map[string]EventList),
			onceEventMap: make(map[string]EventList),
			counter: 1,
		}
	})
	return instance
}

type EventManager struct {
	eventMap map[string]EventList
	onceEventMap map[string]EventList
	counter int
}

// 默认过滤器
func DefaultFilter(...any) bool {
	return true
}

func (slf *EventManager) RegisterEvent(eventName string, callback func(...any), filter func(...any) bool, once bool) int {
	if once {
		return slf.listenOnce(eventName, callback, filter)
	}
	return slf.listen(eventName, callback, filter)
}

func (slf *EventManager) listen(eventName string, callback func(...any), filter func(...any) bool) int {
	eventList, exist := slf.eventMap[eventName]
	if !exist {
		eventList = EventList{}
	}
	evtId := slf.counter
	eventList = append(eventList, Event{evtId: evtId, callback: callback, eventName: eventName, filter: filter})
	slf.eventMap[eventName] = eventList
	slf.counter++
	return evtId
}

func (slf *EventManager) listenOnce(eventName string, callback func(...any), filter func(...any) bool) int {
	eventList, exist := slf.onceEventMap[eventName]
	if !exist {
		eventList = EventList{}
	}
	evtId := slf.counter
	eventList = append(eventList, Event{evtId: evtId, callback: callback, eventName: eventName, filter: filter})
	slf.onceEventMap[eventName] = eventList
	slf.counter++
	return evtId
}

func (slf *EventManager) RemoveEvent(eventName string, evtId interface{}) {
	if evtId == nil {
		slf.removeAll(eventName)
		return
	}
	slf.remove(eventName, evtId.(int))
	slf.removeOnce(eventName, evtId.(int))
}

func (slf *EventManager) remove(eventName string, evtId int) {
	eventList, exist := slf.eventMap[eventName]
	if !exist {
		return
	}
	for i, event := range eventList {
		if event.GetID() != evtId {
			continue
		}
		eventList = append(eventList[:i], eventList[i+1:]...)
		slf.eventMap[eventName] = eventList
		return
	}
}

func (slf *EventManager) removeOnce(eventName string, evtId int) {
	eventList, exist := slf.onceEventMap[eventName]
	if !exist {
		return
	}
	for i, event := range eventList {
		if event.GetID() != evtId {
			continue
		}
		eventList = append(eventList[:i], eventList[i+1:]...)
		slf.onceEventMap[eventName] = eventList
		return
	}
}

func (slf *EventManager) removeAll(eventName string) {
	delete(slf.eventMap, eventName)
	delete(slf.onceEventMap, eventName)
}

func (slf *EventManager) TriggerEvent(eventName string, evtId int, args ...any) {
	slf.execute(eventName, evtId, args...)
	slf.executeOnce(eventName, evtId, args...)
}

func (slf *EventManager) execute(eventName string, evtId int, args ...any) {
	eventList, exist := slf.eventMap[eventName]
	if !exist {
		return
	}
	if evtId == -1 {
		for _, event := range eventList {
			if event.GetID() != evtId {
				continue
			}
			event.Execute(args...)
			return
		}
	}
	for _, event := range eventList {
		event.Execute(args...)
	}
}

func (slf *EventManager) executeOnce(eventName string, evtId int, args ...any) {
	eventList, exist := slf.onceEventMap[eventName]
	if !exist {
		return
	}
	if evtId == -1 {
		for i, event := range eventList {
			if event.GetID() != evtId {
				continue
			}
			event.Execute(args...)
			eventList = append(eventList[:i], eventList[i+1:]...)
			slf.onceEventMap[eventName] = eventList
			return
		}
	}
	for _, event := range eventList {
		event.Execute(args...)
	}
	delete(slf.onceEventMap, eventName)
}


type EventHandler struct {
	EventManager *EventManager
}

func (slf *EventHandler) OnInit() error {
	slf.EventManager = GetEventManager()
	return nil
}

func (slf *EventHandler) Register(eventName string, callback func(...any), filter func(...any) bool, once bool) int {
	return slf.EventManager.RegisterEvent(eventName, callback, filter, once)
}

func (slf *EventHandler) Deregister(eventName string, evtId interface{}) {
	slf.EventManager.RemoveEvent(eventName, evtId)
}

func (slf *EventHandler) OnEvent(eventName string, args ...any) {
	slf.EventManager.TriggerEvent(eventName, -1, args...)
}

func (slf *EventHandler) OnSEvent(eventName string, evtId int, args ...any) {
	slf.EventManager.TriggerEvent(eventName, evtId, args...)
}
