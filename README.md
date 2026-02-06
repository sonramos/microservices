# Microservices - Order, Payment e Shipping

Sistema de microserviços gRPC para gerenciamento de pedidos, pagamentos e envios.

## Arquitetura

O projeto utiliza **Arquitetura Hexagonal** (Ports & Adapters) com comunicação via gRPC.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│    Order    │────▶│   Payment   │     │  Shipping   │
│   :3000     │     │   :3001     │     │   :3002     │
└─────────────┘     └─────────────┘     └─────────────┘
       │                   │                   │
       └───────────────────┴───────────────────┘
                           │
                    ┌──────┴──────┐
                    │    MySQL    │
                    │    :3306    │
                    └─────────────┘
```

## Pré-requisitos

- Docker e Docker Compose
- Go 1.25+ (para desenvolvimento local)
- grpcurl (opcional, para testes)

## Deploy com Docker

### 1. Iniciar todos os serviços

```bash
cd microservices
docker-compose up -d
```

Isso irá:
- Iniciar o MySQL com os bancos `order`, `payment` e `shipping`
- Construir e iniciar os três microserviços

### 2. Verificar status

```bash
docker-compose ps
```

### 3. Parar serviços

```bash
docker-compose down
```

Para remover também os volumes (dados):
```bash
docker-compose down -v
```

## Desenvolvimento Local

### 1. Configurar dependências proto

O projeto depende do módulo `github.com/sonramos/microservices-proto`. Para desenvolvimento local:

```bash
# Os arquivos go.mod já contêm replace directives para paths locais
# Certifique-se de que microservices-proto está no path correto:
# ../microservices-proto/
```

### 2. Iniciar MySQL

```bash
docker run -d --name mysql-dev \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=minhasenha \
  mysql:8.0

# Criar bancos de dados
docker exec mysql-dev mysql -uroot -pminhasenha -e "
  CREATE DATABASE IF NOT EXISTS \`order\`;
  CREATE DATABASE IF NOT EXISTS \`payment\`;
  CREATE DATABASE IF NOT EXISTS \`shipping\`;
"
```

### 3. Iniciar serviços

**Terminal 1 - Payment:**
```bash
cd payment
ENV=development \
APPLICATION_PORT=3001 \
DATA_SOURCE_URL="root:minhasenha@tcp(localhost:3306)/payment?charset=utf8mb4&parseTime=True&loc=Local" \
go run cmd/main.go
```

**Terminal 2 - Shipping:**
```bash
cd shipping
ENV=development \
APPLICATION_PORT=3002 \
DATA_SOURCE_URL="root:minhasenha@tcp(localhost:3306)/shipping?charset=utf8mb4&parseTime=True&loc=Local" \
go run cmd/main.go
```

**Terminal 3 - Order:**
```bash
cd order
ENV=development \
APPLICATION_PORT=3000 \
DATA_SOURCE_URL="root:minhasenha@tcp(localhost:3306)/order?charset=utf8mb4&parseTime=True&loc=Local" \
PAYMENT_SERVICE_URL=localhost:3001 \
SHIPPING_SERVICE_URL=localhost:3002 \
go run cmd/main.go
```

## Variáveis de Ambiente

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `ENV` | Ambiente (development habilita reflection) | `development` |
| `APPLICATION_PORT` | Porta do serviço gRPC | `3000` |
| `DATA_SOURCE_URL` | Connection string MySQL | `root:senha@tcp(host:3306)/db` |
| `PAYMENT_SERVICE_URL` | Endereço do Payment (só Order) | `payment:3001` |
| `SHIPPING_SERVICE_URL` | Endereço do Shipping (só Order) | `shipping:3002` |

## Testando com grpcurl

### Instalar grpcurl
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### Inserir dados de estoque (necessário para validação)
```bash
docker exec mysql mysql -uroot -pminhasenha -e "
  USE \`order\`;
  INSERT INTO stock_items (product_code, quantity, created_at, updated_at) 
  VALUES ('PROD001', 100, NOW(), NOW()), ('PROD002', 50, NOW(), NOW());
"
```

### Criar pedido
```bash
grpcurl -plaintext -d '{
  "costumer_id": 1, 
  "order_items": [
    {"product_code": "PROD001", "unit_price": 10.0, "quantity": 5},
    {"product_code": "PROD002", "unit_price": 20.0, "quantity": 3}
  ]
}' localhost:3000 Order/Create
```

### Testar Payment diretamente
```bash
grpcurl -plaintext -d '{
  "user_id": 1,
  "order_id": 1,
  "total_price": 100.0
}' localhost:3001 Payment/Create
```

### Testar Shipping diretamente
```bash
grpcurl -plaintext -d '{
  "order_id": 1,
  "items": [
    {"product_code": "PROD001", "quantity": 10}
  ]
}' localhost:3002 Shipping/Create
```

## Validações

| Validação | Serviço | Erro |
|-----------|---------|------|
| Quantidade total > 50 itens | Order | `InvalidArgument` |
| Produto não existe no estoque | Order | `NotFound` |
| Total do pagamento > 1000 | Payment | `InvalidArgument` |

## Cálculo de Dias de Entrega

O Shipping calcula os dias de entrega com base na quantidade total:
```
dias = 1 + (quantidade_total / 5)
```

Exemplos:
- 5 itens → 2 dias
- 12 itens → 3 dias
- 25 itens → 6 dias

## Estrutura do Projeto

```
microservices/
├── docker-compose.yml
├── init.sql
├── README.md
├── order/
│   ├── Dockerfile
│   ├── cmd/main.go
│   ├── config/
│   └── internal/
│       ├── adapters/
│       │   ├── db/
│       │   ├── grpc/
│       │   ├── payment/
│       │   └── shipping/
│       ├── application/core/
│       │   ├── api/
│       │   └── domain/
│       └── ports/
├── payment/
│   ├── Dockerfile
│   ├── cmd/main.go
│   ├── config/
│   └── internal/
│       ├── adapters/
│       │   ├── db/
│       │   └── grpc/
│       ├── application/core/
│       │   ├── api/
│       │   └── domain/
│       └── ports/
└── shipping/
    ├── Dockerfile
    ├── cmd/main.go
    ├── config/
    └── internal/
        ├── adapters/
        │   ├── db/
        │   └── grpc/
        ├── application/core/
        │   ├── api/
        │   └── domain/
        └── ports/
```

## Licença

Projeto acadêmico - IFPB Programação Distribuída
