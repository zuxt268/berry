import { Typography } from 'antd';

const { Title, Paragraph } = Typography;

const Home = () => (
  <div style={{ padding: '40px', maxWidth: 800, margin: '0 auto' }}>
    <Title>Hello, World!</Title>
    <Paragraph>プロジェクトの雛形が正常に動作しています。</Paragraph>
  </div>
);

export default Home;
