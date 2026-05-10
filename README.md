# API Clima Go

Uma API REST desenvolvida em Go que agrega dados climáticos de múltiplos provedores externos simultaneamente, retornando a resposta mais rápida em um formato padronizado.

Este projeto foi desenvolvido como estudo de **concorrência em Go**, explorando goroutines, channels, context e padrões de arquitetura idiomáticos da linguagem.

---

## Objetivo

A maioria dos tutoriais de backend ensina operações CRUD. Este projeto segue uma abordagem diferente: ao invés de gerenciar um banco de dados, a API orquestra chamadas para serviços externos, agrega os resultados e entrega uma resposta unificada.

O desafio central é fazer isso de forma eficiente — não sequencialmente, mas de forma concorrente.

---

## Como Funciona

```
Cliente
  |
  v
Handler (Gin)
  |
  v
WeatherService
  |              |
  v              v
OpenWeather   WeatherAPI     <- goroutines disparadas simultaneamente
  |              |
  v              v
  Primeira resposta vence
  |
  v
JSON padronizado
```

Ao invés de chamar uma API, esperar, e depois chamar a outra, os dois provedores são disparados ao mesmo tempo usando goroutines. O que responder primeiro é retornado ao cliente. A goroutine restante é cancelada via context.

---

## Endpoints

**Buscar clima por cidade**
```
GET /weather/:city
```

Exemplo de requisição:
```
GET http://localhost:8080/weather/recife
```

Exemplo de resposta:
```json
{
  "city": "Recife",
  "temperature": 26.02,
  "humidity": 83,
  "condition": "algumas nuvens",
  "source": "OpenWeather"
}
```

O campo `source` identifica qual provedor respondeu primeiro.

**Health check**
```
GET /ping
```

---

## Estrutura do Projeto

```
api-clima-go/
├── cmd/
│   └── main.go              # Ponto de entrada, injeção de dependências
│
├── internal/
│   ├── config/
│   │   └── config.go        # Carrega variáveis de ambiente
│   ├── dto/
│   │   ├── openweather.go   # Mapeia resposta da OpenWeather API
│   │   └── weatherapi.go    # Mapeia resposta da WeatherAPI
│   ├── handler/
│   │   └── weather.go       # Camada HTTP, handlers do Gin
│   ├── model/
│   │   └── weather.go       # Structs internas padronizadas
│   ├── provider/
│   │   ├── provider.go      # Interface WeatherProvider
│   │   ├── openweather.go   # Implementação OpenWeather
│   │   └── weatherapi.go    # Implementação WeatherAPI
│   └── service/
│       └── weather.go       # Lógica de concorrência e agregação
│
├── .env.example
├── go.mod
└── go.sum
```

---

## Conceitos Praticados

**Goroutines**

Cada provedor é chamado dentro de sua própria goroutine, permitindo execução simultânea sem bloqueio.

```go
for _, p := range s.providers {
    go func(p provider.WeatherProvider) {
        result, err := p.GetWeather(ctx, city)
        // ...
    }(p)
}
```

**Channels**

Resultados e erros das goroutines são comunicados de volta através de channels com buffer.

```go
resultChan := make(chan *model.WeatherResponse, len(s.providers))
errorChan  := make(chan error, len(s.providers))
```

**Context com Timeout**

Cada requisição tem um prazo de 3 segundos. Se nenhum provedor responder a tempo, o context é cancelado e um erro é retornado.

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
```

**Interfaces**

Ambos os provedores implementam a mesma interface, tornando a camada de serviço completamente desacoplada de qualquer API específica.

```go
type WeatherProvider interface {
    GetWeather(ctx context.Context, city string) (*model.WeatherResponse, error)
    Name() string
}
```

Adicionar um novo provedor requer apenas implementar essa interface — sem alterações no service ou no handler.

**Padrão DTO**

Cada API externa retorna uma estrutura JSON diferente. Os DTOs mapeiam essas respostas brutas para um modelo interno único, mantendo o restante da aplicação desacoplado dos formatos externos.

---

## Tecnologias

| Tecnologia | Finalidade |
|---|---|
| Go | Linguagem principal |
| Gin | Framework HTTP |
| godotenv | Carregamento de variáveis de ambiente |
| OpenWeatherMap API | Provedor de dados climáticos |
| WeatherAPI | Provedor de dados climáticos |

---

## Como Rodar

**1. Clone o repositório**
```bash
git clone https://github.com/robertoMateus/API-Clima-Go.git
cd api-clima-go
```

**2. Configure as variáveis de ambiente**
```bash
cp .env.example .env
```

Edite o `.env` com suas chaves de API:
```env
OPENWEATHER_API_KEY=sua_chave_aqui
WEATHERAPI_KEY=sua_chave_aqui
SERVER_PORT=8080
```

Ambos os provedores oferecem plano gratuito:
- OpenWeatherMap: https://openweathermap.org/api
- WeatherAPI: https://www.weatherapi.com

**3. Instale as dependências**
```bash
go mod tidy
```

**4. Execute a partir da raiz do projeto**
```bash
go run cmd/main.go
```

**5. Teste**
```bash
curl http://localhost:8080/weather/recife
```

---

## Autor

Desenvolvido por Roberto Mateus como projeto de estudo de backend com foco em concorrencia em Go.