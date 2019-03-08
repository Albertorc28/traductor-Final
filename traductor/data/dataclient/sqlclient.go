package dataclient

import (
	"database/sql"
	"fmt"
	"traductor/data/model"

	_ "github.com/go-sql-driver/mysql" ///El driver se registra en database/sql en su función Init(). Es usado internamente por éste
)

//ListarIdiomas test
func ListarIdiomas() []model.RIdioma {
	db, err := sql.Open("mysql", "ubuntu:ubuntu@tcp(localhost:3306)/Traductor?parseTime=true")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//Consultamos todos los idiomas de la base de datos
	comando := "SELECT * FROM Idioma"
	fmt.Println(comando)
	query, err := db.Query("SELECT * FROM Idioma")

	if err != nil {
		panic(err.Error())
	}
	defer query.Close()

	resultado := make([]model.RIdioma, 0)
	for query.Next() {
		var fila = model.RIdioma{}

		err = query.Scan(&fila.ID, &fila.Idioma)
		if err != nil {
			panic(err.Error())
		}
		resultado = append(resultado, fila)
	}
	return resultado
}

//ListarTraducciones test
func ListarTraducciones(objeto *model.Traduccion) []model.RTraduccion {
	db, err := sql.Open("mysql", "ubuntu:ubuntu@tcp(localhost:3306)/Traductor?parseTime=true")

	if err != nil {
		panic(err.Error())
	}

	//Comprobamos que el Idioma existe
	comando := "SELECT ID FROM Idioma WHERE (Idioma = '" + objeto.Idioma + "')"
	fmt.Println(comando)
	query, err := db.Query("SELECT ID FROM Idioma WHERE (ID = ?)", objeto.Idioma)

	if err != nil {
		panic(err.Error())
	}
	defer query.Close()
	var aux = model.Filtro{}

	for query.Next() {
		err = query.Scan(&aux.ID)
		fmt.Println("ID Idioma: ", aux.ID)
	}
	resultado := make([]model.RTraduccion, 0)

	//Si no existe salimos ya de la funcion
	if aux.ID == 0 {
		fmt.Println("Idioma no encontrado")
		var aux = model.RTraduccion{}
		aux.Fail = true
		resultado = append(resultado, aux)
		return resultado
	}

	//Si el idioma existe, vamos a comprobar que la palabra exista
	comando = "SELECT ID FROM Palabra WHERE (Palabra = '" + objeto.Palabra + "')"
	fmt.Println(comando)
	query, err = db.Query("SELECT ID FROM Palabra WHERE (Palabra = ?)", objeto.Palabra)

	if err != nil {
		panic(err.Error())
	}
	defer query.Close()
	var fila = model.Filtro{}

	for query.Next() {
		err = query.Scan(&fila.ID)
		fmt.Println("ID Palabra: ", fila.ID)
	}

	//Si no existe salimos ya de la funcion
	if fila.ID == 0 {
		fmt.Println("Palabra no encontrada")
		var fila = model.RTraduccion{}
		fila.Fail = true
		resultado = append(resultado, fila)

		//Guardamos ya la plabra en la base de datos
		comando := "INSERT INTO Palabra(Palabra) VALUES ('" + objeto.Palabra + "')"
		fmt.Println(comando)
		insert, err := db.Query("INSERT INTO Palabra(Palabra) VALUES (?)", objeto.Palabra)
		if err != nil {
			panic(err.Error())
		}

		insert.Close()

		return resultado
	}

	//Si el Idioma y la Palabra existen vamos a comprobar si existe una traduccion
	comando = "SELECT t.ID, p.Palabra, i.Idioma, Traduccion, Fecha, FechaConsulta FROM Palabra p INNER JOIN Traduccion t ON p.ID=t.PalabraID INNER JOIN Idioma i ON t.IdiomaID=i.ID WHERE (p.Palabra = '" + objeto.Palabra + "' AND IdiomaID = '" + objeto.Idioma + "')"
	fmt.Println(comando)
	query, err = db.Query("SELECT t.ID, p.Palabra, Traduccion, Fecha, FechaConsulta, i.Idioma FROM Palabra p INNER JOIN Traduccion t ON p.ID=t.PalabraID INNER JOIN Idioma i ON t.IdiomaID=i.ID WHERE (p.Palabra = ? AND IdiomaID = ?)", objeto.Palabra, objeto.Idioma)

	if err != nil {
		panic(err.Error())
	}
	defer query.Close()

	for query.Next() {
		var fila = model.RTraduccion{}

		err = query.Scan(&fila.ID, &fila.Palabra, &fila.Traduccion, &fila.Fecha, &fila.FechaConsulta, &fila.Idioma)
		if err != nil {
			panic(err.Error())
		}
		fila.Ok = true
		resultado = append(resultado, fila)

		//Registramos la fecha de la ultima consulta
		comando := "UPDATE Traduccion SET FechaConsulta = utc_timestamp() WHERE (ID = '" + string(fila.ID) + "')"
		fmt.Println(comando)
		update, err := db.Query("UPDATE Traduccion SET FechaConsulta = utc_timestamp() WHERE (ID = ?)", fila.ID)

		if err != nil {
			panic(err.Error())
		}
		defer update.Close()

	}

	//Si no se encontró esa traduccion, vamos a comprobar si existe en otro idioma
	if len(resultado) == 0 {
		comando := "SELECT t.ID, p.Palabra, i.Idioma, Traduccion, Fecha, FechaConsulta FROM Palabra p INNER JOIN Traduccion t ON p.ID=t.PalabraID INNER JOIN Idioma i ON t.IdiomaID=i.ID WHERE (p.Palabra = '" + objeto.Palabra + "' )"
		fmt.Println(comando)
		query, err := db.Query("SELECT t.ID, p.Palabra, Traduccion, Fecha, FechaConsulta, i.Idioma FROM Palabra p INNER JOIN Traduccion t ON p.ID=t.PalabraID INNER JOIN Idioma i ON t.IdiomaID=i.ID WHERE (p.Palabra = ?)", objeto.Palabra)

		if err != nil {
			panic(err.Error())
		}
		defer query.Close()

		for query.Next() {
			var fila = model.RTraduccion{}

			err = query.Scan(&fila.ID, &fila.Palabra, &fila.Traduccion, &fila.Fecha, &fila.FechaConsulta, &fila.Idioma)
			if err != nil {
				panic(err.Error())
			}
			fila.Ok = false
			resultado = append(resultado, fila)

		}
	}

	//Si no existe ninguna traduccion
	if len(resultado) == 0 {
		var fila = model.RTraduccion{}
		fila.Fail = true
		resultado = append(resultado, fila)
	}

	return resultado
}

