package handlers

import "net/http"

//PathInicio Ruta raíz
const PathInicio string = "/"

//PathJSFiles Ruta a la carpeta de scripts de javascript
const PathJSFiles string = "/js/"

//PathCSSFiles Ruta a la carpeta de estilos css
const PathCSSFiles string = "/css/"

//PathListadoTraducciones Ruta de obtención de las traducciones
const PathListadoTraducciones string = "/traduccion"

//PathListadoIdiomas Ruta de obtención de los idiomas
const PathListadoIdiomas string = "/idioma"

//PathInsertarTraducciones Ruta de insercion de traducciones
const PathInsertarTraducciones string = "/Insertar"

//ManejadorHTTP encapsula como tipo la función de manejo de peticiones HTTP, para que sea posible almacenar sus referencias en un diccionario
type ManejadorHTTP = func(w http.ResponseWriter, r *http.Request)

//Manejadores Lista es el diccionario general de las peticiones que son manejadas por nuestro servidor
var Manejadores map[string]ManejadorHTTP

func init() {
	Manejadores = make(map[string]ManejadorHTTP)
	Manejadores[PathInicio] = IndexFile
	Manejadores[PathJSFiles] = JsFile
	Manejadores[PathCSSFiles] = CssFile
	Manejadores[PathListadoIdiomas] = ListIdioma
	Manejadores[PathInsertarTraducciones] = InsertTraduccion
	Manejadores[PathListadoTraducciones] = ListTraduccion
}
