# Sys-Called — Sistema de Controle de Chamados Internos

Este é o meu projeto para o Desafio Técnico Full Stack (nível básico) da Codificar: um sistema para que funcionários de uma empresa abram chamados internos de suporte, e para que o time de suporte os acompanhe, atenda e distribua entre si.

Neste README explico como rodar o projeto, minhas escolhas tecnológicas e arquiteturais, como atendi cada requisito do desafio, e os trade-offs que assumi conscientemente.

## Como executar o projeto

### Dependências

- Docker e Docker Compose (é tudo que eu preciso pra rodar o projeto inteiro — não instalei nada manualmente).
- Opcionalmente, Go 1.26+ se eu quiser rodar os serviços fora de container (ver segunda opção abaixo).

Não precisei criar nem popular banco manualmente: as migrações (SQL puro, em `internal/infra/postgres/migrations/` de cada serviço) rodam automaticamente como um passo do `docker-compose.yml`, e cada serviço já sobe com dados seed (o `employees-service` registra 3 responsáveis fixos na primeira execução).

### Subindo tudo com Docker Compose (é o que eu recomendo)

```bash
docker compose up --build
```

Isso sobe, nesta ordem: os dois bancos Postgres (um por serviço), o Kafka, as migrações de cada serviço, e por fim os dois serviços. Depois de subir, ficam disponíveis:

| Serviço | URL |
|---|---|
| `ticket-service` | http://localhost:18080 |
| `employees-service` | http://localhost:18081 |
| Postgres do `ticket-service` | `localhost:5433` |
| Postgres do `employees-service` | `localhost:5434` |
| Kafka | `localhost:9092` |

Para derrubar tudo e limpar os dados dos bancos:

```bash
docker compose down -v
```

### Rodando os serviços fora do Docker (Go local + infra via Docker)

Se eu quiser rodar os binários Go diretamente (pra debugar, por exemplo), subo só a infraestrutura:

```bash
docker compose up postgres employees-postgres kafka kafka-init ticket-service-migrate employees-service-migrate
```

E em terminais separados:

```bash
DATABASE_URL="postgres://ticket_service:ticket_service@localhost:5433/ticket_service?sslmode=disable" \
KAFKA_BROKERS="localhost:9092" \
PORT=8080 \
go run ./services/ticket-service/cmd/server
```

```bash
DATABASE_URL="postgres://employees_service:employees_service@localhost:5434/employees_service?sslmode=disable" \
KAFKA_BROKERS="localhost:9092" \
PORT=8081 \
go run ./services/employees-service/cmd/server
```

### Comandos que deixei prontos no `makefile`

```bash
make build   # compila os dois binários em bin/
make test    # roda go test ./... em cada módulo Go
make fmt     # go fmt em cada módulo
make tidy    # go mod tidy em cada módulo
make docker-up    # docker compose up --build -d
make docker-down  # docker compose down
```

## Minhas escolhas tecnológicas e arquiteturais

O enunciado deixou a stack livre e citou Go como uma das tecnologias usadas no dia a dia da Codificar — escolhi Go por isso, e porque acho a linguagem um bom encaixe pra um domínio com bastante regra de negócio e concorrência (várias pessoas mexendo no mesmo chamado).

Fui bem além do mínimo pedido pelo desafio — sei que isso é mais do que o "nível básico" pede, e quero deixar isso explícito em vez de esconder. Decidi usar essa entrega pra mostrar como eu abordaria esse problema em um cenário real de produto que vai crescer, não só como resolver o exercício da forma mais rápida possível. Concretamente, isso significou:

- **Dois serviços independentes** (`ticket-service` e `employees-service`) em vez de uma aplicação única, cada um dono do seu próprio banco. Separei porque "chamados" e "funcionários" são domínios com motivos de mudança bem diferentes — quem mexe em regra de atendimento não deveria precisar tocar em cadastro de funcionário, e vice-versa.
- **Event Sourcing** no `ticket-service`: em vez de guardar só o estado atual de um chamado, guardo cada evento que aconteceu com ele (aberto, editado, atribuído, prioridade mudada, etc.) e reconstruo o estado a partir do histórico. Escolhi isso porque rastreabilidade — "quem fez o quê e quando" — é exatamente o tipo de coisa que a pessoa que pediu o sistema (a área administrativa, no enunciado) provavelmente vai querer no futuro, e é muito mais barato ter isso desde o início do que adicionar depois.
- **DDD tático** (agregados, value objects, eventos de domínio, repositórios como interface) pra manter a regra de negócio isolada de framework HTTP e de banco.
- **Comunicação assíncrona entre os serviços via Kafka + Outbox Pattern**, em vez de um serviço chamar o outro por HTTP direto. Cheguei a implementar a versão síncrona primeiro, percebi que isso deixava a distribuição automática de chamados refém do `employees-service` estar no ar, e troquei pela versão assíncrona — deixei esse processo de decisão registrado no [ADR-005](docs/decisions/ADR-005-comunicacao-assincrona-employees-ticket.md), porque acho que o raciocínio por trás da mudança é mais interessante do que só o resultado final.
- **Testes em todas as camadas**, seguindo TDD como prática (documentei essa decisão à parte no [ADR-002](docs/decisions/ADR-002-TDD.md)) — pra mim isso pesa mais do que ter mais funcionalidades, como o próprio desafio sugere ("Qualidade > Quantidade").

