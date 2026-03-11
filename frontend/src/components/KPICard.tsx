import { Card, Typography } from 'antd';
import { ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons';
import type { KPIMetric } from '../types';
import { formatValue } from '../utils/format';

const { Text } = Typography;

const LOWER_IS_BETTER_KEYWORDS = ['直帰率', '掲載順位', 'ブロック率'];

interface Props {
  metric: KPIMetric;
  size?: 'default' | 'large';
}

const KPICard: React.FC<Props> = ({ metric, size = 'default' }) => {
  const { label, value, previousValue, format } = metric;
  const change = previousValue !== 0 ? ((value - previousValue) / previousValue) * 100 : 0;
  const isUp = change > 0;

  const lowerIsBetter = LOWER_IS_BETTER_KEYWORDS.some((kw) => label.includes(kw));
  const isPositive = lowerIsBetter ? !isUp : isUp;
  const trendBg = change === 0
    ? 'rgba(0,0,0,0.04)'
    : isPositive
      ? 'rgba(82,196,26,0.08)'
      : 'rgba(245,34,45,0.08)';
  const trendColor = change === 0 ? '#8c8c8c' : isPositive ? '#389e0d' : '#cf1322';
  const Arrow = isUp ? ArrowUpOutlined : ArrowDownOutlined;

  const isLarge = size === 'large';

  return (
    <Card
      className="kpi-card"
      variant="borderless"
      style={{
        height: '100%',
        borderRadius: 10,
        boxShadow: '0 1px 4px rgba(0,0,0,0.06)',
        borderLeft: isLarge ? '4px solid #1677ff' : undefined,
      }}
      styles={{ body: { padding: isLarge ? '20px 24px' : '16px 20px' } }}
    >
      <Text type="secondary" style={{ fontSize: isLarge ? 14 : 13, letterSpacing: 0.3 }}>
        {label}
      </Text>
      <div
        style={{
          fontSize: isLarge ? 40 : 28,
          fontWeight: 700,
          lineHeight: 1.2,
          margin: isLarge ? '10px 0 10px' : '6px 0 8px',
          color: '#141414',
        }}
      >
        {formatValue(value, format)}
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 8, flexWrap: 'wrap' }}>
        <span
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: 3,
            background: trendBg,
            color: trendColor,
            fontSize: 12,
            fontWeight: 600,
            padding: '2px 8px',
            borderRadius: 12,
          }}
        >
          <Arrow style={{ fontSize: 10 }} />
          {Math.abs(change).toFixed(1)}%
        </span>
        <Text type="secondary" style={{ fontSize: 12 }}>
          前月 {formatValue(previousValue, format)}
        </Text>
      </div>
    </Card>
  );
};

export default KPICard;
