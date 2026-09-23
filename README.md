# SGA — Sistema de Gestão de Alocação

<p align="center">
  <strong>API REST para gerenciamento de alunos, turmas, salas e alocações acadêmicas</strong>
</p>

<p align="center">
  Desenvolvido em Go utilizando o framework Gin
</p>

---

## Sobre o projeto

O **SGA — Sistema de Gestão de Alocação** é uma API REST desenvolvida para auxiliar instituições de ensino no gerenciamento de:

- alunos;
- turmas;
- salas de aula e laboratórios;
- matrículas;
- capacidade das salas;
- horários;
- conflitos de agenda;
- alocação física das turmas.

O projeto foi desenvolvido em **Go**, utilizando o framework **Gin**, e implementa regras de negócio voltadas à organização acadêmica e ao controle de ocupação dos espaços físicos.

A API possui versionamento através da rota:

```text
/api/v1
```

---

## Objetivo

Instituições de ensino frequentemente precisam lidar com situações como:

- salas com capacidade insuficiente para determinadas turmas;
- alunos matriculados duas vezes na mesma turma;
- duas turmas utilizando a mesma sala no mesmo horário;
- alunos matriculados em turmas com horários conflitantes;
- dificuldade para consultar a ocupação das salas.

O SGA centraliza essas informações e aplica automaticamente regras de validação antes de permitir determinadas operações.

---

## Tecnologias utilizadas

| Tecnologia | Utilização |
|---|---|
| **Go** | Linguagem principal |
| **Gin** | Framework HTTP para construção da API REST |
| **JSON** | Formato utilizado na comunicação da API |
| **Git** | Controle de versão |
| **GitHub** | Hospedagem do código-fonte |
| **Postman** | Testes das requisições HTTP |

---

## Arquitetura do projeto

```text
projetoFinalGoVictorZuccoli/
│
├── main.go
├── models.go
├── handlers.go
├── regras.go
├── dados.go
├── go.mod
├── go.sum
└── README.md
```

### `main.go`

Responsável pela inicialização do servidor, configuração do Gin e definição das rotas da API.

### `models.go`

Contém as estruturas utilizadas pelo sistema:

```text
Aluno
Sala
Turma
Alocacao
MatriculaRequest
AlocacaoRequest
UsoSala
```

### `handlers.go`

Contém os handlers responsáveis por receber as requisições HTTP, validar os dados e retornar as respostas da API.

### `regras.go`

Centraliza as principais regras de negócio, como:

- identificação de alunos, salas e turmas;
- prevenção de matrícula duplicada;
- validação de horários;
- detecção de sobreposição;
- conflito de agenda de alunos;
- conflito na utilização das salas.

### `dados.go`

Armazena temporariamente os dados utilizados pela aplicação.

> Atualmente o projeto utiliza armazenamento **em memória**. Portanto, os dados cadastrados são apagados sempre que o servidor é reiniciado.

---

# Entidades

## Aluno

Representa um aluno cadastrado no sistema.

```json
{
  "matricula": "20260001",
  "nome": "Victor Zuccoli",
  "email": "victor@faculdade.edu.br",
  "ativo": true
}
```

A matrícula funciona como identificador único do aluno.

---

## Sala

Representa uma sala de aula ou laboratório.

```json
{
  "id": "LAB01",
  "nome": "Laboratório de Informática",
  "capacidade": 30,
  "recursos": [
    "Projetor",
    "Computadores",
    "Ar-condicionado"
  ],
  "ativa": true
}
```

Cada sala possui uma capacidade máxima e uma lista de recursos disponíveis.

---

## Turma

Representa uma turma acadêmica.

```json
{
  "id": "TURMA01",
  "nome": "Engenharia de Software - Turma A",
  "disciplina": "Desenvolvimento de APIs",
  "professor": "Professor responsável",
  "quantidade_alunos": 0,
  "alunos": [],
  "status_alocacao": "nao_alocada",
  "ativa": true
}
```

Uma turma pode possuir diversos alunos e posteriormente ser associada a uma sala e horário.

---

## Alocação

Representa a utilização de uma sala por uma turma.

```json
{
  "sala_id": "LAB01",
  "dia_semana": "segunda-feira",
  "horario_inicio": "19:00",
  "horario_termino": "21:00"
}
```

---

# Regras de negócio

O sistema possui validações para impedir inconsistências durante o gerenciamento das turmas.

### Matrícula duplicada

