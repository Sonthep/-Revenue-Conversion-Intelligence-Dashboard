import { useEffect, useState } from 'react';
import { DashboardLayout } from '../components/DashboardLayout';
import { KPICard } from '../components/KPICard';
import { LineChart } from '../components/LineChart';
import { fetchMetric } from '../lib/api';

type TrendPoint = { date: string; value: number };

export default function Home() {
  const [revenue, setRevenue] = useState<string>('0');
  const [conversion, setConversion] = useState<string>('0%');
  const [arpu, setArpu] = useState<string>('0');
  const [mrr, setMrr] = useState<string>('0');
  const [nrr, setNrr] = useState<string>('0%');
  const [churn, setChurn] = useState<string>('0%');
  const [ltv, setLtv] = useState<string>('0');
  const [cac, setCac] = useState<string>('0');
  const [revenueTrend, setRevenueTrend] = useState<TrendPoint[]>([]);
  const [conversionTrend, setConversionTrend] = useState<TrendPoint[]>([]);

  useEffect(() => {
    const load = async () => {
      const revenueMetric = await fetchMetric('revenue');
      const conversionMetric = await fetchMetric('conversion-rate');
      const arpuMetric = await fetchMetric('arpu');
      const mrrMetric = await fetchMetric('mrr');
      const nrrMetric = await fetchMetric('nrr');
      const churnMetric = await fetchMetric('churn-rate');
      const ltvMetric = await fetchMetric('ltv');
      const cacMetric = await fetchMetric('cac');
      const revenueTrendMetric = await fetchMetric('revenue-trend');
      const conversionTrendMetric = await fetchMetric('conversion-trend');

      setRevenue(revenueMetric.value ?? '0');
      setConversion(conversionMetric.value ?? '0%');
      setArpu(arpuMetric.value ?? '0');
      setMrr(mrrMetric.value ?? '0');
      setNrr(nrrMetric.value ?? '0%');
      setChurn(churnMetric.value ?? '0%');
      setLtv(ltvMetric.value ?? '0');
      setCac(cacMetric.value ?? '0');
      setRevenueTrend(Array.isArray(revenueTrendMetric.value) ? revenueTrendMetric.value : []);
      setConversionTrend(Array.isArray(conversionTrendMetric.value) ? conversionTrendMetric.value : []);
    };

    load().catch(() => {
      setRevenue('0');
      setConversion('0%');
      setArpu('0');
      setMrr('0');
      setNrr('0%');
      setChurn('0%');
      setLtv('0');
      setCac('0');
      setRevenueTrend([]);
      setConversionTrend([]);
    });
  }, []);

  const revenueRows = revenueTrend.slice(-14);
  const conversionRows = conversionTrend.slice(-14);

  return (
    <DashboardLayout>
      <div className="section-header">
        <div className="section-title">Executive Snapshot</div>
        <div className="section-subtitle">Core KPIs for the last 30 days</div>
      </div>
      <div className="grid">
        <KPICard title="Revenue" value={revenue} subtitle="Last 30 days" />
        <KPICard title="Conversion Rate" value={conversion} subtitle="Last 30 days" />
        <KPICard title="ARPU" value={arpu} subtitle="Last 30 days" />
        <KPICard title="MRR" value={mrr} subtitle="Current" />
        <KPICard title="NRR" value={nrr} subtitle="Last 30 days" />
        <KPICard title="Churn Rate" value={churn} subtitle="Last 30 days" />
        <KPICard title="LTV" value={ltv} subtitle="Derived" />
        <KPICard title="CAC" value={cac} subtitle="Last 30 days" />
      </div>

      <div className="section-header">
        <div className="section-title">Momentum</div>
        <div className="section-subtitle">Short-term trend lines</div>
      </div>
      <div className="trend-section">
        <LineChart
          title="Revenue Trend"
          subtitle="Last 14 days"
          data={revenueRows}
        />
        <LineChart
          title="Conversion Trend"
          subtitle="Last 14 days"
          data={conversionRows}
          valueSuffix="%"
        />
      </div>
    </DashboardLayout>
  );
}
