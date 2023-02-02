DROP DATABASE IF EXISTS simpeg;
CREATE DATABASE simpeg;
USE simpeg;

CREATE TABLE users (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    photo TEXT NOT NULL
);

CREATE TABLE lembur (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    no_payroll VARCHAR(50) NOT NULL,
    nama VARCHAR(50) NOT NULL,
    jabatan VARCHAR(50) NOT NULL,
    gaji_pokok FLOAT NOT NULL,
    basis FLOAT NOT NULL
);

/*
CREATE TABLE tanggal_lembur (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    tanggal VARCHAR(50) NOT NULL,
    total_jam FLOAT NOT NULL,
    id_lembur VARCHAR(255) NOT NULL,
    FOREIGN KEY (id_lembur) REFERENCES lembur(id) ON DELETE CASCADE ON UPDATE CASCADE
);
*/

ALTER TABLE tanggal_lembur ADD FOREIGN KEY (id_lembur) REFERENCES lembur(id) ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE tanggal_lembur (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    tanggal VARCHAR(50) NOT NULL,
    total_jam FLOAT NOT NULL,
    id_lembur VARCHAR(255) NOT NULL
);

INSERT INTO lembur(id, no_payroll, nama, jabatan, gaji_pokok, basis) VALUES
('1111111', '500408', 'Kusnul Mathatir', 'Dan Jaga', 4750000, 27457),
('11111112', '500511', 'Maralat', 'Patroli Satpam', 3901700, 22553),
('11111113', '500532', 'Amos Sibarani', 'Dan Jaga', 4750000, 27457);

INSERT INTO tanggal_lembur(id, tanggal, total_jam, id_lembur) VALUES
('aaaa','2020-10-15', 8, "1111111"),
('bbbb','2020-10-18', 8, "1111111"),
('cccc','2020-10-20', 8, "1111111"),
('dddd','2020-10-1', 8, "11111112"),
('eeee','2020-10-29', 8, "11111112"),
('ffff','2020-10-1', 8, "11111113"),
('gggg','2020-10-29', 8, "11111113");

CREATE TABLE non_staff (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    no_payroll VARCHAR(255) NOT NULL,
    nama VARCHAR(50) NOT NULL,
    foto VARCHAR(255) NOT NULL,
    status_pph_21 VARCHAR(50) NOT NULL,
    alamat VARCHAR(50) NOT NULL,
    nik VARCHAR(50) NOT NULL,
    npwp VARCHAR(50) NOT NULL,
    gender VARCHAR(50) NOT NULL,
    jabatan VARCHAR(50) NOT NULL,
    golongan VARCHAR(50) NOT NULL,
    expat VARCHAR(50) NOT NULL,
    ptkp VARCHAR(50) NOT NULL,
    status_kerja VARCHAR(50) NOT NULL,
    awal_kerja INT NOT NULL,
    akhir_kerja INT NOT NULL,
    gaji_pokok FLOAT DEFAULT NULL,
    tunjangan_kemahalan FLOAT DEFAULT NULL,
    tunjangan_perumahan FLOAT DEFAULT NULL,
    tunjangan_jabatan FLOAT DEFAULT NULL,
    tunjangan_lain_pph21 FLOAT DEFAULT NULL
);

INSERT INTO non_staff
(
    id, no_payroll, nama, foto, status_pph_21, alamat, nik,npwp,gender,jabatan,golongan,
    expat,ptkp,status_kerja,awal_kerja,akhir_kerja,gaji_pokok,tunjangan_kemahalan,
    tunjangan_perumahan,tunjangan_jabatan,tunjangan_lain_pph21
) VALUES
    ('111111110', '502275', 'TIMOTIUS SEPTRADA S','sammidev.codes/sammi.png', 'Beban Karyawan', 'Pekanbaru', '3171041010710010',
    '142650043216000', 'M', 'KERANI A - PERSONALIA', ' III-4', 'L', 'TK/0', 'Aktif', 1, 12, 3467876, 500000, NULL, NULL, NULL),

    ('111111111','500408','KUSNUL MATHATIR','sammidev.codes/sammi.png','Beban Karyawan','Pekanbaru','3171041010710010',
    '090063744023000','M','DAN JAGA','III-9','L','K/3','Aktif',1,12,4750000,500000,NULL,NULL,NULL),

    ('111111112','500511','MARALAT','sammidev.codes/sammi.png','Beban Karyawan','Pekanbaru','1471112905640001',
    '142650043216000','M','Patroli Satpam','III-6','L','K/3','Aktif',1,12,3901700,500000,NULL,NULL,NULL),

    ('111111113','500532','Amos Sibarani','sammidev.codes/sammi.png','Beban Karyawan','Pekanbaru','1471115905640001',
    '142650043216000','M','DAN JAGA','III-9','L','K/3','Aktif',1,12,4750000,500000,NULL,NULL,NULL);
