'use client';
import { Button } from '../ui/button';
import type { HeroContent } from '../../content/landingContent';

export default function Hero({ content }: { content: HeroContent }) {
  const handleScroll = (e: React.MouseEvent) => {
    e.preventDefault();
    const el = document.getElementById('demo-form');
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  };
  return (
    <section className="relative flex flex-col items-center justify-center min-h-[60vh] py-16 text-center bg-gradient-to-b from-gray-50 to-white">
      {/* Animated background visual placeholder */}
      <div aria-hidden="true" className="absolute inset-0 z-0 flex items-center justify-center">
        {/* TODO: Replace with animated SVG or canvas chain/sensor flow */}
        <div className="w-64 h-64 bg-gradient-to-tr from-blue-200 to-indigo-200 rounded-full opacity-30 blur-2xl animate-pulse" />
      </div>
      <div className="relative z-10 max-w-2xl mx-auto">
        {/* Brand name */}
        <div className="mb-6">
          <span className="inline-block text-2xl md:text-3xl font-extrabold tracking-tight text-blue-700">Scain</span>
        </div>
        <h1 className="text-4xl md:text-5xl font-bold mb-4 leading-tight">
          {content.headline}
        </h1>
        <p className="text-lg md:text-xl text-gray-700 mb-8">
          {content.subheadline}
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Button asChild>
            <a href="#demo-form" onClick={handleScroll} className="font-semibold" aria-label={content.primaryCta.label}>
              {content.primaryCta.label}
            </a>
          </Button>
          <Button asChild variant="outline">
            <a href="/login" className="font-semibold" aria-label="Login to dashboard">
              Login
            </a>
          </Button>
        </div>
      </div>
    </section>
  );
} 