# 🔐 Orcrux

> **Secure Secret Sharing with Elegant UI**

Orcrux is a modern, cross-platform desktop application that implements Shamir's Secret Sharing algorithm with a beautiful, intuitive interface. Split your secrets into multiple shards or recompose them back together with enterprise-grade security and a delightful user experience.

![Orcrux Demo](https://img.shields.io/badge/Status-Production%20Ready-brightgreen)
![Platform](https://img.shields.io/badge/Platform-Cross%20Platform-blue)

<div align="center">

![Orcrux Demo](docs/demo.gif)

*Smooth animations, elegant transitions, and intuitive interface*

</div>

## ✨ Features

### 🔒 **Core Functionality**
- **Secret Splitting**: Divide secrets into multiple shards using Shamir's Secret Sharing
- **Secret Reconstruction**: Reconstruct original secrets from a subset of shards
- **Flexible Configuration**: Customize number of total shards and required shards
- **Multiple Output Formats**: Support for Base64 and Hexadecimal encoding

### 🎨 **User Experience**
- **Beautiful Interface**: Modern, crystal-themed design with smooth animations
- **Real-time Feedback**: Interactive sliders and immediate visual updates
- **Responsive Design**: Adapts to different screen sizes and orientations
- **Smooth Transitions**: Elegant color animations between different modes

### 🛡️ **Security & Reliability**
- **Cryptographically Secure**: Implements proven Shamir's Secret Sharing algorithm
- **Input Validation**: Comprehensive error checking and validation
- **File Operations**: Secure file import/export capabilities
- **Cross-platform**: Built with Go and Wails for native performance

## 🚀 Getting Started

### Prerequisites
- **Go 1.21+** - [Download here](https://golang.org/dl/)
- **Node.js 18+** - [Download here](https://nodejs.org/)
- **Wails CLI** - Install with `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/polymorphl/orcrux.git
   cd orcrux
   ```

2. **Install dependencies**
   ```bash
   # Install Go dependencies
   go mod download
   
   # Install frontend dependencies
   cd frontend
   npm install
   cd ..
   ```

3. **Build and run**
   ```bash
   # Development mode
   wails dev
   
   # Build production binary
   wails build
   ```

## 📖 Usage

### Splitting Secrets

1. **Navigate to the Split tab**
2. **Enter your secret** in the text area
3. **Configure shard parameters**:
   - Set total number of shards (2-255)
   - Set required shards for reconstruction
4. **Choose output format**: Base64 or Hexadecimal
5. **Click Split** to generate your shards
6. **Copy or save** the generated shards

### Reconstructing Secrets

1. **Navigate to the Bind tab**
2. **Add shards** using the + button
3. **Paste your shard data** into each input field
4. **Click Recompose** to reconstruct the secret
5. **View the result** in the output area

### File Operations

- **Import shards** from text files
- **Export results** to your local system
- **Drag and drop** support for easy file handling

## 🏗️ Architecture

### **Backend (Go)**
- **Shamir Implementation**: Custom Go package for secret sharing
- **File Operations**: Secure file handling and validation
- **Wails Integration**: Native desktop app framework

### **Frontend (React + TypeScript)**
- **Modern UI**: Built with React 18 and TypeScript
- **Styling**: Tailwind CSS with custom crystal theme
- **Animations**: Framer Motion for smooth transitions
- **Components**: Modular, reusable UI components

### **Key Technologies**
- **Go** - Backend logic and Shamir's algorithm
- **Wails** - Cross-platform desktop framework
- **React** - Frontend user interface
- **TypeScript** - Type-safe development
- **Tailwind CSS** - Utility-first styling
- **Framer Motion** - Smooth animations

## 🔧 Development

### Project Structure
```
orcrux/
├── app.go              # Main application logic
├── files.go            # File operations
├── main.go             # Entry point
├── shamir/             # Shamir's Secret Sharing implementation
├── frontend/           # React frontend application
│   ├── src/
│   │   ├── components/ # UI components
│   │   ├── lib/        # Utilities and helpers
│   │   └── types/      # TypeScript type definitions
│   └── package.json
└── README.md
```

### Running Tests
```bash
# Run all tests
go test -v

# Run specific test suites
go test -v ./shamir
go test -v ./files

# Run with coverage
go test -cover
```

### Building for Distribution
```bash
# Build for current platform
wails build

# Build for specific platforms
wails build -platform windows/amd64
wails build -platform darwin/universal
wails build -platform linux/amd64
```

## 🎯 Use Cases

### **Personal Security**
- **Password Management**: Split master passwords into multiple parts
- **Recovery Keys**: Distribute recovery keys across trusted contacts
- **Private Documents**: Secure storage of sensitive information

### **Enterprise Applications**
- **Key Management**: Secure distribution of cryptographic keys
- **Access Control**: Multi-party authorization systems
- **Compliance**: Meet regulatory requirements for secret sharing

### **Blockchain & Crypto**
- **Wallet Recovery**: Secure backup of cryptocurrency wallets
- **Private Keys**: Distributed storage of blockchain private keys
- **Multi-sig Wallets**: Enhanced security for digital assets

## 🤝 Contributing

We welcome contributions!

### **How to Contribute**
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### **Development Setup**
1. Follow the installation steps above
2. Make your changes
3. Run tests to ensure everything works
4. Submit your pull request


## 🙏 Acknowledgments

- **Shamir's Secret Sharing**: The cryptographic foundation
- **Wails Team**: For the excellent cross-platform framework

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/polymorphl/orcrux/issues)

---

<div align="center">

**Made with ❤️ and ☕ **

[![GitHub stars](https://img.shields.io/github/stars/polymorphl/orcrux?style=social)](https://github.com/polymorphl/orcrux)
[![GitHub forks](https://img.shields.io/github/forks/polymorphl/orcrux?style=social)](https://github.com/polymorphl/orcrux)
[![GitHub issues](https://img.shields.io/github/issues/polymorphl/orcrux)](https://github.com/polymorphl/orcrux/issues)

</div>
