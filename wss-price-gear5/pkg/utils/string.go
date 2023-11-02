package utils

import "strings"

func OneOfStrings(mainString string, stringList ...string) bool {
	for _, str := range stringList {
		if strings.Compare(mainString, str) == 0 {
			return true
		}
	}

	return false
}

func RemoveUUIDStrikeThrough(uuidStr string) string {
	uuidString := strings.Replace(uuidStr, "-", "", -1)
	return uuidString
}

func AddUUIDStrikeThrough(str string) string {
	strUUID := str[:8] + "-" + str[8:12] + "-" + str[12:16] + "-" + str[16:20] + "-" + str[20:]
	return strUUID
}
