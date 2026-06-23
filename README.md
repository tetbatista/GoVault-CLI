# 🔐 GoVault — Gerenciador de Senhas pelo Terminal

> CLI seguro para guardar e gerenciar credenciais localmente, feito em Go.  
> Tudo fica no seu computador — nenhum dado vai para a nuvem.

---

## ✨ Funcionalidades

- Cadastro de usuário com **senha mestre** protegida por bcrypt
- **Login** com autenticação local
- **Adicionar** credenciais (site, usuário, senha)
- **Listar** todas as credenciais salvas
- **Buscar** por nome do site
- **Deletar** credenciais
- Senhas armazenadas com **criptografia AES-256**
- Banco de dados **SQLite local** (um arquivo `.db` na sua máquina)
- **Sessão com token JWT** — faz login uma vez e usa por X minutos

---

## 🚀 Como rodar (do zero)

### Pré-requisitos

- [Go 1.21+](https://go.dev/dl/) instalado

### Instalação

```bash
# 1. Clone o repositório
git clone https://github.com/tetbatista/govault.git
cd govault

# 2. Instale as dependências
go mod tidy

# 3. Rode o projeto
go run .
```

Pronto. Sem Docker, sem configuração, sem variáveis de ambiente.

---

## 📦 Ou baixe o binário direto

> Em breve na seção de [Releases](https://github.com/tetbatista/govault/releases)

```bash
# Linux/macOS — exemplo após download
chmod +x govault
./govault
```

---

## 🖥️ Como usar

### Primeiro acesso — crie sua conta

```bash
go run . register
```

Você vai criar um usuário e uma **senha mestre**. Guarde bem essa senha — ela protege tudo.

---

### Login

```bash
go run . login
```

Após o login, uma sessão é criada localmente com JWT válido por 30 minutos.

---

### Adicionar uma credencial

```bash
go run . add
```

O sistema vai pedir:
- Nome do site (ex: `github`)
- Usuário/email
- Senha

A senha é criptografada com **AES-256** antes de ser salva.

---

### Listar todas as credenciais

```bash
go run . list
```

---

### Buscar por site

```bash
go run . search github
```

---

### Deletar uma credencial

```bash
go run . delete github
```

---

### Logout

```bash
go run . logout
```

---

## 🔒 Segurança — como funciona por baixo

| Camada | Tecnologia | Detalhe |
|---|---|---|
| Senha mestre | `bcrypt` (custo 14) | Nunca armazenada em texto puro |
| Senhas salvas | `AES-256-GCM` | Chave derivada da senha mestre via PBKDF2 |
| Sessão local | `JWT` (HS256) | Expira em 30 min, salvo em arquivo local |
| Banco de dados | `SQLite` | Arquivo local `.govault.db` na home do usuário |

> **Importante:** a chave de criptografia AES é derivada da sua senha mestre usando PBKDF2 com 100.000 iterações e salt aleatório. Isso significa que **sem a senha mestre, as senhas salvas são irrecuperáveis**.

---

## 🗂️ Estrutura do projeto

```
govault/
├── main.go                  # Ponto de entrada, lê os comandos
├── go.mod
├── go.sum
│
├── cmd/                     # Um arquivo por comando
│   ├── register.go          # Cadastro de usuário
│   ├── login.go             # Login e geração de sessão
│   ├── logout.go            # Logout
│   ├── add.go               # Adicionar credencial
│   ├── list.go              # Listar credenciais
│   ├── search.go            # Buscar credencial
│   └── delete.go            # Deletar credencial
│
├── internal/                # Lógica interna (não exposta)
│   ├── auth/
│   │   ├── hash.go          # bcrypt — hash e verificação de senha
│   │   └── session.go       # JWT — criação e validação de sessão
│   │
│   ├── crypto/
│   │   └── aes.go           # AES-256-GCM — cifrar/decifrar senhas
│   │
│   └── db/
│       ├── database.go      # Conexão e inicialização do SQLite
│       ├── user.go          # Queries de usuário
│       └── credential.go    # Queries de credenciais
│
└── README.md
```

---

## 🧰 Dependências principais

| Pacote | Uso |
|---|---|
| `github.com/mattn/go-sqlite3` | Driver SQLite para Go |
| `golang.org/x/crypto` | bcrypt + PBKDF2 |
| `github.com/golang-jwt/jwt` | Sessão com JWT |

Todas instaladas automaticamente com `go mod tidy`.

---

## 🗺️ Roadmap (próximas melhorias)

- [ ] Gerador de senhas fortes embutido (`govault generate`)
- [ ] Export/import de credenciais em formato criptografado
- [ ] Timeout de sessão configurável
- [ ] Suporte a categorias (trabalho, pessoal, etc.)

---

## 🧑‍💻 Sobre o projeto

Projeto desenvolvido para aprender Go na prática, aplicando:

- Organização de packages (`cmd/`, `internal/`)
- Criptografia real com AES-256 e bcrypt
- Banco de dados local com SQLite
- Autenticação com JWT
- CLI com args e flags nativas do Go

---

## 📄 Licença

MIT — fique à vontade para usar, estudar e modificar.
