# QuantumScholar

<p align="center">
   <img src="assets/Qubitopia-4096x2048.png" alt="Qubitopia Logo" width="400"/>
</p>

<p align="center">
   <img src="https://img.shields.io/badge/License-AGPL%20v3-blue.svg" alt="AGPL v3 License"/>
</p>

<p align="center">
   <img src="https://img.shields.io/badge/AI%20Proctoring-Powered%20by%20QuantumScholar-blueviolet?style=for-the-badge&logo=quantconnect&logoColor=white" alt="AI Proctoring Badge"/>
   <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=for-the-badge" alt="Status Badge"/>
</p>

---

> "Empowering secure, fair, and scalable online assessments with AI."

## Overview

**QuantumScholar** is an AI-based proctoring website designed to ensure secure, fair, and scalable online assessments. Leveraging advanced artificial intelligence, it provides real-time monitoring, identity verification, and automated analysis to detect suspicious activities during online exams.

## Features

- AI-powered live proctoring
- Automated identity verification
- Real-time cheating detection
- Secure user authentication
- Scalable for institutions and organizations
- Detailed reporting and analytics

## Getting Started

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

3. Build the image

   ```sh
   # if using powershell
   build.ps1 myapp my-qs-backend latest

   # if using bash
   chmod +x build.sh
   APP_NAME=myapp IMAGE_REPO=my-qs-backend VERSION_TAG=latest build.sh
   ```

4. Start the application using Docker Compose:

   ```sh
   # update the image for gin to use my-qs-backend

   docker compose up
   ```
