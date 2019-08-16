package ard

//
// Locator
//

type Locator interface {
	Locate(path ...PathElement) (int, int, bool)
}
