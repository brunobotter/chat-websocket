# 💬 ChatApp — Etapa 1: WebSocket Local em Go

Este projeto é o início de um **chat distribuído em tempo real**, feito em **Go**.  
Nesta primeira etapa, o servidor implementa um **broadcast local** via WebSocket —  
todas as mensagens recebidas são enviadas para **todos os clientes conectados**.

---

## 🚀 Tecnologias usadas

- **Go 1.24+**
- **Gorilla WebSocket**
- **HTTP nativo (net/http)**

## ⚙️ Como rodar o servidor

1. Instale as dependências:
   ```bash
   go mod tidy
   go run main.go

   O servidor será iniciado em:
   http://localhost:8080

🧪 Testando via navegador

Abra duas abas no navegador.

Acesse este código JavaScript no console (pressione F12 → aba Console):

const ws = new WebSocket("ws://localhost:8080/ws");

ws.onopen = () => console.log("✅ Conectado ao servidor WebSocket");

ws.onmessage = (event) => console.log("📩 Mensagem recebida:", event.data);

// Para enviar uma mensagem:
// ws.send("Oi, outra aba aqui!");

🧪 Testando via Insomnia (ou Postman)

Abra o Insomnia (ou outro cliente WebSocket).

Crie uma nova requisição WebSocket:

Método: WS

URL: ws://localhost:8080/ws

Clique em Connect.

Envie uma mensagem: