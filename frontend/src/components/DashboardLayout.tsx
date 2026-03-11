import { useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { Layout, Menu, Typography, Grid } from 'antd';
import {
  DashboardOutlined,
  ApiOutlined,
  SettingOutlined,
} from '@ant-design/icons';

const { Sider, Content } = Layout;
const { useBreakpoint } = Grid;

const menuItems = [
  { key: 'dashboard', icon: <DashboardOutlined />, label: 'ダッシュボード' },
  { key: 'connections', icon: <ApiOutlined />, label: '連携設定' },
  { type: 'divider' as const },
  { key: 'settings', icon: <SettingOutlined />, label: '設定' },
];

const DashboardLayout: React.FC = () => {
  const screens = useBreakpoint();
  const isMobile = !screens.lg;
  const [collapsed, setCollapsed] = useState(isMobile);
  const location = useLocation();
  const navigate = useNavigate();

  const selectedKey = location.pathname.split('/')[1] || 'dashboard';

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        collapsible
        collapsed={collapsed}
        onCollapse={setCollapsed}
        breakpoint="lg"
        collapsedWidth={isMobile ? 0 : 80}
        theme="light"
        style={{
          borderRight: '1px solid #f0f0f0',
          position: isMobile ? 'fixed' : 'relative',
          height: isMobile ? '100vh' : 'auto',
          zIndex: 100,
        }}
      >
        <div
          style={{
            padding: '16px 0',
            textAlign: 'center',
            fontWeight: 700,
            fontSize: collapsed ? 18 : 22,
            color: '#1677ff',
            borderBottom: '1px solid #f0f0f0',
            marginBottom: 4,
          }}
        >
          {collapsed ? 'B' : 'Berry'}
        </div>
        <Menu
          mode="inline"
          selectedKeys={[selectedKey]}
          onClick={({ key }) => {
            navigate(`/${key}`);
            if (isMobile) setCollapsed(true);
          }}
          items={menuItems}
          style={{ borderRight: 'none' }}
        />
      </Sider>

      <Layout>
        {isMobile && (
          <Layout.Header
            style={{
              background: '#fff',
              padding: '0 16px',
              borderBottom: '1px solid #f0f0f0',
              display: 'flex',
              alignItems: 'center',
              height: 48,
            }}
          >
            <Typography.Title level={5} style={{ margin: 0, color: '#1677ff' }}>
              Berry
            </Typography.Title>
          </Layout.Header>
        )}
        <Content
          style={{
            padding: isMobile ? 16 : 24,
            margin: 0,
            minHeight: 280,
            overflow: 'auto',
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default DashboardLayout;
