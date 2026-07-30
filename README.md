<h1 align="left">Open Web Drive</h1>

<p align="left"> 
<img alt="backend" src="https://img.shields.io/github/actions/workflow/status/OliverKeefe/open-web-drive/backend-ci.yml?branch=main&label=backend">
<img alt="frontend" src="https://img.shields.io/github/actions/workflow/status/OliverKeefe/open-web-drive/frontend-ci.yml?branch=main&label=frontend">
</p>

## 🛠️ About This Project

<p align="left">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go">
  <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-5.8-3178C6?logo=typescript">
  <img alt="React" src="https://img.shields.io/badge/React-19-61DAFB?logo=react">
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-18-4169E1?logo=postgresql">
  <img alt="Terraform" src="https://img.shields.io/badge/Terraform-1.15-844FBA?logo=terraform">
  <img alt="Kubernetes" src="https://img.shields.io/badge/Kubernetes-123?logo=kubernetes">
  <img alt="License" src="https://img.shields.io/badge/license-MIT-blue.svg">
</p>

**Open Web Drive** is a self-hostable, open-source cloud storage solution that combines the familiar user experience of mainstream cloud platforms with total ownership over your data and infrastructure.
Built for privacy-conscious individuals and homelab enthusiasts, it provides file management, client-side encryption, and integration with your preferred OAuth2 / OIDC identity provider.

<p align="center">
  <img src="docs/media/readme/OWA-Screenshot.png" alt="App Screenshot" width="55%">
</p>

## 🚀 Quick Start

> [!NOTE]
> Auth / Client-Side Encryption are currently undergoing a full re-write.


```shell
git clone https://github.com/OliverKeefe/open-web-drive.git \
cd open-web-drive

# Install Frontend dependencies
cd frontend && npm install && cd ..

# Download Backend dependencies
cd backend && go mod download && cd ..

# Run local build script
chmod +x scripts/build.sh
./scripts/build.sh -d
```

📄 License
Distributed under the MIT License. See `LICENSE.md` for more information.