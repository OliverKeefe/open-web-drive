<h1 align="left">Open Web Drive</h1>

<p align="left"><img alt="License" src="https://img.shields.io/badge/license-MIT-blue.svg" /> 
<img alt="backend" src="https://img.shields.io/github/actions/workflow/status/OliverKeefe/gestalto/backend-ci.yml?branch=main&label=backend">
<img alt="frontend" src="https://img.shields.io/github/actions/workflow/status/OliverKeefe/gestalto/frontend-ci.yml?branch=main&label=frontend">
</p>

## 🛠️ About This Project

**Open Web Drive** is a self-hostable, open-source cloud storage solution that combines the familiar user experience of mainstream cloud platforms with total ownership over your data and infrastructure.
Built for privacy-conscious individuals and homelab enthusiasts, it provides file management, client-side encryption, and integration with your preferred OAuth2 / OIDC identity provider.

<p align="center">
  <img src="docs/media/readme/OWA-Screenshot.png" alt="Screenshot of Gestalt Web Storage" width="45%">
</p>

## 🚀 Quick Start

> [!NOTE]
> Client-Side Encryption and IPFS Integration are currently undergoing a full re-write.

### Prerequisites (Linux / Ubuntu / Debian)

#### AWS CLI & LocalStack
*Skip this step if deploying to another, non-aws cloud provider.*

```shell
# Update package list and install dependencies
sudo apt update

# Install AWS CLI
curl "[https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip](https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip)" -o "awscliv2.zip"
unzip awscliv2.zip && sudo ./aws/install

# Install LocalStack CLI
curl -fsSL [https://github.com/localstack/lstk/releases/download/v0.19.0/lstk_0.19.0_linux_amd64.tar.gz](https://github.com/localstack/lstk/releases/download/v0.19.0/lstk_0.19.0_linux_amd64.tar.gz) | tar -xz
sudo mv lstk /usr/local/bin/

# Verify Installation
aws --version
localstack --version
```

**Minkube**
```shell
curl -LO [https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64](https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64)
sudo install minikube-linux-amd64 /usr/local/bin/minikube && rm minikube-linux-amd64

# Start Minikube
# If minkube fails to start something went wrong, see https://minikube.sigs.k8s.io/docs/start
minikube start 
```

**Install Terraform**

```shell
wget -O - https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(grep -oP '(?<=UBUNTU_CODENAME=).*' /etc/os-release || lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install terraform
```

**Install Docker**
```shell
# Add Docker's official GPG key:
sudo apt update
sudo apt install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/docker.asc
EOF

sudo apt update
```

```shell
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin 
sudo systemctl status docker
```

**Install Go / NPM dependencies and Build Dev Environment**

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