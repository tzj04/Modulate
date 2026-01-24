// src/api/modules.ts
import { client } from "./client";
import { Module } from "../types/module";

export const moduleApi = {
  // Get all modules for the list view
  getAll: () => client<Module[]>("/modules"),

  // Get a single module by its code (e.g., CS1231S)
  getById: (id: string) => client<Module>(`/modules/${id}`),
};
