# Task Controller - Backend

## Tecnologias

- **Lenguaje:** Go
- **Framework:** Gin
- **Base de datos:** PostgreSQL
- **Autenticacion:** JWT
- **Modelo de roles:** RBAC (Role-Based Access Control)
- **Documentacion API:** Swagger (swag + gin-swagger)
- **Migraciones:** golang-migrate
- **Hot Reload:** Air

## Estructura de directorios

```
backend/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point
├── internal/
│   ├── config/
│   │   └── config.go                # Carga y validacion de variables de entorno
│   ├── database/
│   │   ├── database.go              # Conexion a PostgreSQL
│   │   └── seed.go                  # Seed de datos iniciales (roles, permisos, admin)
│   ├── dto/
│   │   ├── auth_request.go          # DTOs de autenticacion
│   │   ├── user_request.go          # DTOs de usuarios
│   │   ├── task_request.go          # DTOs de tareas
│   │   └── task_comment_request.go  # DTOs de comentarios
│   ├── handler/
│   │   ├── handler.go               # Interface RouteRegister
│   │   ├── auth_handler.go          # Endpoints: login, logout, change-password
│   │   ├── user_handler.go          # Endpoints: CRUD users
│   │   ├── task_handler.go          # Endpoints: CRUD tasks
│   │   └── task_comment_handler.go  # Endpoints: comentarios de tareas
│   ├── middleware/
│   │   ├── auth.go                  # JWT validation middleware
│   │   └── rbac.go                  # Middleware de permisos por rol
│   ├── model/
│   │   ├── user.go                  # User struct
│   │   ├── role.go                  # Role struct
│   │   ├── permission.go            # Permission struct + Action/Module constants
│   │   ├── role_permission.go       # RolePermission struct
│   │   ├── user_role.go             # UserRole struct
│   │   ├── task.go                  # Task struct + TaskState constants
│   │   └── task_comment.go          # TaskComment struct
│   ├── repository/
│   │   ├── user_repository.go       # Queries de usuarios
│   │   ├── role_repository.go       # Queries de roles
│   │   ├── permission_repository.go # Queries de permisos
│   │   ├── task_repository.go       # Queries de tareas
│   │   └── task_comment_repository.go # Queries de comentarios
│   ├── server/
│   │   └── server.go                # Configuracion de Gin y registro de rutas
│   └── service/
│       ├── auth_service.go          # Login, logout, cambio de contraseña
│       ├── user_service.go          # CRUD usuarios + reglas de negocio
│       ├── task_service.go          # CRUD tareas + reglas de negocio
│       └── task_comment_service.go  # Comentarios + reglas de negocio
├── migrations/                      # Archivos SQL de migraciones
├── pkg/
│   └── utils/
│       ├── jwt.go                   # Generacion y validacion de JWT
│       └── hash.go                  # Hashing de contraseñas (bcrypt)
├── docs/                            # Documentacion Swagger generada
├── .air.toml                        # Configuracion de Air (hot reload)
├── .env-example                     # Variables de entorno de ejemplo
├── docker-compose.yaml              # PostgreSQL con Docker
├── go.mod
└── go.sum
```

## Capas de la aplicacion

| Capa | Responsabilidad |
|---|---|
| `cmd/` | Entry point. Inicializa dependencias y arranca el servidor. |
| `internal/handler/` | Recibe HTTP requests, valida input, llama al service. |
| `internal/service/` | Logica de negocio (reglas RBAC, validaciones de estado, etc.) |
| `internal/repository/` | Acceso a datos (queries SQL). |
| `internal/model/` | Structs que representan las tablas de BD. |
| `internal/dto/` | Structs de request/response (separados del modelo de BD). |
| `internal/middleware/` | Auth JWT y control de permisos por rol. |
| `internal/config/` | Carga y validacion de variables de entorno. |
| `internal/database/` | Conexion a PostgreSQL y seed de datos iniciales. |
| `pkg/utils/` | Utilidades reutilizables que no dependen del dominio. |

