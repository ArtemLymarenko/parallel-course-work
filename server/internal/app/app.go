package app

type Env string

const (
	EnvProd Env = "production"
	EnvDev  Env = "development"
)

func (e Env) IsProduction() bool {
	return e == EnvProd
}

func (e Env) String() string {
	return string(e)
}

func (e Env) Valid() bool {
	switch e {
	case EnvProd, EnvDev:
		return true
	default:
		return false
	}
}
