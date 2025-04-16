# 🛠️ Portal CLI

**Portal** is a powerful developer tool with two core components:

- **Parser**: Scans your codebase for custom annotations and extracts structured variables.
- **Patcher**: Spins up a secure web-based dashboard where non-technical users can update those variables and automatically push changes to a GitHub repository (via commit or pull request).

---

## 🚀 Features

- 🧠 **Code annotation parsing** via CLI
- 🌐 **Interactive dashboard** for variable editing
- 🔐 Built-in **authentication** for protected access
- 📤 Automatically creates commits or pull requests to your GitHub repo
- ⚙️ Configurable via CLI flags, environment variables, or config files

## 🛠️ Development Instructions

### 🔨 Build the Executable
To build the Portal CLI executable, run the following command:

```bash
go build -o portal ./cmd/
```

### ▶️ Run the Executable
After building, you can run the executable with:

```bash
./portal
```
