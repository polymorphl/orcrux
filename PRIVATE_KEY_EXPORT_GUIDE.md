# 🔐 **Private Key Export Guide for GitHub Actions**

## **The Problem**
The automatic export script couldn't export your private key because macOS protects private keys by default. This is why the code signing isn't working in GitHub Actions.

## **The Solution**
You need to manually export the private key from Keychain Access.

## **Step-by-Step Instructions**

### **1. Open Keychain Access**
- Press `Cmd + Space` and type "Keychain Access"
- Or go to `Applications > Utilities > Keychain Access`

### **2. Find Your Certificate**
- In the left sidebar, click on "login" keychain
- Look for "Apple Development: lucterracher@gmail.com (7GX2PQQM2W)"
- Double-click on it

### **3. Export the Private Key**
- In the certificate window, click the **"Access Control"** tab
- Look for the **private key** entry (it should show the same name)
- Right-click on the private key → **Export**
- Choose a location and save it as `private_key.p12`
- **Important**: When prompted for a password, leave it blank (or set a simple one)

### **4. Convert to PEM Format**
Open Terminal and run:
```bash
# Convert .p12 to .pem (you'll be prompted for the export password)
openssl pkcs12 -in private_key.p12 -out private_key.pem -nodes

# If you set a password, use: openssl pkcs12 -in private_key.p12 -out private_key.pem -nodes -passin pass:YOUR_PASSWORD
```

### **5. Update GitHub Secrets**
1. Go to your GitHub repository
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Edit the `MACOS_PRIVATE_KEY` secret
4. Copy the **entire content** of `private_key.pem` (including the BEGIN and END lines)
5. Save the secret

### **6. Test the Workflow**
Create a new release tag to test:
```bash
./scripts/version.sh bump patch
git push --tags
```

## **What This Will Give You**

✅ **Real Apple Development certificate** for code signing  
✅ **Much better security warnings** for users  
✅ **Professional appearance** for releases  
✅ **Better user experience**  

## **Alternative: Use Current Setup**
If you prefer not to export the private key, the current workflow will:
- Try to use your certificate (will fail without private key)
- Fall back to unsigned app
- Users can right-click to open (standard for open-source apps)

## **Security Note**
- The private key is only used during GitHub Actions builds
- It's never stored in your repository
- It's encrypted in GitHub's secrets system
- This is the standard approach for automated code signing

## **Need Help?**
If you encounter issues:
1. Check that the certificate name matches exactly
2. Ensure you're exporting from the correct keychain
3. Try exporting without a password first
4. Verify the .pem file contains both certificate and private key
