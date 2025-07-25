import { Database, Link, ShieldCheck, BadgeCheck } from 'lucide-react';

export type HeroContent = {
  headline: string;
  subheadline: string;
  primaryCta: { label: string; href: string };
};

export type KPIStripContent = {
  kpis: { label: string; value: string }[];
};

export type HowItWorksContent = {
  title: string;
  steps: { icon: any; title: string; description: string }[];
};

export type OnboardingContent = {
  title: string;
  bullets: string[];
};

export type ComplianceContent = {
  title: string;
  cteCoveragePercent: number;
  exportButtonLabel: string;
  blurb: string;
};

export type CompareContent = {
  title: string;
  columns: string[];
  rows: string[][];
};

export type PublicProofContent = {
  title: string;
  badgeLabel: string;
  explorerLabel: string;
  explorerHref: string;
  blurb: string;
};

export type CTAContent = {
  title: string;
  ctaLabel: string;
  socialProof: { src: string; alt: string }[];
};

export type FooterContent = {
  links: { label: string; href: string; external?: boolean }[];
};

const content = {
  hero: {
    headline: 'Any signal in. Instant, standards-grade traceability out.',
    subheadline: 'Normalize sensors, trackers & ERP events to GS1 EPCIS 2.0, anchor integrity on Hyperledger Fabric, and (optionally) prove it publicly.',
    primaryCta: { label: 'Book a 15‑min traceability audit', href: '/contact' },
  },
  kpi: {
    kpis: [
      { label: 'Lots monitored', value: '1,200+' }, // TODO: confirm
      { label: 'Devices online', value: '340' }, // TODO: confirm
      { label: 'Avg latency', value: '2.2s' }, // TODO: confirm
      { label: 'CTE coverage', value: '98%' }, // TODO: confirm
      { label: 'Anchors published', value: '1,800' }, // TODO: confirm
    ],
  },
  howItWorks: {
    title: 'How Scain Works',
    steps: [
      { icon: Database, title: 'Connect any data source', description: 'DIY ESP32, LoRaWAN, Tive, ERP, and more.' },
      { icon: Link, title: 'Normalize to EPCIS 2.0', description: 'Validate KDE/CTE, ensure standards compliance.' },
      { icon: ShieldCheck, title: 'Hash to Fabric', description: 'Private by default, public anchor optional.' },
      { icon: BadgeCheck, title: 'Dashboards & Recalls', description: 'Trace, recall, and audit in seconds.' },
    ],
  },
  onboarding: {
    title: 'Self-Service IoT Onboarding',
    bullets: [
      'Arduino/MicroPython code snippets',
      'ExpressLink AT command script',
      'LoRaWAN OTAA join form',
      'Generic tracker webhook template',
    ],
  },
  compliance: {
    title: 'Compliance & Audit Export',
    cteCoveragePercent: 98, // TODO: confirm
    exportButtonLabel: 'Export EPCIS JSON + SHA-256 manifest',
    blurb: 'FSMA §204 and SFCR ready. Export full EPCIS event history and manifest for audit or regulatory review.',
  },
  compare: {
    title: 'Compare vs Status Quo',
    columns: ['','Manual spreadsheets','QR-only','Token blockchain','Scain'],
    rows: [
      ['CTE coverage','Low','Medium','Medium','High'],
      ['Sensor fidelity','None','Low','Medium','High'],
      ['Cost predictability','Manual','QR fees','Token risk','Deterministic'],
      ['Public transparency','None','QR scan','Token ledger','Optional'],
      ['DIY onboarding','No','No','No','Yes'],
    ],
  },
  publicProof: {
    title: 'Optional Public Proof',
    badgeLabel: 'Public Proof',
    explorerLabel: 'View on Explorer',
    explorerHref: '/explorer',
    blurb: 'Hybrid model: permissioned privacy by default, with optional public hash anchoring for extra trust. Show the badge when anchored.',
  },
  cta: {
    title: 'Ready to see traceability in action?',
    ctaLabel: 'Request Demo',
    socialProof: [
      { src: '/icons/tive.svg', alt: 'Tive' },
      { src: '/icons/controlant.svg', alt: 'Controlant' },
      { src: '/icons/aws-iot.svg', alt: 'AWS IoT Core LoRaWAN' },
    ],
  },
  footer: {
    links: [
      { label: 'Docs', href: '/docs' },
      { label: 'EPCIS Schema', href: 'https://ref.gs1.org/epcis/', external: true },
      { label: 'Security', href: '/security' },
      { label: 'Careers', href: '/careers' },
    ],
  },
};

export default content; 