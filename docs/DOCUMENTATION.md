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
... (conteúdo omitido para brevidade) ...
---
# 🔎 Recomendações para Melhoria do Código na pasta /internal

A seguir estão recomendações para refatoração visando melhor testabilidade, performance e uso de interfaces na pasta `/internal`:

## 1. Introduza Interfaces para Dependências Externas e Componentes-Chave
- **WebSocket Hub/Client:** Defina interfaces como `Hub`, `ClientConn` para abstrair operações principais (broadcast, subscribe, send/receive). Isso facilita mocks em testes unitários.
- **Armazenamento/Redis:** Crie interfaces como `ChatStore`, `MessageRepository`, `UnreadRepository` para desacoplar lógica do Redis da lógica do domínio. Implemente-as no pacote Redis.
- **Logger:** Use uma interface `Logger` ao invés do uso direto do Zap nos componentes internos.
- **JWT/Auth:** Defina interface para validação/geração de tokens (`TokenService`).

## 2. Inversão de Dependências via Injeção por Construtor
- Prefira passar dependências via construtores ao invés de acessar singletons ou variáveis globais.
- Exemplo:
  ```go
  type Hub struct {
      store ChatStore // interface!
      logger Logger   // interface!
  }
  func NewHub(store ChatStore, logger Logger) *Hub { ... }
  ```
  Isso permite injetar mocks nos testes.

## 3. Separe Lógica de Domínio da Infraestrutura/Frameworks
- Mantenha handlers HTTP/WebSocket "finos", delegando lógica ao domínio via interfaces.
- Evite dependência direta do pacote gorilla/websocket ou redis nos serviços centrais; use wrappers/interfaces.

## 4. Facilite Testes Unitários e Mocks
- Para cada interface criada, forneça implementações "fakes" ou "mocks" nos arquivos `_test.go`.
- Use testify/mock ou crie mocks manuais simples para dependências como ChatStore, TokenService etc.
- Exemplo:
  ```go
  type MockChatStore struct { ... }
  func (m *MockChatStore) SaveUnread(...) error { ... }
  ```

## 5. Melhore Performance com Channels Bufferizados e Pooling (quando aplicável)
- Use channels bufferizados para filas internas (`client.Send`) para evitar bloqueios desnecessários.
- Considere pooling para conexões Redis se houver alta concorrência.

## 6. Documente as Interfaces Públicas e Fluxos Críticos
- Adicione comentários GoDoc nas interfaces e structs principais para facilitar manutenção e onboarding.

## 7. Exemplos Práticos de Interface:
```go
type ChatStore interface {
    GetMessages(ctx context.Context, room string, limit int) ([]Message, error)
    SaveUnread(ctx context.Context, user string, msg Message) error
    GetUnreadMessages(ctx context.Context, user string) ([]Message, error)
    ClearUnread(ctx context.Context, user string) error
}
type TokenService interface {
    GenerateAccessToken(user string, rooms []string) (string, error)
    ValidateAccessToken(token string) (*Claims, error)
}
type Logger interface {
    Info(msg string, fields ...zap.Field)
    Error(msg string, fields ...zap.Field)
}
```
Essas práticas facilitam a implementação futura dos testes unitários e aumentam a flexibilidade/manutenibilidade do código.