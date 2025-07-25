import type { HowItWorksContent } from '../../content/landingContent';
import { LucideIcon } from 'lucide-react';

export default function HowItWorks({ content }: { content: HowItWorksContent }) {
  return (
    <section className="py-16 bg-white">
      <div className="max-w-5xl mx-auto px-4">
        <h2 className="text-2xl md:text-3xl font-bold text-center mb-8">{content.title}</h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {content.steps.map((step: { icon: LucideIcon; title: string; description: string }, i: number) => {
            const Icon = step.icon;
            return (
              <div key={i} className="flex flex-col items-center text-center">
                <div className="mb-4">
                  <Icon className="w-10 h-10 text-blue-600" aria-hidden="true" />
                </div>
                <h3 className="font-semibold text-lg mb-2">{step.title}</h3>
                <p className="text-gray-600 text-sm">{step.description}</p>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
} 