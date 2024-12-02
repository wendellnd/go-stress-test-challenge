# go-stress-test-challenge

#### Execução

1. Crie uma imagem utilizando o `Dockerfile` disponível no repositório

   ```
   docker build -t stress:latest .
   ```

2. Execute o comando para executar o teste de stress utilizando imagem gerada
   ```
   docker run stress:latest --url https://google.com --requests 100 --concurrency 10
   ```
