[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/Nhr6kRao)
[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/3O2Jx-Z_)
# Graded Challenge 3 - P2

Graded Challenge ini dibuat guna mengevaluasi pembelajaran pada Hacktiv8 Program Fulltime Golang khususnya pada pembelajaran Echo Framework

## Assignment Objectives
Graded Challenge 3 ini dibuat guna mengevaluasi pemahaman SQL sebagai berikut:

- Mampu memahami konsep REST API
- Mampu membuat REST API using Echo Framework
- Mampu membuat REST API dengan implementasi database postgresql
- Mampu menggunakan third-party API
- Mampu implemantasi Autentikasi dan Autorisasi menggunakan JWT pada REST API dengan Echo Framework

## Assignment Directions - Hacktivagram API
Anda diminta oleh perusahaan bernama Hacktivagram untuk membuat sebuah aplikasi social media serupa dengan instagram. Dimana pada aplikasi ini user dapat memposting sebuah konten, dan memberikan komen pada konten. Terdapat juga fitur untuk menghapus postingan, komen, yang sudah disubmit oleh user.
Untuk mempermudah proses pengerjaan, Hacktivagram sudah menyediakan ERD awal, sebagai berikut:

![Alt text](erd.png)


### Requirements:
- Database Requirement
  - Buatlah table sesuai dengan ERD yang telah disediakan, improvement atau perubahan struktur table diperbolehkan, selama masih dapat memenuhi kebetuhuna requirement fitur.
  - RESTRICTION
    - Pastikan email dan username setiap user harus unik dan tidak boleh ada yang sama antar user
    - Pastikan semua field pada database tidak boleh kosong/null
  - Pastikan untuk menyertakan query DDL dan Query seeding data pada folder project GC ini, buatlah file dengan nama ddl.sql pada root folder (jika ada)

- Web API untuk user harus memiliki beberapa fitur sebagai berikut:
  - <b>POST</b> /users/register - Menyimpan data user baru
    - Request dari endpoint ini harus meliputi nama, alamat, email, password, dan umur
    - Response dari endpoint ini harus berupa message sukses, dan data user yang berhasil disimpan, jika terdapat kesalahan pada request maka response harus terdiri dari message yang menjelaskan kesalahan pada input request nya.
    - Perlu diperhatikan, untuk tidak mengikutsertakan data sensitif seperti password pada response
  - <b>POST</b> /users/login - Login untuk mendapatkan akses token
    - Request dari endpoint ini harus meliputi email, password
    - Response dari endpoint ini adalah akses token JWT, jika terdapat kesalahan pada request maka response harus terdiri dari message yang menjelaskan kesalahan pada input request nya.
  - Memerlukan login/autentikasi menggunakan JWT, pastikan untuk mengakses setiap endpoint dibawah, perlu disertakan akses token pada headers dengan key `authorization`, jika akses token jwt tidak terautentikasi, maka web api perlu memberikan response error tidak terautentikasi beserta message yang sesuai.
    - <b>POST</b> /posts - Membuat postingan baru milik user yang sedang login (akses token terlampir pada header).
      - Request dari endpoint ini harus meliputi konten dan image_url (pastikan image_url merupakan format url, gunakan regex untuk memvalidasi), jika konten tidak dimasukkan pada request body, maka isi konten dengan random jokes dari 3rd Party API `https://api-ninjas.com/api/jokes`
      - Response dari endpoint ini adalah message sukses dan objek post yang baru berhasil dibuat. Jika terdapat kesalahan pada request maka response harus terdiri dari message yang menjelaskan kesalahan pada input request nya.
    - <b>GET</b> /posts - Menampilkan semua postingan yang
     tersimpan pada database
      - Response dari endpoint ini adalah sebuah array of objek post yang tersimpan pada database
    
    - <b>GET</b> /posts/:id - Menampilkan detail postingan yang tersimpan pada database sesuai dengan param ID
      - Resopnse dari endpoint ini adalah objek post sesuai dengan param id, data lain yang perlu ditampilkan adalah list dari komen pada post(user pembuat komen perlu ditampilkan) tersebut. Jika post dengan ID tidak ditemukan, maka response harus terdiri dari message yang menjelaskan bahwa data post tidak ditemukan

    - <b>DELETE</b> /posts/:id - Menghapus sebuah postingan yang tersimpan pada database sesuai dengan param ID
      - Response dari endpoint ini adalah data post sesuai dengan ID pada parameter endpoint dan message yang yang menjelaskan bahwa proses penghapusan data post berhasil, jika post dengan ID tidak ditemukan, maka response harus terdiri dari message yang menjelaskan bahwa data post tidak ditemukan.
      - Pastikan hanya owner dari post yang dapat melakukan aksi ini. Jika user lain selain owner dari post melakukan aksi ini pastikan untuk membatalkan aksi dan memberikan response error terdiri dari message yang menjelaskan bahwa user tersebut tidak dapat melakukan aksi delete

    - <b>POST</b> /comments - Membuat comment baru pada sebuah post 
      - Request dari endpoint ini harus meliputi konten komen yang akan disubmit.
      - Response dari endpoint ini adalah message sukses dan objek komen yang baru berhasil dibuat. Jika terdapat kesalahan pada request maka response harus terdiri dari message yang menjelaskan kesalahan pada input request nya.
    
    - <b>GET</b> /comments/:id - Menampilkan comment yang tersimpan pada database sesuai dengan param ID
      - Resopnse dari endpoint ini adalah objek comment sesuai dengan param id, data lain yang perlu ditampilkan adalah owner pembuat komen dan data post nya. Jika post dengan ID tidak ditemukan, maka response harus terdiri dari message yang menjelaskan bahwa data post tidak ditemukan

    - <b>DELETE</b> /comments/:id - Menghapus sebuah postingan yang tersimpan pada database sesuai dengan param ID
      - Response dari endpoint ini adalah data komen sesuai dengan ID pada parameter endpoint dan message yang yang menjelaskan bahwa proses penghapusan data komen berhasil, jika post dengan ID tidak ditemukan, maka response harus terdiri dari message yang menjelaskan bahwa data komen tidak ditemukan.
      - Pastikan hanya owner dari komen yang dapat melakukan aksi ini. Jika user lain selain owner dari komen melakukan aksi ini pastikan untuk membatalkan aksi dan memberikan response error terdiri dari message yang menjelaskan bahwa user tersebut tidak dapat melakukan aksi delete

    - <b>GET</b> /activities - Menampilkan user activities yang tersimpan pada database sesuai dengan user yang sedang login
      - Resopnse dari endpoint ini adalah data list aktivitas user pada aplikasi yang tersimpan pada table user_activity_logs

  - Setiap aksi yang dilakukan user harus dicatat dalam table user_activity_logs, dengan deskripsi yang menjelaskan aksi yang dilakukan user, sebagai contoh jika user berhasil membuat post baru, maka anda juga harus menambahkan data ke table user_activity_logs dengan isi deskripsi `user create new POST with ID [new post id]`
  - Setiap endpoint diatas harus menerapkan best practice REST termasuk status code dan http method yang digunakan
  - Setiap endpoint diatas perlu dibuat dokumentasi API menggunakan Swagger dan dapat diakses pada `/swagger/index.html`
