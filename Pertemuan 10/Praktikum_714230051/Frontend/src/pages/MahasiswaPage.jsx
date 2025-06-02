import { TableWithStripedRows } from "../components/organisms/TableWithStripedRows";

export function MahasiswaPage() {
  return (
    <div className="py-6">
      <h1 className="text-2xl font-bold mb-6">Data Mahasiswa</h1>
      <TableWithStripedRows />
    </div>
  );
}