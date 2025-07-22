import { redirect } from 'next/navigation'
import { Loader2 } from 'lucide-react'

export default function RootPage() {
  // Redirect to dashboard
  redirect('/dashboard')

  // This component will never render due to redirect
  // But we include it for type safety and potential loading state
  return (
    <div className="flex h-screen items-center justify-center">
      <div className="flex flex-col items-center gap-4">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <p className="text-muted-foreground">Redirecting to dashboard...</p>
      </div>
    </div>
  )
} 