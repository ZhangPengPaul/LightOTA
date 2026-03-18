import { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Card } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useStore } from '../store';

export default function Products() {
  const { products, fetchProducts, createProduct, updateProduct, deleteProduct, setCurrentProduct } = useStore();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingProduct, setEditingProduct] = useState<any>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  const handleOpenModal = (product?: any) => {
    setEditingProduct(product || null);
    if (product) {
      form.setFieldsValue(product);
      setCurrentProduct(product);
    } else {
      form.resetFields();
      setCurrentProduct(null);
    }
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    await deleteProduct(id);
  };

  const handleSubmit = async () => {
    const values = await form.validateFields();
    if (editingProduct) {
      await updateProduct(editingProduct.id, values);
    } else {
      await createProduct(values);
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
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
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
      width: 150,
      render: (_: any, record: any) => (
        <div style={{ display: 'flex', gap: '8px' }}>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleOpenModal(record)}
          >
            Edit
          </Button>
          <Button
            type="link"
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
          >
            Delete
          </Button>
        </div>
      ),
    },
  ];

  return (
    <Card title="Products" extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => handleOpenModal()}>New Product</Button>}>
      <Table
        columns={columns}
        dataSource={products}
        rowKey="id"
        pagination={{ pageSize: 10 }}
      />
      <Modal
        title={editingProduct ? 'Edit Product' : 'New Product'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: 'Please enter product name' }]}
          >
            <input placeholder="Product name" />
          </Form.Item>
          <Form.Item
            name="description"
            label="Description"
          >
            <textarea placeholder="Product description" rows={3} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
}
