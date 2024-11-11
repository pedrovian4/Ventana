# Documentação do Ventana


## Configurações

O Ventana utiliza um arquivo de configuração JSON que armazena as preferências do usuário, incluindo idioma, caminho para o histórico de conexões, diretório de mensagens, e outras configurações. Abaixo estão as principais configurações e como utilizá-las.

## Estrutura do Arquivo de Configuração
O ventana cria uma configuração padrão no primeira abertura dele.

Exemplo de arquivo de configuração `ventana.json`:

```json

{
  "language": "en",
  "history_file": "/home/usuario/.config/ventana/history.csv",
  "message_directory": "./messages",
  "default_retry_count": 3,
  "default_retry_delay": 5,
  "config_file_path": "/home/usuario/.config/ventana/ventana.json",
  "scripts_directory": "/home/usuario/.config/ventana/scripts"
}

```

##  Explicação das Configurações

- language: Define o idioma da interface do Ventana. Os idiomas suportados são "en" (inglês) e "pt" (português).

- history_file: Caminho para o arquivo de histórico (history.csv). Esse arquivo armazena as conexões feitas anteriormente, permitindo fácil acesso a servidores utilizados recentemente. O caminho padrão é ~/.config/ventana/history.csv.

- message_directory: Diretório onde estão armazenadas as mensagens de interface. O Ventana usa arquivos de mensagem para exibir textos na língua definida na configuração "language". O diretório padrão é ./messages.


- default_retry_count: Define o número padrão de tentativas em caso de falha de conexão com um servidor. O valor padrão é 3, mas pode ser ajustado conforme a necessidade.

- default_retry_delay: Tempo (em segundos) de espera entre cada tentativa de reconexão ao servidor. O valor padrão é 5 segundos.

- config_file_path: Caminho para o próprio arquivo de configuração ventana.json. Este arquivo armazena as preferências e deve estar localizado em ~/.config/ventana/ventana.json.

- scripts_directory: Diretório onde os scripts devem ser armazenados para execução. O Ventana executa scripts escritos em Shell usando o bash, e esses scripts devem ser colocados na pasta ~/.config/ventana/scripts para serem reconhecidos pelo programa.