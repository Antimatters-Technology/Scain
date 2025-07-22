import { Card, CardContent, CardHeader } from '@/components/ui/card'

export function DashboardSkeleton() {
  return (
    <div className="grid gap-4 sm:grid-cols-2">
      {Array.from({ length: 6 }).map((_, i) => (
        <Card key={i} className="relative overflow-hidden">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="h-8 w-8 rounded-lg bg-muted animate-pulse" />
                <div className="space-y-2">
                  <div className="h-4 w-20 bg-muted rounded animate-pulse" />
                  <div className="h-6 w-16 bg-muted rounded animate-pulse" />
                </div>
              </div>
              <div className="space-y-1">
                <div className="h-2 w-2 bg-muted rounded-full animate-pulse ml-auto" />
                <div className="h-3 w-12 bg-muted rounded animate-pulse" />
              </div>
            </div>
          </CardHeader>
        </Card>
      ))}
    </div>
  )
} 