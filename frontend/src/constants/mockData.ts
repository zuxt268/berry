import type { KPIMetric, ChannelSection, DailyDataPoint, ChartLine, PlatformConnection, AIAdvice } from '../types';

// --- Deterministic pseudo-random ---
function rand(seed: number): number {
  const x = Math.sin(seed * 9301 + 49297) * 233280;
  return x - Math.floor(x);
}

function generateDailyData(
  days: number,
  generators: Record<string, (day: number, r: number) => number>,
): DailyDataPoint[] {
  const baseDate = new Date('2026-03-10');
  return Array.from({ length: days }, (_, i) => {
    const date = new Date(baseDate);
    date.setDate(date.getDate() - (days - 1 - i));
    const label = `${date.getMonth() + 1}/${date.getDate()}`;
    const point: DailyDataPoint = { date: label };
    for (const [key, gen] of Object.entries(generators)) {
      point[key] = gen(i, rand(i * 17 + key.length * 31));
    }
    return point;
  });
}

// ============================================================
// ① 成果サマリー
// ============================================================
export const resultSummaryKPIs: KPIMetric[] = [
  { key: 'inquiries', label: 'お問い合わせ数', value: 42, previousValue: 35, format: 'number' },
  { key: 'phoneTaps', label: '電話タップ数', value: 128, previousValue: 104, format: 'number' },
];

// ============================================================
// ② アクセス概況
// ============================================================
export const accessOverviewKPIs: KPIMetric[] = [
  { key: 'sessions', label: 'アクセス数', value: 18542, previousValue: 16230, format: 'number' },
  { key: 'users', label: 'ユーザー数', value: 14821, previousValue: 13105, format: 'number' },
  { key: 'impressions', label: '検索表示回数', value: 85420, previousValue: 74300, format: 'number' },
  { key: 'clicks', label: 'クリック数', value: 5840, previousValue: 5120, format: 'number' },
  { key: 'position', label: '平均掲載順位', value: 14.2, previousValue: 15.8, format: 'decimal' },
];

export const accessDailyData: DailyDataPoint[] = generateDailyData(30, {
  access: (d, r) => Math.round(520 + d * 4 + r * 180),
  users: (d, r) => Math.round(400 + d * 3 + r * 140),
  clicks: (d, r) => Math.round(160 + d * 2 + r * 80),
});

export const accessChartLines: ChartLine[] = [
  { dataKey: 'access', name: 'アクセス数', color: '#1677ff' },
  { dataKey: 'users', name: 'ユーザー数', color: '#13c2c2' },
  { dataKey: 'clicks', name: 'クリック数', color: '#52c41a' },
];

// ============================================================
// ③ 集客チャネル
// ============================================================
export const channelSections: ChannelSection[] = [
  {
    key: 'gbp',
    name: 'Googleビジネスプロフィール',
    connected: true,
    kpis: [
      { key: 'gbpViews', label: '閲覧数', value: 4520, previousValue: 3890, format: 'number' },
      { key: 'gbpCalls', label: '電話タップ', value: 89, previousValue: 72, format: 'number' },
      { key: 'gbpRoutes', label: 'ルート検索', value: 156, previousValue: 134, format: 'number' },
    ],
  },
  {
    key: 'instagram',
    name: 'Instagram',
    connected: true,
    accountLabel: '@example_official',
    kpis: [
      { key: 'igFollowers', label: 'フォロワー数', value: 3240, previousValue: 2980, format: 'number' },
      { key: 'igEngagement', label: 'エンゲージメント率', value: 4.2, previousValue: 3.8, format: 'percent' },
    ],
  },
  {
    key: 'line',
    name: 'LINE公式アカウント',
    connected: false,
    kpis: [
      { key: 'lineFriends', label: '友だち数', value: 0, previousValue: 0, format: 'number' },
      { key: 'lineBlock', label: 'ブロック率', value: 0, previousValue: 0, format: 'percent' },
    ],
  },
];

// ============================================================
// ④ AIアドバイス
// ============================================================
export const aiAdvices: AIAdvice[] = [
  {
    id: '1',
    title: 'Googleビジネスプロフィールの写真を更新しましょう',
    content:
      '最新の写真が3ヶ月以上前です。新しい写真を追加すると、閲覧数が平均20%向上する傾向があります。店舗の外観・内装・商品の写真を5枚以上追加することをおすすめします。',
    priority: 'high',
  },
  {
    id: '2',
    title: 'Instagramの投稿頻度を上げましょう',
    content:
      '先月の投稿数は8件でした。週3回以上の投稿でエンゲージメント率が改善する傾向があります。リール動画の活用も効果的です。',
    priority: 'medium',
  },
  {
    id: '3',
    title: 'サイトの「料金ページ」のSEOを強化しましょう',
    content:
      '「料金」関連キーワードでの検索表示回数が増加していますが、クリック率が2.1%と低めです。タイトルタグとメタディスクリプションの見直しで改善が期待できます。',
    priority: 'high',
  },
  {
    id: '4',
    title: 'LINE公式アカウントを連携しましょう',
    content:
      'LINE公式アカウントを連携すると、友だち数やメッセージ配信の効果を一元管理できます。まずはリッチメニューの設定から始めましょう。',
    priority: 'low',
  },
];

// ============================================================
// 連携設定ページ用
// ============================================================
export const platformConnections: PlatformConnection[] = [
  {
    key: 'ga4',
    name: 'Google Analytics 4',
    connected: true,
    accountLabel: 'example.com',
    description: 'ウェブサイトのアクセス解析データを取得します。',
  },
  {
    key: 'gsc',
    name: 'Google Search Console',
    connected: true,
    accountLabel: 'example.com',
    description: '検索パフォーマンスデータを取得します。',
  },
  {
    key: 'gbp',
    name: 'Googleビジネスプロフィール',
    connected: true,
    description: 'ビジネスプロフィールの閲覧数やアクションを取得します。',
  },
  {
    key: 'instagram',
    name: 'Instagram',
    connected: true,
    accountLabel: '@example_official',
    description: 'Instagramビジネスアカウントのインサイトを取得します。',
  },
  {
    key: 'line',
    name: 'LINE公式アカウント',
    connected: false,
    description: 'LINE公式アカウントの友だち数やメッセージ配信データを取得します。',
  },
];
