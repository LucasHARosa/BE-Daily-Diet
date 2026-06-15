# Daily Diet API — Roadmap

Back-end em Go para o app **Daily Diet**. Cada fase termina com algo funcionando e testável.

---

## Stack

| Camada | Lib |
|---|---|
| Roteamento | chi v5 |
| Banco | PostgreSQL 16 (Docker) |
| Driver | pgx v5 (pgxpool) |
| Queries | sqlc (código gerado a partir de SQL puro) |
| Migrations | goose |
| Autenticação | JWT (golang-jwt/jwt v5) |
| Senha | bcrypt |
| Validação | go-playground/validator v10 |
| Env | godotenv |

---

## Fase 1 — Fundação ✅

> Objetivo: API sobe, banco conecta, `/health` responde.

- [x] Estrutura de pastas (`cmd/`, `internal/`, `db/`)
- [x] Docker Compose com PostgreSQL 16
- [x] `.env` com variáveis de ambiente
- [x] `internal/config/config.go` — carrega `.env`
- [x] `internal/infra/database/postgres.go` — pool de conexões (pgxpool)
- [x] `internal/http/responses/json.go` — helpers `JSON()` e `Error()`
- [x] `cmd/api/main.go` — servidor HTTP com `/health`
- [x] `go mod tidy` — go.sum sincronizado

**Verificar:**
```powershell
docker compose up -d
go run ./cmd/api
# GET http://localhost:3333/health → OK
```

---

## Fase 2 — Autenticação

> Objetivo: usuário cria conta, faz login, recebe JWT, acessa rota protegida.

**Migrations a criar:**
- [ ] `001_create_users.sql`
- [ ] `002_create_refresh_tokens.sql`

**Queries (`db/queries/`):**
- [ ] `users.sql` — CreateUser, GetUserByEmail, GetUserByID
- [ ] `refresh_tokens.sql` — CreateRefreshToken, GetRefreshToken, DeleteRefreshToken

**`sqlc generate` → `internal/infra/postgres/sqlc/`**

**Camadas:**
- [ ] `internal/infra/auth/password_hasher.go` — bcrypt hash + compare
- [ ] `internal/infra/auth/jwt_service.go` — gerar e validar access token
- [ ] `internal/infra/postgres/repositories/user_repository.go`
- [ ] `internal/application/auth/service.go` — Register, Login
- [ ] `internal/http/middleware/auth_middleware.go` — valida JWT, injeta user_id no ctx
- [ ] `internal/http/handlers/auth_handler.go`
- [ ] `internal/http/routes/routes.go`
- [ ] `internal/shared/errors/errors.go` — ErrNotFound, ErrUnauthorized, ErrConflict...

**Endpoints:**
```
POST /users     → { id, name, email, created_at }
POST /sessions  → { access_token, user: { id, name, email } }
GET  /me        → { id, name, email }  ← rota protegida (testa middleware)
```

**Verificar:**
```
POST /users    → cria usuário
POST /sessions → recebe token
GET  /me       → retorna dados com Authorization: Bearer <token>
GET  /me       → retorna 401 sem token
```

---

## Fase 3 — CRUD de Refeições

> Objetivo: usuário autenticado gerencia suas próprias refeições. Toda query usa `WHERE user_id = $1`.

**Migration:**
- [ ] `003_create_meals.sql`

**Queries (`db/queries/meals.sql`):**
- [ ] CreateMeal, GetMealByIDAndUserID, ListMealsByUserFiltered, UpdateMeal, DeleteMeal

**Camadas:**
- [ ] `internal/infra/postgres/repositories/meal_repository.go`
- [ ] `internal/application/meals/service.go` — Create, Update, Delete, GetByID, List
- [ ] `internal/http/handlers/meal_handler.go`

**Endpoints:**
```
POST   /meals
GET    /meals?start=&end=&status=on_diet|off_diet
GET    /meals/:id
PUT    /meals/:id
DELETE /meals/:id
```

**Verificar:**
- Criar, listar, editar, excluir refeição
- Usuário A não vê/edita refeições do usuário B
- Conectar app mobile: substituir mocks da lista de refeições

---

## Fase 4 — Métricas

> Objetivo: calcular sequência, totais e agrupamentos por período.

**Queries (`db/queries/metrics.sql`):**
- [ ] GetMetricsSummary, ListMealsForSequence, ListMealsGroupedByDay

