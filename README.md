# Sys-Called — Sistema de Controle de Chamados Internos

Este é o meu projeto para o Desafio Técnico Full Stack (nível básico) da Codificar: um sistema para que funcionários de uma empresa abram chamados internos de suporte, e para que o time de suporte os acompanhe, atenda e distribua entre si.

Neste README explico como rodar o projeto, minhas escolhas tecnológicas e arquiteturais, como atendi cada requisito do desafio, e os trade-offs que assumi conscientemente.

## Telas

| Login | Lista de chamados |
| --- | --- |
| ![Tela de login](docs/screenshots/01-login.png) | ![Lista de chamados em atendimento](docs/screenshots/02-lista-chamados.png) |
| **Detalhe do chamado (suporte)** | **Notificações (administrador)** |
| ![Detalhe de um chamado com as ações do suporte](docs/screenshots/03-detalhe-chamado.png) | ![Painel de notificações do administrador](docs/screenshots/04-notificacoes.png) |

## Como executar o projeto

### Dependências

- Docker e Docker Compose — é tudo que é preciso pra rodar o projeto inteiro.
- Opcionalmente, Go 1.26+ e Node 24+ pra rodar os serviços e o frontend fora de container (ver mais abaixo).

Não é preciso criar nem popular banco manualmente: as migrações (SQL puro, em `internal/adapters/out/postgres/migrations/` de cada serviço) rodam automaticamente como um passo do `docker-compose.yml`, e o `employees-service` já sobe com um usuário de cada perfil.

### Subindo tudo com Docker Compose (é o que eu recomendo)

```bash
make env                  # cria o .env a partir do .env.example, com uma chave JWT nova
docker compose up --build
```

As credenciais dos bancos, o `CLUSTER_ID` do Kafka, a porta do frontend e a chave que assina os JWT ficam no `.env`, que não é versionado. O `.env.example` traz os valores de desenvolvimento. Sem o `.env`, o compose se recusa a subir e diz qual variável falta. Quem não tiver `make` pode fazer `cp .env.example .env`; aí o `JWT_PRIVATE_KEY` fica vazio e o serviço gera uma chave temporária a cada subida.

Depois de subir, é só abrir **http://localhost:8000** e entrar com um dos usuários de desenvolvimento (a senha de todos é `senha123`):

| Username | Perfil | O que vê ao entrar |
|---|---|---|
| `usuario` | Usuário padrão | Abrir novo chamado, Meus chamados |
| `ana`, `bruno`, `carla` | Suporte | Chamados abertos, Em atendimento, Fechados por mim |
| `admin` | Administrador | Abrir novo chamado, Todos os chamados, Usuários, Suportes |

O compose sobe, nesta ordem: os dois bancos Postgres (um por serviço), o Kafka, as migrações, os dois serviços e o frontend. Cada etapa só começa quando a anterior está saudável: os bancos, o Kafka, os dois serviços (`/healthz`) e o nginx têm healthcheck, e `docker compose ps` mostra o estado de cada um. Os containers de longa duração sobem com `restart: unless-stopped` e com rotação de logs (3 arquivos de 10 MB).

O único ponto de entrada público é o nginx do frontend, que serve a SPA e faz o proxy de `/api/employees` e `/api/tickets` para os serviços. As outras portas ficam presas em `127.0.0.1`, só pra depuração e pra rodar os serviços fora do Docker:

| Serviço | URL |
|---|---|
| Frontend (Vue, servido por nginx) | http://localhost:8000 |
| `ticket-service` | http://localhost:18080 |
| `employees-service` | http://localhost:18081 |
| Postgres do `ticket-service` | `localhost:5433` |
| Postgres do `employees-service` | `localhost:5434` |
| Kafka (listener para o host) | `localhost:9094` |

Para derrubar tudo e limpar os dados dos bancos:

```bash
docker compose down -v
```

> Se você já tinha rodado uma versão anterior do projeto, rode `docker compose down -v` antes de subir de novo: os funcionários antigos no volume não têm username/senha, e não conseguiriam entrar.

### Rodando no Kubernetes (kind + Helm)