Documentei cada decisão maior como um ADR em [`docs/decisions/`](docs/decisions/), incluindo os antipadrões e dívidas técnicas que assumi conscientemente em cada serviço — prefiro deixar isso visível a fingir que não existe:

| ADR | Assunto |
|---|---|
| [ADR-001](docs/decisions/ADR-001-architecture.md) | Arquitetura geral: microsserviços orientados a eventos, DDD, bounded contexts |
| [ADR-002](docs/decisions/ADR-002-TDD.md) | Adoção de TDD |
| [ADR-003](docs/decisions/ADR-003-ticket-service.md) | Implementação do `ticket-service`: camadas, event sourcing, padrões e antipadrões |
| [ADR-004](docs/decisions/ADR-004-employees-service.md) | Implementação do `employees-service`: persistência, outbox |
| [ADR-005](docs/decisions/ADR-005-comunicacao-assincrona-employees-ticket.md) | Comunicação assíncrona entre os dois serviços via Outbox + Kafka |
| [ADR-006](docs/decisions/ADR-006-ci-e-dependabot.md) | Pipeline de CI e atualização de dependências |

## Como atendi cada requisito do desafio

### 2.0 — Cadastro de chamados

Implementei cadastro (`POST /tickets`), edição (`PUT /tickets/:id`), listagem (`GET /tickets`) e visualização (`GET /tickets/:id`). Cada chamado tem título, descrição, prioridade (`Low`/`Medium`/`High`), status (`Open`/`In Progress`/`Closed`), responsável e data/hora de abertura — os campos mínimos pedidos, mais o histórico de respostas trocadas no chamado, que decidi incluir porque achei natural pro caso de uso.

### 3.0 — Responsáveis pelo atendimento

Segui a sugestão do próprio enunciado de não construir um cadastro completo: o `employees-service` sobe com 3 responsáveis fixos (`GET /employees` lista os disponíveis). Optei por fazer disso um serviço de verdade, com seu próprio banco, em vez de uma lista estática dentro do `ticket-service` — o raciocínio completo está no [ADR-004](docs/decisions/ADR-004-employees-service.md).

### 4.0 — Distribuição automática

`POST /tickets/:id/assign/auto` atribui o chamado ao responsável com menos chamados em aberto no momento; `POST /tickets/:id/assign` continua disponível pra atribuição manual. Defini "em aberto" como qualquer chamado com status diferente de `Closed` (ou seja, `Open` e `In Progress` contam como carga de trabalho) — decidi assim porque é o que reflete trabalho pendente de verdade para quem está atendendo.

### 5.0 — Listagem e acompanhamento

`GET /tickets` retorna todos os chamados. Não priorizei filtros/ordenação/busca no backend porque não cheguei a construir a tela que consumiria isso (ver "o que não entreguei" abaixo) — não fazia sentido pra mim construir parâmetros de filtro sem uma UI real pra validar o que faz sentido pro caso de uso.

### 6.0 — Funcionamento da aplicação

Coberto pelas instruções de execução no topo deste README.

## Bibliotecas e referências externas

