### Apresentação

Este projeto visa resolver o desafio de Bancos de Dados proposto pela fullcycle no curso de Go Expert.

Os requisitos são entregar dois sistemas em Go:
- client.go
- server.go
 
Os requisitos de cada aplicação serão descritos à seguir.

### Server.go 

O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente.

Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.

O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.
 

### Client.go

O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
 
O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.

O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}

### Requisitos Gerais

Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.

## Como Executar o Projeto
1. Clone o repositório
`git clone https://github.com/gabriel-ulisses-andrade/goexpert-desafio-client-server-api/tree/main`
`cd goexpert-desafio-client-server-api`

2. Execute o servidor
`cd server`
`go mod tity` 
`go run main.go` 

3. Execute o client
`cd client`
`go mod tity` 
`go run main.go` 