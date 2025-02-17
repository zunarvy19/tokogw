import React from "react";

interface BarangItemProps {
  barang: {id: number; nama:string; harga:number; stok:number}
  onEdit: (barang: {id: number; nama:string; harga:number; stok: number}) => void
  onDelete: (id:number) => void
}

export default function BarangItem ({barang, onEdit, onDelete}: BarangItemProps){
  return(
    <tr>
      <td className="border p-2 text-center">{barang.id}</td>
      <td className="border p-2">{barang.nama}</td>
      <td className="border p-2 text-right">{barang.harga}</td>
      <td className="border p-2 text-center">{barang.stok}</td>
      <td className="border p-2 text-center">
        <button 
        onClick={()=>onEdit(barang)}
        className="bg-yellow-500 text-white px-2 py-1 rounded mr-2"
        >Edit</button>
        <button
        onClick={()=>onDelete(barang.id)}
        className="bg-red-500 text-white px-2 py-1 rounded"
        >
          Delete
        </button>
      </td>
    </tr>
  )
}