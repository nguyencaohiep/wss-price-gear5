package utils

func OneOfInts(mainInt int, intList ...int) bool {
	for _, i := range intList {
		if mainInt == i {
			return true
		}
	}

	return false
}
