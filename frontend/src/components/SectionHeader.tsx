import { Typography } from 'antd';

const { Title, Text } = Typography;

interface Props {
  title: string;
  subtitle?: string;
}

const SectionHeader: React.FC<Props> = ({ title, subtitle }) => (
  <div style={{ marginBottom: 16 }}>
    <Title level={5} style={{ margin: 0 }}>
      {title}
    </Title>
    {subtitle && (
      <Text type="secondary" style={{ fontSize: 12 }}>
        {subtitle}
      </Text>
    )}
  </div>
);

export default SectionHeader;
