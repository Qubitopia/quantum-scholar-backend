# QuantumScholar
**QuantumScholar** is an online Exam portal with AI powerd proctering

<p align="center">
   <img src="https://raw.githubusercontent.com/Qubitopia/quantum-scholar-web/refs/heads/main/src/assets/Qubitopia-4096x2048.png" alt="Qubitopia Logo" width="400"/>
</p>

> "Empowering secure, fair, and scalable online assessments with AI."

<p align="center">
   <img src="https://img.shields.io/badge/License-AGPL%20v3-blue.svg" alt="AGPL v3 License"/>
</p>

---

## Technologies Used

<table>
   <thead>
      <tr>
         <th scope="col">Category</th>
         <th scope="col">Technology</th>
         <th scope="col">Why ?</th>
      </tr>
   </thead>
   <tbody>
      <tr>
         <td>Programming Language</td>
         <td>Go</td>
         <td>Multi-architecture support and static binaries</td>
      </tr>
      <tr>
         <td rowspan="4">Frameworks and Libraries</td>
         <td>Gin</td>
         <td>HTTP web framework</td>
      </tr>
      <tr>
         <td>GORM</td>
         <td>Type-safe ORM</td>
      </tr>
      <tr>
         <td>JWT-Go</td>
         <td>JSON Web Token authentication</td>
      </tr>
      <tr>
         <td>AWS SDK for Go</td>
         <td>Cloudflare R2 compatibility</td>
      </tr>
      <tr>
         <td rowspan="3">Database</td>
         <td>PostgreSQL</td>
         <td>Primary relational database</td>
      </tr>
      <tr>
         <td>Redis</td>
         <td>Rate limiting and caching</td>
      </tr>
      <tr>
         <td>Cloudflare R2</td>
         <td>Object storage for images and video</td>
      </tr>
      <tr>
         <td rowspan="2">Containerization</td>
         <td>Docker</td>
         <td>Container runtime</td>
      </tr>
      <tr>
         <td>Docker Compose</td>
         <td>For easy testing</td>
      </tr>
      <tr>
         <td rowspan="3">CI/CD</td>
         <td>GitHub Actions</td>
         <td>Continuous integration</td>
      </tr>
      <tr>
         <td>GitHub Packages / Container Registry</td>
         <td>Package and image hosting</td>
      </tr>
      <tr>
         <td>Trivy</td>
         <td>Filesystem and container vulnerability scanning</td>
      </tr>
      <tr>
         <td>Email Service</td>
         <td>Oracle Cloud Infrastructure (OCI) Email Delivery</td>
         <td>Affordable email delivery with a generous free tier</td>
      </tr>
   </tbody>
</table>


## Things Contributor [CHETAN-INGALE](https://github.com/CHETAN-INGALE/) worked on
### Backend
- Built high availability API services using Go and Gin framework.
- Implemented JWT-based authentication and authorization.
- Designed and optimized PostgreSQL database schemas.
- Integrated Redis for caching and rate limiting.
- Implemented file storage using Cloudflare R2.
- Integrated payment gateway Razorpay for handling transactions.
- Email service integration using custom domain.

### DevSecOps
- Created Dockerfiles for containerizing the application with Multistage builds.
- Optimized Docker images for multi-architecture support (amd64 & arm64).
- Configured GitHub Actions for CI/CD pipelines.
- Implemented Trivy for filesystem and container vulnerability scanning, and uploaded the SARIF report to GitHub Security.
- Build and published Docker images to [GitHub Container Registry](https://github.com/Qubitopia/quantum-scholar-backend/pkgs/container/quantum-scholar-backend).
- Managed application secrets using GitHub Secrets & environment variables.
- Set up Docker Compose for local development and testing.


## How to Run this project
- Using pre-built container images (Easy and repoducable)
- Build locally (from source)

### Using Pre-built Container Images
1. Clone the repository:

   ```sh
   git clone https://github.com/Qubitopia/quantum-scholar-backend.git

   cd quantum-scholar-backend
   ```

2. Copy the environment template and update your configuration:

   ```sh
   cp template.env .env
   # Edit .env and fill in the required values
   ```

3. Update the `docker-compose.yml` to use the pre-built image from GitHub Container Registry:

   ```yaml
   services:
     gin:
       image: ghcr.io/qubitopia/quantum-scholar-backend:latest
   ```

4. Start the application using Docker Compose:

   ```sh
   docker compose up
   ```

### Build Locally with dependency on Docker
1. Clone the repository:

   ```sh
   git clone https://github.com/Qubitopia/quantum-scholar-backend.git

   cd quantum-scholar-backend
   ```

2. Copy the environment template and update your configuration:

   ```sh
   cp template.env .env
   # Edit .env and fill in the required values
   ```

3. Start the application using Docker Compose:

   ```sh
   docker compose up
   # This will build the Docker image locally and start the application with all dependencies.
   ```

Note: This project was  seperated from its monorepo [Qubitopia/QuantumScholar](https://github.com/Qubitopia/) for better segregation of website, Exam appplicaton and backend services.

Checkout my co-contributor [Araya Bhagat's](https://github.com/Aaru911) work on the [QuantumScholar Web](https://github.com/Qubitopia/quantum-scholar-web)