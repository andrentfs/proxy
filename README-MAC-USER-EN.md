# 🍏 User Guide: Proxy Manager for macOS

To use the **Proxy Manager** on your Mac, follow this guide to set up and run the service correctly.

## 📦 What you need
You should have received two items:
1.  **`proxy-gui.app`**: The main application.
2.  **`config.yaml`**: The configuration file (contains the server address).

---

## 🚀 Step-by-Step to Start

### 1. Organization
Place the `proxy-gui.app` application and the `config.yaml` file in the **same folder** (for example, on your Desktop or in the Documents folder).

### 2. Bypassing the Security Warning (Gatekeeper)
Since this is specialized software, macOS may show a warning saying it "cannot verify the developer." **Don't worry**, this is normal for custom apps.

To open for the first time:
1.  **Right-click** (or hold the `Control` key and click) on the **Proxy Manager** icon.
2.  Select the **Open** option from the menu that appears.
3.  A confirmation window will pop up. Click the **Open** button.
4.  *In the future, you can open it with a double-click as usual.*

### 3. Using the Application
*   **Auto-Connect**: As soon as the app opens, it will try to connect to the server. You will see a status dot change color.
*   **Green Status**: Means you are successfully connected! The tunnel is active.
*   **Orange/Red Status**: The app is trying to connect or there was a network failure. Check your internet connection.

---

## 🛠 Troubleshooting

### "The config.yaml file was not found"
Ensure that the `config.yaml` file is named exactly like that and is in the same folder as the application. If you moved the app to the `Applications` folder, move the `config.yaml` there as well.

### "The application does not open or closes by itself"
Try moving the project folder to a different location (like the `Documents` folder) and repeat the right-click process to open.

---
