export interface KPIMetric {
  key: string;
  label: string;
  value: number;
  previousValue: number;
  format: 'number' | 'percent' | 'decimal' | 'duration';
}

export interface DailyDataPoint {
  date: string;
  [key: string]: string | number;
}

export interface ChartLine {
  dataKey: string;
  name: string;
  color: string;
}

export interface ChannelSection {
  key: string;
  name: string;
  connected: boolean;
  accountLabel?: string;
  kpis: KPIMetric[];
}

export interface PlatformConnection {
  key: string;
  name: string;
  connected: boolean;
  accountLabel?: string;
  description: string;
}

export interface AIAdvice {
  id: string;
  title: string;
  content: string;
  priority: 'high' | 'medium' | 'low';
}
