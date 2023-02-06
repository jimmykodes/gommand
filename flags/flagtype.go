package flags

//go:generate stringer -type=FlagType
type FlagType int

const (
	UnknownFlagType FlagType = iota
	StringFlagType
	BoolFlagType
	DurationFlagType
	IntFlagType
	Int8FlagType
	Int16FlagType
	Int32FlagType
	Int64FlagType
	UintFlagType
	Uint8FlagType
	Uint16FlagType
	Uint32FlagType
	Uint64FlagType
	Float32FlagType
	Float64FlagType
	StringSliceFlagType
	BoolSliceFlagType
	DurationSliceFlagType
	IntSliceFlagType
	Int8SliceFlagType
	Int16SliceFlagType
	Int32SliceFlagType
	Int64SliceFlagType
	UintSliceFlagType
	Uint8SliceFlagType
	Uint16SliceFlagType
	Uint32SliceFlagType
	Uint64SliceFlagType
	Float32SliceFlagType
	Float64SliceFlagType
)