- [Gin](https://github.com/gin-gonic/gin) — escolhi por ser o framework HTTP em Go mais maduro e com melhor ergonomia pra roteamento e binding de JSON que eu conheço; não vi motivo pra escrever isso na mão.
- [pgx/v5](https://github.com/jackc/pgx) — driver Postgres; preferi ele ao `database/sql` genérico + driver `lib/pq` porque tem pool de conexões nativo (`pgxpool`) e é o driver recomendado atualmente pra Postgres em Go.
- [google/uuid](https://github.com/google/uuid) — geração e parsing de UUID, usado nos IDs de agregado.
- [segmentio/kafka-go](https://github.com/segmentio/kafka-go) — cliente Kafka; escolhi essa lib especificamente por ser Go puro (sem cgo), o que manteve meu Dockerfile multi-stage em Alpine simples, sem precisar instalar `librdkafka`.
- [testify](https://github.com/stretchr/testify) — asserções de teste (`assert`/`require`).
- Imagens Docker: `postgres:16-alpine`, `apache/kafka:3.8.0` (modo KRaft, sem Zookeeper), `golang:1.26-alpine` e `alpine:3.20` como base das minhas imagens.

## Trade-offs e o que eu não cheguei a entregar

O desafio pede pra eu documentar isso quando cortar escopo por tempo, então quero ser direto:

- **Não implementei o frontend.** O enunciado pede uma aplicação web com telas de cadastro, edição, listagem e seleção de responsável — eu só entreguei a API que sustentaria essas telas. Prioricei aprofundar a modelagem de domínio, os testes e a arquitetura de comunicação entre os serviços, e não sobrou tempo pra construir a interface nesta entrega. Se eu fosse continuar, o próximo passo seria um front em Vue (consistente com a stack que a Codificar usa) consumindo a API já pronta — os endpoints já retornam tudo que uma tela de listagem/edição precisaria.
- **Não implementei autenticação/autorização** em nenhuma rota — qualquer requisição é aceita. Fora de escopo pra essa entrega, mas seria bloqueador antes de qualquer uso real.
- **`ticket-service` não tem uma projeção de leitura dedicada** — a listagem e a distribuição automática reidratam os chamados a partir dos eventos a cada chamada. Funciona bem no volume de um desafio técnico, mas eu não deixaria assim em produção com um volume real de chamados (detalhei a solução que eu faria — uma projeção própria — no [ADR-003](docs/decisions/ADR-003-ticket-service.md)).
- **O cache de chamados é em memória e local a cada instância** do `ticket-service` — funciona rodando uma única instância (como faço aqui), mas não seria seguro com múltiplas réplicas sem trocar por um cache compartilhado.

Fui mais detalhista sobre essas e outras decisões (incluindo antipadrões que assumi de propósito) nos ADRs de cada serviço, linkados acima.

## Testes

Rodo os testes unitários (sem nenhuma dependência externa) com:

```bash
make test
```

Os testes de integração contra Postgres real ficam atrás de uma variável de ambiente e são pulados automaticamente se ela não estiver definida — assim não preciso de banco no ar só pra rodar a suíte unitária no dia a dia:

```bash
# ticket-service
TICKET_SERVICE_TEST_DATABASE_URL="postgres://ticket_service:ticket_service@localhost:5433/ticket_service?sslmode=disable" \
  go test ./services/ticket-service/internal/infra/postgres/...

# employees-service
EMPLOYEES_SERVICE_TEST_DATABASE_URL="postgres://employees_service:employees_service@localhost:5434/employees_service?sslmode=disable" \
  go test ./services/employees-service/internal/infra/postgres/...
```

## CI/CD

Configurei um pipeline em [`.github/workflows/ci.yml`](.github/workflows/ci.yml) que roda em todo push/PR: build, `go vet`, testes com `-race` e checagem de `go.mod`/`go.sum` pra cada serviço; build das duas imagens Docker; e um smoke test que sobe a stack inteira via `docker compose` e exercita o fluxo real (abrir chamado → distribuir automaticamente → listar), validando que a integração assíncrona entre os serviços funciona de verdade, não só nos testes unitários com fakes. Deixei o [`dependabot.yml`](.github/dependabot.yml) configurado pra manter as dependências Go, as imagens Docker e as próprias GitHub Actions atualizadas. Detalhei as decisões dessa configuração no [ADR-006](docs/decisions/ADR-006-ci-e-dependabot.md).

## Referência de API

### `ticket-service`

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/healthz` | Health check |
| `POST` | `/tickets` | Abre um novo chamado |
| `GET` | `/tickets` | Lista todos os chamados |
| `GET` | `/tickets/:id` | Detalhe de um chamado |
| `PUT` | `/tickets/:id` | Edita título/descrição |
| `POST` | `/tickets/:id/assign` | Atribui manualmente (`{"assignee_id": "..."}`) |
| `POST` | `/tickets/:id/assign/auto` | Atribui automaticamente ao menos ocupado |
| `POST` | `/tickets/:id/priority` | Muda a prioridade (`{"priority": "Low\|Medium\|High"}`) |
| `POST` | `/tickets/:id/start` | Move para "em andamento" |
| `POST` | `/tickets/:id/close` | Fecha o chamado |
| `POST` | `/tickets/:id/responses` | Adiciona uma resposta |

### `employees-service`

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/healthz` | Health check |
| `GET` | `/employees` | Lista os funcionários registrados |

## Estrutura do repositório

```
services/<serviço>/
  cmd/server/main.go     # composição: monta as dependências e sobe o servidor HTTP
  internal/
    domain/               # entidades, value objects, eventos, interfaces de repositório
    application/          # um use case por ação, dependendo só de interfaces do domínio
    infra/                # implementações concretas (Postgres, HTTP, Kafka, cache)
docs/decisions/           # ADRs
.github/                  # CI e Dependabot
docker-compose.yml
makefile
```
