# Task Controller - Frontend

## Tecnologias

| Tecnologia | Descripcion |
|---|---|
| [React](https://react.dev/) 19 | Libreria UI |
| [TypeScript](https://www.typescriptlang.org/) | Tipado estatico |
| [Vite](https://vite.dev/) | Build tool y dev server |
| [Tailwind CSS](https://tailwindcss.com/) 4 | Framework de estilos |
| [shadcn/ui](https://ui.shadcn.com/) | Componentes UI (Radix + Tailwind) |
| [React Router](https://reactrouter.com/) 7 | Enrutamiento |
| [Lucide](https://lucide.dev/) | Iconos |
| [Sonner](https://sonner.emilkowal.dev/) | Notificaciones toast |

## Estructura del proyecto

```
frontend/
├── public/                      # Assets estaticos
├── src/
│   ├── assets/                  # Imagenes y recursos
│   ├── components/
│   │   ├── ui/                  # Componentes shadcn/ui (button, card, dialog, etc.)
│   │   ├── layout.tsx           # Layout principal de la aplicacion
│   │   └── protected-route.tsx  # Wrapper de rutas protegidas por autenticacion
│   ├── context/
│   │   └── auth.tsx             # Context de autenticacion (login, logout, token)
│   ├── lib/
│   │   ├── api.ts               # Cliente HTTP para comunicarse con el backend
│   │   └── utils.ts             # Utilidades (cn helper para clases CSS)
│   ├── pages/
│   │   ├── login.tsx            # Pagina de login
│   │   ├── change-password.tsx  # Pagina de cambio de contraseña
│   │   ├── tasks.tsx            # Listado de tareas
│   │   ├── task-detail.tsx      # Detalle de tarea
│   │   └── users.tsx            # Gestion de usuarios
│   ├── types/
│   │   └── index.ts             # Tipos e interfaces TypeScript
│   ├── App.tsx                  # Componente raiz con rutas
│   ├── main.tsx                 # Entry point
│   └── index.css                # Estilos globales (Tailwind)
├── components.json              # Configuracion de shadcn/ui
├── vite.config.ts               # Configuracion de Vite
├── tsconfig.json                # Configuracion de TypeScript
├── eslint.config.js             # Configuracion de ESLint
└── package.json
```

## Requisitos previos

- [Node.js](https://nodejs.org/) >= 18
- npm (incluido con Node.js)
- Backend corriendo en `http://localhost:8081`

## Instalacion y ejecucion

```bash
# Instalar dependencias
npm install

# Iniciar servidor de desarrollo
npm run dev
```

La aplicacion se levanta en `http://localhost:5173`.

## Comandos disponibles

| Comando | Descripcion |
|---|---|
| `npm run dev` | Inicia el servidor de desarrollo con HMR |
| `npm run build` | Compila TypeScript y genera build de produccion en `dist/` |
| `npm run preview` | Sirve el build de produccion localmente |
| `npm run lint` | Ejecuta ESLint sobre el proyecto |

## Proxy al backend

Vite esta configurado para redirigir las peticiones `/api` al backend:

```
/api/* → http://localhost:8081
```

Esto permite que el frontend haga peticiones a `/api/v1/task-controller/...` sin problemas de CORS durante el desarrollo.

## Agregar componentes de shadcn/ui

Para agregar nuevos componentes de shadcn/ui:

```bash
npx shadcn@latest add <nombre-del-componente>
```

Ejemplo:

```bash
npx shadcn@latest add tooltip
npx shadcn@latest add tabs
```

Los componentes se generan en `src/components/ui/`.

## Alias de imports

El proyecto usa el alias `@` para importar desde `src/`:

```tsx
import { Button } from "@/components/ui/button"
import { useAuth } from "@/context/auth"
import type { Task } from "@/types"
```
