# Boost

[![Go](https://github.com/xgodev/boost/actions/workflows/go.yml/badge.svg)](https://github.com/xgodev/boost/actions/workflows/go.yml)

## Visão Geral

Boost é um framework modular e extensível para desenvolvimento de aplicações em Go, projetado para simplificar a criação de serviços robustos, escaláveis e observáveis. O framework fornece um conjunto abrangente de componentes que podem ser utilizados de forma independente ou combinados para construir aplicações completas.

## Características Principais

- **Arquitetura Modular**: Componentes independentes que podem ser utilizados conforme necessário
- **Injeção de Dependências**: Sistema flexível para gerenciamento de dependências entre componentes
- **Observabilidade Integrada**: Suporte nativo para métricas, logs e rastreamento
- **Adaptadores para Bibliotecas Populares**: Wrappers consistentes para diversas bibliotecas do ecossistema Go
- **Padrões de Design**: Implementações de padrões como Factory, Middleware e Wrapper
- **Suporte a Cloud Events**: Integração nativa com o padrão CloudEvents
- **Extensibilidade**: Facilidade para adicionar suporte a novas bibliotecas e frameworks

## Estrutura do Projeto

O projeto está organizado em pacotes modulares, cada um com responsabilidades específicas:

- **bootstrap**: Componentes para inicialização de aplicações
- **examples**: Exemplos de uso do framework
- **extra**: Funcionalidades adicionais (health checks, middleware, multiserver)
- **factory**: Implementações de padrões de fábrica para diversos componentes
- **fx**: Injeção de dependências e gerenciamento do ciclo de vida
- **model**: Definições de estruturas de dados e interfaces
- **utils**: Utilitários diversos
- **wrapper**: Adaptadores para bibliotecas externas (cache, config, log, publisher)

## Começando

### Pré-requisitos

- Go 1.18 ou superior
- Dependências específicas variam conforme os componentes utilizados

### Instalação

```bash
go get github.com/xgodev/boost
```

### Exemplo Básico

```go
package main

import (
    "context"
    "github.com/xgodev/boost/bootstrap"
    "github.com/xgodev/boost/wrapper/log"
    "github.com/xgodev/boost/wrapper/config"
)

func main() {
    ctx := context.Background()
    
    // Inicializar configuração
    cfg := config.New()
    
    // Inicializar logger
    logger := log.New(cfg)
    
    // Iniciar aplicação
    app := bootstrap.NewApp(cfg, logger)
    app.Start(ctx)
    
    // Aguardar sinal de término
    app.WaitForShutdown()
}
```

## Componentes Principais

### Bootstrap

O pacote Bootstrap fornece componentes para inicialização de aplicações, incluindo suporte para funções serverless e processamento de eventos.

[Documentação do Bootstrap](./bootstrap/README.md)

### Factory

O pacote Factory implementa o padrão de design Factory para diversos componentes, facilitando a criação e configuração de objetos complexos.

[Documentação do Factory](./factory/README.md)

### Wrapper

O pacote Wrapper fornece adaptadores para bibliotecas externas, incluindo cache, configuração, logging e publicação de mensagens.

[Documentação do Wrapper](./wrapper/README.md)

### FX

O pacote FX fornece funcionalidades de injeção de dependências e gerenciamento do ciclo de vida de componentes.

[Documentação do FX](./fx/README.md)

### Extra

O pacote Extra inclui funcionalidades adicionais como health checks, middleware e suporte para múltiplos servidores.

[Documentação do Extra](./extra/README.md)

## Integração com Tecnologias Populares

O Boost oferece integrações com diversas tecnologias populares:

- **Messaging**: NATS, Kafka, Google Cloud Pub/Sub
- **Cache**: Redis, BigCache, FreeCache
- **Logging**: Zap, Zerolog, Logrus
- **Web**: Echo, gRPC
- **Configuração**: Koanf, Environment Variables
- **Observabilidade**: Prometheus, OpenTelemetry

## Contribuição

Contribuições são bem-vindas! Para contribuir:

1. Fork o repositório
2. Crie uma branch para sua feature (`git checkout -b feature/amazing-feature`)
3. Commit suas mudanças (`git commit -m 'Add some amazing feature'`)
4. Push para a branch (`git push origin feature/amazing-feature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob os termos da licença incluída no arquivo [LICENSE](./LICENSE).

## Contato

Para questões, sugestões ou contribuições, por favor abra uma issue no repositório GitHub.

---

Desenvolvido com ❤️ pela comunidade xgodev.
