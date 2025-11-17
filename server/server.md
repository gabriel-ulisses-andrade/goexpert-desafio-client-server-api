# ✅ Checklist do Desafio

## 📡 Consumo da API de câmbio
 - O server.go deve consumir a API: https://economia.awesomeapi.com.br/json/last/USD-BRL
 - O retorno da requisição deve fornecer a cotação USD/BRL ao servidor.

## 🔁 Retorno ao cliente
 - O servidor deve retornar a cotação em JSON no endpoint /cotacao.

## 🕒 Uso de contextos e timeouts
 - Utilizar o package context para controlar timeouts.
 - Timeout máximo de 200ms para chamar a API externa de cotação.
 - Timeout máximo de 10ms para persistir os dados no banco SQLite.
 - Os 3 contextos criados (requisição externa, persistência e request do cliente) devem registrar erro nos logs caso o tempo estoure.

## 🗃️ Persistência no banco
 - Registrar cada cotação recebida em um banco SQLite.
 - A gravação deve respeitar o timeout de 10ms.

## 🌐 Servidor HTTP
 - Criar o endpoint /cotacao.
 - O servidor deve rodar na porta 8080.

## 📤 Entrega
 - Enviar o link do repositório (GitHub, GitLab, etc.) com a solução final para correção.