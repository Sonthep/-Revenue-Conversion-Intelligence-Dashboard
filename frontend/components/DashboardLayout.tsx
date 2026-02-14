import { ReactNode } from 'react';

interface DashboardLayoutProps {
  children: ReactNode;
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  const now = new Date();
  const formatted = now.toISOString().slice(0, 10);

  return (
    <div className="layout">
      <header className="header">
        <div className="header-inner">
          <div>
            <div className="title">Revenue & Conversion Intelligence</div>
            <div className="subtitle">Executive Dashboard MVP</div>
          </div>
          <div className="header-meta">
            <span className="badge badge-accent">Live Preview</span>
            <span className="badge badge-muted">As of {formatted}</span>
          </div>
        </div>
      </header>
      <main className="content">{children}</main>
    </div>
  );
}
