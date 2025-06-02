import { Card } from "@material-tailwind/react";
import { ButtonAtom } from "../atoms/ButtomAtom";
import { TypographyAtom } from "../atoms/TypographyAtom";
import { useMahasiswa } from "../../hooks/useMahasiswa";
import { deleteMahasiswa } from "../../services/mahasiswaServices";
import { useNavigate } from "react-router-dom";
import Swal from "sweetalert2";

const TABLE_HEAD = ["NPM", "Name", "Prodi", "Fakultas", "Minat", "Mata Kuliah", "Aksi"];

export function TableWithStripedRows() {
  const { users, loading, error, retry } = useMahasiswa();
  const navigate = useNavigate(); // <-- Tambahkan ini

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <TypographyAtom variant="h6" color="gray">
          Loading...
        </TypographyAtom>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col justify-center items-center h-64 space-y-4">
        <TypographyAtom variant="h6" color="red">
          Gagal mengambil data mahasiswa.
        </TypographyAtom>
        <ButtonAtom color="red" onClick={retry}>
          Coba Lagi
        </ButtonAtom>
      </div>
    );
  }

  const handleDelete = async (npm) => {
    const result = await Swal.fire({
        title: "Yakin hapus data?",
        text: "Data yang dihapus tidak dapat dikembalikan.",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#dc2626", // merah
        cancelButtonColor: "#6b7280", // abu
        confirmButtonText: "Ya, hapus!",
        cancelButtonText: "Batal"
    });

    if (result.isConfirmed) {
        try {
            await deleteMahasiswa(npm);
            Swal.fire("Terhapus!", "Data berhasil dihapus.", "success");
            retry(); // refresh data
        // eslint-disable-next-line no-unused-vars
        } catch (err) {
            Swal.fire("Gagal!", "Terjadi kesalahan saat menghapus.", "error");
        }
    }
};

  return (
    <Card className="h-full w-100 overflow-x-scroll p-6">
      <div className="flex justify-end p-4">
        <ButtonAtom color="blue" onClick={() => navigate("/mahasiswa/tambah")}>
          Tambah Data
        </ButtonAtom>
      </div>

      <table className="w-full table-auto text-left">
        <thead>
          <tr>
            {TABLE_HEAD.map((head) => (
              <th key={head} className="border-b border-blue-gray-100 bg-blue-gray-50 p-4">
                <TypographyAtom variant="small" color="blue-gray" className="font-normal leading-none opacity-70">
                  {head}
                </TypographyAtom>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <tr key={user._id} className="even:bg-blue-gray-50/50 align-top">
              <td className="p-4">{user.npm}</td>
              <td className="p-4">{user.nama}</td>
              <td className="p-4">{user.prodi}</td>
              <td className="p-4">{user.fakultas}</td>
              <td className="p-4">
                <ul className="list-disc">
                  {user.minat.map((minat, i) => {
                    return <li key={i}>{minat}</li>;
                  })}
                </ul>
              </td>
              <td className="p-4">
                <ul className="list-disc">
                  {user.mata_kuliah.map((matkul) => {
                    return (
                      <li key={matkul.kode}>
                        {matkul.nama} ({matkul.kode}) - Nilai: {matkul.nilai}
                      </li>
                    );
                  })}
                </ul>
              </td>
              <td className="p-4">
                <div className="flex flex-col space-y-2">
                  <ButtonAtom color="green" onClick={() => navigate(`/mahasiswa/edit/${user.npm}`)}>
                    Edit
                  </ButtonAtom>
                  <ButtonAtom color="red" onClick={() => handleDelete(user.npm)}>
                    Hapus
                  </ButtonAtom>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </Card>
  );
}
