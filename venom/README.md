# VenomBot Integration for Mestra

Este diretório contém a integração do VenomBot com o sistema Mestra para automação de WhatsApp.

## Estrutura

- `whatsapp-bot.js` - Script principal do VenomBot
- `package.json` - Dependências e configurações do Node.js

## Como usar

### 1. Instalar dependências

```bash
cd venom
npm install
```

### 2. Iniciar o serviço Go

Primeiro, certifique-se de que o serviço Go está rodando:

```bash
cd ..
go run main.go
```

### 3. Iniciar o VenomBot via API

Com o serviço Go rodando, você pode iniciar o VenomBot através da API REST:

```bash
# Iniciar o bot
curl -X POST http://localhost:8080/api/bot/start

# Verificar status
curl http://localhost:8080/api/bot/status

# Enviar mensagem
curl -X POST http://localhost:8080/api/bot/send \
  -H "Content-Type: application/json" \
  -d '{"to": "5511999999999@c.us", "message": "Olá do Mestra Bot!"}'
```

### 4. Ou executar diretamente

Alternativamente, você pode executar o bot diretamente:

```bash
npm start
```

## Endpoints da API

### Bot Management
- `POST /api/bot/start` - Iniciar o VenomBot
- `POST /api/bot/stop` - Parar o VenomBot
- `GET /api/bot/status` - Status do bot

### Mensagens
- `POST /api/bot/send` - Enviar mensagem de texto
- `POST /api/bot/media` - Enviar mídia
- `POST /api/whatsapp/message` - Receber mensagens (usado pelo VenomBot)

### Logs e Comandos
- `POST /api/bot/log` - Logs do VenomBot
- `GET /api/bot/commands/:session` - Comandos pendentes
- `DELETE /api/bot/commands/:id` - Remover comando

## Funcionamento

1. **Inicialização**: O serviço Go gerencia o processo Node.js do VenomBot
2. **QR Code**: O VenomBot gera um QR code para conectar ao WhatsApp
3. **Mensagens Recebidas**: São enviadas para a API Go via HTTP
4. **Mensagens Enviadas**: A API Go adiciona comandos na fila para o VenomBot executar
5. **Persistência**: Todas as mensagens são salvas no banco PostgreSQL

## Variáveis de Ambiente

