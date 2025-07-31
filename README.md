# 🚀 API Golang MESTRA

API minimalista em Golang para processar mensagens econômicas do WhatsApp e extrair informações financeiras.

## 📋 Funcionalidades

- ✅ Recebe mensagens via POST `/api/processar`
- ✅ Extrai valores, descrições e datas automaticamente
- ✅ Interpreta marcadores temporais (hoje, ontem, anteontem)
- ✅ Salva mensagens originais no SQL Server
- ✅ Retorna JSON estruturado com dados econômicos

## 🛠️ Tecnologias

- **Linguagem**: Go 1.21+
- **Framework**: Gin
- **Banco**: SQL Server
- **ORM**: GORM

## ⚡ Instalação Rápida

1. **Clone e configure**:
```bash
cd mestra-golang
cp .env.example .env
```

2. **Configure o banco no `.env`**:
```env
DB_HOST=localhost
DB_PORT=1433
DB_NAME=apiGo
DB_USER=sa
DB_PASSWORD=SuaSenhaAqui
PORT=8080
```

3. **Instale dependências**:
```bash
go mod tidy
```

4. **Execute**:
```bash
go run main.go
```

## 🔥 Exemplo de Uso

**POST** `http://localhost:8080/api/processar`

```json
{
  "mensagem": "[30/07/2025 12:04] Mãe: Ontem 31 00 panificadora\n[30/07/2025 12:06] Mãe: 7 00 pão de quijo\n[30/07/2025 12:05] Mãe: 28 de julho\n28 00 feijoada",
  "autor": "Mãe",
  "data": "30/07/2025"
}
```

**Resposta**:
```json
[
  {
    "valor": 31.00,
    "descricao": "panificadora",
    "data": "2025-07-29"
  },
  {
    "valor": 7.00,
    "descricao": "pão de quijo",
    "data": "2025-07-30"
  },
  {
    "valor": 28.00,
    "descricao": "feijoada",
    "data": "2025-07-28"
  }
]
```

## 📍 Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/api/processar` | Processa mensagem e extrai dados |
| GET | `/api/mensagens` | Lista mensagens salvas |
| GET | `/api/entradas` | Lista entradas econômicas |
| GET | `/ping` | Teste de saúde da API |

## 🎯 Integração com VenomBot

O VenomBot (Node.js) monitora grupos WhatsApp e envia mensagens para esta API:

```javascript
// Exemplo de integração no whatsapp-bot.js
const response = await fetch('http://localhost:8080/api/processar', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    mensagem: message.body,
    autor: message.author,
    data: new Date().toLocaleDateString('pt-BR')
  })
});
```

## 📊 Regras de Extração

| Marcador | Regra |
|----------|-------|
| "Hoje" | Usa data da mensagem |
| "Ontem" | Data da mensagem - 1 dia |
| "Anteontem" | Data da mensagem - 2 dias |
| "DD de mês" | Usa essa data como referência |
| Sem marcador | Fallback para data da mensagem |

## 🗃️ Estrutura do Banco

```sql
-- Tabelas criadas automaticamente
CREATE TABLE mensagem_whats_apps (
  id INT IDENTITY PRIMARY KEY,
  autor NVARCHAR(255),
  data DATETIME2,
  conteudo NTEXT,
  created_at DATETIME2
);

CREATE TABLE entradas (
  id INT IDENTITY PRIMARY KEY,
  valor FLOAT,
  descricao NVARCHAR(255),
  data DATETIME2,
  created_at DATETIME2
);
```

## 🚀 Deploy

Para produção, defina `GIN_MODE=release` no `.env` e compile:

```bash
go build -o api-mestra main.go
./api-mestra
```

---

**Desenvolvido para o projeto MESTRA** 🎯
