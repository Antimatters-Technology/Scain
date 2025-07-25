import Hero from '../../components/landing/Hero';
import KPIStrip from '../../components/landing/KPIStrip';
import HowItWorks from '../../components/landing/HowItWorks';
import OnboardingSection from '../../components/landing/OnboardingSection';
import ComplianceSection from '../../components/landing/ComplianceSection';
import CompareTable from '../../components/landing/CompareTable';
import PublicProof from '../../components/landing/PublicProof';
import CTA from '../../components/landing/CTA';
import Footer from '../../components/landing/Footer';
import content from '../../content/landingContent';

export default function LandingPage() {
  return (
    <main className="bg-white text-gray-900 min-h-screen flex flex-col">
      <Hero content={content.hero} />
      <KPIStrip content={content.kpi} />
      <HowItWorks content={content.howItWorks} />
      <OnboardingSection content={content.onboarding} />
      <ComplianceSection content={content.compliance} />
      <CompareTable content={content.compare} />
      <PublicProof content={content.publicProof} />
      <CTA content={content.cta} />
      <Footer content={content.footer} />
    </main>
  );
} 