import { useEffect, useState } from 'react';
import { DashboardLayout } from '../components/DashboardLayout';
import { KPICard } from '../components/KPICard';
import { fetchMetric } from '../lib/api';

export default function Home() {
  const [revenue, setRevenue] = useState<string>('0');
  const [conversion, setConversion] = useState<string>('0%');
  const [arpu, setArpu] = useState<string>('0');
  const [mrr, setMrr] = useState<string>('0');

  useEffect(() => {
    const load = async () => {
      const revenueMetric = await fetchMetric('revenue');
      const conversionMetric = await fetchMetric('conversion-rate');
      const arpuMetric = await fetchMetric('arpu');
      const mrrMetric = await fetchMetric('mrr');

      setRevenue(revenueMetric.value ?? '0');
      setConversion(conversionMetric.value ?? '0%');
      setArpu(arpuMetric.value ?? '0');
      setMrr(mrrMetric.value ?? '0');
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
        <KPICard title="ARPU" value={arpu} subtitle="Last 30 days" />
        <KPICard title="MRR" value={mrr} subtitle="Current" />
      </div>
    </DashboardLayout>
  );
}