- Requirement 3rd Party API
  - Kembangkanlah endpoint POST /post, jika tidak terdapat input konten dari user maka, buatkanlah konten secara random dari 3rd Party API berikut `https://api-ninjas.com/api/jokes` 
- Deployment Requirement
  - Buatlah database pada platform Supabase(postgreSQL) dan sambungkan dengan aplikasi anda.
  - Deploy REST API yang sudah anda buat dengan menggunakan platform Heroku, dan pastikan mencantumkan url hasil deployment pada section expected result dan deployment notes.
- Pastikan untuk mengikuti best practice untuk penggunaan environment variable

## Expected Result
- Web API dapat diakses pada _________ (isi dengan url hasil deployment anda).
- Web API memiliki endpoint sebagai berikut
  - <b>POST</b> /users/register
    - request body -> `{ first_name, last_name, address, email, username, password, age }`
  - <b>POST</b> /users/login 
    - request body -> `{ username, password }`
  - <b>POST</b> /posts
    - request headers -> `{ authorization }`
    - request body -> `{ content, image_url }`
  - <b>GET</b> /posts
    - request headers -> `{ authorization }`
  - <b>GET</b> /posts/:id
    - request headers -> `{ authorization }`
  - <b>DELETE</b> /posts/:id
    - request headers -> `{ authorization }`
  - <b>POST</b> /comments
    - request headers -> `{ authorization }`
    - request body -> `{ content, post_id }`
  - <b>GET</b> /comments/:id
    - request headers -> `{ authorization }`
  - <b>DELETE</b> /comments/:id
    - request headers -> `{ authorization }`
  - <b>GET</b> /activities
    - request headers -> `{ authorization }`
  
## Assignment Submission
Push Assigment yang telah Anda buat ke akun Github Classroom Anda masing-masing.

### Assignment Notes:
- Jangan terburu-buru dalam menyelesaikan masalah atau mencoba untuk menyelesaikannya sekaligus.
- Jangan menyalin kode dari sumber eksternal tanpa memahami bagaimana kode tersebut bekerja.
- Jangan menentukan nilai secara hardcode atau mengandalkan asumsi yang mungkin tidak berlaku dalam semua kasus.
- Jangan lupa untuk menangani negative case, seperti input yang tidak valid
- Jangan ragu untuk melakukan refaktor kode Anda, buatlah struktur project anda lebih mudah dibaca dan dikembangkan kedepannya, pisahkanlah setiap bagian kode program pada folder sesuai dengan tugasnya masing-masing.

### Additional Notes
Total Points : 100

Deadline : Diinformasikan oleh instruktur saat briefing GC. Keterlambatan pengumpulan tugas mengakibatkan skor GC 3 menjadi 0.

Informasi yang tidak dicantumkan pada file ini harap dipastikan/ditanyakan kembali kepada instruktur. Kesalahan asumsi dari peserta mungkin akan menyebabkan kesalahan pemahaman requirement dan mengakibatkan pengurangan nilai.

### Deployment Notes
- Deployed url: _________ (isi dengan url hasil deployment anda)
