export interface MetricResponse {
  metric: string;
  value: string;
  updated_at: string;
  cached: boolean;
  time_window: string;
}

export async function fetchMetric(metric: string): Promise<MetricResponse> {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';
  const response = await fetch(`${baseUrl}/api/metrics/${metric}`);
  if (!response.ok) {
    throw new Error('Failed to fetch metric');
  }
  return response.json();
}
