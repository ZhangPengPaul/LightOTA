import { useState } from 'react';
import { Button, Card, Form, Input, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { setApiKey } from '../api/client';
import { useStore } from '../store';

export default function Login() {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { fetchTenants } = useStore();

  const onFinish = async (values: { apiKey: string }) => {
    setLoading(true);
    try {
      setApiKey(values.apiKey);
      localStorage.setItem('api_key', values.apiKey);
      await fetchTenants();
      message.success('Login successful');
      navigate('/tenants');
    } catch (err: any) {
      message.error('Login failed: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: '100vh',
      background: '#f0f2f5',
    }}>
      <Card title="LightOTA Admin Login" style={{ width: 400 }}>
        <Form
          name="login"
          onFinish={onFinish}
          initialValues={{ apiKey: localStorage.getItem('api_key') || '' }}
        >
          <Form.Item
            name="apiKey"
            label="API Key"
            rules={[{ required: true, message: 'Please enter your API Key' }]}
          >
            <Input.Password placeholder="Enter your tenant API Key" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              Login
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}