//InsertTraducciones test
func InsertTraducciones(objeto *model.Insercion) {
	db, err := sql.Open("mysql", "ubuntu:ubuntu@tcp(localhost:3306)/Traductor")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//Consultamos el Idioma en la Base de Datos, si existe sacamos su ID
	comando := "SELECT ID FROM Idioma WHERE (Idioma = '" + objeto.Idioma + "')"
	fmt.Println(comando)
	query, err := db.Query("SELECT ID FROM Idioma WHERE (Idioma = ?)", objeto.Idioma)

	if err != nil {
		panic(err.Error())
	}
	defer query.Close()
	var filaidioma = model.RIdioma{}
	resultado := make([]model.RIdioma, 0)
	for query.Next() {
		err = query.Scan(&filaidioma.ID)
		resultado = append(resultado, filaidioma)
	}

	//Si no existe lo registramos, y despues sacamos su ID
	if len(resultado) == 0 {
		comando := "INSERT INTO Idioma(Idioma) VALUES ('" + objeto.Idioma + "')"
		fmt.Println(comando)
		insert, err := db.Query("INSERT INTO Idioma(Idioma) VALUES (?)", objeto.Idioma)
		if err != nil {
			panic(err.Error())
		}

		insert.Close()

		comando = "SELECT ID FROM Idioma WHERE (Idioma = '" + objeto.Idioma + "')"
		fmt.Println(comando)
		query, err := db.Query("SELECT ID FROM Idioma WHERE (Idioma = ?)", objeto.Idioma)

		if err != nil {
			panic(err.Error())
		}
		defer query.Close()

		for query.Next() {

			err = query.Scan(&filaidioma.ID)

		}
	}
	fmt.Println("ID de Idioma: ", filaidioma.ID)

	//Consultamos la palabra en la Base de Datos, si existe sacamos su ID
	comando = "SELECT ID FROM Palabra WHERE (Palabra = '" + objeto.Palabra + "')"
	fmt.Println(comando)
	query, err = db.Query("SELECT ID FROM Palabra WHERE (Palabra = ?)", objeto.Palabra)

	if err != nil {
		panic(err.Error())
	}
	defer query.Close()
	var filapalabra = model.RIdioma{}
	resultado = make([]model.RIdioma, 0)
	for query.Next() {
		err = query.Scan(&filapalabra.ID)
		resultado = append(resultado, filapalabra)
	}

	//Si no existe lo registramos, y despues sacamos su ID
	if len(resultado) == 0 {
		comando := "INSERT INTO Palabra(Palabra) VALUES ('" + objeto.Palabra + "')"
		fmt.Println(comando)
		insert, err := db.Query("INSERT INTO Palabra(Palabra) VALUES (?)", objeto.Palabra)
		if err != nil {
			panic(err.Error())
		}

		insert.Close()

		comando = "SELECT ID FROM Palabra WHERE (Palabra = '" + objeto.Palabra + "')"
		fmt.Println(comando)
		query, err := db.Query("SELECT ID FROM Palabra WHERE (Palabra = ?)", objeto.Palabra)

		if err != nil {
			panic(err.Error())
		}
		defer query.Close()

		for query.Next() {

			err = query.Scan(&filapalabra.ID)

		}
	}
	fmt.Println("ID de Palabra: ", filapalabra.ID)

	comando = "INSERT INTO Traduccion(PalabraID, Traduccion, IdiomaID, Fecha, FechaConsulta) VALUES ('" + string(filapalabra.ID) + "', '" + objeto.Traduccion + "', '" + string(filaidioma.ID) + "', utc_timestamp(), utc_timestamp())"
	fmt.Println(comando)
	insert, err := db.Query("INSERT INTO Traduccion(PalabraID, Traduccion, IdiomaID, Fecha, FechaConsulta) VALUES (?, ?, ?, utc_timestamp(), utc_timestamp())", filapalabra.ID, objeto.Traduccion, filaidioma.ID)
	if err != nil {
		panic(err.Error())
	}
	insert.Close()
}
