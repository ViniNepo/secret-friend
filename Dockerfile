# Etapa 1: Build da aplicação
FROM golang:1.21 AS builder

WORKDIR /app

# Copiar e instalar dependências
COPY go.mod go.sum ./
RUN go mod tidy

# Copiar o código e construir a aplicação
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o secretfriend .

# Etapa 2: Imagem final com certificados confiáveis
FROM gcr.io/distroless/static:nonroot

# Copiar o binário da etapa de build
COPY --from=builder /app/secretfriend /

# Definir o ponto de entrada
ENTRYPOINT [ "/secretfriend" ]