**Camadas:**
- [ ] `internal/domain/metrics.go` — `CalculateBestStreak(meals) int` (função pura, testável)
- [ ] `internal/application/metrics/service.go` — GetSummary, GetByPeriod
- [ ] `internal/http/handlers/metrics_handler.go`

**Endpoints:**
```
GET /metrics/summary
GET /metrics?start=2026-06-01&end=2026-06-30&groupBy=day|month
```

**Teste unitário:**
```go
// [✓, ✓, ✗, ✓, ✓, ✓] → bestStreak = 3
```

**Verificar:**
- Conectar app mobile: home (card de percentual) e tela de estatísticas

---

## Fase 5 — Plano Alimentar

> Objetivo: estrutura hierárquica plano → dias da semana → refeições planejadas.

**Migrations:**
- [ ] `004_create_food_plans.sql`
- [ ] `005_create_food_plan_days.sql`
- [ ] `006_create_food_plan_meals.sql` (inclui `food_plan_items`)

**Camadas:**
- [ ] Repository + Service + Handler para food plans
- [ ] Lógica de `is_active` — apenas um plano ativo por usuário
- [ ] Integração com refeições: ao criar refeição com `food_plan_meal_id`, salvar snapshot JSONB

**Endpoints:**
```
POST   /food-plans
GET    /food-plans
GET    /food-plans/active
GET    /food-plans/:id
PUT    /food-plans/:id
DELETE /food-plans/:id

POST   /food-plans/:id/days/:weekday/meals
PUT    /food-plan-meals/:id
DELETE /food-plan-meals/:id
```

**Verificar:**
- Conectar app mobile: tela de plano alimentar

---

## Fase 6 — Perfil + Estimativa de Calorias

> Objetivo: dados físicos do usuário e estimativa de calorias via IA.

**Migration:**
- [ ] `007_create_user_profiles.sql`
- [ ] `008_create_calorie_estimations.sql`

**Perfil:**
- [ ] Repository + Service + Handler para user profiles
- [ ] Campos: weight_kg, height_cm, birth_date, body_fat_percentage, basal_calories, activity_level, gym_frequency_per_week

**Endpoints perfil:**
```
GET   /me/profile
PATCH /me/profile
```

**Estimativa de Calorias:**
- [ ] `internal/infra/ai/calorie_estimator.go` — chama API externa (Claude/OpenAI)
- [ ] A chave da API fica só no backend, nunca no app mobile
- [ ] Salva histórico em `calorie_estimations`

**Endpoint:**
```
POST /calorie-estimations
Body: { "description": "..." }
→ { estimatedCalories, confidence, items[], observation }
```

**Verificar:**
- Conectar app mobile: tela de perfil e botão "Estimar calorias" no form de refeição

---

## Fase 7 — Qualidade e Portfólio

- [ ] Testes unitários: `CalculateBestStreak`, validação de inputs, regras de autorização
- [ ] Testes de integração: fluxo register → login → criar refeição → listar → deletar
- [ ] Collection HTTP (arquivo `.http` no projeto com todos os endpoints)
- [ ] Diagrama do banco (dbdiagram.io ou draw.io)
- [ ] README com instruções de setup, stack e prints do app
- [ ] Deploy da API (Railway, Render, Fly.io ou VPS)

---

## Fluxo de dados (padrão de toda request)

```
chi route
  → AuthMiddleware (valida JWT, injeta user_id no ctx)
    → Handler (lê body/params, valida, pega user_id do ctx)
      → UseCase/Service (regras de negócio)
        → Repository interface
          → PostgresRepository (usa sqlc gerado)
            → pgx → PostgreSQL
```

> `user_id` **nunca** vem do body. Sempre do token via contexto.

---

## Comandos do dia a dia

```powershell
# Subir banco
docker compose up -d

# Rodar API
go run ./cmd/api

# Criar migration
goose -dir db/migrations create nome_da_migration sql

# Aplicar migrations
goose -dir db/migrations postgres "postgres://docker:docker@localhost:5432/daily_diet?sslmode=disable" up

# Status das migrations
goose -dir db/migrations postgres "postgres://docker:docker@localhost:5432/daily_diet?sslmode=disable" status

# Gerar código sqlc
sqlc generate

# Rodar testes
go test ./...

# Formatar código
gofmt -w .

# Organizar dependências
go mod tidy
```
