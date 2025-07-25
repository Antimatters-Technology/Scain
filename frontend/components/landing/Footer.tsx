import type { FooterContent } from '../../content/landingContent';

export default function Footer({ content }: { content: FooterContent }) {
  return (
    <footer className="w-full bg-gray-900 text-white py-8 mt-auto">
      <div className="max-w-5xl mx-auto px-4 flex flex-col md:flex-row justify-between items-center gap-4">
        <div className="flex gap-6 mb-2 md:mb-0">
          {content.links.map((link: { label: string; href: string; external?: boolean }, i: number) => (
            <a
              key={i}
              href={link.href}
              className="text-sm hover:underline focus:outline-none focus:ring-2 focus:ring-blue-500"
              aria-label={link.label}
              target={link.external ? '_blank' : undefined}
              rel={link.external ? 'noopener noreferrer' : undefined}
            >
              {link.label}
            </a>
          ))}
        </div>
        <div className="text-xs text-gray-400">&copy; {new Date().getFullYear()} Scain. All rights reserved.</div>
      </div>
    </footer>
  );
} 