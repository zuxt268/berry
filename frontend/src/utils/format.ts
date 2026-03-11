import type { KPIMetric } from '../types';

export function formatValue(value: number, format: KPIMetric['format']): string {
  switch (format) {
    case 'number':
      return value.toLocaleString();
    case 'percent':
      return `${value.toFixed(1)}%`;
    case 'decimal':
      return value.toFixed(1);
    case 'duration': {
      const minutes = Math.floor(value / 60);
      const seconds = Math.round(value % 60);
      return `${minutes}:${seconds.toString().padStart(2, '0')}`;
    }
  }
}
