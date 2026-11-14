### Panduan Tugas Kelompok

Pengantar
Dokumen ini memberi gambaran contoh pengerjaan tugas kelompok, memakai studi kasus fiktif
RuangKu — Sistem Peminjaman Ruang Kampus. Kelompok kalian bebas memilih sistem lain; yang
ditiru adalah bentuk dan alur kerjanya.

Tugas terdiri dari empat langkah:

1. Bentuk kelompok dan bagi peran.
2. Susun rancangan sistem.
3. Buat boilerplate backend dengan Go dan Clean Architecture sesuai modul.
4. Kerjakan bersama di GitHub dengan micro commits; Git Master melakukan merge ke main.

Langkah 1 — Membentuk Kelompok
Setiap kelompok berisi 3 orang. Kelompok 4 orang hanya dibuat jika jumlah mahasiswa tidak habis
dibagi 3. Contoh: kelas 31 orang → 9 kelompok berisi 3 orang + 1 kelompok berisi 4 orang.

Contoh susunan Kelompok 1 — RuangKu:

Anggota Pengalaman Peran

Andi Sudah biasa memakai Git dan Go dalam
tim
Git Master + logika bisnis

Bima Pernah membuat REST API sederhana Database + akses data

Citra Masih belajar Go Entity + routing API

(Dewi, jika 4 orang) Masih belajar Konfigurasi + dokumentasi

Anggota yang berpengalaman tidak mengerjakan semuanya. Tugasnya membimbing dan mereview
pekerjaan anggota lain.

Langkah 2 — Rancangan Sistem
Tulis rancangan di README.md repo sebelum mulai membuat kode.

Nama sistem: RuangKu — Sistem Peminjaman Ruang Kampus.

Deskripsi: RuangKu membantu mahasiswa dan organisasi kemahasiswaan meminjam ruang kelas, lab,
atau aula secara online. Peminjam mengajukan peminjaman, lalu pengelola ruang menyetujui atau
menolaknya. Sistem mencegah dua peminjaman pada ruang dan jam yang sama.

Role pada sistem:

Role Hak akses utama

Admin Mengelola user dan data ruang

---

Role Hak akses utama

Pengelola Ruang Menyetujui atau menolak peminjaman

Peminjam Melihat ruang, mengajukan dan membatalkan peminjaman, melihat
riwayat

Fitur:

No Fitur Role

1 Registrasi dan login Semua

2 Kelola data ruang (tambah, lihat, ubah, hapus) Admin

3 Cek ketersediaan ruang per tanggal Peminjam

4 Ajukan peminjaman Peminjam

5 Setujui / tolak peminjaman Pengelola Ruang

6 Batalkan peminjaman Peminjam

7 Riwayat peminjaman Peminjam, Pengelola Ruang

Rancangan database:

Tabel Kolom utama Relasi

users id, name, email, password_hash, role, created_at Satu user punya banyak
peminjaman

rooms id, code, name, building, capacity, is_active,
created_at
Satu ruang punya banyak
peminjaman

bookings id, user_id, room_id, purpose, start_time,
end_time, status, created_at
Milik satu user dan satu ruang

booking_logs id, booking_id, actor_id, from_status, to_status,
note, created_at
Mencatat setiap perubahan
status peminjaman

Status peminjaman: pending, approved, rejected, cancelled. Lampirkan juga gambar ERD (boleh dari
dbdiagram.io, draw.io, atau tools lain).

Langkah 3 — Boilerplate Backend Go
Boilerplate adalah kerangka proyek yang sudah bisa dijalankan, minimal dengan satu fitur yang
berjalan utuh (misalnya kelola data ruang). Fitur lain cukup berupa kerangka yang nanti dilengkapi.

Aturan wajib:

• Menggunakan bahasa Go.
• Menggunakan Clean Architecture sesuai modul praktikum sebelumnya.
• Tidak menggunakan istilah atau layer controller, handler, dan sejenisnya.

