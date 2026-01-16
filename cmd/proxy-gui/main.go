package main

import (
	"context"
	"fmt"
	"image/color"
	"io"
	"log"
	"net"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/hashicorp/yamux"
	"gopkg.in/yaml.v3"
)

type Config struct {
	RelayHost         string        `yaml:"relay_host"`
	RelayPort         string        `yaml:"relay_port"`
	KeepAliveInterval time.Duration `yaml:"keepalive_interval"`
	KeepAliveTimeout  time.Duration `yaml:"keepalive_timeout"`
	AutoConnect       bool          `yaml:"auto_connect"`
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	// Valores mais conservadores para evitar o erro de deadline reached
	if cfg.KeepAliveInterval == 0 {
		cfg.KeepAliveInterval = 25 * time.Second
	}
	if cfg.KeepAliveTimeout == 0 {
		cfg.KeepAliveTimeout = 60 * time.Second
	}
	// Habilitar auto_connect por padrão se não estiver no yaml
	dataMap := make(map[string]interface{})
	yaml.Unmarshal(data, &dataMap)
	if _, ok := dataMap["auto_connect"]; !ok {
		cfg.AutoConnect = true
	}
	return &cfg, nil
}

var (
	statusColorDisconnected = color.NRGBA{R: 255, G: 165, B: 0, A: 255} // Laranja
	statusColorConnected    = color.NRGBA{R: 0, G: 255, B: 0, A: 255}   // Verde
)

func main() {
	// Configurar log em arquivo
	logFile, err := os.OpenFile("proxy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	}
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	cfg, err := loadConfig()
	if err != nil {
		log.Printf("Erro ao carregar configuração: %v. Usando padrões.", err)
		cfg = &Config{
			RelayHost:         "1.2.3.4",
			RelayPort:         "8080",
			KeepAliveInterval: 10 * time.Second,
			KeepAliveTimeout:  30 * time.Second,
		}
	}

	myApp := app.New()
	myApp.SetIcon(resourceIconPng)
	myWindow := myApp.NewWindow("Proxy Manager")
	myWindow.Resize(fyne.NewSize(300, 180))

	// Componentes visuais do Ícone e Status
	img := canvas.NewImageFromResource(resourceIconPng)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(36, 36)) // Ícone interno menor (reduzido ~50%)

	// Indicador visual (LED) - Agora proporcionalmente maior (40-50% do ícone)
	dot := canvas.NewCircle(statusColorDisconnected)
	dot.StrokeColor = color.NRGBA{255, 255, 255, 200}
	dot.StrokeWidth = 2
	dotBox := container.NewGridWrap(fyne.NewSize(18, 18), dot) // Bolinha maior em relação ao ícone

	// Stack para colocar o LED no canto do ícone
	dotBox.Move(fyne.NewPos(22, 22)) // Posicionamento proporcional
	iconWithStatus := container.NewStack(
		img,
		container.NewWithoutLayout(dotBox),
	)

	statusText := widget.NewLabel("Status: Desconectado")
	statusText.TextStyle = fyne.TextStyle{Bold: true}

	var cancel context.CancelFunc
	var ctx context.Context

	connectBtn := widget.NewButton("Conectar", nil)
	connectBtn.Importance = widget.HighImportance

	connectBtn.OnTapped = func() {
		if connectBtn.Text == "Conectar" {
			addr := fmt.Sprintf("%s:%s", cfg.RelayHost, cfg.RelayPort)
			connectBtn.SetText("Desconectar")
			statusText.SetText("Status: Conectando...")
			log.Printf("Iniciando conexão com %s", addr)

			ctx, cancel = context.WithCancel(context.Background())

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						err := runTunnelWithStatus(ctx, cfg, func(connected bool) {
							fyne.Do(func() {
								if connected {
									statusText.SetText("Status: Conectado")
									dot.FillColor = statusColorConnected
								} else {
									statusText.SetText("Status: Desconectado")
									dot.FillColor = statusColorDisconnected
								}
								dot.Refresh()
								statusText.Refresh()
							})
						})
						if err != nil {
							if ctx.Err() != nil {
								return // Cancelado pelo usuário
							}
							log.Printf("Conexão falhou: %v. Tentando novamente...", err)
							time.Sleep(3 * time.Second)
						}
					}
				}
			}()
		} else {
			if cancel != nil {
				cancel()
			}
			connectBtn.SetText("Conectar")
			statusText.SetText("Status: Desconectado")
			dot.FillColor = statusColorDisconnected
			dot.Refresh()
			log.Println("Conexão interrompida pelo usuário")
		}
	}

	// Layout minimalista e centralizado
	header := widget.NewLabelWithStyle("Serviço de Proxy", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		container.NewPadded(connectBtn),
		layout.NewSpacer(),
		container.NewCenter(
			container.NewVBox(
				container.NewCenter(iconWithStatus),
				container.NewCenter(statusText),
			),
		),
		layout.NewSpacer(),
	)

	myWindow.SetContent(content)

	// System Tray
	if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("Proxy Manager",
			fyne.NewMenuItem("Abrir", func() {
				myWindow.Show()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Sair", func() {
				myApp.Quit()
			}),
		)
		desk.SetSystemTrayMenu(m)
	}

	myWindow.SetCloseIntercept(func() {
		myWindow.Hide()
	})

	// Conexão automática ao iniciar
	if cfg.AutoConnect {
		go func() {
			// Pequeno delay para garantir que a UI esteja pronta
			time.Sleep(100 * time.Millisecond)
			fyne.Do(func() {
				connectBtn.OnTapped()
			})
		}()
	}

	myWindow.ShowAndRun()
}

