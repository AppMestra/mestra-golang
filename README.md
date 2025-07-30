# Mestra WhatsApp Bot - Golang

Bot em Go para integração com WhatsApp que se comunica com a API .NET.

## 🏗️ Estrutura do Projeto

```
mestra-golang/
├── main.go                   # Ponto de entrada
├── config/                   # Configurações (env, banco, etc)
│   ├── database.go
│   └── config.go
├── controllers/              # Handlers das rotas
│   └── message_controller.go
├── models/                   # Modelos das mensagens e usuários
│   └── message.go
├── routes/                   # Definição das rotas
│   └── routes.go
├── services/                 # Lógica de envio/recebimento WhatsApp
│   └── message_service.go
├── utils/                    # Funções auxiliares
│   └── helpers.go
├── uploads/                  # Arquivos recebidos
├── .env.example              # Exemplo de configuração
└── go.mod
```

## ⚙️ Configuração

### 1. Configuração do Ambiente

Copie o arquivo de exemplo e configure suas variáveis:

```bash
cp .env.example .env.local
```

### 2. Configuração do Banco de Dados

No `.env.local`, configure as variáveis:

```bash
# Configurações do Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=seu_usuario
DB_PASSWORD=sua_senha
DB_NAME=mestra_db
DB_SSLMODE=disable
```

### 3. Configurações de APIs

Configure as URLs e chaves no `.env.local`:

```bash
# URL da API .NET
DOTNET_API_URL=http://localhost:5000/api/messages

# Configurações do WhatsApp
WHATSAPP_SESSION_NAME=mestra-bot
JWT_SECRET=sua_chave_jwt_super_secreta
```

## 🚀 Execução

### Local
```bash
go mod tidy
go run main.go
```

### Docker
```bash
docker-compose up mestra-bot
```

## 📡 API Endpoints

### POST /api/v1/send
Envia uma mensagem via WhatsApp
```json
{
  "to": "+5511999999999",
  "message": "Olá! Esta é uma mensagem de teste.",
  "type": "text"
}
```

### POST /api/v1/webhook
Recebe webhooks do WhatsApp (configurado automaticamente)

### GET /api/v1/messages?limit=50
Busca mensagens armazenadas

### GET /health
Health check do serviço

## 🔌 Integração com WhatsApp

Este projeto está preparado para integração com bibliotecas como:
- VenomBot
- WhatsApp Web.js
- Baileys

Para implementar a integração real com WhatsApp, você precisará:

1. Escolher uma biblioteca (recomendo VenomBot para Go)
2. Configurar a sessão do WhatsApp
3. Implementar os handlers de envio e recebimento
4. Configurar os webhooks

## 🔄 Comunicação com API .NET

O serviço automaticamente encaminha mensagens recebidas para a API .NET configurada em `DOTNET_API_URL`.

## 🔒 Segurança

- **Nunca** commite arquivos `.env.local` ou similares
- Use variáveis de ambiente em produção
- Configure JWT com chaves fortes
- Use HTTPS em produção

## 🌍 Deploy

Para deploy em cloud:

1. Configure as variáveis de ambiente no seu provedor (AWS, GCP, Azure, etc.)
2. Use Docker para containerização
3. Configure load balancer e SSL
4. Configure monitoramento e logs

Exemplo para deploy no Heroku:
```bash
heroku config:set DB_HOST=seu_host_producao
heroku config:set DB_USER=seu_usuario_producao
# ... outras variáveis
```
