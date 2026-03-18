import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
});

export function setApiKey(apiKey: string) {
  if (apiKey) {
    api.defaults.headers.common.Authorization = `Bearer ${apiKey}`;
  }
}

export interface Tenant {
  id: string;
  name: string;
  api_key: string;
  external_device_api_url: string;
  created_at: string;
  updated_at: string;
}

export interface Product {
  id: string;
  tenant_id: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Firmware {
  id: string;
  tenant_id: string;
  product_id: string;
  version: string;
  version_code: number;
  changelog: string;
  file_size: number;
  md5: string;
  release_notes: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface UpgradeTask {
  id: string;
  tenant_id: string;
  product_id: string;
  firmware_id: string;
  task_name: string;
  upgrade_type: 'specified' | 'all' | 'gray';
  gray_percent: number;
  target_devices_count: number;
  push_rate: number;
  status: 'created' | 'running' | 'paused' | 'completed' | 'cancelled';
  created_by: string;
  started_at: string | null;
  completed_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface TaskStats {
  total: number;
  success_count: number;
  failed_count: number;
  pending_count: number;
  percent: number;
  status_counts: Record<string, number>;
};

export const apiClient = {
  // Tenant
  listTenants: async (): Promise<{list: Tenant[]; total: number}> => {
    const res = await api.get('/tenants');
    return res.data.data;
  },
  createTenant: async (data: {name: string; external_device_api_url: string}) => {
    const res = await api.post('/tenants', data);
    return res.data;
  },
  updateTenant: async (id: string, data: {name?: string; external_device_api_url?: string}) => {
    const res = await api.put(`/tenants/${id}`, data);
    return res.data;
  },

  // Product
  listProducts: async (params?: {limit?: number; offset?: number}): Promise<{list: Product[]; total: number}> => {
    const res = await api.get('/products', {params});
    return res.data.data;
  },
  createProduct: async (data: {name: string; description: string}) => {
    const res = await api.post('/products', data);
    return res.data;
  },
  updateProduct: async (id: string, data: {name?: string; description?: string}) => {
    const res = await api.put(`/products/${id}`, data);
    return res.data;
  },
  deleteProduct: async (id: string) => {
    const res = await api.delete(`/products/${id}`);
    return res.data;
  },

  // Firmware
  listFirmwares: async (productId: string, params?: {limit?: number; offset?: number}): Promise<{list: Firmware[]; total: number}> => {
    const res = await api.get('/firmwares', {params: {productId, ...params}});
    return res.data.data;
  },
  createFirmware: async (formData: FormData) => {
    const res = await api.post('/firmwares', formData, {
      headers: {'Content-Type': 'multipart/form-data'},
    });
    return res.data;
  },
  deleteFirmware: async (id: string) => {
    const res = await api.delete(`/firmwares/${id}`);
    return res.data;
  },

  // Upgrade Task
  createUpgradeTask: async (data: {
    product_id: string;
    firmware_id: string;
    task_name: string;
    upgrade_type: 'specified' | 'all' | 'gray';
    gray_percent?: number;
    target_device_ids?: string[];
    push_rate?: number;
  }) => {
    const res = await api.post('/upgrade/task', data);
    return res.data;
  },
  listUpgradeTasks: async (params?: {productId?: string; limit?: number; offset?: number}): Promise<{list: UpgradeTask[]; total: number}> => {
    const res = await api.get('/upgrade/tasks', {params});
    return res.data.data;
  },
  getUpgradeTask: async (id: string): Promise<{
    task: UpgradeTask;
  } & TaskStats> => {
    const res = await api.get(`/upgrade/task/${id}`);
    return res.data.data;
  },
};

export default api;
