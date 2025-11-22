## Índice

1. [📟 Principais tecnologias utilizadas](#📟-principais-tecnologias-utilizadas)
2. [💿 Instalação e Execução](#💿-instalação-e-execução)
3. [🌎 Visão Geral](#🌎-visão-geral)
4. [💵 Principais Regras de Negócio](#💵-principais-regras-de-negócio)
5. [📐 Arquitetura e Design](#📐-arquitetura-e-design)
6. [🚀 API - Endpoints HTTP](#🚀-api---endpoints-http)
    - [📡 Endpoints Expostos pela Aplicação](#📡-endpoints-expostos-pela-aplicação)
    - [📡 cURL dos Endpoints](#📡-curl-dos-endpoints)
    - [📟 Endpoints Consumidos pela Aplicação](#📟-endpoints-consumidos-pela-aplicação)
7. [✉️ Comunicação Assíncrona (Mensageria)](#✉️-comunicação-assíncrona-(mensageria))
    - [👂 Consumers](#👂-consumers)
    - [📣 Producers](#📣-producers)
8. [🎲 Modelo de Dados da Aplicação](#🎲-modelo-de-dados-da-aplicação)
9. [🚨 Estratégia de Testes](#🚨-estratégia-de-testes)
10. [🔎 Observabilidade](#🔎-observabilidade)
    - [Logs](#logs)
    - [Métricas](#métricas)
    - [Tracing](#tracing)
11. [🚔 Segurança](#🚔-segurança)



 # 📘 chat-websocket

Sistema de chat em tempo real desenvolvido em Go, com autenticação baseada em JWT, suporte a múltiplas salas e comunicação entre instâncias por meio de Redis Pub/Sub. Permite múltiplos servidores conectados, integrados via WebSocket, com autenticação simples e gerenciamento de mensagens. O balanceamento de carga e o acesso externo são realizados via Nginx, com ambiente pronto para deploy em Docker.

## 📟 Principais tecnologias utilizadas
- Go (Golang)
- Echo Framework
- Redis (Pub/Sub)
- JWT (JSON Web Token)
- WebSocket (gorilla/websocket)
- Nginx
- Docker & Docker Compose

 # 💿 Instalação e Execução

## Requisitos
- Go 1.24+
- Docker (para execução com containers)
- Docker Compose (para orquestração dos serviços)

## Instalação
```bash
# Caso queira rodar localmente sem Docker, instale as dependências Go:
go mod download
```

## Variáveis de Ambiente

- `APP_REDIS_ADDR`: Endereço do Redis utilizado pela aplicação (exemplo: `redis:6379`)
- `APP_SERVER_PORT`: Porta em que o servidor será iniciado (exemplo: `8080`)

*Essas variáveis são referenciadas no arquivo de configuração do Docker Compose e no código de configuração da aplicação.*

## Executando Localmente
```bash
# Para executar a aplicação localmente na porta padrão (ajuste as variáveis conforme necessário):

export APP_REDIS_ADDR=localhost:6379
export APP_SERVER_PORT=8080

# Inicie o Redis em outro terminal se necessário
docker run --name redis -p 6379:6379 redis:7-alpine

# Execute o servidor Go
go run ./cmd/server/main.go
```

## Usando Docker
```bash
# Para subir toda a infraestrutura (aplicação, múltiplas instâncias, Redis e Nginx load balancer):

docker compose up --build
```

---

 # 🌎 Visão Geral

O sistema é uma aplicação modular desenvolvida em **Go** (Golang), com arquitetura inspirada em Clean/Hexagonal Architecture. Seu foco principal é viabilizar **comunicação em tempo real via WebSocket**, com autenticação baseada em JWT, integração com Redis para Pub/Sub, e suporte a múltiplas instâncias por meio de balanceamento com Nginx e Docker.  

### Principais objetivos:
- Permitir troca de mensagens em tempo real entre usuários, agrupados por salas.
- Garantir autenticação simples e segura via JWT, com refresh de tokens.
- Escalar horizontalmente, permitindo múltiplos servidores conectados por Redis.
- Prover persistência leve de mensagens e histórico temporário via Redis.
- Integrar facilmente com infraestrutura moderna (Docker, Nginx).

---

## Módulos Principais

- `auth`: Gerenciamento de autenticação e autorização (JWT).
- `handler`: Camada de controle responsável por lidar com requisições HTTP e conexões WebSocket.
- `dto`: Estruturas de transferência de dados para mensagens e autenticação.
- `websocket`: Gerenciamento de clientes, salas (hubs), transmissão e recebimento de mensagens em tempo real.
- `redis`: Integração com Redis para Pub/Sub, histórico de mensagens, controle de mensagens não lidas.
- `config`: Inicialização e gerenciamento das configurações da aplicação.
- `router`: Definição dos endpoints HTTP e WebSocket.
- `logger`: Implementação de log centralizado para observabilidade.
- `infraestrutura`: Arquivos Docker, Nginx e configurações para deploy e desenvolvimento.

---

> **Nota:**  
O domínio do projeto é **comunicação/chat em tempo real**, direcionado para cenários onde múltiplos usuários interagem em salas, incluindo suporte a mensagens privadas. A análise foi baseada na estrutura de código e arquivos fornecidos, não havendo menção a setores específicos como financeiro ou comércio eletrônico.

 # 💵 Principais Regras de Negócio

**Lista Regras de Negócio:**

<details>
<summary>Autorização de Acesso via Token JWT</summary>

**Regra:**  
- Toda conexão WebSocket exige autenticação via token JWT fornecido no header `Authorization`.  
- O token é validado antes do usuário acessar qualquer sala.  
- Caso o token seja inválido ou ausente, a conexão é rejeitada com erro HTTP 401 (unauthorized).
- O token é limpo do prefixo "Bearer " antes da validação.

**Trecho do código:**
```go
tokenStr := r.Header.Get("Authorization")
if tokenStr == "" {
    http.Error(w, "unauthorized", http.StatusUnauthorized)
    ws.Close()
    return
}
if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
    tokenStr = tokenStr[7:]
}
claims, err := auth.ValidateAccessToken(tokenStr)
if err != nil {
    http.Error(w, "unauthorized", http.StatusUnauthorized)
    ws.Close()
    return
}
```
[Ver linhas 26–42](internal/websocket/client.go#L26-L42)
</details>

<details>
<summary>Permissão de Acesso às Salas de Chat</summary>

**Regra:**  
- O usuário só pode acessar salas (rooms) explicitamente permitidas em seu token JWT.
- O backend verifica se a sala requerida está presente no array `Rooms` das claims do usuário.
- Se não autorizado, rejeita a conexão com erro HTTP 403 (forbidden).

**Trecho do código:**
```go
authorized := false
for _, r := range claims.Rooms {
    if r == room {
        authorized = true
        break
    }
}
if !authorized {
    http.Error(w, "forbidden", http.StatusForbidden)
    ws.Close()
    return
}
```
[Ver linhas 45–56](internal/websocket/client.go#L45-L56)
</details>

<details>
<summary>Histórico das Últimas Mensagens da Sala</summary>

**Regra:**  
- Ao ingressar em uma sala, o usuário recebe automaticamente o histórico das últimas 50 mensagens daquela sala.
- Esse histórico é recuperado do Redis e enviado individualmente ao novo cliente conectado.

**Trecho do código:**
```go
if history, err := store.GetMessages(r.Context(), room, 50); err == nil {
    for _, msg := range history {
        client.Send <- []byte(msg.Content)
    }
}
```
[Ver linhas 64–68](internal/websocket/client.go#L64-L68)
</details>

<details>
<summary>Mensagens Privadas e Não Lidas</summary>

**Regra:**  
- Mensagens com o campo `Target` definido são consideradas privadas e enviadas apenas para o usuário alvo.
- Quando o usuário alvo está offline, a mensagem é armazenada no Redis como "não lida".
- Ao reconectar, o usuário recebe automaticamente todas as mensagens não lidas e, após o envio, elas são removidas do Redis.

**Trecho do código:**  
(Armazenamento de mensagem privada não lida)
```go
if msg.Target != "" {
    for _, clients := range h.Rooms {
        for client := range clients {
            if client.User == msg.Target {
                select {
                case client.Send <- []byte(msg.Content):
                default:
                    close(client.Send)
                    delete(clients, client)
                }
            }
        }
    }
    // Salva mensagem como não lida no Redis
    if h.ChatStore != nil {
        _ = h.ChatStore.SaveUnread(ctx, msg.Target, msg)
    }
    continue
}
```
[Ver linhas 66–83](internal/websocket/hub.go#L66-L83)


(Envio de mensagens não lidas ao conectar)
```go
go func(c *Client) {
    unread, err := h.ChatStore.GetUnreadMessages(ctx, c.User)
    if err != nil {
        h.logger.Error("Falha ao buscar mensagens não lidas", zap.String("user", c.User), zap.Error(err))
        return
    }
    for _, msg := range unread {
        payload, _ := json.Marshal(msg)
        c.Send <- payload
    }
    // Limpa mensagens não lidas depois de enviar
    _ = h.ChatStore.ClearUnread(ctx, c.User)
}(client)
```
[Ver linhas 31–43](internal/websocket/hub.go#L31-L43)
</details>

<details>
<summary>Persistência e Limite do Histórico de Mensagens</summary>

**Regra:**  
- As mensagens enviadas em uma sala são persistidas no Redis.
- É mantido um limite de 50 mensagens por sala; ao exceder esse limite, as mensagens mais antigas são descartadas (LTRIM).
- O histórico expira automaticamente após 6 horas.

**Trecho do código:**
```go
// LPUSH adiciona no início da lista
if err := cw.Client.LPush(ctx, key, payload).Err(); err != nil { ... }
// LTRIM mantém apenas as últimas `maxMessages`
if err := cw.Client.LTrim(ctx, key, 0, int64(maxMessages-1)).Err(); err != nil { ... }
if err := cw.Client.Expire(ctx, key, 6*time.Hour).Err(); err != nil { ... }
```
[Ver linhas 27–41](internal/redis/client.go#L27-L41)
</details>

<details>
<summary>Autenticação e Refresh Token</summary>

**Regra:**  
- O login exige nome de usuário e senha (hardcoded para "1234" neste exemplo).
- Após login bem-sucedido, são gerados e retornados tokens de acesso e refresh.
- O refresh token pode ser utilizado para obter um novo access token.

**Trecho do código:**
```go
user := r.FormValue("user")
password := r.FormValue("password")
if password != "1234" {
    http.Error(w, "invalid credentials", http.StatusUnauthorized)
    return
}
// salas que o usuário pode acessar
rooms := []string{"default", "vip"}
accessToken, _ := auth.GenerateAccessToken(user, rooms)
refreshToken, _ := auth.GenerateRefreshToken(user)
w.Header().Set("Content-Type", "application/json")
w.Write([]byte(`{"access_token":"` + accessToken + `","refresh_token":"` + refreshToken + `"}`))
```
[Ver linhas 109–120](internal/websocket/client.go#L109-L120)
</details>

---

## Categorias das Regras

| Categoria         | Descrição                                                                                      |
|-------------------|------------------------------------------------------------------------------------------------|
| Autorização       | Controle de acesso por JWT e permissão por sala.                                               |
| Persistência      | Armazenamento de histórico com limite de mensagens e expiração automática.                     |
| Processos automáticos | Envio automático de histórico e mensagens não lidas ao reconectar.                         |
| Validações        | Exige token válido e credenciais válidas para login.                                           |
| Privacidade       | Encaminhamento de mensagens privadas e armazenamento de não lidas para usuários offline.       |

---

> Para mais detalhes sobre cada regra e trechos de código associados, utilize os links diretos fornecidos em cada bloco.

 # 📐 Arquitetura e Design

A aplicação segue uma arquitetura modular inspirada em Clean Architecture e Hexagonal, com separação clara de responsabilidades entre camadas de entrada (handlers/controllers), camada de domínio (lógica de negócio), camada de infraestrutura (Redis, WebSocket, Logger) e configuração. Os componentes são organizados em pacotes internos, privilegiando o encapsulamento e facilitando a manutenção e escalabilidade.

## Camadas

- **Handler (Controller):** Responsável por receber e tratar requisições HTTP/WebSocket, validar autenticação, inicializar conexões e interagir com a camada de domínio (Hub). Exemplos: `internal/handler/auth.go`, `internal/handler/websocket.go`.

- **Domain/Hub (Service):** Centraliza a lógica do chat em tempo real, gerenciando clientes conectados, broadcast, mensagens privadas e integração com o armazenamento de mensagens. Exemplo: `internal/websocket/hub.go`.

- **DTO:** Define objetos de transferência de dados utilizados entre camadas, especialmente para autenticação e mensagens. Exemplo: `internal/dto/message.go`.

- **Infraestrutura (Repository/Adapter):** Implementa a persistência e integração externa, principalmente com Redis para Pub/Sub e armazenamento de mensagens. Exemplo: `internal/redis/client.go`, `internal/redis/publisher.go`, `internal/redis/subscriber.go`.

- **Configuração:** Gerencia carregamento e mapeamento das configurações da aplicação. Exemplo: `internal/config/config.go`.

- **Logger:** Fornece logging estruturado centralizado. Exemplo: `internal/logger/logger.go`.

## Diagrama:

```mermaid
graph TD;
  subgraph "Entrada"
    Handler[Handler<br/>(Controllers)]
  end
  subgraph "Domínio"
    Hub[Hub<br/>(Service)]
    DTO[DTO]
  end
  subgraph "Infraestrutura"
    RedisClient[Redis<br/>Client/Adapter]
    Logger[Logger]
    Config[Configuração]
  end

  Handler --> Hub
  Handler --> DTO
  Hub --> RedisClient
  Hub --> DTO
  Handler --> RedisClient
  Handler --> Logger
  Handler --> Config
  RedisClient --> Logger
  Hub --> Logger
```

---

**Observações:**
- O Hub centraliza a lógica do domínio do chat.
- Redis atua como repositório e mecanismo de Pub/Sub.
- Handlers controlam entrada HTTP/WebSocket e interagem com as demais camadas.
- Logger e Configuração são utilizados em diversas camadas para logging e parametrização.
- Não há um "Repository" tradicional com acesso a banco relacional, pois Redis cumpre papel principal de armazenamento.

 # 🚀 API - Endpoints HTTP
## 📡 Endpoints Expostos pela Aplicação
**Lista Endpoints:**

<details>
<summary>Login do usuário (Geração de tokens JWT)</summary>

### Descrição
Permite que um usuário realize login informando nome de usuário e senha. Se as credenciais estiverem corretas, retorna um par de tokens JWT: `access_token` e `refresh_token`. O endpoint é utilizado para autenticação inicial dos usuários no sistema de chat.

### Entrada
- **Verbo HTTP:** POST
- **Caminho da Rota:** `/login`
- **Nome do Método Handler:** `LoginHandler`
- **Payload esperado:**
  - `user` (string, obrigatório): Nome do usuário.
  - `password` (string, obrigatório): Senha do usuário.
- **Exemplo de JSON de entrada:**
  ```json
  {
    "user": "bruno",
    "password": "1234"
  }
  ```

### Processamento
- **Validações:**
  - O campo `user` é recuperado do formulário.
  - O campo `password` é comparado com o valor fixo `"1234"`.
- **Recuperação de dados externos:** Nenhuma integração com banco de dados real neste ponto; apenas validação fixa.
- **Motivo:** Simulação de autenticação para fins de exemplo.
- **Geração de dados:**
  - Gera tokens `access_token` e `refresh_token` com permissão para salas `default` e `vip`.
- **Resposta:**
  - Status HTTP 200 em caso de sucesso.
  - Status HTTP 401 se as credenciais estiverem incorretas.
  - Resposta JSON:
    ```json
    {
      "access_token": "<jwt_access_token>",
      "refresh_token": "<jwt_refresh_token>"
    }
    ```

</details>

<details>
<summary>Refresh de token JWT (Gerar novo access_token)</summary>

### Descrição
Permite que um usuário troque um `refresh_token` válido por um novo `access_token` JWT. Usado para manter sessões ativas sem obrigar o usuário a realizar login novamente.

### Entrada
- **Verbo HTTP:** POST
- **Caminho da Rota:** `/refresh`
- **Nome do Método Handler:** `RefreshHandler`
- **Payload esperado:**
  - Header HTTP `Authorization` no formato: `Bearer <refresh_token>`
- **Exemplo de requisição:**
  ```
  POST /refresh
  Authorization: Bearer <refresh_token>
  ```

### Processamento
- **Validações:**
  - Verifica se o header `Authorization` contém um refresh token JWT válido.
- **Recuperação de dados externos:** Nenhuma integração com banco de dados real; validação apenas do JWT.
- **Motivo:** Troca segura de tokens JWT expirados por válidos.
- **Geração de dados:**
  - Gera novo `access_token` JWT para o usuário autorizado.
- **Resposta:**
  - Status HTTP 200 em caso de sucesso.
  - Status HTTP 401 se o refresh token for inválido.
  - Resposta JSON:
    ```json
    {
      "access_token": "<novo_access_token>"
    }
    ```

</details>

<details>
<summary>Conexão WebSocket autenticada (Entrada no chat)</summary>

### Descrição
Permite que o cliente estabeleça uma conexão WebSocket autenticada para participar do chat em tempo real. O acesso depende do envio de um JWT válido e da permissão do usuário para a sala requisitada.

### Entrada
- **Verbo HTTP:** GET (upgrade para WebSocket)
- **Caminho da Rota:** `/ws`
- **Nome do Método Handler:** `HandleConnections`
- **Payload esperado:**
  - Header HTTP `Authorization` no formato: `Bearer <access_token>`
  - Query parameters:
    - `room` (string, opcional): Nome da sala. Padrão `"default"` se não informado.
    - `user` (string, obrigatório): Nome do usuário.
- **Exemplo de requisição:**
  ```
  GET ws://localhost:8000/ws?room=default&user=bruno
  Authorization: Bearer <access_token>
  ```

### Processamento
- **Validações:**
  - Verifica existência e validade do token JWT enviado no header `Authorization`.
  - Verifica se o usuário tem permissão para acessar a sala especificada no token JWT.
- **Recuperação de dados externos:**
  - Recupera histórico das últimas mensagens da sala via Redis (`store.GetMessages`).
- **Motivo da recuperação:** Disponibilizar contexto/conversa anterior ao usuário recém-conectado.
- **Geração de dados:**
  - Envia mensagem de boas-vindas ao cliente.
  - Inicia o ciclo de leitura/escrita WebSocket com envio/recebimento de mensagens em tempo real.
- **Resposta:**
  - Conexão WebSocket estabelecida e mensagens trocadas em tempo real.
  - Se não autorizado, resposta HTTP de erro (`401 Unauthorized` ou `403 Forbidden`) e conexão fechada.

</details>

 ## 📡 cURL dos Endpoints
**Lista de endpoints:**
<details>
<summary>Login do usuário (gera access_token e refresh_token)</summary>

- **Endpoint:** [POST] /login
- **Base URL:** http://localhost:8000
- **Segurança:** Nenhuma
- **Body (application/x-www-form-urlencoded):**
  ```text
  user=bruno&password=1234
  ```
- **cURL:**
  ```code  copy
  curl -X POST "http://localhost:8000/login" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "user=bruno&password=1234"
  ```
</details>

<details>
<summary>Refresh do access_token usando refresh_token</summary>

- **Endpoint:** [POST] /refresh
- **Base URL:** http://localhost:8000
- **Segurança:** Bearer token (refresh_token)
- **Headers:** Authorization: Bearer &lt;REFRESH_TOKEN&gt;
- **cURL:**
  ```code  copy
  curl -X POST "http://localhost:8000/refresh" \
    -H "Authorization: Bearer <REFRESH_TOKEN>" \
    -H "Accept: application/json"
  ```
</details>

<details>
<summary>Conectar ao chat via WebSocket autenticado (por sala)</summary>

- **Endpoint:** [GET] /ws?room={{room}}&user={{user}}
- **Base URL:** ws://localhost:8000
- **Segurança:** Bearer token (access_token)
- **Headers:** Authorization: Bearer &lt;ACCESS_TOKEN&gt;
- **Query params:** 
  - room: string (opcional, padrão "default")
  - user: string (exemplo: bruno)
- **cURL:**
  ```code  copy
  curl -i -N -H "Authorization: Bearer <ACCESS_TOKEN>" \
    "ws://localhost:8000/ws?room=default&user=bruno"
  ```
  > Observação: O cURL nativo não suporta WebSocket; utilize ferramentas como [wscat](https://github.com/websockets/wscat) para testes reais:
  ```code  copy
  wscat -c "ws://localhost:8000/ws?room=default&user=bruno" \
    -H "Authorization: Bearer <ACCESS_TOKEN>"
  ```
</details>

---

**Observação:**  
- O endpoint `/login` espera dados via `x-www-form-urlencoded`, retornando `access_token` e `refresh_token` no corpo JSON.
- O endpoint `/refresh` deve ser chamado com o `refresh_token` no header de autorização para obter um novo `access_token`.
- O endpoint `/ws` é uma conexão WebSocket protegida por JWT, usada para chat em tempo real. Use ferramentas próprias para WebSocket, pois o cURL padrão não suporta esse protocolo.

Se precisar de exemplos de payloads para envio de mensagens via WebSocket, ou de resposta dos endpoints, solicite!

 ## 📟 Endpoints Consumidos pela Aplicação

**Resumo:**  
Após análise dos arquivos do repositório e do fluxo de execução da aplicação, **não foram identificadas chamadas a endpoints HTTP externos (consumo de APIs via cliente HTTP) no código-fonte Go**. Toda a comunicação da aplicação ocorre por meio de WebSocket, Redis Pub/Sub e manipulação interna entre handlers, controllers e serviços próprios.

### Detalhamento do Processo de Verificação

- Foram inspecionadas as principais funções de handlers, services, integrações Redis, WebSocket e autenticação.
- Não há uso de bibliotecas/clients HTTP (como net/http para requisições externas, http.Client, ou clients de terceiros para consumo de APIs REST).
- Não foram encontradas chamadas do tipo `http.Get`, `http.Post`, `client.Do`, nem uso de variáveis de ambiente ou arquivos de configuração que apontem para URLs de APIs externas consumidas como cliente.
- A única comunicação HTTP identificada no código refere-se à **exposição de endpoints próprios** (login, refresh, websocket), não ao consumo de APIs externas.

---

### Resultado

**Nenhum endpoint externo consumido foi identificado nesta aplicação.**

> Caso no futuro sejam implementadas integrações com APIs externas (exemplo: validação de tokens via serviço externo, integração com sistemas terceiros), este documento deve ser atualizado conforme o padrão estabelecido acima.

 # ✉️ Comunicação Assíncrona (Mensageria)
A aplicação interage com sistemas de mensageria para comunicação desacoplada entre instâncias e componentes, utilizando Redis como mecanismo principal.

## 👂 Consumers
<details>
<summary>RedisRoomSubscriber (Assinante de Sala via Redis Pub/Sub)</summary>

- **Nome Consumidor:** RedisRoomSubscriber
- **Fila/Tópico:** chat:*
- **Tipo de evento esperado:** dto.Message (estrutura da aplicação para mensagens de chat)
- **Implementação:**  
  - Função: `SubscribeAllRooms` (internal/redis/client.go)
  - Detalhe: Escuta todos os canais Redis que seguem o padrão `chat:*`. Ao receber uma mensagem, desserializa o payload (espera um objeto do tipo dto.Message) e encaminha para o Hub WebSocket para broadcast assíncrono aos clientes conectados.
</details>

## 📣 Producers
<details>
<summary>RedisPublisher (Publicador de Mensagens no Chat)</summary>

- **Nome Produtor:** RedisPublisher
- **Fila/Tópico:** chat:{roomID}
- **Mensagem publicada:** dto.Message (estrutura da aplicação para mensagens de chat)
- **Implementação:**  
  - Função: `PublishMessage` (internal/redis/publisher.go)
  - Detalhe: Publica mensagens no canal Redis correspondente à sala (`chat:{roomID}`), permitindo que múltiplas instâncias do servidor recebam e propaguem as mensagens em tempo real.
</details>

<details>
<summary>RedisUnreadPublisher (Gerenciador de Mensagens Não Lidas)</summary>

- **Nome Produtor:** RedisUnreadPublisher
- **Fila/Tópico:** unread:{user}
- **Mensagem publicada:** dto.Message (mensagens privadas não lidas para um usuário)
- **Implementação:**  
  - Função: `SaveUnread` (internal/redis/client.go)
  - Detalhe: Adiciona mensagens privadas enviadas a um usuário à lista Redis específica para mensagens não lidas daquele usuário, permitindo a recuperação posterior quando o usuário se reconectar.
</details>

---

**Observações adicionais:**
- Não foram identificadas outras bibliotecas de mensageria (Kafka, RabbitMQ, SQS, etc.) além do Redis Pub/Sub.
- O fluxo assíncrono está centrado no Pub/Sub do Redis, integrando múltiplas instâncias do servidor Go e garantindo entrega de mensagens em tempo real para clientes conectados via WebSocket.
- O payload das mensagens trafegadas segue a estrutura `dto.Message` definida internamente pela aplicação.

 # 🎲 Modelo de Dados da Aplicação

## 🗄️ Banco de Dados: **Redis (NoSQL, armazenamento em memória com suporte a listas e chaves-valor)**

---

### 📁 Estruturas de Dados Armazenadas no Redis

#### 1. **Lista de Mensagens por Sala**
- **Chave:** `chat:[roomID]` (exemplo: `chat:default`)
- **Tipo:** Lista (`LPUSH`/`LTRIM`/`LRange`)
- **Descrição:** Armazena as mensagens mais recentes de cada sala de chat, com limite configurável (exemplo: 50 mensagens por sala).

| Campo      | Tipo de Dado   | Atributos         | Observações                          |
|------------|----------------|-------------------|--------------------------------------|
| User       | string         |                   | Usuário que enviou a mensagem        |
| Content    | string         |                   | Conteúdo da mensagem                 |
| Timestamp  | datetime       |                   | Data/hora de envio (RFC3339)         |
| RoomID     | string         |                   | Sala associada à mensagem            |
| Target     | string/null    |                   | Usuário-alvo para mensagens privadas |

**Notas:**
- As mensagens são armazenadas serializadas em JSON.
- Cada elemento da lista representa uma instância da estrutura acima.

---

#### 2. **Lista de Mensagens Não Lidas por Usuário**
- **Chave:** `unread:[user]` (exemplo: `unread:bruno`)
- **Tipo:** Lista (`LPUSH`/`LRange`/`Del`)
- **Descrição:** Armazena mensagens privadas recebidas por um usuário que ainda não foram lidas.

| Campo      | Tipo de Dado   | Atributos         | Observações                          |
|------------|----------------|-------------------|--------------------------------------|
| User       | string         |                   | Usuário remetente                    |
| Content    | string         |                   | Conteúdo da mensagem                 |
| Timestamp  | datetime       |                   | Data/hora de envio                   |
| RoomID     | string         |                   | Sala de origem                       |
| Target     | string         |                   | Usuário-alvo (deve ser igual ao dono da chave) |

**Notas:**
- Mensagens são serializadas em JSON.
- Expiração automática configurada para 24 horas nas listas de não lidas.

---

#### 3. **Tokens de Autenticação (JWT)**
- **Armazenamento:** Não persistido no Redis, gerado e validado em memória.
- **Campos relevantes (JWT Claims):**
  - User: string
  - Rooms: lista de strings (salas autorizadas para o usuário)
  - Expiração e demais claims padrão do JWT

---

## 🌳 Estrutura Hierárquica das Entidades

```
Redis
├── chat:[roomID]           # Lista de mensagens por sala
│   └── [
│        {
│          User: string,
│          Content: string,
│          Timestamp: datetime,
│          RoomID: string,
│          Target: string/null
│        }, ...
│      ]
├── unread:[user]           # Lista de mensagens privadas não lidas por usuário
│   └── [
│        {
│          User: string,
│          Content: string,
│          Timestamp: datetime,
│          RoomID: string,
│          Target: string
│        }, ...
│      ]
```

---

### Exemplos Textuais de Relacionamento

**Relacionamento entre Salas e Mensagens:**
- Cada "sala" é identificada por um RoomID utilizado como parte da chave (`chat:[roomID]`) no Redis. 
- A lista correspondente armazena as mensagens enviadas para a sala.
- Não há chaves primárias ou estrangeiras explícitas em Redis, mas o campo `RoomID` nas mensagens serve como referência à sala.

**Relacionamento entre Usuários e Mensagens Privadas Não Lidas:**
- Cada usuário possui uma lista `unread:[user]`.
- As mensagens com o campo `Target` igual ao nome do usuário são adicionadas nesta lista.
- Não existe relacionamento relacional formal, mas a correspondência é feita pela chave e pelo campo `Target`.

---

## ℹ️ Observações

- **Não há tabelas relacionais** — Toda a persistência é feita via listas e chaves nomeadas no Redis, sem esquemas fixos.
- **Não foram identificados índices ou constraints** além dos nomes das chaves (padrão em bancos NoSQL).
- **Tokens JWT** são usados para controle de acesso, mas não são armazenados no banco.
- **Relacionamentos** são implícitos por chave e campos dos objetos serializados.

---

Caso precise detalhamento dos campos das estruturas DTO utilizadas, solicite explicitamente.

 # 🚨 Estratégia de Testes

A análise do repositório da aplicação **chat-websocket** indica as seguintes informações sobre a estratégia de testes:

- **Testes unitários**: A aplicação demonstra intenção de realizar testes unitários, principalmente para lógicas isoladas nas camadas de autenticação, WebSocket e integração com Redis.
- **Testes de integração**: Há menção ao uso de testes de integração, sugerindo a validação do funcionamento conjunto entre módulos (ex: interação entre WebSocket, Redis e autenticação JWT).
- **Testes end-to-end (E2E)**: A documentação menciona planos para implementação de testes E2E utilizando Docker Compose para simular cenários completos com múltiplos serviços.
- **Cobertura em andamento**: O README indica que os testes unitários e integrados ainda estão "em andamento".

## Frameworks Utilizados

- **Testify** (github.com/stretchr/testify): Framework Go amplamente utilizado para assertions e mocks em testes unitários e de integração.
- **Ferramentas nativas do Go**: O ecossistema Go utiliza por padrão o comando `go test` e arquivos com sufixo `_test.go` para definição de casos de teste.

## Estrutura dos Testes

- Não há pastas dedicadas explicitamente à testes (como `test/` ou `__tests__/`).
- Espera-se que os testes estejam distribuídos nos próprios pacotes, em arquivos com o padrão Go (`*_test.go`), por exemplo:
  - `internal/auth/auth_test.go`
  - `internal/handler/auth_test.go`
  - `internal/websocket/hub_test.go`
  - `internal/redis/client_test.go`
- A documentação do projeto cita:
  - "Testes unitários - Em andamento"
  - "Criar testes E2E completos via Docker Compose"

## Executando os Testes

Para rodar os testes unitários e de integração da aplicação Go, utiliza-se geralmente:

```bash
go test ./...
```

Para execução de testes end-to-end (E2E) utilizando Docker Compose (conforme sugerido no README):

```bash
docker compose up --build
# (os testes E2E completos ainda estão em desenvolvimento)
```

## Relatórios de Cobertura

Nenhum relatório de cobertura foi identificado no repositório da aplicação.

## Lacunas Identificadas

- Não foram encontrados arquivos de teste implementados (`*_test.go`) no repositório até o momento.
- Os testes automatizados ainda estão "em andamento", conforme mencionado na documentação.
- Ausência de relatórios de cobertura ou integração com ferramentas como `cover`, `Codecov`, etc.
- Não há exemplos concretos de comandos de execução de testes E2E automatizados.
- Não foram detectados scripts dedicados para rodar testes via Makefile ou arquivos auxiliares.

> [!NOTE]
> A aplicação segue as práticas do ecossistema Go, mas carece da implementação efetiva dos testes automatizados. Recomenda-se priorizar a criação dos arquivos de teste (`*_test.go`) para módulos críticos e a inclusão de relatórios de cobertura para acompanhamento da qualidade do código.

---

 # 🔎 Observabilidade

A aplicação implementa os seguintes mecanismos de observabilidade:

## Logs
- Ferramenta(s) utilizada(s): **Uber Zap** (`go.uber.org/zap`).
- Formato: **Estruturado** (os logs são emitidos em formato estruturado, geralmente JSON, padrão do Zap).
- Integração com sistemas externos: **Não identificado no código** integração direta com stacks como ELK, mas o formato estruturado facilita ingestão por ferramentas externas.
- Configurações de nível de log: **Não explicitamente detalhadas** nos trechos fornecidos, mas o Zap suporta níveis como INFO, DEBUG, WARN, ERROR – e há uso explícito de `.Info`, `.Error`, `.Warn`, `.Debug` no código.
- Exemplos de campos de log (observados nos usos):
  - Mensagens de fluxo ("Iniciando subscriber genérico Redis para todas as salas", "🚀 Servidor iniciado", "❌ Servidor com problema")
  - Contexto de erro (`zap.Error(err)`)
  - Metadados como nomes de usuários (`zap.String("user", c.User)`), nomes de sala (`zap.String("room", roomID)`), canais do Redis, payload das mensagens.

## Métricas
- Ferramenta(s) utilizada(s): **Não detectado** uso de bibliotecas/frameworks de métricas (ex: Prometheus, Micrometer, OpenTelemetry Metrics) no código analisado.
- Endpoint de exposição: **Não identificado**.
- Exemplos de métricas detectadas: **Nenhuma métrica customizada ou padrão identificada**.
- Integração com sistemas externos: **Não identificado**.

## Tracing
- Ferramenta(s) utilizada(s): **Não detectado** uso de tracing distribuído (ex: OpenTelemetry, Zipkin, Jaeger) no código analisado.
- Integração com sistemas externos: **Não identificado**.
- Configuração de amostragem: **Não aplicável**.
- Exemplos de integração: **Não identificado** middleware ou interceptors para tracing.

---

**Resumo objetivo:**
- A aplicação utiliza logs estruturados via Uber Zap em diversos pontos críticos (inicialização, erros, eventos no Redis, fluxo de mensagens).
- Não foram identificados mecanismos implementados para métricas ou tracing distribuído nos arquivos analisados.
- Não há integração explícita com sistemas externos de observabilidade além do suporte a ingestão por logs estruturados.

 # 🚔 Segurança

A análise da aplicação revela a seguinte estratégia de segurança, baseada exclusivamente nas evidências presentes no código-fonte e configurações fornecidas:

## Autenticação

- A aplicação utiliza **JWT (JSON Web Token)** para autenticação de requisições.
  - O token de acesso é esperado no header HTTP `Authorization` no formato `Bearer <token>`.
  - O fluxo de login ocorre via o endpoint HTTP, no handler `LoginHandler`, que valida o usuário e a senha (apenas um caso fixo: senha igual a `"1234"`), e emite tokens de acesso e refresh usando funções do módulo `auth`.
  - O endpoint de refresh (`RefreshHandler`) valida o token de refresh (também via `Authorization: Bearer <token>`) e emite um novo token de acesso.
  - No estabelecimento de conexões WebSocket (função `HandleConnections`), o token JWT é extraído do header HTTP e validado por `auth.ValidateAccessToken`.

## Autorização

- Após validação do JWT, são avaliados os "claims" do token para determinar as permissões de acesso às salas de chat:
  - Apenas usuários cujos claims incluem a sala desejada (`claims.Rooms`) conseguem conectar-se à respectiva sala via WebSocket.
  - Caso o usuário não possua permissão para a sala, a conexão é imediatamente encerrada com status HTTP 403 (`forbidden`).
- Não há anotações, decorators ou middlewares genéricos identificados para autorização; toda a verificação ocorre explicitamente no tratamento das conexões WebSocket e nos handlers HTTP.

## Configurações adicionais

- **CORS:** Não foi identificado no código analisado nenhum mecanismo explícito de configuração ou restrição de CORS.
- **CSRF:** Nenhum mecanismo ou proteção contra CSRF foi identificado, nem menção à sua configuração explícita.
- **Rate Limiting:** Não há qualquer implementação ou configuração visível de rate limiting nos endpoints ou conexões WebSocket.
- **Validação de entrada de dados:** A validação dos dados enviados parece ser limitada à verificação básica de presença do token e permissões, sem validações estruturais avançadas ou sanitização detalhada dos payloads.

---

> [!WARNING]
> - Não foi identificado mecanismo explícito de CORS, CSRF ou rate limiting na aplicação.
> - O controle de autenticação e autorização é realizado exclusivamente via handlers customizados e validação manual dos claims do JWT.
> - Não há uso de middlewares, decorators ou anotações padronizadas para segurança.
> - A senha do usuário está fixa no código (exemplo), o que não é seguro em ambiente real.
> - Recomenda-se análise adicional para identificar eventuais pontos complementares não presentes nesse trecho do código.