import { useEffect, useState } from 'react';
import { DashboardLayout } from '../components/DashboardLayout';
import { KPICard } from '../components/KPICard';
import { fetchMetric } from '../lib/api';

export default function Home() {
  const [revenue, setRevenue] = useState<string>('0');
  const [conversion, setConversion] = useState<string>('0%');

  useEffect(() => {
    const load = async () => {
      const revenueMetric = await fetchMetric('revenue');
      const conversionMetric = await fetchMetric('conversion-rate');

      setRevenue(revenueMetric.value ?? '0');
      setConversion(conversionMetric.value ?? '0%');
    };

    load().catch(() => {
      setRevenue('0');
      setConversion('0%');
    });
  }, []);

  return (
    <DashboardLayout>
      <div className="grid">
        <KPICard title="Revenue" value={revenue} subtitle="Last 30 days" />
        <KPICard title="Conversion Rate" value={conversion} subtitle="Last 30 days" />
        <KPICard title="ARPU" value="0" subtitle="Last 30 days" />
        <KPICard title="MRR" value="0" subtitle="Current" />
      </div>
    </DashboardLayout>
  );
}
