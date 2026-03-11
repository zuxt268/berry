import { Card, Tag, Button, Space, Typography } from 'antd';
import type { AIAdvice } from '../types';

const { Text } = Typography;

const priorityConfig = {
  high:   { color: '#f5222d', bg: 'rgba(245,34,45,0.06)',  tagColor: 'red',    label: '高' },
  medium: { color: '#fa8c16', bg: 'rgba(250,140,22,0.06)', tagColor: 'orange', label: '中' },
  low:    { color: '#1677ff', bg: 'rgba(22,119,255,0.06)', tagColor: 'blue',   label: '低' },
} as const;

interface Props {
  advice: AIAdvice;
}

const AIAdviceCard: React.FC<Props> = ({ advice }) => {
  const p = priorityConfig[advice.priority];

  return (
    <Card
      className="advice-card"
      variant="borderless"
      style={{
        marginBottom: 12,
        borderRadius: 10,
        boxShadow: '0 1px 4px rgba(0,0,0,0.06)',
        borderLeft: `3px solid ${p.color}`,
        background: p.bg,
      }}
      styles={{ body: { padding: '16px 20px' } }}
    >
      <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 6 }}>
        <Tag color={p.tagColor} style={{ margin: 0, fontSize: 11, lineHeight: '18px' }}>
          優先度 {p.label}
        </Tag>
      </div>
      <Text strong style={{ fontSize: 14, display: 'block', marginBottom: 6 }}>
        {advice.title}
      </Text>
      <Text type="secondary" style={{ fontSize: 13, display: 'block', lineHeight: 1.7, marginBottom: 14 }}>
        {advice.content}
      </Text>
      <Space size="small">
        <Button size="small">自分でやる</Button>
        <Button size="small" type="primary">
          当社に依頼
        </Button>
      </Space>
    </Card>
  );
};

export default AIAdviceCard;