Um aluno não pode ser matriculado mais de uma vez na mesma turma.

Nesse caso, a API retorna:

```text
409 Conflict
```

---

### Capacidade da sala

Quando uma turma já possui uma sala alocada, a inclusão de um novo aluno não pode ultrapassar a capacidade máxima da sala.

Exemplo:

```text
Capacidade da sala: 30 alunos
Alunos matriculados: 30

Nova matrícula → recusada
```

Resposta:

```text
422 Unprocessable Entity
```

---

### Conflito de horário da sala

Uma mesma sala não pode receber duas turmas no mesmo período.

Exemplo:

```text
Turma A
19:00 ───────────── 21:00

Turma B
        20:00 ───────────── 22:00
```

Existe sobreposição entre os horários.

A regra utilizada é:

```text
novo_inicio < horario_final_existente
E
novo_final > horario_inicio_existente
```

Se as duas condições forem verdadeiras, existe conflito.

Resposta:

```text
409 Conflict
```

---

### Conflito de agenda do aluno

Um aluno não pode participar de duas turmas com horários sobrepostos no mesmo dia.

Exemplo:

```text
Aluno: Victor

Turma A
segunda-feira
19:00 - 21:00

Turma B
segunda-feira
20:00 - 22:00
```

A matrícula ou alocação será recusada.

Resposta:

```text
409 Conflict
```

---

### Validação de horário

Os horários devem utilizar o formato:

```text
HH:MM
```

Exemplo válido:

```text
19:00
```

Além disso:

```text
horario_inicio < horario_termino
```

---

# Endpoints

A URL base da API é:

```text
http://localhost:8080/api/v1
```

## Monitoramento

| Método | Endpoint | Descrição |
|---|---|---|
| `GET` | `/health` | Verifica o funcionamento da API |

---

## Alunos

| Método | Endpoint | Descrição |
|---|---|---|
| `POST` | `/alunos` | Cadastra um aluno |
| `GET` | `/alunos` | Lista todos os alunos |
| `GET` | `/alunos/:id` | Busca um aluno pela matrícula |

---

## Salas

| Método | Endpoint | Descrição |
|---|---|---|
| `POST` | `/salas` | Cadastra uma sala |
| `GET` | `/salas` | Lista todas as salas |
| `GET` | `/salas/:id/grade` | Consulta a utilização de uma sala |

---

## Turmas

| Método | Endpoint | Descrição |
|---|---|---|
| `POST` | `/turmas` | Cadastra uma turma |
| `GET` | `/turmas` | Lista todas as turmas |
| `POST` | `/turmas/:id/alunos` | Matricula um aluno em uma turma |
| `GET` | `/turmas/:id/alunos` | Lista os alunos matriculados |
| `POST` | `/turmas/:id/alocar` | Aloca uma sala para uma turma |

---

# Como executar o projeto

## 1. Pré-requisitos

É necessário possuir o **Go 1.22 ou superior** instalado.

Para verificar:

```bash
go version
```

Também é recomendado possuir o Git instalado:

```bash
git --version
```

---

## 2. Clonar o repositório

Execute:

```bash
git clone https://github.com/zZuccoli/projetoFinalGoVictorZuccoli.git
```

Entre na pasta:

```bash
cd projetoFinalGoVictorZuccoli
```

---

## 3. Instalar as dependências

Execute:

```bash
go mod tidy
```

O projeto utiliza o Gin:

```text
github.com/gin-gonic/gin
```

As demais dependências são gerenciadas automaticamente pelo Go através dos arquivos:

```text
go.mod
go.sum
```

---

## 4. Executar a API

Execute:

```bash
go run .
```

O servidor será iniciado em:

```text
http://localhost:8080
```

---

## 5. Verificar se a API está funcionando

Acesse no navegador ou Postman:

```http
GET http://localhost:8080/api/v1/health
```

Resposta esperada:

```json
{
  "status": "healthy",
  "timestamp": "2026-09-23T18:00:00-03:00",
  "version": "1.0.0"
}
```

Se essa resposta for apresentada, a API está funcionando corretamente.

---

# Exemplos de utilização

## Cadastrar um aluno

```http
POST /api/v1/alunos
```

Body:

```json
{
  "matricula": "20260001",
  "nome": "Victor Zuccoli",
  "email": "victor@faculdade.edu.br"
}
```

Resposta:

```text
201 Created
```

---

## Cadastrar uma sala

