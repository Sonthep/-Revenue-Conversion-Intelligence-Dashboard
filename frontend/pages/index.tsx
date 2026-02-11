import { useEffect, useState } from 'react';
import { DashboardLayout } from '../components/DashboardLayout';
import { KPICard } from '../components/KPICard';
import { fetchMetric } from '../lib/api';

export default function Home() {
  const [revenue, setRevenue] = useState<string>('0');
  const [conversion, setConversion] = useState<string>('0%');
  const [arpu, setArpu] = useState<string>('0');
  const [mrr, setMrr] = useState<string>('0');
  const [nrr, setNrr] = useState<string>('0%');
  const [churn, setChurn] = useState<string>('0%');
  const [ltv, setLtv] = useState<string>('0');

  useEffect(() => {
    const load = async () => {
      const revenueMetric = await fetchMetric('revenue');
      const conversionMetric = await fetchMetric('conversion-rate');
      const arpuMetric = await fetchMetric('arpu');
      const mrrMetric = await fetchMetric('mrr');
      const nrrMetric = await fetchMetric('nrr');
      const churnMetric = await fetchMetric('churn-rate');
      const ltvMetric = await fetchMetric('ltv');

      setRevenue(revenueMetric.value ?? '0');
      setConversion(conversionMetric.value ?? '0%');
      setArpu(arpuMetric.value ?? '0');
      setMrr(mrrMetric.value ?? '0');
      setNrr(nrrMetric.value ?? '0%');
      setChurn(churnMetric.value ?? '0%');
      setLtv(ltvMetric.value ?? '0');
    };

    load().catch(() => {
      setRevenue('0');
      setConversion('0%');
      setArpu('0');
      setMrr('0');
      setNrr('0%');
      setChurn('0%');
      setLtv('0');
    });
  }, []);

  return (
    <DashboardLayout>
      <div className="grid">
        <KPICard title="Revenue" value={revenue} subtitle="Last 30 days" />
        <KPICard title="Conversion Rate" value={conversion} subtitle="Last 30 days" />
        <KPICard title="ARPU" value={arpu} subtitle="Last 30 days" />
        <KPICard title="MRR" value={mrr} subtitle="Current" />
        <KPICard title="NRR" value={nrr} subtitle="Last 30 days" />
        <KPICard title="Churn Rate" value={churn} subtitle="Last 30 days" />
        <KPICard title="LTV" value={ltv} subtitle="Derived" />
      </div>
    </DashboardLayout>
  );
}
