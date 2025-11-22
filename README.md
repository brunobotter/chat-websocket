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

### 1. Login
POST http://localhost:8000/login
Content-Type: application/json

```json
{
  "user": "bruno",
  "password": "1234"
}
```

### 2. Refresh token
POST http://localhost:8000/refresh
Authorization: Bearer <refresh_token>

### 3. Conecta ao chat da sala
GET ws://localhost:8000/ws?room=default&user=bruno
Authorization: Bearer <access_token>

```json
{
    "content": "Olá, mundo!"
}
```

---

## 🛣️ Próximos passos
- Testes unitários (em andamento)
- Integrar banco de dados real (usuários, permissões, histórico)
- Adicionar logs estruturados em todas as rotas
- Implementar middleware de autenticação JWT no Echo
- Criar testes E2E completos via Docker Compose
- Adicionar métricas e monitoramento (Prometheus + Grafana)
- Refatoração do código para melhor organização e manutenibilidade
