import { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Card } from 'antd';
import { PlusOutlined, EditOutlined } from '@ant-design/icons';
import { useStore } from '../store';

export default function Tenants() {
  const { tenants, fetchTenants, createTenant, updateTenant } = useStore();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingTenant, setEditingTenant] = useState<any>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchTenants();
  }, [fetchTenants]);

  const handleOpenModal = (tenant?: any) => {
    setEditingTenant(tenant || null);
    if (tenant) {
      form.setFieldsValue(tenant);
    } else {
      form.resetFields();
    }
    setModalVisible(true);
  };

  const handleSubmit = async () => {
    const values = await form.validateFields();
    if (editingTenant) {
      await updateTenant(editingTenant.id, values);
    } else {
      await createTenant(values);
    }
    setModalVisible(false);
    form.resetFields();
  };

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'API Key',
      dataIndex: 'api_key',
      key: 'api_key',
    },
    {
      title: 'External API URL',
      dataIndex: 'external_device_api_url',
      key: 'external_device_api_url',
    },
    {
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 100,
      render: (_: any, record: any) => (
        <Button
          type="link"
          size="small"
          icon={<EditOutlined />}
          onClick={() => handleOpenModal(record)}
        >
          Edit
        </Button>
      ),
    },
  ];

  return (
    <Card title="Tenants" extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => handleOpenModal()}>New Tenant</Button>}>
      <Table
        columns={columns}
        dataSource={tenants}
        rowKey="id"
        pagination={false}
      />
      <Modal
        title={editingTenant ? 'Edit Tenant' : 'New Tenant'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: 'Please enter tenant name' }]}
          >
            <input placeholder="Tenant name" />
          </Form.Item>
          <Form.Item
            name="external_device_api_url"
            label="External Device API URL"
          >
            <input placeholder="https://your-device-api.com/api" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
}
