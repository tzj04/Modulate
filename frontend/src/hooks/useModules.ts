// src/hooks/useModules.ts
import { useState, useEffect } from "react";
import { moduleApi } from "../api/modules";
import { Module } from "../types/module";

export const useModules = () => {
  const [modules, setModules] = useState<Module[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchModules = async () => {
      try {
        setLoading(true);
        const data = await moduleApi.getAll();
        setModules(data);
      } catch (err: any) {
        setError(err.message || "Failed to fetch modules");
      } finally {
        setLoading(false);
      }
    };

    fetchModules();
  }, []);

  return { modules, loading, error };
};
