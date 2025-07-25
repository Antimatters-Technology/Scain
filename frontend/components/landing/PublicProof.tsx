'use client';
import type { PublicProofContent } from '../../content/landingContent';
import { useState } from 'react';

export default function PublicProof({ content }: { content: PublicProofContent }) {
  const [open, setOpen] = useState(false);
  return (
    <section className="py-16 bg-white">
      <div className="max-w-3xl mx-auto px-4 text-center">
        <div className="inline-flex items-center gap-2 mb-4">
          <span className="inline-block px-3 py-1 bg-green-100 text-green-700 rounded-full font-semibold text-xs">
            {content.badgeLabel}
          </span>
          <button
            type="button"
            onClick={() => setOpen(true)}
            className="text-blue-600 underline text-xs font-medium hover:text-blue-800"
          >
            {content.explorerLabel}
          </button>
        </div>
        <h2 className="text-xl md:text-2xl font-bold mb-2">{content.title}</h2>
        <p className="text-gray-700 text-sm mb-4">{content.blurb}</p>
        {open && (
          <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg shadow-lg p-6 max-w-md w-full text-left">
              <h3 className="font-bold text-lg mb-2">Sample Transaction</h3>
              <div className="mb-2 text-xs text-gray-600">Hash: <span className="font-mono">0x1234abcd5678ef90</span></div>
              <div className="mb-4 text-green-700 font-semibold">Status: Confirmed</div>
              <button
                className="px-4 py-2 bg-blue-600 text-white rounded font-semibold text-sm hover:bg-blue-700 transition"
                onClick={() => setOpen(false)}
              >
                Close
              </button>
            </div>
          </div>
        )}
      </div>
    </section>
  );
} 