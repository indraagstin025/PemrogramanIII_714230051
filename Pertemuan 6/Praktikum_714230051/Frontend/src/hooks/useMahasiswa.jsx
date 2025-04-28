import { useEffect, useState, useRef } from "react";
import axios from "axios";

export const useMahasiswa = () => {
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const fetchedRef = useRef(false); // supaya cache tidak fetch 2x

    const fetchData = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await axios.get("http://127.0.0.1:3000/api/mahasiswa");
            setUsers(response.data.data || []); // ambil .data array, kalau ada
            console.log("Fetched mahasiswa:", response.data.data);
        } catch (err) {
            console.error("Error fetching mahasiswa:", err);
            setError(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (!fetchedRef.current) {
            fetchData();
            fetchedRef.current = true; // tandain sudah fetch
        }
    }, []);

    const retry = () => {
        fetchedRef.current = false; // reset cache, agar bisa fetch ulang
        fetchData();
    };

    return { users, loading, error, retry };
};