## Flujo de una request

```
Cliente HTTP
    │
    ▼
  Router (server.go)
    │
    ▼
  Middleware (auth.go → rbac.go)
    │
    ▼
  Handler (valida input, llama al service)
    │
    ▼
  Service (aplica reglas de negocio)
    │
    ▼
  Repository (ejecuta queries SQL)
    │
    ▼
  PostgreSQL
```

## Endpoints

### Auth

| Metodo | Ruta | Descripcion | Auth |
|---|---|---|---|
| POST | `/api/v1/task-controller/auth/login` | Login de usuario | No |
| POST | `/api/v1/task-controller/auth/logout` | Logout de usuario | No |
| PUT | `/api/v1/task-controller/auth/change-password` | Cambio de contraseña | Si |

### Users

| Metodo | Ruta | Descripcion | Permiso |
|---|---|---|---|
| GET | `/api/v1/task-controller/users` | Listar usuarios | read:user |
| POST | `/api/v1/task-controller/users` | Crear usuario | create:user |
| PUT | `/api/v1/task-controller/users/:id` | Actualizar usuario | update:user |
| DELETE | `/api/v1/task-controller/users/:id` | Eliminar usuario | delete:user |

### Tasks

| Metodo | Ruta | Descripcion | Permiso |
|---|---|---|---|
| GET | `/api/v1/task-controller/tasks` | Listar tareas | read:task |
| GET | `/api/v1/task-controller/tasks/:id` | Ver tarea | read:task |
| POST | `/api/v1/task-controller/tasks` | Crear tarea | create:task |
| PUT | `/api/v1/task-controller/tasks/:id` | Actualizar tarea | update:task |
| PATCH | `/api/v1/task-controller/tasks/:id/state` | Actualizar estado | update:task |
| DELETE | `/api/v1/task-controller/tasks/:id` | Eliminar tarea | delete:task |

### Task Comments

| Metodo | Ruta | Descripcion | Permiso |
|---|---|---|---|
| GET | `/api/v1/task-controller/tasks/:id/comments` | Listar comentarios | read:task |
| POST | `/api/v1/task-controller/tasks/:id/comments` | Crear comentario | update:task |

## Reglas de negocio

### Admin
- CRUD completo de tareas y usuarios
- Solo puede asignar tareas a usuarios con rol **Ejecutor**
- No puede crear usuarios con rol **Admin** (solo Ejecutor o Auditor)
- No puede actualizar ni eliminar tareas en estado distinto a **Asignado**

### Ejecutor
- Solo ve las tareas asignadas a el
- Puede actualizar el estado de sus tareas (con transiciones validas)
- No puede actualizar el estado si la tarea esta vencida
- Puede agregar comentarios solo en tareas vencidas asignadas a el

### Auditor
- Solo puede ver todas las tareas (lectura)

### General
- Todo usuario nuevo inicia con `must_change_password = true`
- El middleware RBAC bloquea el acceso a recursos si el usuario debe cambiar su contraseña
- El login devuelve el flag `must_change_password` para que el frontend redirija al cambio

## Diagramas

### Modelo RBAC - Permisos por rol

```mermaid
graph LR
    subgraph Roles
        Admin
        Ejecutor
        Auditor
    end

    subgraph "Permisos User"
        CU[create:user]
        RU[read:user]
        UU[update:user]
        DU[delete:user]
    end

    subgraph "Permisos Task"
        CT[create:task]
        RT[read:task]
        UT[update:task]
        DT[delete:task]
    end

    subgraph "Permisos Role"
        CR[create:role]
        RR[read:role]
        UR[update:role]
        DR[delete:role]
    end

    Admin --> CU & RU & UU & DU
    Admin --> CT & RT & UT & DT
    Admin --> CR & RR & UR & DR

    Ejecutor --> RT & UT

    Auditor --> RT
```

