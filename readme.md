# Cirrus

![Status](https://img.shields.io/badge/Status-In%20Development-orange)
![CI](https://github.com/toryjarvis/cirrus/actions/workflows/ci.yml/badge.svg)
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-20232A?style=flat&logo=react&logoColor=61DAFB)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=flat&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-DC382D?style=flat&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-007ACC?style=flat&logo=typescript&logoColor=white)
![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?style=flat&logo=kubernetes&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=flat&logo=github-actions&logoColor=white)

**Cirrus** is a general purpose link management platform with detailed click analytics, workspace organization, and a high-throughput redirect pipeline. Built with a service-oriented Go backend, a Redis-buffered event system, and a React dashboard for real-time analytics.

## Features

- **Link Management**: Create, organize, and manage shortened links across workspaces
- **Click Analytics**: Real-time tracking of clicks, referrers, geography, and device data
- **Fast Redirects**: Redis-buffered click events keep the redirect path non-blocking
- **Secure Authentication**: JWT-based auth with bcrypt password hashing
- **Workspace Organization**: Group and manage links by project or context

## Tech Stack

- **Frontend**: React, TypeScript
- **Backend**: Go (Fiber)
- **Database**: PostgreSQL
- **Cache / Buffer**: Redis
- **Authentication**: JWT
- **Infrastructure**: Docker, Kubernetes, GitHub Actions

## Current Phase: Foundations and Core Backend Development

Last Updated: May 2026

Sprint 2: Core link services set and other backend preparations

## Contributions

While this is a personal portfolio project, contributions are always welcome. Feel free to fork the repository, submit issues, or open pull requests.

### About The Author

**Victor "Tory" Jarvis**  
Full Stack Software Engineer | [GitHub](https://github.com/toryjarvis) | [LinkedIn](https://www.linkedin.com/in/victorjarvis)
