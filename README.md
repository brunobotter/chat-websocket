# 💬 Chat WebSocket em Go

Este projeto é um **sistema de chat em tempo real** desenvolvido em **Go**, com suporte a **autenticação JWT**, **Redis Pub/Sub** e **load balancing com Nginx**.  
A aplicação é escalável, permitindo múltiplas instâncias de servidor comunicando-se via Redis.

---

## 🚀 Tecnologias utilizadas

- **Go 1.24+**
- **Echo Framework** — roteamento HTTP rápido e minimalista
- **Redis** — Pub/Sub e persistência leve das mensagens
- **JWT (JSON Web Token)** — autenticação e autorização (bem simples)
- **WebSocket (gorilla/websocket)** — comunicação em tempo real
- **Nginx** — proxy reverso e balanceamento de carga
- **Docker & Docker Compose** — ambiente de desenvolvimento e deploy
- **Testify** — testes unitários e integrados (Em andamento)


## ⚙️ Como rodar o projeto

### 🐳 Rodando com Docker Compose

```bash
docker compose up --build
```

1. Login
POST http://localhost:8000/login
Content-Type: application/json

{
  "user": "bruno",
  "password": "1234"
}

2. Refresh token
POST http://localhost:8000/refresh
Authorization: Bearer <refresh_token>

3. Conecta ao chat da sala
GET ws://localhost:8000/ws?room=default&user=bruno
Authorization: Bearer <access_token>

{
    "content": "Olá, mundo!"
}

---

## 📁 Estrutura recomendada para facilitar testes e manutenibilidade

Para facilitar testes unitários, performance e manutenção, recomenda-se:

- Organizar a pasta `/internal` em subpacotes por domínio (ex: `internal/user`, `internal/chat`, `internal/auth`)
- Definir interfaces para dependências externas (ex: repositórios, cache, serviços de autenticação)
- Injetar dependências via construtores (dependency injection)
- Evitar lógica em handlers/controllers; delegar para serviços testáveis via interface
- Utilizar mocks para interfaces nos testes unitários (ex: com Testify/mock)
- Separar modelos de dados das regras de negócio (DTOs vs entidades)
- Adotar padrões como Repository, Service e UseCase para clareza e testabilidade
- Escrever testes unitários para cada serviço isoladamente, cobrindo casos de sucesso e erro
- Priorizar funções puras sempre que possível para facilitar o teste isolado

### Exemplo de interface para repositório:
```go
// internal/user/repository.go
type UserRepository interface {
    FindByUsername(ctx context.Context, username string) (*User, error)
    Save(ctx context.Context, user *User) error
}
```
### Exemplo de injeção de dependência:
```go
// internal/user/service.go
type UserService struct {
    repo UserRepository
}
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```
---
## Próximos passos
- [ ] Testes unitários - Em andamento  
- [ ] Integrar banco de dados real (usuários, permissões, histórico)
- [ ] Adicionar logs estruturados em todas as rotas
- [ ] Implementar middleware de autenticação JWT no Echo
- [ ] Criar testes E2E completos via Docker Compose
- [ ] Adicionar métricas e monitoramento (Prometheus + Grafana)
- [ ] Refatoração para uso extensivo de interfaces e injeção de dependências na pasta `/internal`
