package typhon4g

// Prop defines the interface of properties.
type Prop interface {
	// Str get the string value of key specified by name.
	Str(name string) string
	// Str get the string value of key specified by name or defaultValue when value is empty or missed.
	StrOr(name, defaultValue string) string

	// Bool get the bool value of key specified by name.
	Bool(name string) bool
	// BoolOr get the bool value of key specified by name or defaultValue when value is empty or missed.
	BoolOr(name string, defaultValue bool) bool

	// Int get the int value of key specified by name.
	Int(name string) int
	// IntOr get the int value of key specified by name or defaultValue when value is empty or missed.
	IntOr(name string, defaultValue int) int

	// Int32 get the int32 value of key specified by name.
	Int32(name string) int32
	// Int32Or get the int32 value of key specified by name or defaultValue when value is empty or missed.
	Int32Or(name string, defaultValue int32) int32

	// Int64 get the int64 value of key specified by name.
	Int64(name string) int64
	// Int64Or get the int64 value of key specified by name or defaultValue when value is empty or missed.
	Int64Or(name string, defaultValue int64) int64

	// Float32 get the float32 value of key specified by name.
	Float32(name string) float32
	// Float32Or get the float32 value of key specified by name or defaultValue when value is empty or missed.
	Float32Or(name string, defaultValue float32) float32

	// Float64 get the float64 value of key specified by name.
	Float64(name string) float64
	// Float64Or get the float64 value of key specified by name or defaultValue when value is empty or missed.
	Float64Or(name string, defaultValue float64) float64
}
