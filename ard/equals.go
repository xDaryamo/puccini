package ard

func Equals(a interface{}, b interface{}) bool {
	switch a.(type) {
	case Map:
		if bMap, ok := b.(Map); ok {
			aMap := a.(Map)

			// Does B have keys that A doesn't have?
			for key, _ := range bMap {
				if _, ok := MapValue(aMap, key); !ok {
					return false
				}
			}

			// Are all values in A equal to those in B?
			for key, aValue := range aMap {
				if bValue, ok := MapValue(bMap, key); ok {
					if !Equals(aValue, bValue) {
						return false
					}
				} else {
					return false
				}
			}
		}

	case List:
		if bList, ok := b.(List); ok {
			aList := a.(List)

			// Must have same lengths
			if len(aList) != len(bList) {
				return false
			}

			for index, aValue := range aList {
				bValue := bList[index]
				if !Equals(aValue, bValue) {
					return false
				}
			}
		}

	default:
		return a == b
	}

	return true
}
