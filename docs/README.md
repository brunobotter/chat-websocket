<p align="center" margin="20 0">
  <img src="https://raw.githubusercontent.com/brunobotter/chat-websocket/main/.github/logo.png" alt="logo do time" width="30%" style="max-width:100%;"/>
</p>

# chat-websocket
[![Status do Projeto](https://img.shields.io/badge/Status-Em%20desenvolvimento-brightgreen.svg)]()
[![Go](https://img.shields.io/badge/Go-1.24%2B-blue.svg)]()
[![Redis](https://img.shields.io/badge/Redis-7.x-purple.svg)]()
[![Licença](https://img.shields.io/badge/Licença-Proprietária-red.svg)]()

## Sumário
1. [**Descrição do Projeto**](#descrição-do-projeto)
2. [**Como Usar e Pré-requisitos**](#como-usar-e-pré-requisitos)
3. [**Estrutura do Repositório**](#estrutura-do-repositório)
4. [**Como Executar Localmente**](#como-executar-localmente)
5. [**Como Executar com Docker**](#como-executar-com-docker)
6. [**Testes**](#testes)
7. [**Como Contribuir**](#como-contribuir)
8. [**Equipe Responsável e Contato**](#equipe-responsável-e-contato)
9. [**Referências e Links Úteis**](#referências-e-links-úteis)
10. [**Licenciamento**](#licenciamento)

---

## Descrição do Projeto

### O que é?
Sistema de chat em tempo real desenvolvido em Go, com autenticação baseada em JWT, suporte a múltiplas salas e comunicação entre instâncias por meio de Redis Pub/Sub. Permite múltiplos servidores conectados, integrados via WebSocket, com autenticação simples e gerenciamento de mensagens. O balanceamento de carga e o acesso externo são realizados via Nginx, com ambiente pronto para deploy em Docker.

### Funcionalidades Principais
- Envio e recebimento de mensagens em tempo real via WebSocket.
- Autenticação JWT com refresh token.
- Suporte a múltiplas salas de chat.
- Armazenamento temporário de mensagens no Redis.
- Balanceamento de carga com Nginx.

### Arquitetura
O projeto segue os princípios de arquitetura definidos abaixo:
- **API**: Endpoints HTTP para login, refresh de token e conexão WebSocket autenticada.
- **Application**: Camada de handlers para autenticação e WebSocket, com integração a Redis e gerenciamento de sessões.
- **Domain**: Estruturas DTO para mensagens, usuários e tokens JWT.
- **Infrastructure**: Integração via Redis Pub/Sub, arquivos de configuração, scripts Docker, Nginx e logs centralizados.

---

## Como Usar e Pré-requisitos

### Pré-requisitos
Para utilizar e desenvolver neste projeto, você precisará de:

#### Software Necessário
- **Go 1.24+**
- **Redis 7.x**
- **Docker**
- **Docker Compose**
- IDE de sua preferência:
  - VSCode
  - GoLand
  - Vim

#### Acessos Necessários
Solicite ao administrador da infraestrutura:
- Acesso à instância Redis (local ou remota)
- Permissão para rodar containers Docker localmente

#### Credenciais de API
- Não há consumo de APIs externas neste projeto.
- Os tokens JWT são gerados pelo próprio backend no fluxo de autenticação.

---

## Estrutura do Repositório

```
cmd/
  server/
    main.go
    dev.yaml
internal/
  auth/
    authorzation.go
  config/
    config.go
    mapping.go
  dto/
    auth.go
    message.go
  handler/
    auth.go
    websocket.go
  logger/
    logger.go
  redis/
    client.go
    publisher.go
    subscriber.go
  router/
    router.go
  websocket/
    client.go
    hub.go
    publisher.go
docker-compose.yml
dockerfile
nginx.conf
README.md
go.mod
go.sum
```

---

## Como Executar Localmente

### Configuração Inicial
1. **Instale as dependências Go:**
   ```bash
   go mod download
   ```
2. **Configure as variáveis de ambiente:**
   ```bash
   export APP_REDIS_ADDR=localhost:6379
   export APP_SERVER_PORT=8080
   ```
3. **Suba o Redis local (se necessário):**
   ```bash
   docker run --name redis -p 6379:6379 redis:7-alpine
   ```

### Executando a Aplicação

```bash
go run ./cmd/server/main.go
```
A aplicação estará disponível em `http://localhost:8080`.

---

## Como Executar com Docker

### Usando Docker Compose

```bash
docker compose up --build
```
A aplicação estará disponível em `http://localhost:8000`.

---

## Testes

O projeto está em evolução para incluir testes automatizados.

### Executar Todos os Testes Unitários e Integrados
Para rodar todos os testes (unitários e integrados), execute:
```bash
go test ./...
```
### Executar Testes Unitários Específicos (exemplo)
Para rodar apenas um arquivo ou pacote específico:
```bash
go test ./internal/handler -run ^TestNomeDoTeste$
```
### Escrevendo Testes Unitários (exemplo básico)
Crie arquivos terminando com `_test.go` e funções iniciando com `Test`:
```go
func TestMinhaFuncao(t *testing.T) {
    resultado := MinhaFuncao()
    if resultado != "esperado" {
        t.Errorf("resultado inesperado: %v", resultado)
    }
}
```
*Testes unitários e integrados estão sendo implementados gradualmente.*

---

## Como Contribuir

Para contribuir com o projeto:
1. Faça um fork do repositório.
2. Crie uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`).
3. Faça commit das suas alterações.
4. Faça push para a branch (`git push origin feature/nova-funcionalidade`).
5. Abra um Pull Request.

### Diretrizes de Contribuição
- Siga os padrões de código Go.
- Adicione testes para novas funcionalidades sempre que possível.
- Mantenha a documentação atualizada.
- Certifique-se de que todos os testes passam antes de submeter o PR.

Contato: **bruno.botter@gmail.com**

---

## Equipe Responsável e Contato

### Squad Responsável

**chat-websocket-team**

### Contatos

- **E-mail da Equipe**: bruno.botter@gmail.com

### Suporte

Para dúvidas ou problemas:
1. Abra uma issue no repositório GitHub.
2. Entre em contato por e-mail.

---

## Referências e Links Úteis

### Documentação Técnica e Recursos

- [Go Documentation](https://golang.org/doc/)
- [Redis Documentation](https://redis.io/documentation)
- [Echo Framework](https://echo.labstack.com/)
- [Gorilla WebSocket](https://pkg.go.dev/github.com/gorilla/websocket)
- [Docker Documentation](https://docs.docker.com/)
- [JWT.io](https://jwt.io/)

---

## Licenciamento

Este projeto é de **uso exclusivamente interno** do time chat-websocket.  
Todos os direitos reservados.  
Licença: Proprietária - Uso interno apenas.

---

**Status do Projeto**: 🚧 Em desenvolvimento  
*Última atualização: 2024-06*