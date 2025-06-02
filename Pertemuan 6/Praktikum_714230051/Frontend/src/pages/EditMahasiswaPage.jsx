import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { TypographyAtom } from "../components/atoms/TypographyAtom";
import { MahasiswaForm } from "../components/organisms/MahasiswaForm";
import { getMahasiswaByNpm, updateMahasiswa } from "../services/mahasiswaServices";
import Swal from "sweetalert2";

export function EditMahasiswaPage() {
  const { npm } = useParams();
  const navigate = useNavigate();
  const [initialData, setInitialData] = useState(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const data = await getMahasiswaByNpm(npm);
        setInitialData(data);
      // eslint-disable-next-line no-unused-vars
      } catch (err) {
        Swal.fire("Error", "Gagal mengambil data mahasiswa.", "error");
      }
    }
    fetchData();
  }, [npm]);

  const handleSubmit = async (data) => {
    try {
      await updateMahasiswa(npm, data);
      await Swal.fire({
        icon: "success",
        title: "Berhasil!",
        text: "Data mahasiswa berhasil diperbarui.",
        confirmButtonColor: "#22c55e",
      });
      navigate("/mahasiswa");
    // eslint-disable-next-line no-unused-vars
    } catch (err) {
      Swal.fire("Error", "Gagal memperbarui data mahasiswa.", "error");
    }
  };

  const handleCancel = () => {
    navigate("/mahasiswa");
  };

  if (!initialData) return <div>Loading...</div>;

  return (
    <div className="py-6 px-4">
      <TypographyAtom variant="h5" className="mb-4">
        Edit Data Mahasiswa
      </TypographyAtom>
      <MahasiswaForm onSubmit={handleSubmit} initialValues={initialData} isEditing onCancel={handleCancel} />
    </div>
  );
}