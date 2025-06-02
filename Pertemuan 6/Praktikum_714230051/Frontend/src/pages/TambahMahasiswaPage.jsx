import { TypographyAtom } from "../components/atoms/TypographyAtom";
import { MahasiswaForm } from "../components/organisms/MahasiswaForm";
import { postMahasiswa } from "../services/mahasiswaServices";
import { useNavigate } from "react-router-dom";
import Swal from "sweetalert2";

export function TambahMahasiswaPage() {
  const navigate = useNavigate();

  const handleSubmit = async (data) => {
    try {
      await postMahasiswa(data);
      Swal.fire({
        icon: "success",
        title: "Berhasil!",
        text: "Data mahasiswa berhasil disimpan.",
        showConfirmButton: false,
        timer: 1500,
      });
      navigate("/mahasiswa");
    } catch (error) {
      console.error("Gagal menyimpan data mahasiswa:", error);

      let errorMessage = "Gagal menyimpan data mahasiswa. Silakan coba lagi.";

      // Cek apakah error berasal dari respons server (Axios)
      if (error.response && error.response.data && error.response.data.message) {
        // Jika ada pesan error dari backend (misal: "NPM sudah terdaftar")
        errorMessage = error.response.data.message;
      }

      Swal.fire({
        icon: "error",
        title: "Gagal!",
        text: errorMessage, // Tampilkan pesan error dari backend
        showConfirmButton: true,
      });
    }
  };

  const handleCancel = () => {
    navigate("/mahasiswa");
  };

  return (
    <div className="py-6 px-4"> {/* PASTIKAN penulisan className sudah benar */}
      <TypographyAtom variant="h5" className="mb-4">
        Tambah Data Mahasiswa
      </TypographyAtom>
      <MahasiswaForm onSubmit={handleSubmit} onCancel={handleCancel} />
    </div>
  );
}
