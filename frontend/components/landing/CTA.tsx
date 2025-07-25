'use client';
import { Button } from '../ui/button';
import type { CTAContent } from '../../content/landingContent';
import { useState } from 'react';

const logos = [
  {
    name: 'Tive',
    svg: (
      <svg width="80" height="32" viewBox="0 0 80 32" fill="none" xmlns="http://www.w3.org/2000/svg" aria-label="Tive logo">
        <defs>
          <linearGradient id="tive-bg" x1="0" y1="0" x2="80" y2="32" gradientUnits="userSpaceOnUse">
            <stop stopColor="#2563eb" />
            <stop offset="1" stopColor="#1A56DB" />
          </linearGradient>
        </defs>
        <rect width="80" height="32" rx="8" fill="url(#tive-bg)"/>
        <text x="40" y="22" textAnchor="middle" fill="#fff" fontSize="18" fontFamily="Arial" fontWeight="bold">Tive</text>
      </svg>
    ),
  },
  {
    name: 'Controlant',
    svg: (
      <svg width="80" height="32" viewBox="0 0 80 32" fill="none" xmlns="http://www.w3.org/2000/svg" aria-label="Controlant logo">
        <defs>
          <linearGradient id="controlant-bg" x1="0" y1="0" x2="80" y2="32" gradientUnits="userSpaceOnUse">
            <stop stopColor="#22d3ee" />
            <stop offset="1" stopColor="#0E7490" />
          </linearGradient>
        </defs>
        <rect width="80" height="32" rx="8" fill="url(#controlant-bg)"/>
        <text x="40" y="22" textAnchor="middle" fill="#fff" fontSize="16" fontFamily="Arial" fontWeight="bold">Controlant</text>
      </svg>
    ),
  },
  {
    name: 'AWS IoT Core LoRaWAN',
    svg: (
      <svg width="80" height="32" viewBox="0 0 80 32" fill="none" xmlns="http://www.w3.org/2000/svg" aria-label="AWS IoT Core LoRaWAN logo">
        <defs>
          <linearGradient id="aws-bg" x1="0" y1="0" x2="80" y2="32" gradientUnits="userSpaceOnUse">
            <stop stopColor="#fbbf24" />
            <stop offset="1" stopColor="#F59E42" />
          </linearGradient>
        </defs>
        <rect width="80" height="32" rx="8" fill="url(#aws-bg)"/>
        <text x="40" y="22" textAnchor="middle" fill="#fff" fontSize="15" fontFamily="Arial" fontWeight="bold">AWS IoT</text>
      </svg>
    ),
  },
];

export default function CTA({ content }: { content: CTAContent }) {
  const [submitted, setSubmitted] = useState(false);
  return (
    <section id="demo-form" className="py-16 bg-gray-50">
      <div className="max-w-2xl mx-auto px-4 text-center">
        <h2 className="text-2xl md:text-3xl font-bold mb-4">{content.title}</h2>
        {submitted ? (
          <div className="py-8 text-green-700 font-semibold text-lg">Thank you! We’ll be in touch soon.</div>
        ) : (
          <form
            className="flex flex-col gap-4 items-center"
            aria-label="Contact form"
            action="https://formspree.io/f/xdoqzqzq"
            method="POST"
            target="_blank"
            onSubmit={() => setSubmitted(true)}
          >
            <input
              type="email"
              name="email"
              placeholder="Email"
              required
              className="w-full max-w-xs px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              aria-label="Email"
            />
            <textarea
              name="pain"
              placeholder="What’s your biggest traceability pain?"
              required
              className="w-full max-w-xs px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              aria-label="Pain point"
            />
            <Button type="submit" className="w-full max-w-xs font-semibold">
              {content.ctaLabel}
            </Button>
          </form>
        )}
        <div className="mt-10">
          <div className="mb-2 text-xs font-semibold text-gray-500 tracking-wide uppercase">Trusted by</div>
          <div className="flex flex-wrap justify-center items-center gap-6">
            {logos.map((logo, i) => (
              <div
                key={i}
                className="flex flex-col items-center transition-transform hover:scale-105 drop-shadow-lg bg-white rounded-lg p-2"
                style={{ minWidth: 90 }}
              >
                {logo.svg}
                <span className="text-[11px] mt-1 text-gray-500 font-normal tracking-wide">{logo.name}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
} 