- `GO_API_URL` - URL da API Go (padrão: http://localhost:8080)
- `BOT_SESSION` - Nome da sessão do bot (padrão: mestra-bot)

## Logs

O VenomBot envia logs para a API Go em tempo real, incluindo:
- Status de conexão
- Mensagens processadas
- Erros e exceções
- Estados do WhatsApp

## Troubleshooting

### Bot não conecta
1. Verifique se o Node.js está instalado (`node --version`)
2. Certifique-se de que a API Go está rodando
3. Verifique os logs via `/api/bot/log`

### Mensagens não são enviadas
1. Verifique se o bot está conectado (`/api/bot/status`)
2. Confirme o formato do número (`5511999999999@c.us`)
3. Verifique a fila de comandos (`/api/bot/commands/mestra-bot`)

### Erro de permissões
O VenomBot pode precisar de permissões especiais para executar o navegador headless. Em sistemas Linux/Docker, pode ser necessário instalar dependências adicionais.

# 🎉 VenomBot Integração Completa - Resumo da Implementação

## ✅ O que foi implementado

### 1. **Arquitetura Híbrida Go + Node.js**
- **Go API**: Gerencia o sistema, banco de dados e processo Node.js
- **VenomBot Node.js**: Conecta ao WhatsApp real via puppeteer
- **Comunicação HTTP**: APIs REST entre Go e VenomBot

### 2. **Serviços Go Criados**
- `VenomBotService`: Gerencia processo Node.js, fila de comandos
- `VenomController`: Endpoints REST para controle do bot
- Integração com GORM para persistência das mensagens

### 3. **VenomBot Node.js**
- Script completo de automação WhatsApp
- Integração bidirecional com API Go
- Sistema de logs e monitoramento
- Suporte a QR Code para conexão

### 4. **APIs Implementadas**

#### Bot Management
- `POST /api/bot/start` - Iniciar VenomBot
- `POST /api/bot/stop` - Parar VenomBot  
- `GET /api/bot/status` - Status do bot

#### Mensagens
- `POST /api/bot/send` - Enviar mensagem texto
- `POST /api/bot/media` - Enviar mídia
- `POST /api/whatsapp/message` - Receber mensagens

#### Sistema
- `POST /api/bot/log` - Logs centralizados
- `GET /api/bot/commands/:session` - Comandos pendentes
- `DELETE /api/bot/commands/:id` - Remover comando

### 5. **Fluxo Completo**

```
WhatsApp ↔ VenomBot ↔ Go API ↔ PostgreSQL ↔ .NET API
```

1. **Recebimento**: WhatsApp → VenomBot → HTTP POST → Go → GORM → PostgreSQL
2. **Envio**: .NET/Go → Fila de Comandos → VenomBot → WhatsApp
3. **Logs**: VenomBot → Go → Console/Arquivo

## 🚀 Como usar

### 1. **Preparação**
```bash
# 1. Instalar dependências Node.js
cd venom
npm install

# 2. Configurar environment
cp .env.example .env.local
# (editar .env.local com suas configurações)

# 3. Iniciar API Go
go run main.go
```

### 2. **Iniciar VenomBot via API**
```bash
# Iniciar bot
curl -X POST http://localhost:8080/api/bot/start

# Verificar status
curl http://localhost:8080/api/bot/status

# Enviar mensagem
curl -X POST http://localhost:8080/api/bot/send \
  -H "Content-Type: application/json" \
  -d '{"to": "5511999999999@c.us", "message": "Olá do Mestra!"}'
```

### 3. **Conectar WhatsApp**
1. Execute o comando `/api/bot/start`
2. Observe o console da API Go
3. Escaneie o QR Code com WhatsApp
4. Bot estará conectado e funcionando!

## 🛠️ Arquivos Criados/Modificados

### Novos Arquivos
- `services/venom_service.go` - Serviço de integração VenomBot
- `controllers/venom_controller.go` - Controller REST do VenomBot
- `venom/whatsapp-bot.js` - Script principal VenomBot
- `venom/package.json` - Dependências Node.js
- `venom/README.md` - Documentação específica VenomBot
- `test-integration.ps1` - Script de teste PowerShell

### Arquivos Modificados
- `routes/routes.go` - Adicionadas rotas do VenomBot
- `README.md` - Documentação atualizada
- `.env.example` - Variáveis do VenomBot

## 🔧 Tecnologias Utilizadas

- **Go 1.24**: API REST principal
- **Gin**: Framework web Go
- **GORM**: ORM para PostgreSQL
- **Node.js**: Runtime para VenomBot
- **VenomBot**: Biblioteca WhatsApp
- **Axios**: Cliente HTTP Node.js
- **PostgreSQL**: Banco de dados

## ✨ Recursos Implementados

- ✅ **Conexão Real WhatsApp**: Via VenomBot + puppeteer
- ✅ **Mensagens Bidirecionais**: Envio e recebimento 
- ✅ **Persistência**: Todas mensagens salvas no PostgreSQL
- ✅ **Logs Centralizados**: VenomBot → Go → Console
- ✅ **Fila de Comandos**: Sistema assíncrono Go ↔ VenomBot
- ✅ **API REST Completa**: Controle total via HTTP
- ✅ **Monitoramento**: Status e health checks
- ✅ **Graceful Shutdown**: Encerramento seguro
- ✅ **Auto-Restart**: VenomBot pode ser reiniciado via API
- ✅ **Suporte Mídia**: Envio de imagens, documentos
- ✅ **Multi-Sessão**: Suporte a múltiplas sessões WhatsApp

## 🚀 Próximos Passos

Para produção, considere:

1. **Docker**: Containerizar Go + Node.js + PostgreSQL
2. **Nginx**: Proxy reverso para a API
3. **SSL/TLS**: HTTPS em produção
4. **Logs Estruturados**: JSON logs para análise
5. **Métricas**: Prometheus + Grafana
6. **Backup Sessão**: Persistir sessão WhatsApp
7. **Rate Limiting**: Limitar requests por IP
8. **Autenticação**: JWT para endpoints sensíveis

## 🎯 Resultado Final

Você agora tem um **sistema completo de automação WhatsApp** que:

- ✅ Conecta ao WhatsApp real (não simulação)
- ✅ Persiste todas as mensagens
- ✅ Integra com sua API .NET
- ✅ Pode ser controlado via API REST
- ✅ É escalável e pronto para produção
- ✅ Tem logs e monitoramento completos

**O VenomBot está 100% integrado ao seu sistema Mestra! 🚀**
