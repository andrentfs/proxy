# 🚀 K8s-Egress-Relay: Navegação Segura e IPs Residenciais para seus Crawlers

Este projeto nasceu de uma necessidade real: **Como fazer com que meus crawlers rodando no Kubernetes naveguem utilizando meu endereço de IP local (residencial), evitando bloqueios e captchas de sites sensíveis?**

A solução é um **Túnel SOCKS5 Reverso** de alta performance, construído em Go, que permite que o tráfego de saída (egress) de qualquer Pod no K8s seja roteado através da sua máquina física.

---

## 🛠 O que ele faz?

Imagine que você tem um crawler pesado rodando em um cluster GKE, AWS ou Azure. Frequentemente, os IPs desses datacenters são marcados como "suspeitos" por muitos sites. 

Com o **K8s-Egress-Relay**, você abre uma ponte:
1. O **Relay** roda no seu cluster Kubernetes como um Proxy SOCKS5.
2. O **Agente Local (GUI)** roda no seu computador pessoal (Windows/Mac/Linux).
3. Uma conexão persistente (via protocolo Yamux) é estabelecida entre os dois.
4. Seus Pods no K8s passam a "enxergar" a internet através dos olhos (e IP) da sua casa ou escritório.

---

## ✨ Diferenciais deste Projeto

- **Interface Gráfica Premium**: Desenvolvida com Fyne, focada em UX minimalista com indicadores de status dinâmicos.
- **Auto-Connect**: Abra o app e pronto, o túnel é estabelecido automaticamente.
- **Resiliência de Rede**: Implementação customizada de KeepAlive para evitar quedas em conexões instáveis.
- **Configuração via YAML**: Sem complicações, altere IPs e portas em um arquivo simples.
- **Sessão Dinâmica**: O servidor Relay gerencia reconexões de forma transparente, sem derrubar as requisições ativas dos seus Pods.

---

## 🚀 Como usar

### 1. No Kubernetes
Faça o build do `relay` usando o Dockerfile incluso e suba para seu repositório:
```bash
docker build --platform linux/amd64 -t seu-repositorio/proxy-relay:v1 .
docker push seu-repositorio/proxy-relay:v1
```
Aplique os manifestos da pasta `/k8s` no seu cluster.

### 2. Na sua Máquina (Local)
Configure o seu `config.yaml`:
```yaml
relay_host: "1.2.3.4"  # O External-IP do seu Service egress-proxy
relay_port: "8080"
auto_connect: true
```

Execute o binário ou rode via Go:
```bash
go run cmd/proxy-gui/main.go
```

---

## � Compilação e Build (Multiplataforma)

Este projeto utiliza `fyne-cross` para facilitar a criação de binários para diferentes sistemas operacionais sem complicação.

### Pré-requisitos
- Go instalado
- Docker rodando (necessário para o `fyne-cross`)

### Passo a Passo
1. Instale o fyne-cross:
```bash
go install github.com/fyne-io/fyne-cross@latest
```

2. Gere os builds:

* **Windows**:
```bash
fyne-cross windows -app-id com.proxy.manager -arch amd64 ./cmd/proxy-gui/
```

* **Linux**:
```bash
fyne-cross linux -app-id com.proxy.manager -arch amd64 ./cmd/proxy-gui/
```

* **macOS**:
```bash
fyne-cross darwin -app-id com.proxy.manager -arch amd64 ./cmd/proxy-gui/
```

Os executáveis serão gerados na pasta `fyne-cross/bin/`.

> **🍏 Guia para o Usuário Final (Mac)**: Preparamos um guia simplificado para você enviar junto com o app para usuários de macOS. Veja em: [README-MAC-USER.md](./README-MAC-USER.md).

> **💡 Dica de Distribuição**: Você pode disponibilizar diretamente os arquivos `.exe` (Windows), `.tar.gz` (Linux) ou `.app` (Mac). Lembre-se apenas de que o usuário final deve colocar o arquivo `config.yaml` na mesma pasta do executável.

---

Para que seu crawler ou qualquer aplicação dentro de um Pod utilize o túnel, você deve configurar as variáveis de ambiente de proxy utilizando o protocolo `socks5://`.

Exemplo de configuração no seu `deployment.yaml` ou `pod.yaml`:

```yaml
spec:
  containers:
  - name: seu-crawler
    image: sua-imagem:latest
    env:
    - name: HTTP_PROXY
      value: "socks5://egress-proxy.my-namespace.svc.cluster.local:1080"
    - name: HTTPS_PROXY
      value: "socks5://egress-proxy.my-namespace.svc.cluster.local:1080"
    - name: ALL_PROXY
      value: "socks5://egress-proxy.my-namespace.svc.cluster.local:1080"
```

> **Nota**: O endereço `egress-proxy.my-namespace.svc.cluster.local` deve ser alterado para o nome do serviço e namespace que você definiu nos manifestos do Relay.

---

## 🏗 Arquitetura

O projeto utiliza o protocolo **Yamux** para multiplexar conexões sobre um único canal TCP. Isso significa que, mesmo se você estiver navegando em centenas de sites simultaneamente através do túnel, apenas uma conexão física é mantida, garantindo baixa latência e economia de recursos.

---

## 🤝 Contribua!

Este projeto é open-source e focado na comunidade de engenharia de dados e web scraping. 
- Sentiu falta de alguma feature? Abra uma Issue!
- Quer melhorar a estabilidade? Faça um Pull Request!

---

## 💡 Sobre o Autor

Desenvolvi esta ferramenta para resolver problemas reais de extração de dados em escala, onde a reputação de IP é o fator crítico de sucesso. Se isso te ajudou, me dê um alô no [LinkedIn](https://www.linkedin.com/in/andrentfs/)!

---

*Feito com ❤️ e Go.*
