import axios from "axios";

const BASE_URL = "http://localhost:5000/barang";

export const fetchBarang = async () => {
  const response = await axios.get(BASE_URL);
  return response.data;
};

export const addBarang = async (barang: { nama: string; harga: number; stok: number }) => {
  await axios.post(BASE_URL, barang);
};

export const updateBarang = async (id: number, barang: { nama: string; harga: number; stok: number }) => {
  await axios.put(`${BASE_URL}?id=${id}`, barang);
};

export const deleteBarang = async (id: number) => {
  await axios.delete(`${BASE_URL}?id=${id}`);
};
