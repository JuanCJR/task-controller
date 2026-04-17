# Task Controller

Sistema de gestion de tareas con control de acceso basado en roles (RBAC).

## Tecnologias

| Componente        | Tecnologia                |
| ----------------- | ------------------------- |
| Backend           | Go + Gin                  |
| Frontend          | React + TypeScript + Vite |
| UI                | shadcn/ui + Tailwind CSS  |
| Base de datos     | PostgreSQL 16             |
| Autenticacion     | JWT                       |
| Contenedores      | Docker + docker-compose   |
| Documentacion API | Swagger                   |

## Requisitos previos

- [Go](https://go.dev/dl/) >= 1.22
- [Node.js](https://nodejs.org/) >= 18
- [Docker](https://docs.docker.com/get-docker/) y [Docker Compose](https://docs.docker.com/compose/install/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) (para migraciones)

### Instalacion de herramientas Go

```bash
# golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# swag (generador de Swagger)
go install github.com/swaggo/swag/cmd/swag@latest

# air (hot reload)
go install github.com/air-verse/air@latest

# Asegurar que go/bin este en el PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

## Paso a paso para levantar el proyecto

### 1. Clonar el repositorio

```bash
git clone https://github.com/JuanCJR/task-controller.git
cd task-controller
```

### 2. Levantar la base de datos con Docker

```bash
cd backend

# Copiar variables de entorno
cp .env-example .env

# Editar .env con tus valores (o dejar los de ejemplo para desarrollo)

# Levantar PostgreSQL
docker-compose up -d
```

Verificar que el contenedor esta corriendo:

```bash
docker ps
```

Deberia aparecer `task_controller_db` con estado `Up`.

### 3. Ejecutar las migraciones

```bash
# Desde el directorio backend/
migrate -path migrations -database "postgres://YOUR_DB_USER:YOUR_DB_PASSWORD@localhost:5432/task_controller_db?sslmode=disable" up
```

Reemplazar `YOUR_DB_USER` y `YOUR_DB_PASSWORD` con los valores definidos en el `.env`.

Ejemplo con los valores por defecto del `.env-example`:

```bash
migrate -path migrations -database "postgres://your_db_user:your_db_password@localhost:5432/task_controller_db?sslmode=disable" up
```

### 4. Levantar el backend

```bash
# Desde el directorio backend/

# Instalar dependencias Go
go mod download

# Opcion A: Ejecutar directamente
go run cmd/server/main.go

# Opcion B: Ejecutar con hot reload (recomendado para desarrollo)
air
```

El servidor inicia en `http://localhost:8081` (o el puerto definido en `APP_PORT`).

En el primer inicio, si `EXECUTE_SEED=true`, se crearan automaticamente:

- 3 roles: Admin, Ejecutor, Auditor
- 12 permisos (CRUD para modulos user, task, role)
- Asignacion de permisos a roles
- Usuario admin con las credenciales definidas en `.env`

### 5. Levantar el frontend

```bash
# Desde el directorio frontend/

# Instalar dependencias
npm install

# Iniciar servidor de desarrollo
npm run dev
```

El frontend inicia en `http://localhost:5173` (puerto por defecto de Vite).

### 6. Verificar que todo funciona

- **Backend API:** `http://localhost:8081/api/v1/task-controller/auth/login`
- **Swagger UI:** `http://localhost:8081/docs/index.html`
- **Frontend:** `http://localhost:5173`

## Comandos utiles

### Backend

```bash
# Desde el directorio backend/

# Compilar
go build -o tmp/server ./cmd/server


# Generar/actualizar documentacion Swagger
swag init -g cmd/server/main.go

# Crear nueva migracion
migrate create -ext sql -dir migrations -seq nombre_de_la_migracion

# Aplicar migraciones
migrate -path migrations -database "postgres://USER:PASS@HOST:PORT/DB?sslmode=disable" up

# Revertir ultima migracion
migrate -path migrations -database "postgres://USER:PASS@HOST:PORT/DB?sslmode=disable" down 1
```

### Frontend

```bash
# Desde el directorio frontend/

# Desarrollo
npm run dev

# Build de produccion
npm run build

# Lint
npm run lint

# Preview del build
npm run preview
```

### Base de datos

```bash
# Desde el directorio backend/

# Iniciar PostgreSQL
docker-compose up -d

# Detener PostgreSQL
docker-compose down

# Detener y eliminar datos (volumen)
docker-compose down -v

# Ver logs de PostgreSQL
docker-compose logs -f db

# Conectarse a la base de datos
docker exec -it task_controller_db psql -U YOUR_DB_USER -d task_controller_db
```

## Estructura del proyecto

```
task-controller/
├── backend/           # API REST en Go (Gin)
│   ├── cmd/           # Entry point
│   ├── internal/      # Codigo interno (config, handlers, services, repos, models)
│   ├── migrations/    # Archivos SQL de migraciones
│   ├── pkg/           # Utilidades compartidas
│   ├── docs/          # Swagger generado
│   └── docker-compose.yaml
├── frontend/          # SPA en React + TypeScript (Vite)
│   ├── src/
│   └── package.json
└── README.md          # Este archivo
```

Para mas detalle sobre la arquitectura del backend, ver [backend/README.md](backend/README.md).
