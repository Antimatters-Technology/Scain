import type { ComplianceContent } from '../../content/landingContent';

export default function ComplianceSection({ content }: { content: ComplianceContent }) {
  return (
    <section className="py-16 bg-white">
      <div className="max-w-4xl mx-auto px-4 flex flex-col md:flex-row gap-8 items-center">
        {/* CTE coverage bar */}
        <div className="flex-1 flex flex-col items-center">
          <div className="w-full max-w-xs h-4 bg-gray-200 rounded-full overflow-hidden mb-2">
            <div
              className="h-full bg-green-500 transition-all"
              style={{ width: `${content.cteCoveragePercent}%` }}
              aria-label="CTE coverage bar"
            />
          </div>
          <span className="text-xs text-gray-600 mb-2">CTE Coverage: {content.cteCoveragePercent}%</span>
          <a
            href="/sample-epcis.json"
            download
            className="px-4 py-2 bg-blue-600 text-white rounded shadow font-semibold text-sm hover:bg-blue-700 transition"
            aria-label="Export EPCIS JSON + SHA-256 manifest"
          >
            {content.exportButtonLabel}
          </a>
        </div>
        <div className="flex-1">
          <h2 className="text-xl md:text-2xl font-bold mb-2">{content.title}</h2>
          <p className="text-gray-700 text-sm">{content.blurb}</p>
        </div>
      </div>
    </section>
  );
} 