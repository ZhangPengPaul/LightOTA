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
  apiKey: string;
  externalDeviceAPIUrl: string;
  createdAt: string;
  updatedAt: string;
}

export interface Product {
  id: string;
  tenantId: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
}

export interface Firmware {
  id: string;
  tenantId: string;
  productId: string;
  version: string;
  versionCode: number;
  changelog: string;
  fileSize: number;
  md5: string;
  releaseNotes: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UpgradeTask {
  id: string;
  tenantId: string;
  productId: string;
  firmwareId: string;
  taskName: string;
  upgradeType: 'specified' | 'all' | 'gray';
  grayPercent: number;
  targetDevicesCount: number;
  pushRate: number;
  status: 'created' | 'running' | 'paused' | 'completed' | 'cancelled';
  createdBy: string;
  startedAt: string | null;
  completedAt: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface TaskStats {
  total: number;
  successCount: number;
  failedCount: number;
  pendingCount: number;
  percent: number;
  statusCounts: Record<string, number>;
};

export const apiClient = {
  // Tenant
  listTenants: async (): Promise<{list: Tenant[]; total: number}> => {
    const res = await api.get('/tenants');
    return res.data.data;
  },
  createTenant: async (data: {name: string; externalDeviceApiUrl: string}) => {
    const res = await api.post('/tenants', data);
    return res.data;
  },
  updateTenant: async (id: string, data: {name?: string; externalDeviceApiUrl?: string}) => {
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
    productId: string;
    firmwareId: string;
    taskName: string;
    upgradeType: 'specified' | 'all' | 'gray';
    grayPercent?: number;
    targetDeviceIds?: string[];
    pushRate?: number;
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
