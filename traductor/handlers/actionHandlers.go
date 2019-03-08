package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	client "traductor/data/dataclient"
	"traductor/data/model"
)

//ListIdioma Función que devuelve los idiomas de la base de datos
func ListIdioma(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming request from " + r.URL.EscapedPath())
	if r.URL.Path != PathListadoIdiomas {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	lista := client.ListarIdiomas()

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	respuesta, _ := json.Marshal(&lista)
	fmt.Fprint(w, string(respuesta))

}

//ListTraduccion Función que devuelve las traducciones de la base de datos dado un filtro
func ListTraduccion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming request from " + r.URL.EscapedPath())
	if r.URL.Path != PathListadoTraducciones {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	defer r.Body.Close()
	bytes, e := ioutil.ReadAll(r.Body)

	if e == nil {
		var filtro model.Traduccion
		e = json.Unmarshal(bytes, &filtro)
		//Ponemos en mayuscula la primera letra de la Palabra
		filtro.Palabra = strings.Title(filtro.Palabra)

		if e == nil {
			lista := client.ListarTraducciones(&filtro)
			w.WriteHeader(http.StatusOK)

			w.Header().Add("Content-Type", "application/json")

			respuesta, _ := json.Marshal(&lista)
			fmt.Fprint(w, string(respuesta))
		} else {
			fmt.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "La petición no pudo ser parseada")
			fmt.Fprintln(w, e.Error())
			return
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, e)
	}
}

//InsertTraduccion Función que inserta una petición en la base de datos local
func InsertTraduccion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming request from " + r.URL.EscapedPath())
	if r.URL.Path != PathInsertarTraducciones {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	defer r.Body.Close()
	bytes, e := ioutil.ReadAll(r.Body)

	if e == nil {
		var peticion model.Insercion
		enTexto := string(bytes)
		fmt.Println("En texto: " + enTexto)
		_ = json.Unmarshal(bytes, &peticion)

		//Ponemos en mayuscula la primera letra de cada palabra
		peticion.Palabra = strings.Title(peticion.Palabra)
		peticion.Idioma = strings.Title(peticion.Idioma)
		peticion.Traduccion = strings.Title(peticion.Traduccion)

		if peticion.Palabra == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "La petición está vacía")
			return
		}

		w.WriteHeader(http.StatusOK)

		w.Header().Add("Content-Type", "application/json")

		respuesta, _ := json.Marshal(peticion)
		fmt.Fprint(w, string(respuesta))

		go client.InsertTraducciones(&peticion)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, e)
	}
}