A mesma stack também sobe num cluster Kubernetes local, empacotada num chart Helm em [`infra/helm/sys-called`](infra/helm/sys-called). Precisa de Docker, [`kind`](https://kind.sigs.k8s.io/), `kubectl` e `helm`. O script [`infra/scripts/k8s.sh`](infra/scripts/k8s.sh) faz tudo:

```bash
./infra/scripts/k8s.sh up       # cria o cluster kind, faz o build das 3 imagens, carrega no cluster e instala o chart
```

Quando termina, o frontend responde em **http://localhost:30080**, com os mesmos usuários do compose. Os outros comandos:

```bash
./infra/scripts/k8s.sh deploy     # novo build com uma tag nova e helm upgrade no cluster existente
./infra/scripts/k8s.sh status     # pods, services e volumes
./infra/scripts/k8s.sh logs ticket-service
./infra/scripts/k8s.sh uninstall  # remove o release, mas mantém volumes e secret
./infra/scripts/k8s.sh down       # apaga o cluster inteiro
```

Argumentos extras do `up` e do `deploy` vão direto pro `helm` (por exemplo `--set ingress.enabled=true`). `CLUSTER_NAME`, `NAMESPACE`, `RELEASE` e `TAG` podem ser trocados por variável de ambiente.

O que o chart cria:

| Recurso | Detalhes |
|---|---|
| `ticket-postgres`, `employees-postgres` | StatefulSets com volume persistente, um banco por serviço, como no compose |
| `kafka` | StatefulSet em modo KRaft (sem Zookeeper), com volume persistente |
| `employees-service`, `ticket-service` | Deployments com readiness e liveness em `/healthz`, rodando sem root, com sistema de arquivos só leitura e sem capabilities |
| `frontend` | Deployment do nginx, exposto por NodePort (30080) ou, com `ingress.enabled=true`, por um Ingress |
| `sys-called-secrets` | senhas dos bancos e a chave Ed25519 dos JWT |

Algumas decisões:

- **Secrets gerados pelo próprio chart.** Se as senhas e a chave não vierem nos values, o chart gera senhas aleatórias e uma chave Ed25519 na primeira instalação. Nas atualizações seguintes, ele reaproveita o secret que já existe no cluster (via `lookup`), então os bancos continuam acessíveis e os tokens emitidos continuam válidos. O secret tem `helm.sh/resource-policy: keep` e sobrevive a um `helm uninstall`.
- **As migrações viajam dentro da imagem de cada serviço** (em `/migrations`) e rodam num init container antes do serviço subir. Assim o SQL não é duplicado no chart, e a versão da migração é sempre a da imagem que vai rodar. Outro init container garante que o tópico `employees.events` existe.
- **Os nomes dos Services são fixos** (`employees-service`, `ticket-service`, `kafka`...), porque o `nginx.conf` do frontend aponta pra eles. Isso significa um release por namespace.
- **O `ticket-service` fica com uma réplica**, porque o cache de chamados é em memória (ver os trade-offs). O `employees-service` e o frontend podem escalar pelos values.

### Compose, kind e produção: em que porta fica cada um

O compose e o kind são duas formas de subir a mesma aplicação, e cada um é uma cópia independente, com os próprios bancos. Podem rodar ao mesmo tempo (`make up-all`) porque não disputam porta:

| Ambiente | Endereço |
|---|---|
| Docker Compose | http://localhost:8000 |
| Kubernetes local (kind + Helm) | http://localhost:30080 |

Em cada um deles, tudo entra por uma porta só: o nginx do frontend serve as telas e repassa `/api/employees` e `/api/tickets` pros serviços. O navegador nunca fala direto com os serviços.

Em produção ficaria só o Kubernetes, atrás de um único endereço HTTPS:

```
https://chamados.suaempresa.com (443) ──▶ Ingress (TLS)
                                           │
                                           ▼
                                  Service frontend (nginx)
                                    ├── /                → telas (Vue)
                                    ├── /api/employees/  → employees-service
                                    └── /api/tickets/    → ticket-service
```

O Ingress recebe o tráfego na 443 (a 80 só redireciona) e cuida do certificado, por exemplo com cert-manager e Let's Encrypt. O chart já cria o Ingress com `--set ingress.enabled=true --set ingress.host=chamados.suaempresa.com`; falta só a configuração de TLS nele. O NodePort 30080 é um atalho do kind local e não seria usado, e bancos, Kafka e serviços continuam acessíveis só dentro do cluster.

### Rodando fora do Docker (Go e Node locais + infraestrutura via Docker)

Subo só a infraestrutura:

```bash
docker compose up postgres employees-postgres kafka kafka-init ticket-service-migrate employees-service-migrate
```

E em terminais separados:

```bash
DATABASE_URL="postgres://employees_service:employees_service@localhost:5434/employees_service?sslmode=disable" \
KAFKA_BROKERS="localhost:9094" \
PORT=8081 \
go run ./services/employees-service/cmd/server
```

```bash
DATABASE_URL="postgres://ticket_service:ticket_service@localhost:5433/ticket_service?sslmode=disable" \
KAFKA_BROKERS="localhost:9094" \
EMPLOYEES_JWKS_URL="http://localhost:8081/.well-known/jwks.json" \
PORT=8080 \
go run ./services/ticket-service/cmd/server
```

```bash
cd frontend && npm install && npm run dev
```

O Vite sobe em http://localhost:5173 com um proxy igual ao do nginx (`/api/employees` e `/api/tickets`). Se os serviços estiverem em outras portas, `EMPLOYEES_API_URL` e `TICKETS_API_URL` mudam o destino do proxy.

Variáveis de ambiente dos serviços:

| Serviço | Variável | Para quê |
|---|---|---|
| ambos | `DATABASE_URL`, `KAFKA_BROKERS`, `PORT` | conexão com banco, broker e porta HTTP |
| `employees-service` | `JWT_PRIVATE_KEY` (opcional) | chave Ed25519 em PEM (PKCS8) que assina os tokens. No compose ela vem do `.env` (o `make env` gera uma). Sem ela, o serviço gera uma chave temporária a cada subida e avisa no log |
| `ticket-service` | `EMPLOYEES_JWKS_URL` | de onde buscar as chaves públicas pra validar os tokens |

### Comandos que deixei prontos no `makefile`

```bash
make test           # testes Go (os dois serviços) + testes do frontend
make test-go        # só os testes Go
make test-frontend  # só os testes do frontend
make dev-frontend   # sobe o Vite em modo desenvolvimento
make build          # compila os dois binários Go em bin/
make fmt            # go fmt em cada módulo
make tidy           # go mod tidy em cada módulo
make env            # cria o .env com uma chave JWT nova (não sobrescreve um .env existente)
make docker-up      # cria o .env se faltar e roda docker compose up --build -d
make docker-ps      # estado e saúde de cada container
make docker-down    # docker compose down
make k8s-up         # ./infra/scripts/k8s.sh up
make k8s-deploy     # ./infra/scripts/k8s.sh deploy
make k8s-status     # ./infra/scripts/k8s.sh status
make k8s-down       # ./infra/scripts/k8s.sh down
make helm-lint      # helm lint + helm template do chart
make up-all         # sobe o compose (localhost:8000) e o kind (localhost:30080) juntos
make down-all       # derruba o compose e apaga o cluster kind
```

## Minhas escolhas tecnológicas e arquiteturais

O enunciado deixou a stack livre e citou Go e Vue.js como tecnologias do dia a dia da Codificar, eu escolhi Go no backend e Vue no frontend por isso, e porque acho Go um bom encaixe pra um domínio com bastante regra de negócio e concorrência (várias pessoas mexendo no mesmo chamado).

Fui bem além do mínimo pedido pelo desafio — sei que isso é mais do que o "nível básico" pede, e quero deixar isso explícito em vez de esconder. Decidi usar essa entrega pra mostrar como eu abordaria esse problema num cenário real de produto que vai crescer. Concretamente, isso significou:

- **Dois serviços independentes** (`ticket-service` e `employees-service`), cada um dono do seu próprio banco. Separei porque "chamados" e "funcionários" são domínios com motivos de mudança bem diferentes — quem mexe em regra de atendimento não deveria precisar tocar em cadastro de funcionário, e vice-versa.
- **Event Sourcing** no `ticket-service`: em vez de guardar só o estado atual de um chamado, guardo cada evento que aconteceu com ele (aberto, editado, atribuído, prioridade mudada, etc.) e reconstruo o estado a partir do histórico. Cada evento também registra **quem** executou a ação, então o sistema tem uma trilha de auditoria completa desde o primeiro dia — exatamente o que a área administrativa do enunciado ia querer ("ninguém sabe o que ficou pra fazer").
- **DDD tático** (agregados, value objects, eventos de domínio) pra concentrar a regra de negócio no domínio — inclusive as regras de visibilidade e de permissão, que são políticas do próprio agregado de chamado.
- **Arquitetura hexagonal (portas e adaptadores)** dentro de cada serviço, pra manter domínio e casos de uso isolados de framework HTTP, banco, broker e biblioteca de JWT — detalho abaixo.
- **Comunicação assíncrona entre os serviços via Kafka + Outbox Pattern**, em vez de um serviço chamar o outro por HTTP. Cheguei a implementar a versão síncrona primeiro, percebi que isso deixava a distribuição automática de chamados refém do `employees-service` estar no ar, e troquei pela assíncrona: o `employees-service` grava o funcionário e o evento `EmployeeRegistered` na tabela de outbox na mesma transação, um relay publica no Kafka, e o `ticket-service` mantém sua própria lista de responsáveis a partir desses eventos (só funcionários com perfil de suporte entram nela).
- **Autenticação própria com JWT e perfis** (usuário, suporte, administrador) — detalho numa seção própria abaixo.
- **TDD como prática em todo o código**, backend e frontend, no ciclo vermelho → verde → refatoração: escrevo o teste, vejo ele falhar pelo motivo certo, implemento o mínimo, e só então refatoro com os testes verdes como rede de segurança. Pra mim isso pesa mais do que ter mais funcionalidades, como o próprio desafio sugere ("Qualidade > Quantidade").

### Arquitetura interna: portas e adaptadores

Os dois serviços seguem a mesma organização, e a regra é uma só: as dependências apontam sempre para dentro.

- **`domain/`** — agregados, value objects, eventos, políticas e erros. Não conhece nenhuma outra camada nem biblioteca de infraestrutura. É o domínio que decide quais eventos acontecem (por exemplo, `employee.Register` já registra o `EmployeeRegistered`) e quem pode fazer o quê com um chamado (`ticket.CanEdit`, `ticket.CanWork`, `IsVisibleTo`...).
- **`application/`** — os casos de uso e as **portas de saída** (`application/port/out`), que são as interfaces de que os casos de uso precisam: `EventStore`, `TicketCache`, `ResponsibleDirectory`, `TokenVerifier` no `ticket-service`; `EmployeeRepository`, `PasswordHasher`, `TokenIssuer`, `SigningKeys`, `RefreshTokenStore`, `OutboxStore`, `EventPublisher` no `employees-service`. No `ticket-service` separei escrita e leitura em `command/` e `query/` — com event sourcing isso deixa o caminho pronto pra, no futuro, as queries lerem de uma projeção sem tocar nos commands.
- **`adapters/in/`** — quem dispara os casos de uso: HTTP (Gin), o middleware de autenticação, o consumer Kafka do `ticket-service` e o relay do outbox do `employees-service`. Adaptadores de entrada nunca falam com adaptadores de saída; o middleware de autenticação, por exemplo, passa pelo caso de uso `Authenticate` em vez de chamar direto o verificador de JWKS.
- **`adapters/out/`** — implementações das portas: Postgres, cache em memória, publisher Kafka, bcrypt, emissor de JWT, verificador JWKS.
- **`cmd/server/main.go`** — o único lugar que conhece tudo: instancia os adaptadores, monta os casos de uso e os injeta nos adaptadores de entrada.

Pra essa regra não depender só de disciplina, cada serviço tem um teste de arquitetura (`internal/architecture_test.go`) que lê os imports de todos os pacotes e falha se o domínio importar aplicação ou adaptadores, se a aplicação importar adaptadores, se um adaptador importar outro, ou se domínio/aplicação importarem Gin, pgx, kafka-go ou golang-jwt. Ele roda junto com o resto da suíte no `make test` e no CI.

### Autenticação e autorização

A autenticação fica centralizada no `employees-service`, e o `ticket-service` só valida os tokens:

- **Senhas** guardadas com bcrypt. Login com usuário inexistente e com senha errada respondem igual (`401 invalid credentials`) e levam o mesmo tempo — no caso do usuário inexistente comparo a senha contra um hash fictício, pra não dar pra descobrir quais usernames existem medindo o tempo de resposta.
- **Access token JWT de 15 minutos**, assinado com **Ed25519 (EdDSA)**. Escolhi assinatura assimétrica pra que o `ticket-service` valide os tokens só com a chave pública, sem nenhum segredo compartilhado entre os serviços. O token carrega `sub` (id do funcionário), `name`, `role`, `must_change_password`, `iss`, `aud` e `exp`, e o cabeçalho leva o `kid` da chave.
- **JWKS** em `GET /.well-known/jwks.json`: o `ticket-service` baixa as chaves públicas uma vez e guarda em cache; se aparece um token com `kid` desconhecido (troca de chave), ele busca de novo — no máximo uma vez a cada 5 segundos, pra que tokens com `kid` aleatório não virem uma forma de sobrecarregar o `employees-service`. O verificador exige `EdDSA` explicitamente (o que fecha o ataque clássico de "confusão de algoritmo" com HS256), além de emissor, audiência e expiração.
- **Refresh token opaco de 7 dias**, com **rotação a cada uso** e **detecção de reuso**: cada refresh gera um token novo e revoga o anterior; se um token já revogado for apresentado de novo (sinal de que foi roubado), todas as sessões daquele funcionário são revogadas. No banco fica só o hash SHA-256 do refresh token.
- **No navegador**, o access token fica só em memória (nunca em `localStorage`, que um XSS conseguiria ler) e o refresh token fica num cookie `HttpOnly; Secure; SameSite=Strict`, restrito ao caminho de autenticação. Quando a página é recarregada, o frontend recupera a sessão chamando o refresh; quando o access token expira no meio do uso, o cliente da API renova a sessão uma vez e repete a chamada.

- **Contas.** Qualquer pessoa pode criar uma conta, sempre com perfil de usuário padrão (só o administrador define quem é suporte ou administrador). A conta nasce **aguardando aprovação**: o login responde `403` com uma mensagem clara até o administrador aprovar. Senhas precisam ter pelo menos 8 caracteres.
- **Esqueci a senha.** Sem servidor de e-mail, o pedido vai pro administrador: ele recebe uma notificação e gera uma **senha temporária**, que aparece uma única vez na tela pra ele entregar pessoalmente. Quem entra com uma senha temporária é levado direto pra tela de troca de senha e não sai dela até escolher uma nova. O pedido responde sempre `202`, exista o username ou não, pra não revelar quais usernames existem.
- **Notificações.** O `employees-service` publica `EmployeeSignedUp` e `PasswordResetRequested` no Kafka (pelo outbox); o `ticket-service` consome esses eventos e grava notificações pros administradores, junto com as notificações de chamados. O frontend consulta a cada 30 segundos.

A autorização dos chamados fica no `ticket-service`, como regra de domínio:

| Perfil | Vê | Pode |
|---|---|---|
| Usuário | só os chamados que ele abriu | abrir chamado, responder; editar os seus só enquanto ninguém começou o atendimento |
| Suporte | os chamados abertos e os atribuídos a ele (inclusive os que ele fechou) | atribuir (a si ou a outro), distribuir automaticamente, mudar prioridade, responder; iniciar e fechar os atribuídos a ele |
| Administrador | todos | editar, atribuir, distribuir automaticamente, mudar prioridade, responder — **não** inicia nem fecha chamados, isso é trabalho do suporte responsável |

Um chamado que o ator não pode ver responde **404, igual a um chamado inexistente** — assim a API não revela que ele existe. Se ele pode ver mas não pode fazer a ação, a resposta é **403**. O autor de cada resposta é sempre quem está autenticado, nunca um campo enviado no corpo da requisição.

### Frontend

Vue 3 + Vite + TypeScript, com Vue Router e Tailwind CSS (o desafio pediu um framework CSS moderno; o Tailwind me deixou montar as telas rápido sem escrever CSS à mão). Em produção ele é servido por um nginx, que também faz o proxy de `/api/employees` e `/api/tickets` pros serviços — assim frontend e API ficam na mesma origem e o cookie `SameSite=Strict` funciona sem nenhuma configuração de CORS.

As telas:

- **Login**, com links pra "Criar conta" e "Recuperar senha". Uma conta ainda não aprovada recebe a mensagem "Sua conta ainda aguarda a aprovação do administrador".
- **Criar conta** — nome, username, senha e confirmação; ao terminar, explica que a conta precisa ser aprovada.
- **Recuperar senha** — pede o username e avisa que o administrador vai gerar uma senha temporária.
- **Trocar senha** — obrigatória pra quem entrou com senha temporária.
- **Sino de notificações** no canto superior direito, com o número de não lidas: o suporte recebe chamados novos, atribuições e mensagens; o usuário, mensagens e mudanças nos seus chamados; o administrador, chamados novos, contas aguardando aprovação e pedidos de nova senha (esses levam pra tela de Usuários).
- **Usuários** (administrador) — todos os funcionários com perfil e situação ("Ativo", "Aguardando aprovação", "Pediu nova senha"), com os botões **Aprovar** e **Gerar senha temporária**.
- **Suportes** (administrador) — cada atendente com a quantidade de chamados abertos, em andamento e fechados.
- **Tela inicial por perfil**, com o menu de cada um (tabela no começo deste README).
- **Abrir novo chamado** — num modal que sobe na própria tela (inicial ou lista), sem trocar de página: título, descrição e prioridade.
- **Lista de chamados** — com filtro por status sincronizado com a URL (os itens do menu do suporte são essa mesma lista filtrada), responsável pelo nome e link pro detalhe.
- **Detalhe do chamado** — dados (inclusive quando foi aberto e fechado), histórico de respostas, e só as ações que o usuário pode executar (editar, atribuir, distribuir automaticamente, mudar prioridade, iniciar, fechar, responder). Fechar abre um modal que mostra há quanto tempo o chamado está aberto e pede o relatório do que foi feito; o chamado fechado exibe esse relatório com o tempo total até o fechamento. O backend continua sendo quem garante as permissões; esconder botões é só pra não oferecer ação que vai ser recusada.

Os testes do frontend usam Vitest + Testing Library, testando pelo que o usuário vê (textos, campos, botões, links) e não pela estrutura interna dos componentes.

## Como atendi cada requisito do desafio

### 2.0 — Cadastro de chamados

Pela interface: abrir (tela "Abrir novo chamado"), editar (botão "Editar" no detalhe), listar (lista de chamados) e visualizar (detalhe). Pela API: `POST /tickets`, `PUT /tickets/:id`, `GET /tickets` e `GET /tickets/:id`. Cada chamado tem título, descrição, prioridade (baixa/média/alta), status (aberto/em andamento/fechado), responsável e data/hora de abertura — os campos mínimos pedidos, mais quem abriu o chamado e o histórico de respostas trocadas, que achei natural pro caso de uso.

### 3.0 — Responsáveis pelo atendimento

Segui a sugestão do enunciado de não construir um cadastro completo: o `employees-service` sobe com três atendentes (`ana`, `bruno`, `carla`), além de um administrador e um usuário padrão. Optei por fazer disso um serviço de verdade, com banco próprio, porque funcionários são outro domínio. Os atendentes aparecem pelo nome no seletor "Atribuir a" do detalhe do chamado.

### 4.0 — Distribuição automática

O botão "Distribuir automaticamente" (`POST /tickets/:id/assign/auto`) atribui o chamado ao atendente com menos chamados em aberto; a atribuição manual continua disponível. Defini "em aberto" como qualquer chamado com status diferente de fechado — ou seja, aberto e em andamento contam como carga de trabalho — porque é o que reflete trabalho pendente de verdade pra quem está atendendo.

### 5.0 — Listagem e acompanhamento

A lista mostra título, status, prioridade, responsável e data de abertura, com filtro por status. Como cada perfil já recebe do backend só o que pode ver, a mesma tela serve de "Meus chamados" pro usuário, de "Todos os chamados" pro administrador e de "Chamados abertos / Em atendimento / Fechados por mim" pro suporte.

### 6.0 — Funcionamento da aplicação

Um `make env` seguido de `docker compose up --build` sobe tudo, inclusive o frontend, com banco migrado e usuários de desenvolvimento prontos — detalhado no topo deste README.

## Bibliotecas e referências externas

Backend (Go):

- [Gin](https://github.com/gin-gonic/gin) — framework HTTP; escolhi pela maturidade e pela ergonomia de roteamento e binding de JSON.
- [pgx/v5](https://github.com/jackc/pgx) — driver Postgres, com pool de conexões nativo (`pgxpool`).
- [segmentio/kafka-go](https://github.com/segmentio/kafka-go) — cliente Kafka em Go puro (sem cgo), o que manteve os Dockerfiles em Alpine simples.
- [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) — emissão e validação de JWT com EdDSA.
- [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — hash de senhas.
- [google/uuid](https://github.com/google/uuid) — IDs dos agregados.
- [testify](https://github.com/stretchr/testify) — asserções de teste.

Frontend:

- [Vue 3](https://vuejs.org/), [Vue Router](https://router.vuejs.org/) e [Vite](https://vite.dev/), com TypeScript.
- [Tailwind CSS](https://tailwindcss.com/) — estilização.
- [Vitest](https://vitest.dev/), [Testing Library (Vue)](https://testing-library.com/docs/vue-testing-library/intro/), [user-event](https://testing-library.com/docs/user-event/intro) e [jsdom](https://github.com/jsdom/jsdom) — testes.

Imagens Docker: `postgres:16-alpine`, `apache/kafka:3.8.0` (modo KRaft, sem Zookeeper), `golang:1.26-alpine` e `alpine:3.20` nos serviços, `node:24-alpine` e `nginx:1.27-alpine` no frontend.

## Trade-offs e o que eu não cheguei a entregar

O desafio pede pra documentar isso quando cortar escopo por tempo, então quero ser direto:

- **A recuperação de senha passa pelo administrador**, e não por e-mail, porque um envio de e-mail de verdade exigiria um servidor de e-mail. A senha temporária é entregue pessoalmente.
- **As notificações chegam por consulta a cada 30 segundos**, não em tempo real. WebSocket ou Server-Sent Events deixariam a entrega instantânea, mas pra esse volume a consulta periódica é bem mais simples e suficiente.
- **Ainda não dá pra mudar o perfil de um funcionário pela tela** (promover um usuário a suporte, por exemplo); os atendentes e o administrador vêm do seed.
- **Um access token continua válido até expirar (no máximo 15 minutos)**, mesmo depois do logout ou de uma troca de chave. É a natureza do JWT sem estado; o que é revogado de imediato é o refresh token, então a sessão não é renovada.
- **A chave de assinatura dos JWT fica no `.env`**, fora do repositório, e por isso os tokens continuam válidos depois de um restart do `employees-service`. Em produção, ela viria de um gerenciador de secrets, não de um arquivo.
- **O compose e o chart ainda são para desenvolvimento**: HTTP sem TLS, um único broker Kafka e um Postgres por serviço sem réplica nem backup, rodando dentro do próprio cluster. Para produção eu colocaria TLS no Ingress (com cert-manager, por exemplo), usaria bancos e broker gerenciados e passaria os secrets por um gerenciador externo, publicando as imagens num registry em vez de carregá-las direto no kind.
- **O `ticket-service` não tem uma projeção de leitura dedicada** — a listagem e a distribuição automática reidratam os chamados a partir dos eventos a cada chamada. Funciona bem no volume de um desafio técnico, mas eu não deixaria assim em produção. A solução seria uma projeção própria atualizada a partir dos eventos, consumida pelos casos de uso em `application/query/`.
- **O filtro por status da lista é feito no navegador**, porque o backend já devolve só o que cada perfil vê e o volume é pequeno; com a projeção de leitura acima, ele passaria pro backend.
- **O cache de chamados é em memória e local a cada instância** do `ticket-service` — correto com uma instância (como aqui), mas precisaria virar um cache compartilhado com várias réplicas.
- **Eventos gravados antes da autenticação existir não sabem quem abriu o chamado**; nenhum usuário padrão os vê — só o administrador e, pela regra normal, o suporte (quando abertos ou atribuídos a ele). Num banco novo isso não acontece.

## Testes

```bash
make test
```

Roda os testes unitários e de arquitetura dos dois serviços Go e os testes do frontend, sem nenhuma dependência externa.

Os testes de integração contra Postgres real ficam atrás de uma variável de ambiente e são pulados se ela não estiver definida:

```bash
# ticket-service
TICKET_SERVICE_TEST_DATABASE_URL="postgres://ticket_service:ticket_service@localhost:5433/ticket_service?sslmode=disable" \
  go test ./services/ticket-service/internal/adapters/out/postgres/...

# employees-service
EMPLOYEES_SERVICE_TEST_DATABASE_URL="postgres://employees_service:employees_service@localhost:5434/employees_service?sslmode=disable" \
  go test ./services/employees-service/internal/adapters/out/postgres/...
```

Eles usam IDs únicos por execução, então podem rodar com a stack no ar. Um cuidado: os funcionários de teste criados pelo `employees-service` são publicados no Kafka como qualquer outro, e aparecem como atendentes no `ticket-service` se ele estiver rodando — pra uma demonstração limpa, rode-os contra um banco dedicado ou faça `docker compose down -v` depois.

## CI/CD

O pipeline em [`.github/workflows/ci.yml`](.github/workflows/ci.yml) roda em todo push/PR:

- **Go** (por serviço): checagem de `go.mod`/`go.sum`, build, `go vet` e testes com `-race`.
- **Frontend**: `npm ci`, testes, type-check e build.
- **Imagens Docker** dos dois serviços e do frontend.
- **Smoke test**: cria o `.env` com `make env`, sobe a stack inteira via `docker compose up --wait` (que só segue quando todos os healthchecks passam) e exercita o fluxo real — confere que a API recusa chamadas sem token, faz login, abre um chamado, distribui automaticamente e lista; e confere que o frontend é servido e que o login pelo proxy do nginx devolve o cookie de refresh no caminho certo. Isso valida de verdade a integração assíncrona e a autenticação entre os serviços, não só os testes com fakes.

- **Helm**: `helm lint` e `helm template` do chart, e `shellcheck` nos scripts.
- **Smoke test no Kubernetes**: cria um cluster kind no próprio runner, roda `infra/scripts/k8s.sh up` e repete o fluxo de login, abertura, distribuição automática e listagem pelo proxy do frontend. Isso garante que o chart continua instalável e funcionando.

O [`dependabot.yml`](.github/dependabot.yml) mantém atualizadas as dependências Go e npm, as imagens Docker e as próprias GitHub Actions.

## Referência de API

As rotas abaixo são as dos serviços. Pelo frontend (nginx ou Vite), elas ficam sob `/api/employees/...` e `/api/tickets/...`.

### `employees-service`

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/healthz` | Health check |
| `POST` | `/auth/login` | Login (`{"username", "password"}`). Devolve o access token e grava o refresh token em cookie |
| `POST` | `/auth/refresh` | Troca o cookie de refresh por uma sessão nova (rotação) |
| `POST` | `/auth/logout` | Revoga a sessão e apaga o cookie |
| `GET` | `/.well-known/jwks.json` | Chaves públicas pra validar os tokens |
| `POST` | `/auth/signup` | Cria uma conta de usuário padrão aguardando aprovação (`{"name", "username", "password"}`). `409` se o username já existe, `400` se a senha tiver menos de 8 caracteres |
| `POST` | `/auth/password-reset-requests` | Pede ao administrador uma senha nova (`{"username"}`). Responde sempre `202` |
| `POST` | `/auth/change-password` | Troca a senha de quem está autenticado (`{"current_password", "new_password"}`) |
| `GET` | `/employees` | Lista os funcionários com perfil e situação (só administrador) |
| `POST` | `/employees/:id/approve` | Aprova uma conta (só administrador) |
| `POST` | `/employees/:id/temporary-password` | Gera uma senha temporária e devolve `{"temporary_password"}` uma única vez (só administrador) |

### `ticket-service`

Todas as rotas, menos `/healthz`, exigem `Authorization: Bearer <access token>` e respeitam a tabela de perfis acima.

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/healthz` | Health check |
| `POST` | `/tickets` | Abre um chamado (`{"title", "description", "priority"}`) |
| `GET` | `/tickets` | Lista os chamados que o usuário pode ver |
| `GET` | `/tickets/:id` | Detalhe de um chamado, com as respostas |
| `PUT` | `/tickets/:id` | Edita título/descrição |
| `POST` | `/tickets/:id/assign` | Atribui manualmente (`{"assignee_id"}`) |
| `POST` | `/tickets/:id/assign/auto` | Atribui ao atendente menos ocupado |
| `POST` | `/tickets/:id/priority` | Muda a prioridade (`{"priority": "Low\|Medium\|High"}`) |
| `POST` | `/tickets/:id/start` | Move para "em andamento" |
| `POST` | `/tickets/:id/close` | Fecha o chamado com o relatório do que foi feito (`{"resolution"}`, obrigatório) |
| `POST` | `/tickets/:id/responses` | Adiciona uma resposta (`{"content"}`); o autor é quem está autenticado |
| `GET` | `/responsibles` | Lista os atendentes com nome |
| `GET` | `/responsibles/workload` | Quantidade de chamados abertos, em andamento e fechados por atendente (só administrador) |
| `GET` | `/notifications` | Notificações de quem está autenticado, com o número de não lidas |
| `POST` | `/notifications/read` | Marca as notificações como lidas |

## Estrutura do repositório

```
services/<serviço>/
  cmd/server/main.go        # composição: monta adaptadores e casos de uso e sobe o servidor
  internal/
    domain/                 # agregados, value objects, eventos, políticas e erros de domínio
    application/
      port/out/             # portas de saída: interfaces que os casos de uso usam
      port/out/outtest/     # implementações falsas das portas, pros testes
      command/  query/      # (ticket-service) casos de uso de escrita e de leitura
      auth/                 # (ticket-service) autenticação do ator
      *.go                  # (employees-service) casos de uso
    adapters/
      in/                   # HTTP, middleware de autenticação, consumer Kafka, relay do outbox
      out/                  # Postgres (+ migrations), cache, Kafka, bcrypt, JWT, JWKS
    architecture_test.go    # garante a regra de dependência entre as camadas
frontend/
  src/
    views/                  # telas (cada uma com seu .spec.ts)
    components/             # componentes reutilizáveis (campos, cabeçalho, ações do chamado)
    auth/                   # login, sessão, restauração e saída
    tickets/                # cliente da API, permissões, status, prioridades, formatação
    router/                 # rotas, guarda de navegação, menus por perfil
  nginx.conf                # serve a SPA e faz o proxy da API
infra/
  helm/sys-called/          # chart Helm da stack inteira
  kind/cluster.yaml         # cluster kind local, expondo o frontend na porta 30080
  scripts/k8s.sh            # sobe, atualiza, inspeciona e derruba a stack no Kubernetes
.github/                    # CI e Dependabot
.env.example                # variáveis do docker compose (copiado pro .env por make env)
docker-compose.yml
makefile
```
