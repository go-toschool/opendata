# Aries

<img src="aries.png" alt="Aries" align="left" width="160" />

Este servicio está encargado de verificar la valides de los token enviamos por un proveedor, además de comprobar por la valides de las credenciales de un usuario, especificamente, del token de autenticación compartido por el usuario, previa autorización de este.

Esta validación de hace por medio de un request HTTP hacia nuestro servidor principal Sigiriya, quien retorna la valides de este token.

Además se valida el token de acceso de dicho partner consultando el Servicio de [Capricornius](https://github.com/Finciero/opendata/capricornius).

Una vez validadas las llaves de acceso, se realiza un request hacia el Servicio de [Geminis](https://github.com/Finciero/opendata/gamini) quien realiza las consultas de saldo y transacciones de dicho usuario.
