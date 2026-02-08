import { ReactNode } from 'react';

interface DashboardLayoutProps {
  children: ReactNode;
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <div className="layout">
      <header className="header">
        <div className="title">Revenue & Conversion Intelligence</div>
        <div className="subtitle">Executive Dashboard MVP</div>
      </header>
      <main className="content">{children}</main>
    </div>
  );
}
