# 🍏 Guia de Uso: Proxy Manager para macOS

Como fazer para usar o **Proxy Manager** no seu Mac, siga este guia para configurar e rodar o serviço corretamente.

## 📦 O que você precisa
Você deve ter recebido dois itens:
1.  **`proxy-gui.app`**: O aplicativo principal.
2.  **`config.yaml`**: O arquivo de configuração (contém o endereço do servidor).

---

## 🚀 Passo a Passo para Iniciar

### 1. Organização
Coloque o aplicativo `proxy-gui.app` e o arquivo `config.yaml` na **mesma pasta** (por exemplo, na sua Mesa ou na pasta Documentos).

### 2. Contornando o Aviso de Segurança (Gatekeeper)
Como este é um software especializado, o macOS pode mostrar um aviso dizendo que "não pode verificar o desenvolvedor". **Não se preocupe**, isso é normal para apps customizados.

Para abrir pela primeira vez:
1.  Clique com o **botão direito do mouse** (ou segure a tecla `Control` e clique) no ícone do **Proxy Manager**.
2.  Selecione a opção **Abrir** no menu que aparecer.
3.  Uma janela de confirmação surgirá. Clique no botão **Abrir**.
4.  *Nas próximas vezes, você poderá abrir com dois cliques normalmente.*

### 3. Usando o Aplicativo
*   **Conexão Automática**: Assim que o app abrir, ele tentará se conectar ao servidor. Você verá uma bolinha de status mudar de cor.
*   **Status Verde**: Significa que você está conectado com sucesso! O túnel está ativo.
*   **Status Laranja/Vermelho**: O app está tentando conectar ou houve uma falha de rede. Verifique sua internet.

---

## 🛠 Solução de Problemas

### "O arquivo config.yaml não foi encontrado"
Certifique-se de que o arquivo `config.yaml` tem exatamente esse nome e está na mesma pasta do aplicativo. Se você moveu o app para a pasta `Aplicativos`, mova o `config.yaml` para lá também.

### "O aplicativo não abre ou fecha sozinho"
Tente mover a pasta do projeto para um local diferente (como a pasta `Documentos`) e repita o processo de clicar com o botão direito para abrir.