```http
POST /api/v1/salas
```

Body:

```json
{
  "id": "LAB01",
  "nome": "Laboratório de Informática",
  "capacidade": 30,
  "recursos": [
    "Projetor",
    "Computadores",
    "Ar-condicionado"
  ]
}
```

Resposta:

```text
201 Created
```

---

## Cadastrar uma turma

```http
POST /api/v1/turmas
```

Body:

```json
{
  "id": "TURMA01",
  "nome": "Turma A",
  "disciplina": "Desenvolvimento de APIs",
  "professor": "Professor responsável"
}
```

Resposta:

```text
201 Created
```

---

## Matricular aluno em uma turma

```http
POST /api/v1/turmas/TURMA01/alunos
```

Body:

```json
{
  "aluno_id": "20260001"
}
```

Resposta:

```text
201 Created
```

---

## Alocar uma sala

```http
POST /api/v1/turmas/TURMA01/alocar
```

Body:

```json
{
  "sala_id": "LAB01",
  "dia_semana": "segunda-feira",
  "horario_inicio": "19:00",
  "horario_termino": "21:00"
}
```

Se todas as regras de negócio forem atendidas:

```text
200 OK
```

---

# Códigos HTTP utilizados

| Código | Significado | Utilização |
|---:|---|---|
| `200` | OK | Consulta ou operação realizada com sucesso |
| `201` | Created | Recurso criado com sucesso |
| `400` | Bad Request | Dados enviados são inválidos |
| `404` | Not Found | Aluno, turma ou sala não encontrados |
| `409` | Conflict | Duplicidade ou conflito de agenda |
| `422` | Unprocessable Entity | Regra de capacidade ou estado não atendida |

---

# Fluxo principal do sistema

```text
Cadastrar Sala
      │
      ▼
Cadastrar Aluno
      │
      ▼
Cadastrar Turma
      │
      ▼
Matricular Aluno
      │
      ▼
Alocar Sala e Horário
      │
      ▼
┌──────────────────────────┐
│ Motor de validação       │
│                          │
│ ✓ Capacidade             │
│ ✓ Sala disponível        │
│ ✓ Horário válido         │
│ ✓ Conflito de sala       │
│ ✓ Conflito de aluno      │
│ ✓ Duplicidade            │
└─────────────┬────────────┘
              │
              ▼
      Alocação realizada
```

---

# Exemplo de cenário

Considere duas turmas:

```text
TURMA01
Sala: LAB01
Segunda-feira
19:00 - 21:00
```

e:

```text
TURMA02
Sala: LAB01
Segunda-feira
20:00 - 22:00
```

Ao tentar realizar a segunda alocação, o sistema identifica a sobreposição:

```text
19:00 ─────────── 21:00
          20:00 ─────────── 22:00
          ↑
       conflito
```

A operação é recusada com:

```text
409 Conflict
```

Isso impede que duas turmas utilizem a mesma sala simultaneamente.

---

# Monitoramento

O endpoint:

```http
GET /api/v1/health
```

permite verificar continuamente a disponibilidade da aplicação.

Ele informa:

```json
{
  "status": "healthy",
  "timestamp": "horário atual do servidor",
  "version": "1.0.0"
}
```

---

# Persistência de dados

Nesta versão, o projeto utiliza armazenamento em memória:

```go
var alunos = []Aluno{}
var salas = []Sala{}
var turmas = []Turma{}
```

Isso mantém a implementação simples e focada nas regras de negócio da API.

Ao encerrar o programa:

```text
Ctrl + C
```

os dados cadastrados durante aquela execução são descartados.

Uma evolução futura do projeto poderia utilizar um banco de dados como PostgreSQL ou MySQL.

---

# Possíveis evoluções

O projeto pode futuramente receber:

- persistência utilizando banco de dados;
- autenticação de usuários;
- perfis de acesso;
- atualização e exclusão de registros;
- paginação;
- filtros de consulta;
- documentação Swagger/OpenAPI;
- testes automatizados;
- Docker;
- logs estruturados;
- interface web para gerenciamento acadêmico.

---

# Autor

**Victor Zuccoli**

Projeto desenvolvido para atividade acadêmica de construção de APIs utilizando **Go** e **Gin**.

---

<p align="center">
  <strong>SGA — Sistema de Gestão de Alocação</strong>
</p>

<p align="center">
  Gestão acadêmica com validação de capacidade, matrícula e conflitos de horário.
</p>
