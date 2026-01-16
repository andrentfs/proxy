# 🚀 K8s-Egress-Relay: Secure Navigation and Residential IPs for your Crawlers

This project was born out of a real-world need: **How to make crawlers running on Kubernetes navigate using my local (residential) IP address to avoid blocks and captchas from sensitive websites?**

The solution is a high-performance **Reverse SOCKS5 Tunnel**, built in Go, which allows outbound traffic (egress) from any Pod in K8s to be routed through your physical machine.

---

## 🛠 What does it do?

Imagine you have a high-traffic crawler running on a GKE, AWS, or Azure cluster. Frequently, datacenter IPs are flagged as "suspicious" by many websites.

With **K8s-Egress-Relay**, you open a bridge:
1. The **Relay** runs in your Kubernetes cluster as a SOCKS5 Proxy.
2. The **Local Agent (GUI)** runs on your personal computer (Windows/Mac/Linux).
3. A persistent connection (via Yamux protocol) is established between the two.
4. Your Pods in K8s start "seeing" the internet through the eyes (and IP) of your home or office.

---

## ✨ Project Highlights

- **Premium Graphical Interface**: Developed with Fyne, focused on minimalist UX with dynamic status indicators.
- **Auto-Connect**: Open the app and you're ready; the tunnel is established automatically.
- **Network Resilience**: Custom KeepAlive implementation to prevent drops on unstable connections.
- **YAML Configuration**: Simple setup—change IPs and ports in a straightforward file.
- **Dynamic Session**: The Relay server manages reconnections transparently without dropping active requests from your Pods.

---

## 🛠 Compilation and Build (Multi-platform)

This project uses `fyne-cross` to easily create binaries for different operating systems.

### Prerequisites
- Go installed
- Docker running (required for `fyne-cross`)

### Build Steps
1. Install fyne-cross:
```bash
go install github.com/fyne-io/fyne-cross@latest
```

2. Generate the builds:

* **Windows**:
```bash
fyne-cross windows -app-id com.proxy.manager -arch amd64 ./cmd/proxy-gui/
```

* **Linux**:
```bash
fyne-cross linux -app-id com.proxy.manager -arch amd64 ./cmd/proxy-gui/
```

* **macOS** (If on a Mac):
```bash
fyne package -os darwin --appID com.proxy.manager --icon Icon.png --src ./cmd/proxy-gui/
```

The executables will be generated in the `fyne-cross/bin/` folder (or root for macOS).

> **🍏 Mac End-User Guide**: We've prepared a simplified guide for you to send along with the app to macOS users. See: [README-MAC-USER-EN.md](./README-MAC-USER-EN.md).

---

## 🔌 How to make your Pod use the Proxy?

To have your crawler or any application within a Pod use the tunnel, you must configure the proxy environment variables using the `socks5://` protocol.

Example configuration in your `deployment.yaml` or `pod.yaml`:

```yaml
spec:
  containers:
  - name: your-crawler
    image: your-image:latest
    env:
    - name: HTTP_PROXY
      value: "socks5://egress-proxy.my-namespace.svc.cluster.local:1080"
    - name: HTTPS_PROXY
      value: "socks5://egress-proxy.my-namespace.svc.cluster.local:1080"
    - name: ALL_PROXY
      value: "socks5://egress-proxy.my-namespace.svc.cluster.local:1080"
```

> **Note**: The address `egress-proxy.my-namespace.svc.cluster.local` should be adjusted to the service name and namespace you defined in the Relay manifests.

---

## 🏗 Architecture

The project uses the **Yamux** protocol to multiplex connections over a single TCP channel. This means that even if you are browsing hundreds of websites simultaneously through the tunnel, only one physical connection is maintained, ensuring low latency and resource efficiency.

---

## 🤝 Contribute!

This project is open-source and focused on the data engineering and web scraping community.
- Missing a feature? Open an Issue!
- Want to improve stability? Submit a Pull Request!

---

## 💡 About the Author

I developed this tool to solve real scalability issues in data extraction, where IP reputation is the critical success factor. If this helped you, say hi on [LinkedIn](https://www.linkedin.com/in/andrentfs/)!

---

*Made with ❤️ and Go.*
