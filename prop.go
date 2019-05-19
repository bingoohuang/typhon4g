package typhon4g

type Prop interface {
	Str(name string) string
	StrDefault(name, defaultValue string) string

	Bool(name string) bool
	BoolDefault(name string, defaultValue bool) bool

	Int(name string) int
	IntDefault(name string, defaultValue int) int

	Int32(name string) int32
	Int32Default(name string, defaultValue int32) int32

	Int64(name string) int64
	Int64Default(name string, defaultValue int64) int64

	Float32(name string) float32
	Float32Default(name string, defaultValue float32) float32

	Float64(name string) float64
	Float64Default(name string, defaultValue float64) float64
}
