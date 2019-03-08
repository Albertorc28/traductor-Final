package model

import "time"

//RIdioma struct
type RIdioma struct {
	ID     int
	Idioma string
}

//RTraduccion struct
type RTraduccion struct {
	ID            int
	Palabra       string
	Traduccion    string
	Fecha         time.Time
	FechaConsulta time.Time
	Idioma        string
	Ok            bool
	Fail          bool
}
