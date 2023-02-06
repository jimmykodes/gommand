// Code generated by "stringer -type=FlagType"; DO NOT EDIT.

package flags

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnknownFlagType-0]
	_ = x[StringFlagType-1]
	_ = x[BoolFlagType-2]
	_ = x[DurationFlagType-3]
	_ = x[IntFlagType-4]
	_ = x[Int8FlagType-5]
	_ = x[Int16FlagType-6]
	_ = x[Int32FlagType-7]
	_ = x[Int64FlagType-8]
	_ = x[UintFlagType-9]
	_ = x[Uint8FlagType-10]
	_ = x[Uint16FlagType-11]
	_ = x[Uint32FlagType-12]
	_ = x[Uint64FlagType-13]
	_ = x[Float32FlagType-14]
	_ = x[Float64FlagType-15]
	_ = x[StringSliceFlagType-16]
	_ = x[BoolSliceFlagType-17]
	_ = x[DurationSliceFlagType-18]
	_ = x[IntSliceFlagType-19]
	_ = x[Int8SliceFlagType-20]
	_ = x[Int16SliceFlagType-21]
	_ = x[Int32SliceFlagType-22]
	_ = x[Int64SliceFlagType-23]
	_ = x[UintSliceFlagType-24]
	_ = x[Uint8SliceFlagType-25]
	_ = x[Uint16SliceFlagType-26]
	_ = x[Uint32SliceFlagType-27]
	_ = x[Uint64SliceFlagType-28]
	_ = x[Float32SliceFlagType-29]
	_ = x[Float64SliceFlagType-30]
}

const _FlagType_name = "UnknownFlagTypeStringFlagTypeBoolFlagTypeDurationFlagTypeIntFlagTypeInt8FlagTypeInt16FlagTypeInt32FlagTypeInt64FlagTypeUintFlagTypeUint8FlagTypeUint16FlagTypeUint32FlagTypeUint64FlagTypeFloat32FlagTypeFloat64FlagTypeStringSliceFlagTypeBoolSliceFlagTypeDurationSliceFlagTypeIntSliceFlagTypeInt8SliceFlagTypeInt16SliceFlagTypeInt32SliceFlagTypeInt64SliceFlagTypeUintSliceFlagTypeUint8SliceFlagTypeUint16SliceFlagTypeUint32SliceFlagTypeUint64SliceFlagTypeFloat32SliceFlagTypeFloat64SliceFlagType"

var _FlagType_index = [...]uint16{0, 15, 29, 41, 57, 68, 80, 93, 106, 119, 131, 144, 158, 172, 186, 201, 216, 235, 252, 273, 289, 306, 324, 342, 360, 377, 395, 414, 433, 452, 472, 492}

func (i FlagType) String() string {
	if i < 0 || i >= FlagType(len(_FlagType_index)-1) {
		return "FlagType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FlagType_name[_FlagType_index[i]:_FlagType_index[i+1]]
}
