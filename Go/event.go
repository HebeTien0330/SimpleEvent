/**
 * @Author: tangchengqin
 * @Date: 2024/9/18 11:47:49
 * @LastEditors: tangchengqin
 * @LastEditTime: 2024/9/18 16:27:23
 * Description: 
 * Copyright: Copyright (©)}) 2024 Clarify. All rights reserved.
 */
package event

type Event struct {
	evtId int
	eventName string
	callback func(...any)
	filter func(...any) bool
}

func (slf *Event) GetID() int {
	return slf.evtId
}

func (slf *Event) GetName() string {
	return slf.eventName
}

func (slf *Event) Execute(args ...any) {
	if slf.filter != nil && !slf.filter(args) {
		return
	}
	slf.callback(args)
}
