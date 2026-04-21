# Finance Management

Sistema simples de controle financeiro pessoal feito em Go, MariaDB e frontend estático com HTML, CSS e Bootstrap.

O projeto permite criar usuários, registrar receitas e despesas, gerenciar categorias, definir orçamento por categoria, visualizar relatório mensal e acompanhar quanto da renda dos próximos meses já está comprometida por transações parceladas.

## Tecnologias

- Go 1.25
- MariaDB
- GORM
- Gorilla Mux
- Docker e Docker Compose
- HTML, CSS e Bootstrap

## Requisitos

- Docker
- Docker Compose
- Go 1.25, caso rode a API fora do container
- [Air](https://github.com/air-verse/air), caso use o modo desenvolvimento

## Configuração

Crie o arquivo `.env`:

```bash
cp .env.example .env
```

Exemplo de variáveis:

```env
DB_USER=root
DB_PASSWORD=root
DB_HOST=db
DB_PORT=3306
DB_NAME=finance_management
```

## Rodando com Docker

Suba os serviços:

```bash
make up-dettached
```

Rode as migrations:

```bash
make migration-up
```

Acesse:

```text
http://localhost:8080
```

## Modo Desenvolvimento

Neste modo, apenas o banco roda no Docker e a API roda localmente com Air.

Suba o banco:

```bash
make up-service SERVICE=db
```

Rode as migrations:

```bash
make migration-up
```

Baixe as dependências:

```bash
make set-up
```

Suba a API com Air:

```bash
make air-up SERVICE=api
```

Acesse:

```text
http://localhost:8080
```

## Testes

```bash
make test-unit
```

Se tiver problema de permissão com o cache do Go:

```bash
GOCACHE=/tmp/go-build-tcc go test ./...
```

## Frontend

O frontend fica em:

```text
frontend/
```

Ele é servido pela própria API. Em modo desenvolvimento, depois de alterar HTML/CSS/JS, normalmente basta dar refresh no navegador:

```text
Ctrl + F5
```

Telas atuais:

- Criar usuário
- Nova transação
- Gerenciar categorias
- Relatório mensal
- Projeção financeira

## Rotas Disponíveis

```text
POST   /api/usuarios
POST   /api/transacoes
GET    /api/usuarios/{usuario_id}/categorias
POST   /api/categorias
PUT    /api/categorias/{id}
DELETE /api/categorias/{id}
POST   /api/orcamentos
GET    /api/usuarios/{usuario_id}/orcamentos?mes=YYYY-MM-DD
GET    /api/usuarios/{usuario_id}/orcamentos/total?mes=YYYY-MM-DD
PUT    /api/orcamentos/{id}
GET    /api/usuarios/{usuario_id}/relatorios/mensal?mes=YYYY-MM-DD
GET    /api/usuarios/{usuario_id}/relatorios/gastos?mes=YYYY-MM-DD
GET    /api/usuarios/{usuario_id}/projecao/comprometimento?mes=YYYY-MM-DD&meses=4
```

## Exemplos de Uso

Os exemplos abaixo assumem que a API está rodando em:

```text
http://localhost:8080
```

### Criar Usuário

```bash
curl -X POST http://localhost:8080/api/usuarios \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Joao",
    "email": "joao@email.com"
  }'
```

Ao criar um usuário, o sistema cria automaticamente as categorias:

- Alimentação
- Transporte
- Lazer
- Moradia
- Receita

### Listar Categorias do Usuário

```bash
curl http://localhost:8080/api/usuarios/1/categorias
```

### Criar Categoria

```bash
curl -X POST http://localhost:8080/api/categorias \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Saúde",
    "usuario_id": 1
  }'
```

### Editar Categoria

```bash
curl -X PUT http://localhost:8080/api/categorias/6 \
  -H "Content-Type: application/json" \
  -d '{
    "id": 6,
    "nome": "Farmácia",
    "usuario_id": 1
  }'
```

### Excluir Categoria

```bash
curl -X DELETE http://localhost:8080/api/categorias/6
```

Evite excluir categorias que já estejam sendo usadas por transações ou orçamentos, pois o banco pode bloquear por chave estrangeira.

## Transações

### Criar Receita

Use a categoria `Receita`.

```bash
curl -X POST http://localhost:8080/api/transacoes \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": 1,
    "categoria_id": 5,
    "valor": 3000,
    "data": "2026-04-21T00:00:00Z",
    "descricao": "Salário",
    "tipo": "pix",
    "parcelas": 1
  }'
```

Receitas não são parceladas. Mesmo que `parcelas` seja maior que `1`, o sistema registra a receita apenas no mês escolhido.

### Criar Despesa

```bash
curl -X POST http://localhost:8080/api/transacoes \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": 1,
    "categoria_id": 1,
    "valor": 100,
    "data": "2026-04-21T00:00:00Z",
    "descricao": "Mercado",
    "tipo": "debito",
    "parcelas": 1
  }'
```

### Criar Despesa Parcelada

```bash
curl -X POST http://localhost:8080/api/transacoes \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": 1,
    "categoria_id": 3,
    "valor": 150,
    "data": "2026-04-21T00:00:00Z",
    "descricao": "Compra parcelada",
    "tipo": "cartao",
    "parcelas": 3
  }'
```

Esse exemplo cria uma transação em cada mês:

```text
2026-04-21
2026-05-21
2026-06-21
```

O campo `valor` representa o valor de cada parcela.

## Orçamentos

### Criar Orçamento por Categoria

```bash
curl -X POST http://localhost:8080/api/orcamentos \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": 1,
    "categoria_id": 1,
    "limite": 500,
    "mes": "2026-04-01T00:00:00Z"
  }'
```

O campo `mes` é normalizado para o primeiro dia do mês.

### Listar Orçamentos do Mês

```bash
curl "http://localhost:8080/api/usuarios/1/orcamentos?mes=2026-04-01"
```

### Total Planejado do Mês

```bash
curl "http://localhost:8080/api/usuarios/1/orcamentos/total?mes=2026-04-01"
```

### Editar Orçamento

```bash
curl -X PUT http://localhost:8080/api/orcamentos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": 1,
    "categoria_id": 1,
    "limite": 650,
    "mes": "2026-04-01T00:00:00Z"
  }'
```

## Relatórios

### Relatório Mensal

```bash
curl "http://localhost:8080/api/usuarios/1/relatorios/mensal?mes=2026-04-01"
```

Retorna:

- total de receitas
- total de despesas
- saldo atual

### Progresso de Gastos

```bash
curl "http://localhost:8080/api/usuarios/1/relatorios/gastos?mes=2026-04-01"
```

Retorna:

- total gasto no mês
- receita do mês
- percentual da receita já utilizada

### Projeção de Comprometimento

```bash
curl "http://localhost:8080/api/usuarios/1/projecao/comprometimento?mes=2026-04-01&meses=4"
```

Essa rota mostra quanto da renda dos próximos meses já está comprometida por transações futuras, como compras parceladas.

Importante:

- Receitas aparecem apenas no mês em que foram registradas.
- Despesas parceladas aparecem nos meses subsequentes.
- Se um mês futuro não tiver receita registrada, a receita daquele mês fica zerada.

## Dicas

Se uma migration falhar durante o desenvolvimento e você não tiver dados importantes, o caminho mais simples é recriar o banco:

```bash
make down
sudo rm -rf build/docker/mariadb/_dbdata
make up-dettached
make migration-up
```

Se estiver usando Air e alterar apenas arquivos do frontend, use:

```text
Ctrl + F5
```
