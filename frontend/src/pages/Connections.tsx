import { Typography, Card, Tag, Button, Row, Col, Space } from 'antd';
import {
  CheckCircleFilled,
  MinusCircleFilled,
  LinkOutlined,
  DisconnectOutlined,
} from '@ant-design/icons';
import { platformConnections } from '../constants/mockData';
import type { PlatformConnection } from '../types';

const { Title, Text } = Typography;

const ConnectionCard: React.FC<{ platform: PlatformConnection }> = ({ platform }) => {
  const { connected } = platform;

  return (
    <Card
      className="connection-card"
      variant="borderless"
      style={{
        borderRadius: 10,
        boxShadow: '0 1px 4px rgba(0,0,0,0.06)',
        background: connected ? '#fff' : '#fafafa',
        borderLeft: connected ? '3px solid #52c41a' : '3px solid #e0e0e0',
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: 12 }}>
        <div style={{ flex: 1, minWidth: 0 }}>
          <Space align="center" style={{ marginBottom: 8 }}>
            {connected ? (
              <CheckCircleFilled style={{ color: '#52c41a', fontSize: 16 }} />
            ) : (
              <MinusCircleFilled style={{ color: '#d9d9d9', fontSize: 16 }} />
            )}
            <Text
              strong
              style={{ fontSize: 15, color: connected ? '#141414' : '#8c8c8c' }}
            >
              {platform.name}
            </Text>
            <Tag color={connected ? 'green' : 'default'} style={{ fontSize: 11 }}>
              {connected ? '連携済み' : '未連携'}
            </Tag>
          </Space>

          <Text
            type="secondary"
            style={{ display: 'block', fontSize: 13, marginBottom: connected && platform.accountLabel ? 6 : 0 }}
          >
            {platform.description}
          </Text>

          {connected && platform.accountLabel && (
            <Text style={{ fontSize: 13, color: '#1677ff', fontWeight: 500 }}>
              {platform.accountLabel}
            </Text>
          )}
        </div>

        <div style={{ flexShrink: 0 }}>
          {connected ? (
            <Button icon={<DisconnectOutlined />} danger size="small">
              連携解除
            </Button>
          ) : (
            <Button type="primary" icon={<LinkOutlined />}>
              連携する
            </Button>
          )}
        </div>
      </div>
    </Card>
  );
};

const Connections: React.FC = () => (
  <>
    <Title level={4} style={{ marginBottom: 24 }}>
      連携設定
    </Title>
    <Row gutter={[16, 16]}>
      {platformConnections.map((p) => (
        <Col xs={24} md={12} key={p.key}>
          <ConnectionCard platform={p} />
        </Col>
      ))}
    </Row>
  </>
);

export default Connections;
