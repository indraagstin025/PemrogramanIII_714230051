import axios from "axios";

const API_BASE_URL = "http://127.0.0.1:3000/api/mahasiswa";

export const getAllMahasiswa = async () => {
    try {
        const response = await axios.get(API_BASE_URL);
        return response.data.data || [];
    } catch (error) {
        console.error("Error fetching all mahasiswa:", error.response ? error.response.data : error.message);
        throw error; // Lempar error agar bisa ditangkap di komponen
    }
}

export const postMahasiswa = async (payload) => {
    try {
        console.log("Sending payload:", payload); // Log payload yang dikirim
        const response = await axios.post(API_BASE_URL, payload);
        console.log("Response from server (postMahasiswa):", response.data); // Log respons sukses
        return response.data;
    } catch (error) {
        // Log detail error dari server jika ada
        if (error.response) {
            // Server merespons dengan status selain 2xx
            console.error("Server responded with error (postMahasiswa):", error.response.data);
            console.error("Status:", error.response.status);
            console.error("Headers:", error.response.headers);
        } else if (error.request) {
            // Permintaan dibuat tapi tidak ada respons
            console.error("No response received (postMahasiswa):", error.request);
        } else {
            // Kesalahan lain
            console.error("Error setting up request (postMahasiswa):", error.message);
        }
        throw error; // Penting: tetap lempar error agar ditangkap oleh `try...catch` di `TambahMahasiswaPage`
    }
}

export const getMahasiswaByNpm = async (npm) => {
    const response = await axios.get(`http://127.0.0.1:3000/api/mahasiswa/${npm}`);
    return response.data.data;
};

export const updateMahasiswa = async (npm, payload) => {
    const response = await axios.put(`http://127.0.0.1:3000/api/mahasiswa/${npm}`, payload);
    return response.data;
};

export const deleteMahasiswa = async (npm) => {
  const response = await axios.delete(`http://127.0.0.1:3000/api/mahasiswa/${npm}`);
  return response.data;
};
