# ✉️ emailN

API em **Go (Golang)** para **criação e envio de campanhas de e-mail** — simples, performática e extensível 🚀

---

## 🧭 Visão geral

O **emailN** permite criar campanhas de e-mail (com lista de contatos e conteúdo), gerenciar o ciclo de vida (criar, obter, deletar, iniciar) e enviar e-mails via um **worker** que processa campanhas marcadas como iniciadas.

### 🧩 Arquitetura

* 📦 `cmd/api`: API HTTP (`chi`) — endpoints para gerenciar campanhas
* ⚙️ `cmd/worker`: worker que busca campanhas a enviar e dispara via SMTP
* 🧠 `internal/domain`: entidades e regras de negócio
* 🌐 `internal/endpoints`: handlers HTTP
* 🗄️ `internal/infra`: repositório (GORM) e adaptador SMTP
* 🧪 `internal/test`: mocks e testes unitários

---

## ⚙️ Requisitos

* 🐹 Go **1.20+**
* 🐘 Banco de dados (PostgreSQL recomendado)
* 📧 Servidor SMTP (ex: Gmail, SendGrid, Mailgun)
* 🧰 `make` (opcional)

---

## 🔑 Variáveis de ambiente

Copie o arquivo `.env.EXAMPLE` para `.env` e preencha:

```env
# Banco
DATABASE_DSN="postgres://user:pass@localhost:5432/emailn?sslmode=disable"

# SMTP
EMAIL_SMTP="smtp.exemplo.com"
EMAIL_PORT=587
EMAIL_USER="seu-email@example.com"
EMAIL_PASSWORD="sua-senha"

# App
PORT=3000

# Se usar Keycloak/JWT
KEYCLOAK_URL=""
JWT_PUBLIC_KEY=""
```

> ⚠️ **Importante:** nunca envie suas credenciais reais para o repositório.
> Use variáveis locais, secrets no CI/CD ou serviços como Vault / AWS Secrets Manager.

---

## 💻 Instalação local

```bash
# Clonar repositório
git clone https://github.com/seu-usuario/emailN.git
cd emailN

# Criar arquivo de configuração
cp .env.EXAMPLE .env
# Editar o arquivo com suas credenciais

# Rodar banco (Docker ou local)
# Executar API
cd cmd/api
go run main.go
```

---

## 🌐 Endpoints principais

| Método   | Endpoint                 | Descrição                  |
| -------- | ------------------------ | -------------------------- |
| `POST`   | `/campaigns`             | Cria uma nova campanha     |
| `GET`    | `/campaigns/{id}`        | Consulta campanha por ID   |
| `DELETE` | `/campaigns/delete/{id}` | Deleta campanha            |
| `PATCH`  | `/campaigns/start/{id}`  | Inicia o envio da campanha |

### 📨 Exemplo de corpo (POST `/campaigns`)

```json
{
  "name": "Minha Campanha",
  "content": "<h1>Olá</h1>",
  "contacts": [
    {"email": "a@a.com"},
    {"email": "b@b.com"}
  ]
}
```

> 🔒 As rotas sob `/campaigns` requerem autenticação (middleware `Auth`).

---

## ⚙️ Como funciona o envio

O **worker** (`cmd/worker`) busca campanhas com status `started` e utiliza
`internal/infra/database/mail.SendMail` para enviar os e-mails via **SMTP**
(atualmente com o pacote `gomail`).

### 🧠 Recomendações:

* Envio em **lotes (batches)** com **timeout** configurado
* Implementar **retries com backoff exponencial**
* Monitorar métricas e logs de envio

---

## 🧪 Testes

Rodar todos os testes unitários:

```bash
go test ./... -count=1
```

---

## ✅ Boas práticas antes do deploy

* 🔹 Remover bins (`*.exe`, `tmp/`)
* 🔹 Atualizar `.env.EXAMPLE` com todas as variáveis
* 🔹 Configurar logger (`zap`, `logrus`, etc.)
* 🔹 Usar Docker + CI/CD (ex: GitHub Actions)
* 🔹 Validar entradas (ex: `go-playground/validator`)
* 🔹 Implementar **rate-limiting / batch sending** para evitar bloqueios SMTP

---

## 🚀 Roadmap / Ideias futuras

* ☁️ Suporte a provedores externos (SendGrid, Mailgun) via API
* 🧩 Templates HTML com `html/template` + sanitização
* 📊 Dashboard para acompanhar status de envios
* 🔁 Retry + DLQ (Dead Letter Queue) para falhas
* 📎 Suporte a anexos

---

## 🤝 Contribuição

1. Fork 🍴
2. Crie uma branch `feature/<nome>`
3. Abra um **Pull Request** com descrição e testes
4. Aguardamos sua contribuição! 💪

---

## 📄 Licença

Escolha uma licença (ex: **MIT**) e adicione o arquivo `LICENSE` ao projeto.
