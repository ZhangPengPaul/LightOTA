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
  } = useStore();
  const [modalVisible, setModalVisible] = useState(false);
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

  const handleOpenModal = () => {
    form.resetFields();
    setModalVisible(true);
  };

  const handleSubmit = async () => {
    const values = await form.validateFields();
    await createUpgradeTask(values);
    message.success('Upgrade task created');
    setModalVisible(false);
    form.resetFields();
  };

  const handleRowClick = (record: UpgradeTask) => {
    fetchUpgradeTask(record.id);
  };

  const columns = [
    {
      title: 'Name',
      dataIndex: 'taskName',
      key: 'taskName',
    },
    {
      title: 'Product',
      dataIndex: 'productId',
      key: 'productId',
    },
    {
      title: 'Type',
      dataIndex: 'upgradeType',
      key: 'upgradeType',
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
      dataIndex: 'targetDevicesCount',
      key: 'targetDevicesCount',
      width: 100,
    },
    {
      title: 'Created At',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
    },
  ];

  return (
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 350px', gap: '16px' }}>
      <Card title="Upgrade Tasks" extra={<Button type="primary" icon={<PlusOutlined />} onClick={handleOpenModal}>New Task</Button>}>
        <Table
          columns={columns}
          dataSource={upgradeTasks}
          rowKey="id"
          onRow={(record) => ({
            onClick: () => handleRowClick(record),
          })}
          pagination={{ pageSize: 10 }}
        />
      </Card>

      {currentTask && (
        <Card title={`Task: ${currentTask.taskName}`}>
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            <div>
              <div style={{ marginBottom: 8 }}>
                <strong>Status:</strong>{' '}
                <Tag color={statusColors[currentTask.status]}>{currentTask.status}</Tag>
              </div>
              <div><strong>Type:</strong> {currentTask.upgradeType}</div>
              <div><strong>Target Devices:</strong> {currentTask.targetDevicesCount}</div>
              <div><strong>Firmware:</strong> {currentTask.firmwareId}</div>
            </div>

            {currentTask.percent !== undefined && (
              <div>
                <div style={{ marginBottom: 8 }}>
                  <strong>Progress: {currentTask.percent}%</strong>
                </div>
                <Progress percent={currentTask.percent} status="active" />
              </div>
            )}

            {currentTask.successCount !== undefined && (
              <div>
                <div><strong>Success: {currentTask.successCount}</strong></div>
                <div><strong>Failed: {currentTask.failedCount}</strong></div>
                <div><strong>Pending: {currentTask.pendingCount}</strong></div>
              </div>
            )}

            {currentTask.startedAt && <div><strong>Started:</strong> {currentTask.startedAt}</div>}
            {currentTask.completedAt && <div><strong>Completed:</strong> {currentTask.completedAt}</div>}
          </Space>
        </Card>
      )}

      <Modal
        title="Create Upgrade Task"
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
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
                  {f.version} (code: {f.versionCode})
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
            name="pushRate"
            label="Push Rate (per second)"
            initialValue={10}
          >
            <Input type="number" min={1} max={1000} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
