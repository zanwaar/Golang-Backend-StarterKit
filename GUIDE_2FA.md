# Panduan Integrasi 2FA (Two-Factor Authentication)

Dokumen ini menjelaskan cara kerja dan integrasi fitur Two-Factor Authentication (2FA) menggunakan protokol TOTP (Time-based One-Time Password) di backend ini. Fitur ini kompatibel dengan aplikasi seperti **Google Authenticator**, **Authy**, atau **Microsoft Authenticator**.

---

## 🏗️ Alur Kerja (Workflow)

1.  **Inisiasi (Setup)**: User meminta token rahasia (Secret) dan QR Code dari server.
2.  **Verifikasi (Verify)**: User memindai QR Code dan memasukkan kode 6-digit untuk mengaktifkan 2FA secara permanen.
3.  **Login**: Jika 2FA aktif, user wajib menyertakan kode 6-digit saat login.

---

## 🔌 Endpoint API

### 1. Setup 2FA
Meminta server untuk membuatkan Secret Key baru. 2FA belum aktif pada tahap ini.

-   **URL**: `POST /api/2fa/setup`
-   **Auth**: Bearer Token (Login Required)
-   **Response Sukses (200)**:

```json
{
  "success": true,
  "message": "2FA Setup Initiated",
  "data": {
    "secret": "JBSWY3DPEHPK3PXP",
    "qr_code_url": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
  }
}
```

> **Untuk Frontend**: Tampilkan `qr_code_url` sebagai gambar (`<img src="...">`) agar user bisa scan. Simpan `secret` jika ingin ditampilkan sebagai opsi manual input.

### 2. Verifikasi & Aktivasi
User memasukkan kode dari aplikasi Authenticator untuk mengkonfirmasi bahwa mereka telah berhasil menyimpan Secret Key.

-   **URL**: `POST /api/2fa/verify`
-   **Auth**: Bearer Token
-   **Body**:
```json
{
  "code": "123456"
}
```
-   **Response Sukses (200)**:
```json
{
  "success": true,
  "message": "2FA Verified and Enabled"
}
```

> **Note**: Setelah langkah ini sukses, `is_two_fa_enabled` di database user akan berubah menjadi `true`.

### 3. Login dengan 2FA
Jika user sudah mengaktifkan 2FA, login biasa akan gagal (atau bisa diatur untuk meminta kode). Backend ini mendukung pengiriman kode langsung di body login.

-   **URL**: `POST /api/login`
-   **Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "two_fa_code": "123456" // OPTIONAL: Wajib jika 2FA aktif
}
```

#### Skenario Login:
1.  **User tanpa 2FA**: Cukup kirim `email` dan `password`.
2.  **User dengan 2FA (Kode Kosong)**:
    -   Response Login akan error: `401 Unauthorized` dengan pesan `"2FA code required"`.
    -   Frontend harus meminta user memasukkan kode 2FA.
3.  **User dengan 2FA (Kode Salah)**:
    -   Response Error: `"invalid 2FA code"`.

### 4. Cek Status 2FA User
Untuk mengetahui apakah user yang sedang login sudah mengaktifkan 2FA atau belum.

-   **URL**: `GET /api/me`
-   **Response**:
```json
{
  "data": {
    "id": "...",
    "email": "...",
    "is_two_fa_enabled": true // true = aktif, false = belum
  }
}
```

---

## 🧪 Testing Manual (Postman/cURL)

Berikut langkah-langkah untuk mencoba fitur ini:

1.  **Login** sebagai user biasa untuk mendapatkan Token.
2.  **Setup**: Hit endpoint `/api/2fa/setup`. Copy string Base64 dari `qr_code_url`.
3.  **Tampilkan QR**: Paste string Base64 ke [Base64 Image Viewer](https://jaredwinick.github.io/base64-image-viewer/) atau tools sejenis untuk melihat QR Code.
4.  **Scan**: Buka Google Authenticator di HP, scan QR Code tersebut.
5.  **Verifikasi**: Masukkan kode yang muncul di HP ke endpoint `/api/2fa/verify`.
6.  **Re-Login**: Coba login ulang.
    -   Tanpa kode: Gagal.
    -   Dengan kode: Sukses.

---

## 📚 Library
Fitur ini dibangun menggunakan library standar industri:
- [github.com/pquerna/otp](https://github.com/pquerna/otp)