### Transiciones de estado de tareas

```mermaid
stateDiagram-v2
    [*] --> Asignado : Admin crea tarea

    Asignado --> Iniciado : Ejecutor inicia

    Iniciado --> EnEspera : Ejecutor pausa
    Iniciado --> FinalizadaExito : Ejecutor completa
    Iniciado --> FinalizadaError : Ejecutor reporta error

    EnEspera --> Iniciado : Ejecutor retoma
    EnEspera --> FinalizadaExito : Ejecutor completa
    EnEspera --> FinalizadaError : Ejecutor reporta error

    FinalizadaExito --> [*]
    FinalizadaError --> [*]

    state "En espera" as EnEspera
    state "Finalizada con exito" as FinalizadaExito
    state "Finalizada con error" as FinalizadaError
```

### Flujo de creacion de tarea (Admin)

```mermaid
flowchart TD
    A[Admin envia POST /tasks] --> B{assigned_to tiene rol Ejecutor?}
    B -- No --> C[400: tasks can only be assigned to Ejecutor]
    B -- Si --> D[Crear tarea con estado Asignado]
    D --> E[201: tarea creada]
```

### Flujo de actualizacion de estado (Ejecutor)

```mermaid
flowchart TD
    A[Ejecutor envia PATCH /tasks/:id/state] --> B{Tarea asignada a este usuario?}
    B -- No --> C[400: solo puedes actualizar tus tareas]
    B -- Si --> D{Tarea vencida?}
    D -- Si --> E[400: no se puede actualizar tarea vencida]
    D -- No --> F{Transicion de estado valida?}
    F -- No --> G[400: transicion no permitida]
    F -- Si --> H[Actualizar estado]
    H --> I[200: estado actualizado]
```

### Flujo de comentarios en tareas (Ejecutor)

```mermaid
flowchart TD
    A[Ejecutor envia POST /tasks/:id/comments] --> B{Tarea asignada a este usuario?}
    B -- No --> C[400: solo puedes comentar tus tareas]
    B -- Si --> D{Tarea tiene fecha de expiracion?}
    D -- No --> E[400: tarea sin fecha de expiracion]
    D -- Si --> F{Tarea vencida?}
    F -- No --> G[400: solo puede comentar tareas vencidas]
    F -- Si --> H[Crear comentario]
    H --> I[201: comentario creado]
```

### Flujo de creacion de usuario (Admin)

```mermaid
flowchart TD
    A[Admin envia POST /users] --> B{Rol solicitado es Admin?}
    B -- Si --> C[400: no se puede crear usuarios Admin]
    B -- No --> D{Rol existe?}
    D -- No --> E[400: rol no encontrado]
    D -- Si --> F[Crear usuario + asignar rol en transaccion]
    F --> G[must_change_password = true]
    G --> H[201: usuario creado]
```

### Flujo de autenticacion y acceso a recursos

```mermaid
sequenceDiagram
    actor U as Usuario
    participant H as Handler
    participant Auth as Auth Middleware
    participant RBAC as RBAC Middleware
    participant S as Service
    participant DB as PostgreSQL

    U->>H: Request con Bearer token
    H->>Auth: Validar JWT
    alt Token invalido
        Auth-->>U: 401 Unauthorized
    end
    Auth->>RBAC: userID en contexto
    RBAC->>DB: must_change_password?
    alt Debe cambiar password
        RBAC-->>U: 403 password change required
    end
    RBAC->>DB: Tiene permiso action:module?
    alt Sin permiso
        RBAC-->>U: 403 insufficient permissions
    end
    RBAC->>H: Continuar
    H->>S: Ejecutar logica de negocio
    S->>DB: Query
    DB-->>S: Resultado
    S-->>H: Respuesta
    H-->>U: JSON response
```

## Variables de entorno

Copiar `.env-example` a `.env` y configurar:

```bash
cp .env-example .env
```

