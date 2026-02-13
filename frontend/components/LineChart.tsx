type TrendPoint = { date: string; value: number };

interface LineChartProps {
  title: string;
  subtitle?: string;
  data: TrendPoint[];
  valueSuffix?: string;
}

function toPoints(data: TrendPoint[], width: number, height: number) {
  if (data.length === 0) {
    return "";
  }

  const values = data.map((d) => d.value);
  const min = Math.min(...values);
  const max = Math.max(...values);
  const range = max - min || 1;
  const stepX = width / Math.max(data.length - 1, 1);

  return data
    .map((d, i) => {
      const x = i * stepX;
      const y = height - ((d.value - min) / range) * height;
      return `${x.toFixed(1)},${y.toFixed(1)}`;
    })
    .join(" ");
}

function latestValue(data: TrendPoint[]) {
  if (data.length === 0) {
    return "0";
  }
  const value = data[data.length - 1].value;
  return value.toFixed(2);
}

export function LineChart({ title, subtitle, data, valueSuffix }: LineChartProps) {
  const width = 520;
  const height = 160;
  const points = toPoints(data, width, height);

  return (
    <div className="chart-card">
      <div className="chart-header">
        <div>
          <div className="chart-title">{title}</div>
          {subtitle && <div className="chart-subtitle">{subtitle}</div>}
        </div>
        <div className="chart-metric">
          {latestValue(data)}{valueSuffix ?? ""}
        </div>
      </div>
      <div className="chart-body">
        {data.length === 0 ? (
          <div className="chart-empty">No data</div>
        ) : (
          <svg className="chart-svg" viewBox={`0 0 ${width} ${height}`} preserveAspectRatio="none">
            <polyline className="chart-line" points={points} />
          </svg>
        )}
      </div>
      {data.length > 0 && (
        <div className="chart-axis">
          <span>{data[0].date}</span>
          <span>{data[data.length - 1].date}</span>
        </div>
      )}
    </div>
  );
}
