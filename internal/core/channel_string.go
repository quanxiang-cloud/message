// Code generated by "stringer -type Channel"; DO NOT EDIT.

package core

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Letter-0]
	_ = x[Email-1]
}

const _Channel_name = "LetterEmail"

var _Channel_index = [...]uint8{0, 6, 11}

func (i Channel) String() string {
	if i < 0 || i >= Channel(len(_Channel_index)-1) {
		return "Channel(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Channel_name[_Channel_index[i]:_Channel_index[i+1]]
}