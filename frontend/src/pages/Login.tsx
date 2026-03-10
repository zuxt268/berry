import { Button, Card, Typography, Space } from 'antd';
import { GoogleOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

const Login = () => {
  const handleUserLogin = () => {
    window.location.href = '/api/users/auth/google/login';
  };

  const handleOperatorLogin = () => {
    window.location.href = '/api/operators/auth/google/login';
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      }}
    >
      <Card
        style={{
          width: 400,
          textAlign: 'center',
          borderRadius: 12,
          boxShadow: '0 8px 24px rgba(0, 0, 0, 0.15)',
        }}
      >
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <div>
            <Title level={2} style={{ marginBottom: 4 }}>
              Market Pilot
            </Title>
            <Text type="secondary">アカウントにログイン</Text>
          </div>

          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <Button
              type="primary"
              icon={<GoogleOutlined />}
              size="large"
              block
              onClick={handleUserLogin}
            >
              Google でログイン
            </Button>

            <Button
              icon={<GoogleOutlined />}
              size="large"
              block
              onClick={handleOperatorLogin}
            >
              管理者としてログイン
            </Button>
          </Space>
        </Space>
      </Card>
    </div>
  );
};

export default Login;