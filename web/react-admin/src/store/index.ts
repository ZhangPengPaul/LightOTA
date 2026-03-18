import { create } from 'zustand';
import type { Tenant, Product, Firmware, UpgradeTask, TaskStats } from '../api/client';
import { apiClient } from '../api/client';

interface AppState {
  // Tenants
  tenants: Tenant[];
  fetchTenants: () => Promise<void>;
  createTenant: (data: {name: string; external_device_api_url: string}) => Promise<void>;
  updateTenant: (id: string, data: {name?: string; external_device_api_url?: string}) => Promise<void>;

  // Products
  products: Product[];
  currentProduct: Product | null;
  fetchProducts: () => Promise<void>;
  createProduct: (data: {name: string; description: string}) => Promise<void>;
  updateProduct: (id: string, data: {name?: string; description?: string}) => Promise<void>;
  deleteProduct: (id: string) => Promise<void>;
  setCurrentProduct: (product: Product | null) => void;

  // Firmwares
  firmwares: Firmware[];
  fetchFirmwares: (product_id: string) => Promise<void>;
  createFirmware: (formData: FormData, product_id: string) => Promise<void>;
  deleteFirmware: (id: string) => Promise<void>;

  // Upgrade Tasks
  upgradeTasks: UpgradeTask[];
  currentTask: (UpgradeTask & TaskStats) | null;
  fetchUpgradeTasks: (product_id?: string) => Promise<void>;
  createUpgradeTask: (data: {
    product_id: string;
    firmware_id: string;
    task_name: string;
    upgrade_type: 'specified' | 'all' | 'gray';
    gray_percent?: number;
    target_device_ids?: string[];
    push_rate?: number;
  }) => Promise<{taskId: string}>;
  fetchUpgradeTask: (id: string) => Promise<void>;
  setCurrentTask: (task: (UpgradeTask & TaskStats) | null) => void;
}

export const useStore = create<AppState>((set, get) => ({
  tenants: [],
  fetchTenants: async () => {
    const res = await apiClient.listTenants();
    set({ tenants: res.list });
  },
  createTenant: async (data) => {
    await apiClient.createTenant(data);
    await get().fetchTenants();
  },
  updateTenant: async (id, data) => {
    await apiClient.updateTenant(id, data);
    await get().fetchTenants();
  },

  products: [],
  currentProduct: null,
  fetchProducts: async () => {
    const res = await apiClient.listProducts();
    set({ products: res.list });
  },
  createProduct: async (data) => {
    await apiClient.createProduct(data);
    await get().fetchProducts();
  },
  updateProduct: async (id, data) => {
    await apiClient.updateProduct(id, data);
    await get().fetchProducts();
  },
  deleteProduct: async (id) => {
    await apiClient.deleteProduct(id);
    await get().fetchProducts();
  },
  setCurrentProduct: (product) => {
    set({ currentProduct: product });
  },

  firmwares: [],
  fetchFirmwares: async (product_id) => {
    const res = await apiClient.listFirmwares(product_id);
    set({ firmwares: res.list });
  },
  createFirmware: async (formData, productId) => {
    await apiClient.createFirmware(formData);
    await get().fetchFirmwares(productId);
  },
  deleteFirmware: async (id) => {
    await apiClient.deleteFirmware(id);
    const productId = get().currentProduct?.id;
    if (productId) {
      await get().fetchFirmwares(productId);
    }
  },

  upgradeTasks: [],
  currentTask: null,
  fetchUpgradeTasks: async (productId) => {
    const res = await apiClient.listUpgradeTasks(productId ? { productId } : undefined);
    set({ upgradeTasks: res.list });
  },
  createUpgradeTask: async (data) => {
    const res = await apiClient.createUpgradeTask(data);
    await get().fetchUpgradeTasks(data.product_id);
    return res.data;
  },
  fetchUpgradeTask: async (id) => {
    const res = await apiClient.getUpgradeTask(id);
    set({ currentTask: { ...res.task, ...res } });
  },
  setCurrentTask: (task) => {
    set({ currentTask: task });
  },
}));
