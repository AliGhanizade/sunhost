# SunHost

SunHost is a university web development project created for the Web Design course at Sajjad University.

The project demonstrates how a simple hosting company website can be combined with a Go backend to provide authentication, system monitoring, and an administration panel.

---

## Features

- User Registration
- User Login
- User Activity Logs
- Admin Dashboard
- Live System Monitoring
- System Summary API
- Static Website Pages
- SQLite Database

---

## Technologies

- Go
- Gin Framework
- SQLite
- HTML5
- CSS3
- JavaScript

---

## Project Structure

```
config/
controllers/
model/
public/
router/
main.go
```

---

## Installation

Clone the repository

```bash
git clone https://github.com/YourUsername/sunhost.git
```

Install dependencies

```bash
go mod download
```

Run the project

```bash
go run .
```

The application will be available at

```
http://localhost:8080
```

---

## Available Pages

```
/
├── /host
├── /host/products
├── /host/blog
├── /host/about
├── /host/contact
├── /host/login
├── /panel
```

---

## API

### Authentication

```
POST /api/users/register
POST /api/users/login
```

### Logs

```
GET /api/users/logs
```

### System

```
GET /api/system/summary
GET /api/system/stream
```

---

## Educational Purpose

This project was developed as the final project for the **Web Design** course at **Sajjad University**.

The goal was to practice:

- Backend development using Go
- REST API development
- Authentication
- SQLite database integration
- MVC-style project organization
- Basic web interface development

---

## Project Status

> **Note**
>
> This repository contains one of my earlier projects and is shared primarily for educational and portfolio purposes.
>
> The project is **not actively maintained** and was never intended to become a production-ready application.
>
> Some parts of the codebase may be unfinished, unoptimized, or not follow the best coding practices. Documentation is also incomplete, and the commit history does not reflect the actual development process.
>
> I decided to publish this project on GitHub to document my learning journey and preserve my previous work rather than leave it archived on my local machine.
>
> Feedback and suggestions are always welcome.

---

## License

This project is intended for educational purposes.


---

This repository represents where I was, not where I am.