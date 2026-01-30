# Golang Backend Starter Kit

**Starter Kit** Backend REST API yang tangguh, bersih, dan modular. Dibangun dengan **Go**, **Gin**, dan **GORM**.
Proyek ini menyediakan fondasi rock-solid untuk memulai aplikasi apapun, lengkap dengan manajemen user, otentikasi JWT, RBAC, dan Logging terstruktur.

## 🚀 Fitur

- **Otentikasi Pengguna**: Registrasi dan Login aman menggunakan JWT.
- **Layanan Email**: Verifikasi Email dan alur Lupa/Reset Kata Sandi (SMTP).
- **Keamanan**:
  - `bcrypt` untuk hashing kata sandi.
  - Rate Limiting Middleware untuk mencegah penyalahgunaan.
  - CORS Middleware untuk pembagian sumber daya lintas asal.
  - **Kebijakan & Izin**: Sistem otorisasi berbasis kebijakan (Policy-based authorization).
- **Database**: Integrasi PostgreSQL menggunakan GORM ORM.
- **Pencarian & Paginasi Cerdas**: Pencarian Full-Text dan Like otomatis dengan respons terstandar.
- **Dokumentasi**: Dokumen API otomatis via Swagger.
- **Ketahanan**: Penanganan Error Global dan Logging Permintaan.
- **Hot Reload**: Pengembangan real-time menggunakan `air`.

## 🛠️ Tech Stack

- **Bahasa**: [Go](https://go.dev/) (1.24+)
- **Framework**: [Gin Web Framework](https://gin-gonic.com/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **ORM**: [GORM](https://gorm.io/)
- **Dokumentasi**: [Swagger](https://github.com/swaggo/swag)
- **Konfigurasi**: [Godotenv](https://github.com/joho/godotenv)
- **Email**: [Gomail.v2](https://github.com/go-gomail/gomail)

## 📂 Struktur Proyek

```bash
.
├── config/         # Konfigurasi dan Pemuatan Environment
├── controller/     # API Controllers (Handlers)
├── docs/           # File Dokumentasi Swagger
├── dto/            # Data Transfer Objects (Request/Response structs)
├── entity/         # Model Database (GORM Structs)
├── middleware/     # Auth, Logger, RateLimiter, CORS, Policy
├── migrations/     # Migrasi Database Otomatis
├── repository/     # Layer Akses Database (Data Access Object)
├── routes/         # Definisi Rute (Endpoints)
├── service/        # Layer Logika Bisnis
├── utils/          # Utilitas (Email, JWT, Response, Pagination, Seeder, dll.)
├── logs/           # Log Aplikasi (Rotated & Structured)
├── main.go         # Titik Masuk Aplikasi
└── go.mod          # Manajemen Dependensi
```

## ⚡ Memulai (Getting Started)

### Prasyarat

- Go 1.24 atau lebih baru terinstal.
- PostgreSQL berjalan secara lokal atau dapat diakses dari jauh.

### Instalasi

1. **Clone repositori:**
   ```bash
   git clone <repository-url>
   cd <repository-folder>
   ```

2. **Install dependensi:**
   ```bash
   make deps
   # atau `go mod tidy`
   ```

3. **Atur Variabel Lingkungan:**
   Buat file `.env` di direktori root dan tambahkan konfigurasi berikut:

   ```env
   # Konfigurasi Server
   PORT=8080

   # Konfigurasi Database
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=passwordmu
   DB_NAME=namadatabasemu
   DB_PORT=5432
   DB_SSLMODE=disable
   DB_TIMEZONE=Asia/Jakarta

   # Konfigurasi JWT
   JWT_SECRET=rahasia_super_kunci_kamu

   # Konfigurasi SMTP (Untuk Pengiriman Email)
   SMTP_HOST=sandbox.smtp.mailtrap.io
   SMTP_PORT=2525
   SMTP_USER=user_smtp_kamu
   SMTP_PASS=password_smtp_kamu
   ```

### 🏃 Menjalankan Aplikasi

Gunakan perintah `make` untuk kemudahan:

1. **Jalankan Migrasi Database:**
   ```bash
   make db-migrate
   ```

2. **Jalankan Seeder Database (Isi Data Awal):**
   ```bash
   make db-seed
   ```
   *Ini akan membuat akun Super Admin (`superadmin@example.com` / `password`) dan data dummy lainnya.*

3. **Jalankan Server dengan Hot Reload (Mode Pengembangan):**
   ```bash
   make dev
   ```
   *Server otomatis restart saat ada perubahan kode.*

4. **Jalankan Biasa:**
   ```bash
   make run
   ```

5. **Generate Dokumentasi Swagger:**
   ```bash
   make swagger
   ```

6. **Audit Keamanan (Security Audit):**
   ```bash
   make audit
   ```
   *Memeriksa kerentanan keamanan pada dependensi dan standar library Go.*


7. **Release (Production Mode):**
   ```bash
   make run-release
   ```
   *Menjalankan aplikasi tanpa log debug framework (GIN-debug) dan menyimpan log ke file.*

## 📚 Dokumentasi API

Proyek ini menggunakan **Swagger** untuk dokumentasi API.

1. Jalankan aplikasi.
2. Buka browser dan navigasi ke:
   
   ```
   http://localhost:8080/swagger/index.html
   ```
## 🔗 Panduan Pengembangan

Untuk panduan teknis mendalam mengenai cara menambah model baru, pagination, dan kebijakan akses, silakan baca:
👉 **[PANDUAN PENGEMBANGAN (GUIDE_PENGEMBANGAN.md)](./GUIDE_PENGEMBANGAN.md)**

Untuk tips performa database dan indexing:
👉 **[PANDUAN OPTIMASI DATABASE (GUIDE_OPTIMASI_DATABASE.md)](./GUIDE_OPTIMASI_DATABASE.md)**

## 🤝 Berkontribusi

1. Fork Proyek
2. Buat Feature Branch (`git checkout -b feature/FiturKeren`)
3. Commit Perubahan (`git commit -m 'Menambahkan FiturKeren'`)
4. Push ke Branch (`git push origin feature/FiturKeren`)
5. Buka Pull Request
