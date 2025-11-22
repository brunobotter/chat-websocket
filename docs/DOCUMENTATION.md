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
...
não há pastas dedicadas explicitamente à testes (como `test/` ou `__tests__/`).
espera-se que os testes estejam distribuídos nos próprios pacotes, em arquivos com o padrão Go (`*_test.go`), por exemplo:
n  - `internal/auth/auth_test.go`
n  - `internal/handler/auth_test.go`
n  - `internal/websocket/hub_test.go`
n  - `internal/redis/client_test.go`
nA documentação do projeto cita:
n  - "Testes unitários - Em andamento"
n  - "Criar testes E2E completos via Docker Compose"
n## Executando os Testes
Para rodar os testes unitários e de integração da aplicação Go, utiliza-se geralmente:
n```bash
go test ./...
n```
nPara execução de testes end-to-end (E2E) utilizando Docker Compose (conforme sugerido no README):
n```bash
docker compose up --build # (os testes E2E completos ainda estão em desenvolvimento)
n```
n## Relatórios de Cobertura...
n> A aplicação segue as práticas do ecossistema Go, mas carece da implementação efetiva dos testes automatizados. Recomenda-se priorizar a criação dos arquivos de teste (`*_test.go`) para módulos críticos e a inclusão de relatórios de cobertura para acompanhamento da qualidade do código.