---

Contoh struktur folder (hanya ilustrasi):

ruangku-backend/
├── cmd/
│ └── api/
│ └── main.go
├── internal/
│ ├── domain/
│ │ ├── user.go
│ │ ├── room.go
│ │ └── booking.go
│ ├── usecase/
│ │ ├── room_usecase.go
│ │ └── booking_usecase.go
│ ├── repository/
│ │ ├── room_repository.go
│ │ └── booking_repository.go
│ └── delivery/
│ └── http/
│ └── room_delivery.go
├── pkg/
│ ├── config/
│ └── database/
├── migrations/
│ └── 001_init.sql
├── .env.example
├── .gitignore
├── go.mod
└── README.md
Penting: struktur di atas hanya contoh untuk memberi gambaran. Untuk struktur folder
yang benar, silakan sesuaikan dengan modul pada praktikum sebelumnya — termasuk
nama folder, nama layer, dan pembagian file.

Langkah 4 — GitHub, Micro Commits, dan Git Master
Setiap anggota bekerja di branch sendiri, melakukan commit kecil-kecil, lalu membuat Pull Request
(PR). Hanya Git Master yang melakukan merge ke branch main. Repo dibuat public.

Apa itu micro commit? Satu commit berisi satu perubahan kecil yang bisa dijelaskan dalam satu
kalimat, bukan satu commit besar berisi seluruh proyek.

Kurang baik Baik

update feat(domain): add Room entity

upload semua file feat(repository): add query to get all rooms

fix fix(usecase): reject room with zero capacity

revisi final docs(readme): add database design

Alur kerja:

1. Git Master membuat repo public, mengundang anggota sebagai collaborator, dan melakukan
commit awal (inisialisasi proyek, .gitignore, README.md).
2. Git Master mengaktifkan branch protection pada main, supaya semua perubahan wajib lewat PR.

---

3. Setiap anggota membuat branch sesuai bagiannya, misalnya feat/domain, feat/repository-room.
4. Anggota melakukan beberapa micro commit di branch-nya, lalu membuat PR ke main.
5. Git Master mereview PR (apakah kode sesuai layer, bisa di-build), memberi komentar bila perlu,
lalu melakukan merge.
6. Setelah merge, semua anggota memperbarui main di laptop masing-masing sebelum lanjut
bekerja.

Contoh pembagian branch dan commit:

Anggota Branch Contoh commit

Andi (Git Master) main (setup awal),
feat/usecase-room
chore: init go module · docs(readme): add project
description · feat(usecase): validate room input

Citra feat/domain, feat/delivery-
room
feat(domain): add Room entity · feat(domain): add
Booking entity · feat(delivery): add room routes

Bima feat/migration,
feat/repository-room
feat(migration): create rooms table · feat(repository): add
query to get all rooms · feat(repository): add query to
create room

Dewi (jika ada) chore/config, docs/api chore(config): load env variables · docs(readme): add
how to run project

Tips penilaian: merge PR dengan opsi Create a merge commit (bukan Squash) agar commit setiap
anggota tetap terlihat. Pastikan email Git di laptop sama dengan email akun GitHub, supaya kontribusi
tercatat di Insights → Contributors.

Checklist Pengumpulan

• ☐ Kelompok 3 orang (4 hanya jika sisa), campuran anggota berpengalaman dan yang masih
belajar
• ☐ README.md berisi nama sistem, deskripsi, fitur, role, rancangan database, anggota dan
perannya
• ☐ Struktur folder sesuai Clean Architecture pada modul praktikum; tidak ada controller / handler
• ☐ Minimal satu fitur berjalan utuh, dan proyek bisa dijalankan tanpa error
• ☐ Repo GitHub public
• ☐ Setiap anggota punya beberapa micro commit dengan pesan yang jelas
• ☐ Semua perubahan masuk main lewat Pull Request yang di-merge oleh Git Master
• ☐ File .env asli tidak ikut ter-commit