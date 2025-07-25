import type { OnboardingContent } from '../../content/landingContent';

export default function OnboardingSection({ content }: { content: OnboardingContent }) {
  return (
    <section className="py-16 bg-gray-50">
      <div className="max-w-5xl mx-auto px-4 flex flex-col md:flex-row gap-8 items-center">
        {/* Styled illustration for Claim Wizard */}
        <div className="flex-1 flex justify-center">
          <div className="w-72 h-44 bg-gradient-to-tr from-blue-100 to-indigo-100 rounded-lg shadow-inner flex flex-col items-center justify-center text-blue-700 text-base font-semibold">
            <svg width="48" height="48" fill="none" viewBox="0 0 48 48" aria-hidden="true">
              <rect x="8" y="12" width="32" height="24" rx="4" fill="#3b82f6" fillOpacity="0.15" />
              <rect x="14" y="18" width="20" height="4" rx="2" fill="#3b82f6" fillOpacity="0.3" />
              <rect x="14" y="25" width="12" height="4" rx="2" fill="#3b82f6" fillOpacity="0.3" />
              <rect x="14" y="32" width="8" height="2" rx="1" fill="#3b82f6" fillOpacity="0.3" />
            </svg>
            <span className="mt-2">Claim Wizard Demo</span>
          </div>
        </div>
        <div className="flex-1">
          <h2 className="text-xl md:text-2xl font-bold mb-4">{content.title}</h2>
          <ul className="list-disc pl-5 space-y-2 text-gray-700">
            {content.bullets.map((b: string, i: number) => (
              <li key={i}>{b}</li>
            ))}
          </ul>
        </div>
      </div>
    </section>
  );
} 