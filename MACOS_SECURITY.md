# macOS Security & Code Signing Guide

This guide explains how to handle the macOS security warning "cannot verify it doesn't contain any malware" and how to properly sign your Orcrux app.

## 🚨 **The Security Warning Explained**

### **Why This Happens**
When you try to open the DMG file, macOS shows:
> "Orcrux cannot be opened because it is from an unidentified developer"

This is **macOS Gatekeeper** protecting you from potentially malicious software.

### **Root Causes**
1. **No Code Signing**: App isn't signed with a developer certificate
2. **No Notarization**: Apple hasn't verified the app is safe
3. **Security Policy**: macOS blocks unsigned apps by default

## 🔧 **Solutions (In Order of Preference)**

### **1. Code Signing + Notarization (Professional)**
**Cost**: $99/year Apple Developer Account  
**Result**: No warnings, trusted by macOS

```bash
# Step 1: Code sign the app
./scripts/code-sign.sh sign 'Developer ID Application: Your Name'

# Step 2: Create DMG
make create-dmg

# Step 3: Notarize with Apple
./scripts/code-sign.sh notarize your@email.com app_password
```

### **2. Code Signing Only (Development)**
**Cost**: Free (with Xcode)  
**Result**: Reduced warnings, but still flagged

```bash
# Sign with your personal certificate
./scripts/code-sign.sh sign 'Mac Developer: Your Name'

# Create DMG
make create-dmg
```

### **3. Temporary Workaround (Testing)**
**Cost**: Free  
**Result**: Works for testing, not for distribution

```bash
# Right-click DMG → "Open" → "Open anyway"
# Or use Terminal:
sudo spctl --master-disable  # Disable Gatekeeper temporarily
```

## 🛠️ **How to Code Sign**

### **Prerequisites**
1. **Xcode Command Line Tools** (free)
   ```bash
   xcode-select --install
   ```

2. **Developer Certificate** (free with Xcode)
   ```bash
   # List available certificates
   ./scripts/code-sign.sh list
   ```

### **Step-by-Step Signing**

#### **Option A: Personal Development Certificate**
```bash
# 1. Build the app
make build-darwin

# 2. Sign with development certificate
./scripts/code-sign.sh sign 'Mac Developer: Your Name'

# 3. Create DMG
make create-dmg
```

#### **Option B: Apple Developer Certificate ($99/year)**
```bash
# 1. Build and sign
make build-darwin
./scripts/code-sign.sh sign 'Developer ID Application: Your Name'

# 2. Create DMG
make create-dmg

# 3. Notarize with Apple (modern tool)
./scripts/code-sign.sh notarize your@email.com app_password
```

## 🔍 **Available Commands**

### **Code Signing Script**
```bash
./scripts/code-sign.sh help                    # Show all commands
./scripts/code-sign.sh list                    # List available certificates
./scripts/code-sign.sh sign <identity>         # Sign the app
./scripts/code-sign.sh notarize <id> <pass>   # Submit for notarization
./scripts/code-sign.sh check <request_id>      # Check notarization status
```

### **Makefile Commands**
```bash
make sign-app      # Show code signing help
make notarize      # Show notarization help
make create-dmg    # Create DMG after signing
```

## 📱 **What Users See**

### **Unsigned App**
- ❌ **Warning**: "Cannot verify it doesn't contain any malware"
- ❌ **Blocked**: App won't open by default
- ⚠️ **Workaround**: Right-click → "Open" → "Open anyway"

### **Code Signed App**
- ⚠️ **Reduced warning**: "From an unidentified developer"
- ✅ **Can open**: With user permission
- 🔒 **More trusted**: Less scary for users

### **Notarized App**
- ✅ **No warning**: Opens normally
- ✅ **Fully trusted**: macOS recognizes it as safe
- 🎯 **Professional**: Ready for distribution

## 💰 **Cost Comparison**

| Solution | Cost | User Experience | Distribution |
|----------|------|-----------------|--------------|
| **Unsigned** | Free | ❌ Scary warnings | ❌ Not recommended |
| **Code Signed** | Free | ⚠️ Reduced warnings | ⚠️ Development only |
| **Notarized** | $99/year | ✅ No warnings | ✅ Production ready |

## 🚀 **Quick Start for Development**

### **1. Build and Sign (Free)**
```bash
# Build the app
make build-darwin

# Sign with development certificate
./scripts/code-sign.sh sign 'Mac Developer: Your Name'

# Create DMG
make create-dmg
```

### **2. Test the DMG**
- Double-click the DMG
- Accept the reduced security warning
- Drag app to Applications folder

### **3. For Distribution**
- Get Apple Developer Account ($99/year)
- Use `Developer ID Application` certificate
- Submit for notarization
- Users get no warnings

## 🔒 **Security Best Practices**

### **For Development**
- Use development certificates for testing
- Don't distribute unsigned apps
- Test with reduced security settings

### **For Production**
- Always code sign with Apple Developer certificate
- Submit for notarization
- Keep certificates secure
- Rotate certificates regularly

## 📚 **Additional Resources**

- [Apple Developer Documentation](https://developer.apple.com/support/code-signing/)
- [Code Signing Guide](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/Introduction/Introduction.html)
- [Modern Notarization Guide](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [notarytool Documentation](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution#use-notarytool)

---

**Remember**: Code signing is about user trust and security. Even free development certificates significantly improve the user experience! 🎯
