import * as React from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { cn } from '@/lib/utils'

type SensorStatus = 'good' | 'warning' | 'critical' | 'offline'

interface CardSensorProps {
  title: string
  value: number
  unit: string
  status: SensorStatus
  icon?: React.ReactNode
  className?: string
}

const statusStyles: Record<SensorStatus, string> = {
  good: 'ring-green-500/20 bg-green-50 dark:bg-green-950/20',
  warning: 'ring-amber-500/20 bg-amber-50 dark:bg-amber-950/20',
  critical: 'ring-red-500/20 bg-red-50 dark:bg-red-950/20',
  offline: 'ring-gray-500/20 bg-gray-50 dark:bg-gray-950/20',
}

const statusIndicators: Record<SensorStatus, string> = {
  good: 'bg-green-500',
  warning: 'bg-amber-500 animate-pulse',
  critical: 'bg-red-500 animate-pulse',
  offline: 'bg-gray-400',
}

const valueColors: Record<SensorStatus, string> = {
  good: 'text-green-600 dark:text-green-400',
  warning: 'text-amber-600 dark:text-amber-400',
  critical: 'text-red-600 dark:text-red-400',
  offline: 'text-gray-500 dark:text-gray-400',
}

export function CardSensor({ 
  title, 
  value, 
  unit, 
  status, 
  icon, 
  className 
}: CardSensorProps) {
  return (
    <Card className={cn(
      'relative overflow-hidden ring-1 transition-all hover:shadow-md',
      statusStyles[status],
      className
    )}>
      <CardContent className="p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            {icon && (
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-background/80">
                {icon}
              </div>
            )}
            <div>
              <p className="text-sm font-medium text-muted-foreground">
                {title}
              </p>
              <div className="flex items-baseline gap-1">
                <span className={cn(
                  'text-2xl font-bold font-display',
                  valueColors[status]
                )}>
                  {value}
                </span>
                <span className="text-sm text-muted-foreground font-medium">
                  {unit}
                </span>
              </div>
            </div>
          </div>
          
          {/* Status Indicator */}
          <div className="flex flex-col items-end gap-1">
            <div className={cn(
              'h-2 w-2 rounded-full',
              statusIndicators[status]
            )} />
            <span className="text-xs capitalize text-muted-foreground">
              {status}
            </span>
          </div>
        </div>
        
        {/* Background Pattern */}
        <div className="absolute -bottom-6 -right-6 h-16 w-16 rotate-12 opacity-5 text-muted-foreground/20">
          {icon}
        </div>
      </CardContent>
    </Card>
  )
} 