<h1 align="left">Open Web Drive</h1>

<p align="left"> 

</p>

## 🛠️ About This Project

<p align="left">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go">
  <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-5.8-3178C6?logo=typescript">
  <img alt="backend" src="https://img.shields.io/github/actions/workflow/status/OliverKeefe/open-web-drive/backend-ci.yml?branch=main&label=backend">
  <img alt="frontend" src="https://img.shields.io/github/actions/workflow/status/OliverKeefe/open-web-drive/frontend-ci.yml?branch=main&label=frontend">
  <img alt="License" src="https://img.shields.io/badge/license-MIT-white.svg">
</p>

**Open Web Drive** is a self-hostable, open-source file storage solution that provides total ownership of your data.
Built for privacy-conscious individuals, it provides file management, client-side encryption, OAuth2.0 / OIDC and IPFS integration.

**Why Did I Build This?** Aside from an excuse to finally learn Go and React.js, I mostly just wanted to be able to host 
a file storage app like OneDrive or Google Drive where I could host it on my own infrastructure; be that a deployment to a public cloud / hyperscalar tenant 
that I rent or an OpenStack / Bare Metal home lab.


<h2 align="center">UI Screenshots</h2>
<table>
  <tr>
    <td width="50%"><img src="docs/media/readme/Screenshot from 2026-08-05 17-01-55.png" width="100%" alt="Image 1"></td>
    <td width="50%"><img src="docs/media/readme/Screenshot from 2026-08-05 17-02-57.png" width="100%" alt="Image 2"></td>
  </tr>
  <tr>
    <td><img src="docs/media/readme/Screenshot from 2026-08-05 22-06-36.png" width="100%" alt="Image 3"></td>
    <td><img src="docs/media/readme/Screenshot from 2026-08-05 17-02-40.png" width="100%" alt="Image 4"></td>
  </tr>
</table>

> [!NOTE]
> 🏗️ Auth / Client-Side Encryption are currently undergoing a full re-write.

## 🚀 Quick Start
**📋 Prerequisites**

**[Terraform](https://developer.hashicorp.com/terraform/install)**

**[Minikube](https://minikube.sigs.k8s.io/docs/start/?arch=%2Flinux%2Fx86-64%2Fstable%2Fbinary+download)**

**[Docker](https://docs.docker.com/engine/install/)**

**[Go](https://go.dev/dl/)**

**[Node.js](https://nodejs.org/en/download)**

**[AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)**

**[LocalStack](https://docs.localstack.cloud/aws/getting-started/installation/)**

**Clone the repository**
```shell
git clone https://github.com/OliverKeefe/open-web-drive.git 
```

Once you've cloned the repository, change your working directory to ~/open-web-drive and
execute the following setup and dev build scripts.
**Setup Dev Environment**
```shell
chmod +x /scripts/setup.sh && ./scripts/setup.sh
```

**Run in Dev Mode**
```shell
chmod +x ./scripts/build.sh &&./scripts/build.sh
```

📄 License
Distributed under the MIT License. See `LICENSE.md` for more information.
