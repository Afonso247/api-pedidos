# api-pedidos

Este é uma aplicação web criado em Go como linguagem, e Redis como banco de dados remoto. O app recebe e responde a HTTP requests, como o objetivo de criar, gravar, atualizar, e deletar pedidos no banco de dados.

## Índice

- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Pré-requisitos](#pré-requisitos)
- [Instalação](#instalação)
- [Como Usar](#como-usar)
- [Contato](#contato)

## Tecnologias Utilizadas

- **Golang**: Linguagem de programação principal utilizada no desenvolvimento do projeto.
- **Redis**: Sistema de banco de dados NoSQL, baseado em memória, que é utilizado principalmente como um banco do tipo "chave-valor".
- **Bibliotecas**: Chi, para o roteamento HTTP. Viper, para a configuração do servidor e porta (Ex: localhost:8080)

## Pré-requisitos

Antes de começar, certifique-se de que você atendeu aos seguintes requisitos:

- Golang instalado (versão 1.20.5 ou acima)
- Redis instalado e em execução
- (Opcional) Postman instalado, para testar a API e realizar as requisições HTTP

## Instalação

1. Clone o repositório:

    ```bash
    git clone https://github.com/usuario/nome-do-projeto.git
    cd nome-do-projeto
    ```

2. Instale as dependências do projeto:

    ```bash
    go mod tidy
    ```

3. Inicie o Redis (caso não esteja utilizando Docker):

    ```bash
    redis-server
    ```

4. Crie um arquivo `config.json` dentro da pasta do projeto. Dentro do arquivo, insira o seguinte código:

    ```json
    {
    "redis_address": "endereço_redis",
    "server_port": 8080
    }  
    ```
    Edite a string `"endereço_redis"` pelo seu endereço Redis (Ex: `"localhost:6379"`)


5. (Opcional) Utilize Docker para configurar o ambiente:

    ```bash
    docker-compose up
    ```

## Como Usar

Siga os passos abaixo para executar o projeto e utilizar suas funcionalidades:

1. **Inicie o projeto:**

    Execute o comando abaixo para iniciar o servidor Go:

    ```bash
    go run main.go
    ```

    Certifique-se de que o comando está sendo executado dentro da pasta da aplicação.

    O servidor estará rodando na porta definida pela variável de ambiente `server_port` ou na porta padrão `8080`.

2. **Acesse o serviço:**

    Utilize o terminal ou o Postman para acessar e realizar as requisções HTTP.

    Faça um GET request para o servidor(Ex: `localhost:8080/`). O servidor deve retornar uma resposta 200 para você.

3. **Utilize as principais rotas da API:**

    A seguir estão algumas das principais rotas disponíveis na API que você pode testar:

    - **`GET /`**: Verifica se o servidor está ativo.

        ```bash
        curl -X GET http://localhost:8080/
        ```

    - **`GET /pedidos`**: Retorna todos os dados de `pedidos` armazenados no Redis.

        ```bash
        curl -X GET http://localhost:8080/pedidos
        ```

    - **`POST /pedidos`**: Envia pedidos para serem armazenados no Redis.

        ```bash
        curl -X POST http://localhost:8080/pedidos -d '{"cliente_id":"'$(uuidgen)'","line_items":[{"item_id":"'$(uuidgen)'","quantidade":10,"preco":100}]}'
        ```

    - **`GET /pedidos/{id}`**: Recupera um pedido armazenado no Redis por um id específico.

        ```bash
        curl -X GET http://localhost:8080/pedidos/example
        ```

    - **`PUT /pedidos/{id}`**: Atualiza um pedido armazenado no Redis. Somente o status do pedido pode ser atualizado.

        ```bash
        curl -X PUT -d '{"status":"enviado"}' -sS "localhost:8080/pedidos/example" | jq
        ```

        Se o status for "enviado", atualiza o campo `EnviadoEm` para o horário atual UTC-3, caso ainda não tenha sido configurado.

        Se o status for "concluido", atualiza o campo `ConcluidoEm` para o horário atual UTC-3, caso o campo `EnviadoEm` já tenha sido atualizado anteriormente.

    - **`DELETE /pedidos/{id}`**: Deleta um pedido armazenado no Redis por um id específico.

        ```bash
        curl -X DELETE http://localhost:8080/pedidos/example
        ```

4. **Parar o servidor:**

    Para parar o servidor, pressione `CTRL + C` no terminal onde o comando `go run main.go` foi executado.

5. **Logs e Erros:**

    Os logs serão exibidos no terminal onde o servidor foi iniciado. Caso ocorra algum erro, as mensagens de erro serão registradas lá. Verifique o terminal para mais detalhes.

6. **Redis:**

    Certifique-se de que o servidor Redis esteja em execução. Caso contrário, o projeto não funcionará corretamente. Você pode verificar o status do Redis com o comando:

    ```bash
    redis-cli ping
    ```

    Se o retorno for `PONG`, o Redis está ativo.

## Contato

[Afonso Henrique](https://github.com/Afonso247)  
📧 **Email**: [afonsoh.dev@gmail.com](mailto:email@seuemail.com)   
💼 **LinkedIn**: [Afonso Henrique](https://linkedin.com/in/afonsoh247)

Sinta-se à vontade para entrar em contato caso tenha alguma dúvida ou sugestão!

Obrigado!
