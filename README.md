# Base comercial B2B (Go + Fiber + Templ + HTMX + PocketBase)

Primera versión navegable para vender servicios web B2B por industria.

## Stack
- Go
- Fiber
- Templ
- HTMX
- PocketBase (opcional en local; el proyecto funciona con seed en memoria)

## Rutas públicas
- `GET /`
- `GET /soluciones`
- `GET /soluciones/:slug`
- `GET /casos`
- `GET /casos/:industry`
- `GET /diagnostico`
- `POST /diagnostico`
- `GET /p/:slug`

## Rutas admin mínimas
- `GET /admin`
- `GET /admin/prospectos`
- `GET /admin/propuestas`
- `POST /admin/prospectos/:id/status` (HTMX para cambio de estado)

## Colecciones PocketBase esperadas
- `industries`
- `solutions`
- `case_studies`
- `prospects`
- `proposals`
- `events`

## Datos seed incluidos
- Agrícola: sitio web + ERP simple + tablets.
- Centro comercial: directorio + tótem + pantallas.
- Educación: web + tótem + pantallas + WhatsApp + test.
- Legal/profesional: sitio web + agendamiento + WhatsApp.

## Cómo correr localmente
1. Copiar variables de entorno:
   ```bash
   cp .env.example .env
   ```
2. Instalar dependencias:
   ```bash
   go mod tidy
   ```
3. Levantar modo desarrollo:
   ```bash
   make dev
   ```
4. Abrir `http://localhost:3000`.

## Build
```bash
make build
```

## Notas
- `POST /diagnostico` usa HTMX y guarda prospectos en memoria; si PocketBase está configurado también intenta persistir en `prospects`.
- `GET /p/:slug` incrementa `opened_count`, actualiza `last_opened_at` y crea evento `proposal_view` en memoria; si PocketBase está configurado también envía evento.
