import { Row, Col, Typography, Tag, Divider, Space } from 'antd';
import type { ChannelSection } from '../types';
import KPICard from '../components/KPICard';
import DailyChart from '../components/DailyChart';
import AIAdviceCard from '../components/AIAdviceCard';
import SectionHeader from '../components/SectionHeader';
import UnconnectedCard from '../components/UnconnectedCard';
import {
  resultSummaryKPIs,
  accessOverviewKPIs,
  accessDailyData,
  accessChartLines,
  channelSections,
  aiAdvices,
} from '../constants/mockData';

const { Title, Text } = Typography;

const ChannelBlock: React.FC<{ channel: ChannelSection }> = ({ channel }) => (
  <div style={{ marginBottom: 24 }}>
    <Space align="center" style={{ marginBottom: 12 }}>
      <Text strong style={{ fontSize: 14, color: channel.connected ? '#141414' : '#bfbfbf' }}>
        {channel.name}
      </Text>
      {channel.accountLabel && (
        <Tag color="blue" style={{ fontSize: 11 }}>
          {channel.accountLabel}
        </Tag>
      )}
      {!channel.connected && <Tag>未連携</Tag>}
    </Space>

    {channel.connected ? (
      <Row gutter={[12, 12]}>
        {channel.kpis.map((kpi) => (
          <Col xs={12} sm={8} md={6} key={kpi.key}>
            <KPICard metric={kpi} />
          </Col>
        ))}
      </Row>
    ) : (
      <UnconnectedCard />
    )}
  </div>
);

const Dashboard: React.FC = () => (
  <>
    <Title level={4} style={{ marginBottom: 24 }}>
      ダッシュボード
    </Title>

    {/* ① 成果サマリー */}
    <SectionHeader title="成果サマリー" subtitle="今月のコンバージョン成果" />
    <Row gutter={[16, 16]} style={{ marginBottom: 32 }}>
      {resultSummaryKPIs.map((kpi) => (
        <Col xs={12} md={12} key={kpi.key}>
          <KPICard metric={kpi} size="large" />
        </Col>
      ))}
    </Row>

    <Divider style={{ margin: '8px 0 24px' }} />

    {/* ② アクセス概況 */}
    <SectionHeader title="アクセス概況" subtitle="GA4 + Search Console" />
    <Row gutter={[12, 12]} style={{ marginBottom: 16 }}>
      {accessOverviewKPIs.map((kpi) => (
        <Col xs={12} sm={8} md={4} key={kpi.key}>
          <KPICard metric={kpi} />
        </Col>
      ))}
    </Row>
    <DailyChart title="日次アクセス推移" data={accessDailyData} lines={accessChartLines} />

    <Divider style={{ margin: '24px 0' }} />

    {/* ③ 集客チャネル */}
    <SectionHeader title="集客チャネル" />
    {channelSections.map((ch) => (
      <ChannelBlock key={ch.key} channel={ch} />
    ))}

    <Divider style={{ margin: '8px 0 24px' }} />

    {/* ④ AIアドバイス */}
    <SectionHeader title="AIアドバイス" subtitle="データに基づく改善提案" />
    {aiAdvices.map((advice) => (
      <AIAdviceCard key={advice.id} advice={advice} />
    ))}
  </>
);

export default Dashboard;