| Variable | Descripcion | Ejemplo |
|---|---|---|
| `APP_PORT` | Puerto del servidor HTTP | `8081` |
| `JWT_SECRET` | Clave secreta para JWT | `your_jwt_secret_key` |
| `TOKEN_EXPIRATION` | Expiracion del token en horas | `3600` |
| `DB_HOST` | Host de PostgreSQL | `localhost` |
| `DB_PORT` | Puerto de PostgreSQL | `5432` |
| `DB_USER` | Usuario de PostgreSQL | `your_db_user` |
| `DB_PASSWORD` | Contraseña de PostgreSQL | `your_db_password` |
| `DB_NAME` | Nombre de la base de datos | `task_controller_db` |
| `EXECUTE_SEED` | Ejecutar seed al iniciar | `true` |
| `DEFAULT_ADMIN_EMAIL` | Email del admin inicial | `admin@taskcontroller.com` |
| `DEFAULT_ADMIN_PASSWORD` | Contraseña del admin inicial | `admin` |

## Migraciones

El proyecto usa [golang-migrate](https://github.com/golang-migrate/migrate) para gestionar migraciones SQL.

### Instalacion del CLI

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Crear una nueva migracion

```bash
migrate create -ext sql -dir migrations -seq nombre_de_la_migracion
```

Esto genera dos archivos:
- `XXXXXX_nombre_de_la_migracion.up.sql` — aplica el cambio
- `XXXXXX_nombre_de_la_migracion.down.sql` — revierte el cambio

### Orden de migraciones

| Orden | Migracion | Dependencias |
|---|---|---|
| 1 | `create_users_table` | Ninguna |
| 2 | `create_roles_table` | Ninguna |
| 3 | `create_permissions_table` | Ninguna |
| 4 | `create_role_permissions_table` | roles, permissions |
| 5 | `create_user_roles_table` | users, roles |
| 6 | `create_tasks_table` | users |
| 7 | `create_task_comments_table` | tasks, users |

### Comandos de migracion

```bash
# Aplicar todas las migraciones pendientes
migrate -path migrations -database "postgres://USER:PASSWORD@HOST:PORT/DB_NAME?sslmode=disable" up

# Revertir la ultima migracion
migrate -path migrations -database "postgres://USER:PASSWORD@HOST:PORT/DB_NAME?sslmode=disable" down 1

# Revertir todas las migraciones
migrate -path migrations -database "postgres://USER:PASSWORD@HOST:PORT/DB_NAME?sslmode=disable" down

# Ver version actual
migrate -path migrations -database "postgres://USER:PASSWORD@HOST:PORT/DB_NAME?sslmode=disable" version

# Forzar version (en caso de migracion corrupta)
migrate -path migrations -database "postgres://USER:PASSWORD@HOST:PORT/DB_NAME?sslmode=disable" force VERSION
```

## Swagger

La documentacion de la API se genera con [swag](https://github.com/swaggo/swag).

### Instalacion

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generar documentacion

```bash
swag init -g cmd/server/main.go
```

Esto genera/actualiza los archivos en el directorio `docs/`:
- `docs.go` — codigo Go con la spec embebida
- `swagger.json` — especificacion OpenAPI en JSON
- `swagger.yaml` — especificacion OpenAPI en YAML

Ejecutar este comando cada vez que se modifiquen los comentarios godoc en los handlers.

### Acceder a Swagger UI

Con el servidor corriendo:

```
http://localhost:8081/docs/index.html
```

## Hot Reload con Air

El proyecto usa [Air](https://github.com/air-verse/air) para hot-reload durante el desarrollo.

### Instalacion

```bash
go install github.com/air-verse/air@latest
```

Asegurarse de que `$HOME/go/bin` este en el `PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Uso

Desde el directorio `backend/`:

```bash
air
```

Air observa cambios en archivos `.go`, recompila y reinicia el servidor automaticamente. La configuracion se encuentra en `.air.toml`.