func runTunnelWithStatus(ctx context.Context, cfg *Config, onStatusChange func(bool)) error {
	addr := fmt.Sprintf("%s:%s", cfg.RelayHost, cfg.RelayPort)
	dialer := net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		onStatusChange(false)
		return err
	}
	defer conn.Close()

	// Configurar Yamux com KeepAlive para evitar timeout de I/O
	yamuxCfg := yamux.DefaultConfig()
	yamuxCfg.KeepAliveInterval = cfg.KeepAliveInterval
	yamuxCfg.ConnectionWriteTimeout = cfg.KeepAliveTimeout // Tempo limite para o write do keepalive

	session, err := yamux.Client(conn, yamuxCfg)
	if err != nil {
		onStatusChange(false)
		return err
	}
	defer session.Close()

	onStatusChange(true)

	// Monitorar cancelamento do contexto
	go func() {
		<-ctx.Done()
		session.Close()
	}()

	for {
		stream, err := session.Accept()
		if err != nil {
			onStatusChange(false)
			return err
		}
		go handleSOCKS5(stream)
	}
}

func handleSOCKS5(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 256)
	if _, err := io.ReadFull(conn, buf[:2]); err != nil || buf[0] != 0x05 {
		return
	}
	nMethods := int(buf[1])
	if _, err := io.ReadFull(conn, buf[:nMethods]); err != nil {
		return
	}
	conn.Write([]byte{0x05, 0x00})
	if _, err := io.ReadFull(conn, buf[:4]); err != nil || buf[1] != 0x01 {
		return
	}
	var target string
	switch buf[3] {
	case 0x01:
		if _, err := io.ReadFull(conn, buf[:4]); err != nil {
			return
		}
		target = net.IP(buf[:4]).String()
	case 0x03:
		if _, err := io.ReadFull(conn, buf[:1]); err != nil {
			return
		}
		length := int(buf[0])
		if _, err := io.ReadFull(conn, buf[:length]); err != nil {
			return
		}
		target = string(buf[:length])
	case 0x04:
		if _, err := io.ReadFull(conn, buf[:16]); err != nil {
			return
		}
		target = net.IP(buf[:16]).String()
	}
	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		return
	}
	targetPort := (uint16(buf[0]) << 8) | uint16(buf[1])
	targetAddr := net.JoinHostPort(target, fmt.Sprintf("%d", targetPort))
	targetConn, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
	if err != nil {
		conn.Write([]byte{0x05, 0x01, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	}
	defer targetConn.Close()
	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	done := make(chan bool, 2)
	go func() { io.Copy(targetConn, conn); done <- true }()
	go func() { io.Copy(conn, targetConn); done <- true }()
	<-done
}
