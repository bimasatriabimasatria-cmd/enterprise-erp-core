# 🚀 Enterprise ERP Core API

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go)
![Fiber](https://img.shields.io/badge/Fiber-v2-172554?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-336791?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Containerized-2496ED?style=for-the-badge&logo=docker)

Sistem *Backend* ERP (Enterprise Resource Planning) berskala industri yang dirancang dengan arsitektur **Multi-Tenant**. Dibangun untuk memberikan performa tinggi, keamanan maksimal, dan isolasi data antar perusahaan.

## ✨ Fitur Utama

* **🏢 Arsitektur Multi-Tenant:** Isolasi data otomatis menggunakan Tenant ID. Data antar perusahaan tidak akan pernah bocor atau tercampur.
* **🔐 Keamanan Tingkat Lanjut:** Enkripsi *password* menggunakan `Bcrypt` dan otentikasi sesi menggunakan **JSON Web Tokens (JWT)**.
* **📚 Dokumentasi Interaktif:** Dilengkapi dengan antarmuka **Swagger UI** terintegrasi yang di-*generate* secara dinamis untuk memudahkan pengujian API.
* **🐳 Cloud-Ready & Dockerized:** Dibungkus dalam kontainer Alpine Linux super ringan, siap diorbitkan ke platform *Cloud* mana pun.
* **⚙️ CI/CD Pipeline:** Otomatisasi pengerahan (Deployment) dari GitHub ke platform *Cloud* secara *real-time*.

## 🛠️ Tech Stack

* **Bahasa:** Golang (Go)
* **Framework Web:** Go Fiber (Berdasarkan Fasthttp)
* **ORM:** GORM
* **Database:** PostgreSQL (Hosted on Supabase via Connection Pooler)
* **Dokumentasi API:** Swaggo / Swagger
* **Deployment:** Docker & Back4App Containers

## 🌐 Dokumentasi API (Live)

API ini telah mengudara dan dapat diakses publik. Silakan uji coba *endpoint* secara langsung melalui Swagger UI:

👉 **[Buka Halaman Swagger Live](#)** *(Ganti tanda # dengan URL .b4a.run/swagger/index.html milik Anda)*

*(Gunakan email `admin@enterprise.com` dan password `admin123` untuk mendapatkan Token JWT).*

## 📦 Menjalankan Sistem Secara Lokal

Jika Anda ingin menjalankan sistem ini di mesin lokal Anda, ikuti langkah berikut:

1. **Clone repositori:**
   ```bash
   git clone [https://github.com/USERNAME_ANDA/enterprise-erp-core.git](https://github.com/USERNAME_ANDA/enterprise-erp-core.git)