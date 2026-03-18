import { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Input, Upload, Card, message } from 'antd';
import { PlusOutlined, DeleteOutlined, UploadOutlined } from '@ant-design/icons';
import { useStore } from '../store';
import type { UploadProps } from 'antd';

export default function Firmwares() {
  const { firmwares, currentProduct, products, fetchProducts, setCurrentProduct, fetchFirmwares, createFirmware, deleteFirmware } = useStore();
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchProducts();
    if (currentProduct?.id) {
      fetchFirmwares(currentProduct.id);
    }
  }, [currentProduct?.id, fetchFirmwares, fetchProducts]);

  const handleOpenModal = () => {
    form.resetFields();
    setModalVisible(true);
  };

  const handleSubmit = async (values: any) => {
    if (!currentProduct) {
      message.error('Please select a product first');
      return;
    }

    const file = fileList?.[0]?.originFileObj;
    if (!file) {
      message.error('Please select a firmware file');
      return;
    }

    const formData = new FormData();
    formData.append('file', file);
    formData.append('product_id', currentProduct.id);
    formData.append('version', values.version);
    formData.append('versionCode', values.versionCode);
    formData.append('changelog', values.changelog || '');
    formData.append('release_notes', values.releaseNotes || '');

    await createFirmware(formData, currentProduct.id);
    setModalVisible(false);
    form.resetFields();
    setFileList([]);
  };

  const handleDelete = async (id: string) => {
    await deleteFirmware(id);
  };

  const [fileList, setFileList] = useState<any>([]);

  const uploadProps: UploadProps = {
    beforeUpload: () => {
      return true;
    },
    onChange: ({ fileList: newFileList }) => {
      setFileList(newFileList);
    },
    fileList,
    maxCount: 1,
  };

  const columns = [
    {
      title: 'Version',
      dataIndex: 'version',
      key: 'version',
    },
    {
      title: 'Version Code',
      dataIndex: 'version_code',
      key: 'version_code',
      width: 100,
    },
    {
      title: 'File Size',
      dataIndex: 'file_size',
      key: 'file_size',
      width: 120,
      render: (size: number) => `${(size / 1024 / 1024).toFixed(2)} MB`,
    },
    {
      title: 'MD5',
      dataIndex: 'md5',
      key: 'md5',
      width: 200,
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
          danger
          size="small"
          icon={<DeleteOutlined />}
          onClick={() => handleDelete(record.id)}
        >
          Delete
        </Button>
      ),
    },
  ];

  return (
    <Card
      title={currentProduct ? `Firmwares - ${currentProduct.name}` : 'Firmwares'}
      extra={
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenModal}>
          Upload Firmware
        </Button>
      }
    >
      <div style={{ marginBottom: 16 }}>
        <label style={{ marginRight: 8 }}>Select Product:</label>
        <select 
          value={currentProduct?.id || ''} 
          onChange={(e) => {
            const product = products.find(p => p.id === e.target.value);
            setCurrentProduct(product || null);
          }}
          style={{ padding: 4, minWidth: 200 }}
        >
          <option value="">- Please select product -</option>
          {products.map(product => (
            <option key={product.id} value={product.id}>
              {product.name}
            </option>
          ))}
        </select>
      </div>
      {!currentProduct ? (
        <div style={{ textAlign: 'center', padding: '30px', color: '#999' }}>
          Please select a product from the dropdown above
        </div>
      ) : (
        <Table
          columns={columns}
          dataSource={firmwares}
          rowKey="id"
          pagination={{ pageSize: 10 }}
        />
      )}
      <Modal
        title="Upload New Firmware"
        open={modalVisible}
        onOk={() => handleSubmit(form.getFieldsValue())}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="version"
            label="Version"
            rules={[{ required: true, message: 'Please enter version' }]}
          >
            <Input placeholder="1.0.0" />
          </Form.Item>
          <Form.Item
            name="versionCode"
            label="Version Code"
            rules={[{ required: true, message: 'Please enter version code' }]}
          >
            <Input type="number" placeholder="1" />
          </Form.Item>
          <Form.Item
            name="changelog"
            label="Changelog"
          >
            <Input.TextArea placeholder="Changes in this version" rows={3} />
          </Form.Item>
          <Form.Item
            name="releaseNotes"
            label="Release Notes"
          >
            <Input.TextArea placeholder="Release notes" rows={3} />
          </Form.Item>
          <Form.Item label="Firmware File">
            <Upload {...uploadProps}>
              <Button icon={<UploadOutlined />}>Select File</Button>
            </Upload>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
}
