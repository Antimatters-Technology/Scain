import type { KPIStripContent } from '../../content/landingContent';

export default function KPIStrip({ content }: { content: KPIStripContent }) {
  return (
    <section className="w-full bg-gray-900 text-white py-6">
      <div className="max-w-5xl mx-auto flex flex-wrap justify-center gap-6">
        {content.kpis.map((kpi: { label: string; value: string }) => (
          <div
            key={kpi.label}
            className="flex flex-col items-center min-w-[100px] transition-transform hover:scale-105 hover:bg-gray-800 rounded-lg px-2 py-1"
          >
            <span className="text-2xl md:text-3xl font-bold" aria-label={kpi.label}>
              {kpi.value}
            </span>
            <span className="text-xs md:text-sm text-gray-300 mt-1">{kpi.label}</span>
          </div>
        ))}
      </div>
    </section>
  );
} 