# 🧪 Little Alchemy 2 Recipe Finder
🚀 **Little Alchemy 2 Recipe Finder** adalah proyek web yang dibuat untuk menemukan resep elemen pada game Little Alchemy 2 menggunakan algoritma pencarian (BFS, DFS, dan Bidirectional Search).

📌 **Tujuan**: Mempelajari implementasi algoritma pencarian dalam konteks praktis dan pengembangan aplikasi web.

## 🛠 Prasyarat  
Karena proyek ini menggunakan Go, React, dan beberapa teknologi lainnya, pastikan kamu sudah menginstal beberapa alat berikut sebelum menjalankannya:  

### 1️⃣ Go (Backend)  
Jika belum punya, silakan download dan install Go di [sini](https://go.dev/dl/).  

### 2️⃣ Node.js & npm (Frontend)  
Pastikan Node.js dan npm sudah terinstall. Jika belum, download dan install dari [sini](https://nodejs.org/).  

Setelah install, cek apakah sudah terpasang dengan menjalankan perintah berikut di terminal:  
```sh
node -v
npm -v
```

### 3️⃣ Docker (Opsional)  
Docker diperlukan jika ingin menjalankan aplikasi dalam container. Jika belum punya, download dan install dari [sini](https://www.docker.com/products/docker-desktop).  

---

## 🏗 Menjalankan Aplikasi Secara Lokal  
Setelah semua persiapan selesai, ikuti langkah-langkah berikut untuk menjalankan proyek ini di lokal:  

### 1️⃣ Clone Repository  
```sh
git clone https://github.com/username/little-alchemy-recipe-finder.git
cd little-alchemy-recipe-finder
```

### 2️⃣ Menjalankan Backend  
```sh
cd backend
go mod download
go run main.go
```
Backend akan berjalan di `http://localhost:8081`.

### 3️⃣ Menjalankan Frontend  
Buka terminal baru, lalu:
```sh
cd frontend
npm install
npm run dev
```
Frontend akan berjalan di `http://localhost:5173` atau port lain yang tersedia.

---

## 🐳 Menjalankan dengan Docker  
Proyek ini juga bisa dijalankan menggunakan Docker:

### 1️⃣ Menjalankan dengan Docker Compose
Di folder utama proyek, jalankan:
```sh
docker-compose up -d
```

Aplikasi akan berjalan di `http://localhost`.

### 2️⃣ Membangun Image Docker Secara Terpisah

**Backend:**
```sh
cd backend
docker build -t alchemy-backend .
docker run -p 8081:8081 -v ../database:/app/database alchemy-backend
```

**Frontend:**
```sh
cd frontend
docker build -t alchemy-frontend .
docker run -p 80:80 alchemy-frontend
```

---

## 🚢 Deployment  
Proyek ini dapat di-deploy ke berbagai platform:

### 1️⃣ Digital Ocean App Platform
1. Buat akun di Digital Ocean
2. Buat aplikasi baru dan hubungkan ke repository GitHub
3. Konfigurasikan build frontend dan backend
4. Deploy aplikasi

### 2️⃣ Railway.app (Rekomendasi untuk Mahasiswa)
1. Buat akun di Railway.app
2. Hubungkan repository GitHub
3. Tambahkan variabel lingkungan
4. Deploy aplikasi

### 3️⃣ Heroku atau AWS
Lihat panduan lengkap di `deployment-guide.md` untuk instruksi detail tentang deployment ke Heroku atau AWS.

---

## 📁 Struktur Proyek  
Berikut adalah struktur folder utama dalam proyek ini:  

```
📂 Tubes2Stima/
│-- 📂 backend/                # Kode backend (Go)
│   ├── 📂 controllers/     # Controller untuk API
│   ├── 📂 services/        # Layanan dan algoritma pencarian
│   ├── 📂 utils/           # Utilitas dan helper
│   ├── Dockerfile          # Docker config untuk backend
│   ├── .env                # Variabel lingkungan
│   ├── go.mod              # Dependensi Go
│   └── main.go             # Entry point backend
│-- 📂 frontend/            # Kode frontend (React)
│   ├── 📂 public/          # Asset publik
│   ├── 📂 src/             # Source code React
│   │   ├── 📂 assets/      # Gambar dan resource
│   │   ├── 📂 components/  # Komponen React
│   │   ├── 📂 hooks/       # Custom hooks
│   │   ├── 📂 pages/       # Halaman aplikasi
│   │   └── 📂 utils/       # Fungsi utilitas
│   ├── Dockerfile          # Docker config untuk frontend
│   └── nginx.conf          # Konfigurasi nginx
│-- 📂 database/            # Database SQLite dan JSON mapper
│-- docker-compose.yml      # Konfigurasi Docker Compose
│-- .gitignore              # File yang diabaikan Git
└-- deployment-guide.md     # Panduan deployment
```

---

## 🔧 Build untuk Production  
Jika ingin membuat build untuk production:

**Frontend:**
```sh
cd frontend
npm run build
```

**Backend:**
```sh
cd backend
go build -o main .
```

Atau gunakan Docker Compose untuk build dan jalankan sekaligus:
```sh
docker-compose up --build -d
```

---

## ❓ Troubleshooting  

### ❌ Masalah Database
- Pastikan database SQLite (`database/alchemy.db`) tersedia dan dapat diakses.
- Pastikan folder `/database` di-mount dengan benar jika menggunakan Docker.

### ❌ Masalah Port
- Jika port 8081 atau 80 sudah digunakan, ubah di `.env` atau `docker-compose.yml`.
- Untuk frontend development, jika port 5173 sudah digunakan, Vite akan secara otomatis mencari port lain.

### ❌ Masalah CORS
- Jika mengalami masalah CORS, pastikan URL API diatur dengan benar di frontend.
- Dalam mode development, gunakan URL relatif (`/api/search`) atau ubah sesuai kebutuhan.

---

## 📌 Ringkasan  
- Pastikan sudah install **Go, Node.js, dan Docker (opsional)**.
- Gunakan `docker-compose up -d` untuk menjalankan aplikasi dengan Docker.
- Untuk development, jalankan backend dan frontend secara terpisah.
- Ikuti `deployment-guide.md` untuk petunjuk deployment ke berbagai platform.

🎉 **Selamat mencoba! 🚀**
