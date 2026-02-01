# 🚀 Panduan Deployment Ubuntu (Production)

Dokumen ini berisi langkah-langkah detail untuk mendeploy aplikasi **Golang Backend Starter Kit** ke server Ubuntu.

## 📋 Prasyarat
- Server dengan OS Ubuntu (disarankan 22.04 LTS atau terbaru)
- Akses Root atau Sudo
- Domain yang sudah diarahkan ke IP Server (Opsional untuk testing)

---

## 1. Persiapan Awal
Update server dan instal dependensi dasar.

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y git wget curl nginx postgresql postgresql-contrib
```

## 2. Instalasi Go (Build di Server)
Jika Anda ingin melakukan kompilasi koding langsung di server.

```bash
# Download Go (Cek versi terbaru di golang.org)
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz

# Tambahkan ke Environment PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verifikasi
go version
```

## 3. Konfigurasi Database PostgreSQL
1. Masuk sebagai user postgres:
   ```bash
   sudo -u postgres psql
   ```
2. Jalankan perintah SQL berikut:
   ```sql
   CREATE DATABASE mydb;
   CREATE USER myuser WITH PASSWORD 'aman_banget_123';
   GRANT ALL PRIVILEGES ON DATABASE mydb TO myuser;
   \q
   ```

## 4. Setup Aplikasi
1. Clone repository:
   ```bash
   git clone https://github.com/username/repo-anda.git /var/www/golang-api
   cd /var/www/golang-api
   ```
2. Konfigurasi Environment:
   ```bash
   cp .env.example .env
   nano .env
   ```
   **Pastikan menyesuaikan nilai berikut:**
   - `ENVIRONMENT=production`
   - `PORT=8080`
   - `DB_HOST=localhost`
   - `DB_NAME=mydb`
   - `DB_USER=myuser`
   - `DB_PASSWORD=aman_banget_123`
   - `JWT_SECRET=GantiDenganStringSatuParagrafYangUnik`

3. Build Binary:
   ```bash
   make build
   ```

## 5. Menjalankan Database Migration
Jalankan migrasi untuk membuat tabel-tabel di database production.

```bash
./bin/api-gin-production -migrate
```

## 6. Konfigurasi Systemd (Auto-run)
Agar aplikasi berjalan otomatis saat server restart dan bisa memulihkan diri jika crash.

1. Buat file service:
   ```bash
   sudo nano /etc/systemd/system/golang-api.service
   ```
2. Tempelkan konfigurasi ini:
   ```ini
   [Unit]
   Description=Golang API Service
   After=network.target postgresql.service

   [Service]
   Type=simple
   User=ubuntu
   WorkingDirectory=/var/www/golang-api
   ExecStart=/var/www/golang-api/bin/api-gin-production
   Restart=always
   RestartSec=5
   EnvironmentFile=/var/www/golang-api/.env
   StandardOutput=append:/var/www/golang-api/logs/production.log
   StandardError=append:/var/www/golang-api/logs/error.log

   [Install]
   WantedBy=multi-user.target
   ```
3. Aktifkan Service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable golang-api
   sudo systemctl start golang-api
   ```

## 7. Konfigurasi Nginx (Reverse Proxy)
1. Buat konfigurasi baru:
   ```bash
   sudo nano /etc/nginx/sites-available/api.domain.com
   ```
2. Isi dengan:
   ```nginx
   server {
       listen 80;
       server_name api.domain.com; # Ganti dengan domain atau IP

       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   ```
3. Hubungkan ke sites-enabled dan restart Nginx:
   ```bash
   sudo ln -s /etc/nginx/sites-available/api.domain.com /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl restart nginx
   ```

## 8. Keamanan (SSL & Firewall)
1. Instal Certbot (SSL):
   ```bash
   sudo apt install certbot python3-certbot-nginx -y
   sudo certbot --nginx -d api.domain.com
   ```
2. Aktifkan Firewall:
   ```bash
   sudo ufw allow 'Nginx Full'
   sudo ufw allow ssh
   sudo ufw enable
   ```

---

## 🛠️ Perintah Berguna
- **Cek Status Aplikasi:** `sudo systemctl status golang-api`
- **Cek Log Real-time:** `journalctl -u golang-api -f`
- **Restart Aplikasi:** `sudo systemctl restart golang-api`
- **Update Kode:**
  ```bash
  git pull origin main
  make build
  sudo systemctl restart golang-api
  ```
