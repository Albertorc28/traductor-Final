package model

//Insercion struct
type Insercion struct {
	Palabra    string
	Traduccion string
	Idioma     string
}

//Filtro struct
type Filtro struct {
	ID int
}

//Traduccion struct
type Traduccion struct {
	Palabra string
	Idioma  string
}
