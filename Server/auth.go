package server

import (
	"../fastJSON"
	"./common"
)

type personType int

const (
	isVisitorType personType = 0
	isStaffType   personType = 1
	unknownType   personType = 2
)

func getAuthType(msg *fastjson.Value) (res personType) {
	ok := isVisitor(msg)
	if ok {
		return isVisitorType
	}

	ok = isStaff(msg)
	if ok {
		return isStaffType
	}

	return unknownType
}

func isStaff(jsonMsg *fastjson.Value) bool {
	ok := jsonMsg.Exists("command")
	if !ok {
		return false
	}

	if jsonMsg.GetString("command") == common.CommandStaffAuth {
		return true
	}

	return false
}

func isVisitor(jsonMsg *fastjson.Value) bool {
	ok := jsonMsg.Exists("type")
	if !ok {
		return false
	}

	if jsonMsg.GetString("type") == common.CommandVisitorAuth {
		return true
	}

	return false
}