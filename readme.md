# 🌊 AquaScan API Project

<div align="center">
  <img src="https://placehold.co/600x300/0f172a/06b6d4?text=AquaScan+API" alt="Project Banner">
</div>

<div align="center">
  <br />
  <p>
    RESTful API built with <b>Golang</b> and <b>Gin Framework</b> for authentication, profile management, and file upload system.  
    Designed to be clean, consistent, and production-ready.
  </p>
</div>

<!-- Badges -->
<div align="center">
  <img src="https://img.shields.io/badge/Golang-1.22%2B-blue?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Gin-Framework-green?style=for-the-badge&logo=gin" alt="Gin">
  <img src="https://img.shields.io/badge/JWT-Authentication-orange?style=for-the-badge&logo=jsonwebtokens" alt="JWT">
  <img src="https://img.shields.io/badge/License-MIT-purple?style=for-the-badge" alt="License">
</div>

---

## ✨ Features

- 🔐 **JWT Authentication** – secure login with token-based system.
- 👤 **Profile Endpoint** – protected route accessible only with valid token.
- 📂 **File Upload** – upload files to `uploads/temp/`.
- 📡 **RESTful API** – consistent response format.

---

## 🛠️ Built With

- [Go](https://go.dev/) `v1.22+`
- [Gin](https://gin-gonic.com/) (Web Framework)
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) (Password Hashing)
- [JWT](https://jwt.io/) (Token Authentication)

---

## 🚀 Getting Started

### Prerequisites

- Install **Go 1.22+**
- Install **Git**

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/kakonoomoidee/aquascan-api.git
    cd aquascan-api
    ```

2. Install dependencies:

    ```bash
    go mod tidy
    ```

3. Run the server:

    ```bash
    go run main.go
    ```

The API will run at:  
👉 `http://localhost:8080`

---

## 📖 API Documentation

### 1. **Login**

- **Endpoint:** `POST /login`
- **Description:** Authenticate user and receive JWT token.
- **Request:**
    ```json
    {
      "email": "kaisar@kerajaan.com",
      "password": "password123"
    }
    ```
- **Response:**
    ```json
    {
      "status": "success",
      "message": "Login berhasil",
      "data": {
        "token": "<jwt_token>"
      }
    }
    ```

---

### 2. **Profile**

- **Endpoint:** `GET /profile`
- **Description:** Get user profile (requires JWT token).
- **Headers:**
    ```
    Authorization: Bearer <jwt_token>
    ```
- **Response:**
    ```json
    {
      "status": "success",
      "message": "Ini adalah halaman profile rahasia",
      "data": {
        "user_id": 1
      }
    }
    ```

---

### 3. **Upload File**

- **Endpoint:** `POST /upload`
- **Description:** Upload file to `uploads/temp/`.
- **Headers:**
    ```
    Authorization: Bearer <jwt_token>
    ```
- **Form Data:**
    ```
    file: <your_file>
    ```
- **Response:**
    ```json
    {
      "status": "success",
      "message": "File berhasil diupload",
      "data": {
        "file_path": "uploads/temp/example.png"
      }
    }
    ```

---

## 🎮 Testing with Postman

1. **Login** – request `POST /login` → copy JWT token.
2. **Profile** – request `GET /profile` → add header:
3. **Upload** – request `POST /upload` → select `form-data`, key: `file`, type: `File`.

---

## 📂 Project Structure

serber_aquascan/
├── controllers/ # Request handlers
├── middleware/ # JWT authentication
├── routes/ # API routes
├── services/ # Business logic (JWT, etc.)
├── utils/ # Helpers (response, etc.)
├── uploads/temp/ # Uploaded files
└── main.go # Entry point


---

## 📜 License

Distributed under the MIT License.  
See `LICENSE` for more information.

---

## 📬 Contact

Your Name - [@Kakonoomoidee](https://github.com/kakonoomoidee) - kakonoomoidee@gmail.com  

Project Link: [https://github.com/kakonoomoidee/aquascan-api](https://github.com/kakonoomoidee/aquascan-api)