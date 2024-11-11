# Ventana

Ventana é uma ferramenta de linha de comando para facilitar o gerenciamento de múltiplos servidores. Este repositório inclui o código-fonte e a documentação para desenvolvedores e contribuidores. Abaixo, você encontrará informações sobre a estrutura do projeto e instruções para contribuir.

**Ventana** is a command-line tool designed to facilitate the management of multiple servers. This repository includes the source code and documentation for developers and contributors. Below, you’ll find information about the project structure and instructions for contributing.

---

## 📁 Estrutura do Projeto | Project Structure

O repositório está organizado da seguinte forma:

```

├── cmd
│   └── ventana
│       └── main.go
├── go.mod
├── go.sum
├── install.sh
├── messages
│   ├── en.json
│   ├── pt.json
│   └── ru.json
├── pkg
│   ├── arts
│   │   └── welcome.go
│   ├── config
│   │   ├── colors.go
│   │   └── config.go
│   ├── history
│   │   └── history.go
│   ├── localization
│   │   ├── language.go
│   │   └── localization.go
│   ├── menu
│   │   └── menu.go
│   └── server
│       ├── script_manager.go
│       ├── server_manager.go
│       ├── share_manager.go
│       └── syntax_manager.go
├── README.MD
└── ventana-docs
    ├── assets
    │   ├── image.png
    │   └── ventana.gif
    ├── _config.yml
    ├── en
    │   └── index.md
    ├── index.md
    └── pt
        └── index.md
```

---

## 📚 Documentação | Documentation

A documentação completa para usuários e desenvolvedores pode ser encontrada nos seguintes locais:

- **Documentação em Português**: [docs/pt/index.md](./ventana-docs/pt/index.md)
- **English Documentation**: [docs/en/index.md](./ventana-docs/en/index.md)

---

## 🛠️ Instruções para Contribuir | Contributing Instructions

### Português

Para contribuir com o Ventana, siga as etapas abaixo:

 **Fork este repositório** e crie uma nova branch para sua funcionalidade ou correção:
   ```bash
   git checkout -b feature/minha-feature
```


### English

To contribute to Ventana, please follow these:

Fork this repository and create a new branch for your feature or bugFix:
```bash
git checkout -b feature/my-feature
```


