import { Card, Button, Typography } from 'antd';
import { LinkOutlined } from '@ant-design/icons';

const { Text } = Typography;

interface Props {
  message?: string;
}

const UnconnectedCard: React.FC<Props> = ({
  message = 'このサービスはまだ連携されていません',
}) => (
  <Card
    variant="borderless"
    style={{
      background: '#fafafa',
      border: '1.5px dashed #e0e0e0',
      borderRadius: 10,
      textAlign: 'center',
    }}
    styles={{ body: { padding: '32px 16px' } }}
  >
    <Text type="secondary" style={{ display: 'block', marginBottom: 16, fontSize: 13 }}>
      {message}
    </Text>
    <Button type="primary" ghost icon={<LinkOutlined />}>
      連携する
    </Button>
  </Card>
);

export default UnconnectedCard;
