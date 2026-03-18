import { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Input, Radio, Select, Progress, Card, Space, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useStore } from '../store';
import type { UpgradeTask } from '../api/client';

const statusColors: Record<string, string> = {
  created: 'blue',
  running: 'orange',
  completed: 'green',
  paused: 'default',
  cancelled: 'red',
};

export default function UpgradeTasks() {
  const {
    upgradeTasks,
    currentTask,
    products,
    firmwares,
    fetchUpgradeTasks,
    fetchUpgradeTask,
    createUpgradeTask,
    fetchProducts,
    fetchFirmwares,
    setCurrentTask,
  } = useStore();
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchUpgradeTasks();
    fetchProducts();
  }, [fetchUpgradeTasks, fetchProducts]);

  const handleProductChange = (productId: string) => {
    if (productId) {
      fetchFirmwares(productId);
    }
  };

  const handleOpenCreateModal = () => {
    form.resetFields();
    setCreateModalVisible(true);
  };

  const handleSubmit = async () => {
    const values = await form.validateFields();
    let target_device_ids: string[] | undefined = undefined;
    if (values.upgradeType === 'specified' && values.targetDeviceIds) {
      target_device_ids = values.targetDeviceIds
        .split('\n')
        .map((id: string) => id.trim())
        .filter((id: string) => id.length > 0);
    }
    const mappedValues = {
      product_id: values.productId,
      firmware_id: values.firmwareId,
      task_name: values.taskName,
      upgrade_type: values.upgradeType,
      gray_percent: values.grayPercent,
      target_device_ids: target_device_ids,
      push_rate: values.pushRate,
    };
    await createUpgradeTask(mappedValues);
    message.success('Upgrade task created');
    setCreateModalVisible(false);
    form.resetFields();
  };

  const handleRowClick = async (record: UpgradeTask) => {
    await fetchUpgradeTask(record.id);
    setDetailModalVisible(true);
  };

  const handleCloseDetail = () => {
    setDetailModalVisible(false);
    setCurrentTask(null);
  };

  const columns = [
    {
      title: 'Name',
      dataIndex: 'task_name',
      key: 'task_name',
      width: 150,
    },
    {
      title: 'Product',
      dataIndex: 'product_id',
      key: 'product_id',
      width: 200,
      render: (productId: string) => {
        const product = products.find(p => p.id === productId);
        return product?.name || productId;
      },
    },
    {
      title: 'Type',
      dataIndex: 'upgrade_type',
      key: 'upgrade_type',
      width: 100,
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => <Tag color={statusColors[status]}>{status}</Tag>,
    },
    {
      title: 'Devices',
      dataIndex: 'target_devices_count',
      key: 'target_devices_count',
      width: 100,
    },
    {
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 220,
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 100,
      render: (_: any, record: UpgradeTask) => (
        <Button type="link" size="small" onClick={() => handleRowClick(record)}>
          View
        </Button>
      ),
    },
  ];

  return (
    <div>
      <Card title="Upgrade Tasks" extra={<Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreateModal}>New Task</Button>}>
        <Table
          columns={columns}
          dataSource={upgradeTasks}
          rowKey="id"
          pagination={{ pageSize: 10 }}
        />
      </Card>

      <Modal
        title="Create Upgrade Task"
        open={createModalVisible}
        onOk={handleSubmit}
        onCancel={() => setCreateModalVisible(false)}
        width={500}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="productId"
            label="Product"
            rules={[{ required: true, message: 'Please select product' }]}
          >
            <Select placeholder="Select product" onChange={handleProductChange}>
              {products.map((p) => (
                <Select.Option key={p.id} value={p.id}>
                  {p.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="firmwareId"
            label="Target Firmware"
            rules={[{ required: true, message: 'Please select firmware' }]}
          >
            <Select placeholder="Select firmware">
              {firmwares.map((f) => (
                <Select.Option key={f.id} value={f.id}>
                  {f.version} (code: {f.version_code})
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="taskName"
            label="Task Name"
            rules={[{ required: true, message: 'Please enter task name' }]}
          >
            <Input placeholder="March 2025 release" />
          </Form.Item>

          <Form.Item
            name="upgradeType"
            label="Upgrade Type"
            rules={[{ required: true }]}
            initialValue="gray"
          >
            <Radio.Group>
              <Radio value="specified">Specified Devices</Radio>
              <Radio value="all">All Devices</Radio>
              <Radio value="gray">Gray Percentage</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item
            noStyle
            shouldUpdate={(prev, curr) => prev.upgradeType !== curr.upgradeType}
          >
            {({ getFieldValue }) =>
              getFieldValue('upgradeType') === 'gray' && (
                <Form.Item
                  name="grayPercent"
                  label="Gray Percentage"
                  rules={[{ required: true }]}
                  initialValue={10}
                >
                  <Input type="number" min={1} max={100} />
                </Form.Item>
              )
            }
          </Form.Item>

          <Form.Item
            noStyle
            shouldUpdate={(prev, curr) => prev.upgradeType !== curr.upgradeType}
          >
            {({ getFieldValue }) =>
              getFieldValue('upgradeType') === 'specified' && (
                <Form.Item
                  name="targetDeviceIds"
                  label="Target Device IDs"
                  rules={[{ required: true, message: 'Please enter target device IDs' }]}
                >
                  <Input.TextArea 
                    placeholder="One device ID per line" 
                    rows={6}
                  />
                </Form.Item>
              )
            }
          </Form.Item>

          <Form.Item
            name="pushRate"
            label="Push Rate (per second)"
            initialValue={10}
          >
            <Input type="number" min={1} max={1000} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={currentTask?.task_name}
        open={detailModalVisible}
        onCancel={handleCloseDetail}
        footer={null}
        width={500}
      >
        {currentTask && (
          <Space direction="vertical" size="large" style={{ width: '100%', wordBreak: 'break-word' }}>
            <div>
              <div style={{ marginBottom: 8 }}>
                <strong>Status:</strong>{' '}
                <Tag color={statusColors[currentTask.status]}>{currentTask.status}</Tag>
              </div>
              <div><strong>Type:</strong> {currentTask.upgrade_type}</div>
              <div><strong>Target Devices:</strong> {currentTask.target_devices_count}</div>
              <div><strong>Firmware:</strong> {currentTask.firmware_id}</div>
            </div>

            {currentTask.percent !== undefined && (
              <div>
                <div style={{ marginBottom: 8 }}>
                  <strong>Progress: {currentTask.percent}%</strong>
                </div>
                <Progress percent={currentTask.percent} status="active" />
              </div>
            )}

            {currentTask.success_count !== undefined && (
              <div>
                <div><strong>Success: {currentTask.success_count}</strong></div>
                <div><strong>Failed: {currentTask.failed_count}</strong></div>
                <div><strong>Pending: {currentTask.pending_count}</strong></div>
              </div>
            )}

            {currentTask.started_at && <div><strong>Started:</strong> {currentTask.started_at}</div>}
            {currentTask.completed_at && <div><strong>Completed:</strong> {currentTask.completed_at}</div>}
          </Space>
        )}
      </Modal>
    </div>
  );
}
