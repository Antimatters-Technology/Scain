import * as React from 'react'
import { cn, formatRelativeTime } from '@/lib/utils'
import { CheckCircle, Truck, Package, Factory, Store } from 'lucide-react'

interface TimelineEvent {
  id: string
  type: 'harvested' | 'processed' | 'shipped' | 'received' | 'sold'
  location: string
  timestamp: string
  details?: string
}

interface TimelineLotProps {
  lotCode: string
  events?: TimelineEvent[]
  className?: string
}

const mockEvents: TimelineEvent[] = [
  {
    id: '1',
    type: 'harvested',
    location: 'Farm A, Maharashtra',
    timestamp: '2024-01-18T08:00:00Z',
    details: 'Organic tomatoes harvested',
  },
  {
    id: '2',
    type: 'processed',
    location: 'Processing Plant B, Pune',
    timestamp: '2024-01-18T14:30:00Z',
    details: 'Washed and packaged',
  },
  {
    id: '3',
    type: 'shipped',
    location: 'Distribution Center, Mumbai',
    timestamp: '2024-01-19T06:15:00Z',
    details: 'Shipped via cold chain',
  },
  {
    id: '4',
    type: 'received',
    location: 'Warehouse C, Delhi',
    timestamp: '2024-01-19T18:45:00Z',
    details: 'Temperature verified',
  },
]

const eventIcons = {
  harvested: Factory,
  processed: Package,
  shipped: Truck,
  received: Store,
  sold: CheckCircle,
}

const eventColors = {
  harvested: 'text-green-600 bg-green-50 ring-green-500/20',
  processed: 'text-blue-600 bg-blue-50 ring-blue-500/20',
  shipped: 'text-purple-600 bg-purple-50 ring-purple-500/20',
  received: 'text-orange-600 bg-orange-50 ring-orange-500/20',
  sold: 'text-primary bg-primary/10 ring-primary/20',
}

export function TimelineLot({ 
  lotCode, 
  events = mockEvents, 
  className 
}: TimelineLotProps) {
  return (
    <div className={cn('relative', className)}>
      {/* Lot Code Header */}
      <div className="mb-4 flex items-center gap-2">
        <h4 className="font-medium text-foreground">{lotCode}</h4>
        <span className="text-xs text-muted-foreground">
          {events.length} events
        </span>
      </div>

      {/* Timeline */}
      <div className="relative">
        {/* Connecting Line */}
        <div className="absolute left-6 top-0 h-full w-px bg-border" />
        
        {/* Events */}
        <div className="space-y-4">
          {events.map((event, index) => {
            const Icon = eventIcons[event.type]
            const isLast = index === events.length - 1
            
            return (
              <div key={event.id} className="relative flex items-start gap-4">
                {/* Icon */}
                <div className={cn(
                  'flex h-12 w-12 items-center justify-center rounded-full ring-4 ring-background',
                  eventColors[event.type]
                )}>
                  <Icon className="h-5 w-5" />
                </div>
                
                {/* Content */}
                <div className="flex-1 pb-4">
                  <div className="flex items-start justify-between">
                    <div>
                      <h5 className="font-medium capitalize text-foreground">
                        {event.type}
                      </h5>
                      <p className="text-sm text-muted-foreground">
                        {event.location}
                      </p>
                      {event.details && (
                        <p className="mt-1 text-xs text-muted-foreground">
                          {event.details}
                        </p>
                      )}
                    </div>
                    <time className="text-xs text-muted-foreground">
                      {formatRelativeTime(event.timestamp)}
                    </time>
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      </div>
      
      {/* Footer */}
      <div className="mt-4 flex items-center justify-between rounded-lg bg-muted/50 p-3">
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <div className="h-1.5 w-1.5 rounded-full bg-primary" />
          EPCIS 2.0 Compliant
        </div>
        <button className="text-xs font-medium text-primary hover:text-primary/80">
          View Full Trace →
        </button>
      </div>
    </div>
  )
} 