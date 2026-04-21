# Plano de Testes da API

Este guia assume que a API esta rodando em:

```text
http://localhost:8080
```

Se a porta configurada no `.env` for outra, ajuste os exemplos.

## 1. Subir a Aplicacao

Copie o arquivo de ambiente, se ainda nao existir:

```bash
cp .env.example .env
```

Suba os containers:

```bash
make up-dettached
```

Rode as migrations:

```bash
make migration-up
```

## 2. Criar Usuario

```bash
curl -X POST http://localhost:8080/api/usuarios \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Joao",
    "email": "joao@email.com"
  }'
```

Resposta esperada:

```json
{
  "data": {
    "id": 1,
    "nome": "Joao",
    "email": "joao@email.com"
  }
}
```

Ao criar um usuario, o sistema tambem cria automaticamente as categorias:

```text
Alimentacao
Transporte
Lazer
Moradia
Receita
```

## 3. Validar Email Duplicado

Rode novamente a mesma request de criacao de usuario:

```bash
curl -X POST http://localhost:8080/api/usuarios \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Joao 2",
    "email": "joao@email.com"
  }'
```

Resposta esperada:

```text
HTTP 400
```

O erro deve informar que ja existe um usuario com esse email.

## 4. Listar Categorias do Usuario

Troque `1` pelo ID do usuario criado:

```bash
curl http://localhost:8080/api/usuarios/1/categorias
```

Resposta esperada:

```json
{
  "data": [
    {
      "id": 1,
      "nome": "Alimentação",
      "usuario_id": 1
    },
    {
      "id": 2,
      "nome": "Transporte",
      "usuario_id": 1
    }
  ]
}
```

Anote os IDs das categorias, principalmente:

```text
Receita
Alimentacao
Transporte
Lazer
Moradia
```

## 5. Criar Categoria

```bash
curl -X POST http://localhost:8080/api/categorias \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Saude",
    "usuario_id": 1
  }'
```

Resposta esperada:

```json
{
  "data": {
    "id": 6,
    "nome": "Saude",
    "usuario_id": 1
  }
}
```

## 6. Excluir Categoria

Troque `6` pelo ID da categoria que deseja excluir:

```bash
curl -X DELETE http://localhost:8080/api/categorias/6
```

Resposta esperada:

```json
{
  "data": {
    "message": "categoria deleted successfully"
  }
}
```

Evite excluir categorias que ja estejam sendo usadas em transacoes ou orcamentos, porque o banco pode bloquear por chave estrangeira.

## 7. Criar Transacao de Receita

Use o `categoria_id` da categoria `Receita`.

```bash
curl -X POST http://localhost:8080/api/transacoes \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": 1,
    "categoria_id": 5,
    "valor": 3000,
    "data": "2026-04-21T00:00:00Z",
    "descricao": "Salario",
    "tipo": "pix",
    "parcelas": 1
  }'
```

Resposta esperada:

```json
{
  "data": [
    {
      "id": 1,
      "usuario_id": 1,
      "categoria_id": 5,
      "valor": 3000,
      "data": "2026-04-21T00:00:00Z",
      "descricao": "Salario",
      "tipo": "pix",
      "parcelas": 1
    }
  ]
}
```

## 8. Criar Transacao de Despesa

Use uma categoria que nao seja `Receita`, por exemplo `Alimentacao`.

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

Resposta esperada: uma transacao criada dentro do array `data`.

## 9. Criar Transacao Parcelada

Se `parcelas` for maior que `1`, o sistema cria uma transacao por mes subsequente.

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

Resposta esperada: tres transacoes criadas, uma em cada mes:

```text
2026-04-21
2026-05-21
2026-06-21
```

Atencao: atualmente o sistema trata `valor` como o valor de cada parcela. Entao `valor: 150` com `parcelas: 3` cria tres transacoes de `150`.

## 10. Testar Validacoes

### Usuario sem Email

```bash
curl -X POST http://localhost:8080/api/usuarios \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Sem Email"
  }'
```

Resposta esperada:

```text
HTTP 400
```

## 11. Criar Orcamento por Categoria

Use uma categoria de despesa, por exemplo `Alimentacao`.

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

Resposta esperada:

```json
{
  "data": {
    "id": 1,
    "usuario_id": 1,
    "categoria_id": 1,
    "limite": 500,
    "mes": "2026-04-01T00:00:00Z"
  }
}
```

O sistema normaliza o campo `mes` para o primeiro dia do mes. Entao, mesmo se a data enviada for `2026-04-21T00:00:00Z`, o orcamento sera salvo como `2026-04-01`.

## 12. Criar Outro Orcamento no Mesmo Mes

```bash
curl -X POST http://localhost:8080/api/orcamentos \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": 1,
    "categoria_id": 2,
    "limite": 300,
    "mes": "2026-04-01T00:00:00Z"
  }'
```

Resposta esperada: um segundo orcamento criado.

## 13. Listar Orcamentos do Mes

Use `mes` no formato `YYYY-MM-DD`.

```bash
curl "http://localhost:8080/api/usuarios/1/orcamentos?mes=2026-04-01"
```

Resposta esperada:

```json
{
  "data": [
    {
      "id": 1,
      "usuario_id": 1,
      "categoria_id": 1,
      "limite": 500,
      "mes": "2026-04-01T00:00:00Z"
    },
    {
      "id": 2,
      "usuario_id": 1,
      "categoria_id": 2,
      "limite": 300,
      "mes": "2026-04-01T00:00:00Z"
    }
  ]
}
```

## 14. Buscar Total Planejado do Mes

```bash
curl "http://localhost:8080/api/usuarios/1/orcamentos/total?mes=2026-04-01"
```

Resposta esperada:

```json
{
  "data": {
    "total_planejado": 800
  }
}
```

## 15. Editar Orcamento

Troque `1` pelo ID do orcamento que deseja editar.

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

Resposta esperada:

```json
{
  "data": {
    "id": 1,
    "usuario_id": 1,
    "categoria_id": 1,
    "limite": 650,
    "mes": "2026-04-01T00:00:00Z"
  }
}
```

### Transacao com Valor Invalido

```bash
curl -X POST http://localhost:8080/api/transacoes \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": 1,
    "categoria_id": 1,
    "valor": 0,
    "data": "2026-04-21T00:00:00Z",
    "descricao": "Invalida",
    "tipo": "pix",
    "parcelas": 1
  }'
```

Resposta esperada:

```text
HTTP 400
```

## Rotas Disponiveis Hoje

```text
POST   /api/usuarios
POST   /api/transacoes
GET    /api/usuarios/{usuario_id}/categorias
POST   /api/categorias
DELETE /api/categorias/{id}
POST   /api/orcamentos
GET    /api/usuarios/{usuario_id}/orcamentos?mes=YYYY-MM-DD
GET    /api/usuarios/{usuario_id}/orcamentos/total?mes=YYYY-MM-DD
PUT    /api/orcamentos/{id}
```

## Proximos Endpoints a Expor

As services e repositories ja tem base para estes fluxos, mas ainda faltam rotas HTTP:

```text
GET    /api/usuarios/{usuario_id}/relatorios/mensal
GET    /api/usuarios/{usuario_id}/relatorios/gastos
PUT    /api/categorias/{id}
```
