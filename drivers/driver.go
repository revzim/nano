package drivers

type (
	Driver interface {
		Init() error
		Close()
	}
)
