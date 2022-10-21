
## ⚙️ Configuración del entorno windows

Para comenzar a trabar con el back de sercide, es necesario tener una serie de instalaciones realizadas previamente
* Instalar **Ubuntu** `20.04` en Windows Store. Si no sabes como mira este [tutorial](https://www.youtube.com/watch?v=GdkSR8FKoRg)
* [Descargar](https://docs.docker.com/desktop/windows/install/) e instalar Docker-desktop 🐳. Si no sabes como mira este [tutorial](https://www.youtube.com/watch?v=mHdaxzgXvnQ)
* [Descargar](https://golang.org/dl/) e instalar la versión `1.18` de **Go** o superior en windows y en ubuntu. Si no sabes como mira este [tutorial para windows](https://www.youtube.com/watch?v=yf6afhvxLws) y este [tutorial para ubuntu](https://www.youtube.com/watch?v=4zVJBltNwD0)
* Instalar openapi-codegen `v1.9.0` para la generación de código **Go** a partir de la especificación **Yaml** en ubuntu
```bash
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0
```
* [Descargar](https://www.googleadservices.com/pagead/aclk?sa=L&ai=DChcSEwiU3PHqu6D4AhUNkGgJHXqJDZoYABAAGgJ3Zg&ae=2&ohost=www.google.com&cid=CAASJORolaWSzkY_1PLI0WNx1u6fXH67hb3kXqO_fOIC7OtBhuMiqQ&sig=AOD64_0wdIWsZGhLnpaadH4_r2GJN19e-Q&q&adurl&ved=2ahUKEwifsOrqu6D4AhVWP-wKHUF5AFIQ0Qx6BAgDEAE) e instalar **goland**
## ⚡️ Quick start
* Descargar el proyecto del repositorio
```bash
git clone https://github.com/SERCIDE/medidas-backend
```
* Levantar los servicios de docker compose para pruebas en local
```bash
docker compose up -d
```
* Configurar los ficheros `.env` y `service_account.json`. Si no los tienes, tendrás que solicitarlos
* Ejecutar el servicio completo
```bash
go run cmd/api/api.go
```