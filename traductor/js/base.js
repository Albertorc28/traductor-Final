$(document).ready(function() {
    
    ActualizarIdiomas();
    $('#txtTexto').keyup(function(event) {
        if (event.keyCode === 13) {
            $("#btnEnviar").click();
        }
    });

    $("#btnNuevaTraduccion").click(function() {
        document.getElementById("idiomas").style.display = 'none';
        document.getElementById("txtIdioma").style.display = 'block';
        document.getElementById("txtTraduccion").style.display = 'block';
        document.getElementById("btnVolver").style.display = 'block';
        document.getElementById("btnNuevaTraduccion").style.display = 'none';
        document.getElementById("txtTexto").disabled = false;
        $("#txtTraduccion").val('');


    }); 

    $("#btnVolver").click(function() {
        document.getElementById("txtTraduccion").style.display = 'none';
        document.getElementById("sidiomas").disabled = false;
        document.getElementById("txtTexto").disabled = false;
        document.getElementById("btnVolver").style.display = 'none';
        document.getElementById("btnNuevaTraduccion").style.display = 'block';
        document.getElementById("idiomas").style.display = 'block';
        document.getElementById("txtIdioma").style.display = 'none';
        $("#Traduccion tbody").children().remove();
        $('#Traduccion thead').children().remove();
        $('#aviso').children().remove();
    }); 

    $("#btnEnviar").click(function() {

        var tvisible = $("#txtTraduccion:visible").length > 0;
        var ivisible = $("#txtIdioma:visible").length > 0;

        //Si el campo traduccion está oculto, entonces consultamos la traducción
        if (tvisible == false){
            var texto = $("#txtTexto").val();
            var id = $("#idiomas select").val();

            //Comprobamos si el campo de Texto está vacio, para controlar la entrada
            if (texto == "") {
                $('#aviso').html('<p>Ningún campo puede estar vacío.<pr>');
            }
            
            else {
                var envio = {
                    palabra: texto,
                    idioma: id
                };
                console.log(envio);
    
                $.post({
                    url:"/traduccion",
                    data: JSON.stringify(envio),
                    success: function(data, status, jqXHR) {
                        console.log(data);
                    },
                    dataType: "json"
    
                }).done(function(data) {
                    console.log("Petición realizada");
                    MostrarTraduccion(data);
                    
                }).fail(function(data) {
                    console.log("Petición fallida");
                
                }).always(function(data){
                    console.log("Petición completa");
                });
            }
            

        } 

        //Si el campo Traduccion está visible, entonces insertamos la traduccion en la Base de datos
        else {
        
            //Si el campo de idioma está habilitado para introducirlo manualmente cogemos su valor
            if (ivisible == true){
                var idioma = $("#txtIdioma").val();
            }
            //Si no cogemos el de el desplegable
            else {
                var idioma = $("#idiomas option:selected").text();
            }

            var texto = $("#txtTexto").val();
            var traduccion = $("#txtTraduccion").val();

            //Comprobamos que ningun campo esté vacío
            if (texto == "" || traduccion == "" || idioma == "") {
                $('#aviso').html('<p>Ningún campo puede estar vacío.<pr>');
            }
            else {
                var envio = {
                    palabra: texto,
                    traduccion: traduccion,
                    idioma: idioma
                };
                console.log(envio);
    
                $.post({
                    url:"/Insertar",
                    data: JSON.stringify(envio),
                    success: function(data, status, jqXHR) {
                        console.log(data);
                    },
                    dataType: "json"
    
                }).done(function(data) {
                    console.log("Petición realizada");
                    $('#aviso').html('<p>La traducción ha sido añadida correctamente.<pr>');
                    ActualizarIdiomas();
                
                }).fail(function(data) {
                    console.log("Petición fallida");
                
                }).always(function(data){
                    console.log("Petición completa");
                });
            }
        }
    });      
    
});

function ActualizarIdiomas() {
    $.ajax({
        url: "/idioma",
        method: "POST",
        dataType: "json",
        contentType: "application/json",
        success: function(data) {
            if(data != null)
                console.log(data.length + " idiomas obtenidos");
            ListarIdiomas(data);
        },
        error: function(data) {
            console.log(data);
        }
    });
}

function ListarIdiomas(array) {
    var select = $("#idiomas select");
    select.children().remove();
    if(array != null && array.length > 0) {

        for(var x = 0; x < array.length; x++) {
            select.append(
                "<option value='" + array[x].ID + 
                "'>" + array[x].Idioma + 
                "</option>");
        }
    } else {
        select.append('<option value=#> No hay registros </option>');
        
    }
}

function MostrarTraduccion(array) {
    $("#Traduccion tbody").children().remove();
    $('#Traduccion thead').children().remove();
    $('#aviso').children().remove();
    //Comprobamos el valor de Fail, si es true, no habrá ninguna traduccion
    if(array[0].Fail != true) {
        
        //Comprobamos el valor de Ok, Si es false, no ha encontrado la traduccion deseada, y se las mostraremos en otros idiomas
        if (array[0].Ok == false){
            document.getElementById("txtTraduccion").style.display = 'block';
            document.getElementById("sidiomas").disabled = true;
            document.getElementById("txtTexto").disabled = true;
            $('#aviso').html("<p>La palabra "+array[0].Palabra+" no ha sido encontrada en el idoma seleccionado. Si quiere puede añadir la traducción. <br> Aquí tiene la traducción a otros idomas.</p>");
        } 
        $('#Traduccion thead').append( "<th>Palabra</th><th>Idioma</th><th>Traduccion</th><th>Fecha de Registro</th><th>Fecha de Ultima Consulta</th>");
        for(var x = 0; x < array.length; x++) {
            $("#Traduccion tbody").append(
                "<tr><td>" + array[x].Palabra + 
                "</td><td>" + array[x].Idioma + 
                "</td><td>" + array[x].Traduccion + 
                "</td><td>" + moment(array[x].Fecha).format("DD-MM-YY HH:mm:ssZ") + 
                "</td><td>" + moment(array[x].FechaConsulta).format("DD-MM-YY HH:mm:ssZ") + 
                "</td></tr>");
                document.getElementById("btnVolver").style.display = 'block';
        }
    //Si no hay ninguna traduccion 
    } else {
        $('#Traduccion thead').empty();
        document.getElementById("txtTraduccion").style.display = 'block';
        document.getElementById("sidiomas").disabled = true;
        document.getElementById("txtTexto").disabled = true;
        document.getElementById("btnVolver").style.display = 'block';
        $('#aviso').html('<p>No Existe esa traducción, si quiere puede añadirla.<pr>');
      
        
    }
}


