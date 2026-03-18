import { Layout, Menu } from 'antd';
import { Outlet, Link, useLocation } from 'react-router-dom';
import { ApartmentOutlined, ContainerOutlined, RocketOutlined, CloudUploadOutlined } from '@ant-design/icons';
import './App.css';

const { Header, Content, Sider } = Layout;

function App() {
  const location = useLocation();

  const menuItems = [
    {
      key: '/tenants',
      icon: <ApartmentOutlined />,
      label: <Link to="/tenants">Tenants</Link>,
    },
    {
      key: '/products',
      icon: <ContainerOutlined />,
      label: <Link to="/products">Products</Link>,
    },
    {
      key: '/firmwares',
      icon: <CloudUploadOutlined />,
      label: <Link to="/firmwares">Firmwares</Link>,
    },
    {
      key: '/tasks',
      icon: <RocketOutlined />,
      label: <Link to="/tasks">Upgrade Tasks</Link>,
    },
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', background: '#001529' }}>
        <div style={{ color: 'white', fontSize: 18, fontWeight: 'bold' }}>
          LightOTA Admin
        </div>
      </Header>
      <Layout>
        <Sider width={200} theme="light">
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            style={{ height: '100%', borderRight: 0 }}
          />
        </Sider>
        <Layout>
          <Content style={{ padding: 16, background: '#f0f2f5' }}>
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
}

export default App;
