'use client'

import * as React from 'react'
import { cn } from '@/lib/utils'
import { AlertTriangle, Download, Search, X, File } from 'lucide-react'

interface RecallDrillModalProps {
  className?: string
}

export function RecallDrillModal({ className }: RecallDrillModalProps) {
  const [isOpen, setIsOpen] = React.useState(false)
  const [selectedLot, setSelectedLot] = React.useState('')
  const [isSearching, setIsSearching] = React.useState(false)

  const handleRecallDrill = async () => {
    if (!selectedLot) return
    
    setIsSearching(true)
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 2000))
    setIsSearching(false)
    
    // In real implementation, this would open the recall trace results
    alert(`Recall drill initiated for ${selectedLot}. Full trace completed in 2 seconds.`)
  }

  const handleExportPDF = () => {
    // In real implementation, this would generate and download PDF
    alert('PDF export feature would download the recall trace report here.')
  }

  if (!isOpen) {
    return (
      <button
        onClick={() => setIsOpen(true)}
        className={cn(
          'flex items-center gap-2 rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2',
          className
        )}
      >
        <AlertTriangle className="h-4 w-4" />
        Recall Drill
      </button>
    )
  }

  return (
    <>
      {/* Backdrop */}
      <div 
        className="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm"
        onClick={() => setIsOpen(false)}
      />
      
      {/* Modal */}
      <div className="fixed left-1/2 top-1/2 z-50 w-full max-w-lg -translate-x-1/2 -translate-y-1/2 transform rounded-xl bg-background p-6 shadow-2xl ring-1 ring-border">
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
              <AlertTriangle className="h-5 w-5 text-red-600 dark:text-red-400" />
            </div>
            <div>
              <h2 className="text-lg font-semibold text-foreground">
                Recall Drill
              </h2>
              <p className="text-sm text-muted-foreground">
                FSMA §204 2-second trace requirement
              </p>
            </div>
          </div>
          <button
            onClick={() => setIsOpen(false)}
            className="rounded-lg p-1 text-muted-foreground hover:bg-muted hover:text-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        {/* Content */}
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-2">
              Lot Code
            </label>
            <input
              type="text"
              value={selectedLot}
              onChange={(e) => setSelectedLot(e.target.value)}
              placeholder="Enter lot code (e.g., LOT-2024-001)"
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2"
            />
          </div>

          <div className="rounded-lg bg-amber-50 dark:bg-amber-950/20 p-4">
            <div className="flex items-start gap-3">
              <AlertTriangle className="h-5 w-5 text-amber-600 dark:text-amber-400 mt-0.5" />
              <div className="text-sm">
                <p className="font-medium text-amber-800 dark:text-amber-200">
                  Compliance Notice
                </p>
                <p className="text-amber-700 dark:text-amber-300 mt-1">
                  This recall drill must complete within 2 seconds as required by FSMA §204. 
                  All affected products will be identified with full supply chain traceability.
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex items-center justify-between mt-6 pt-4 border-t border-border">
          <div className="flex items-center gap-2">
            <button
              onClick={handleExportPDF}
              disabled={!selectedLot}
              className="flex items-center gap-2 rounded-lg border border-border px-3 py-2 text-sm font-medium hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <File className="h-4 w-4" />
              Export PDF
            </button>
          </div>
          
          <div className="flex items-center gap-3">
            <button
              onClick={() => setIsOpen(false)}
              className="rounded-lg px-4 py-2 text-sm font-medium text-muted-foreground hover:text-foreground"
            >
              Cancel
            </button>
            <button
              onClick={handleRecallDrill}
              disabled={!selectedLot || isSearching}
              className="flex items-center gap-2 rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSearching ? (
                <>
                  <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                  Searching...
                </>
              ) : (
                <>
                  <Search className="h-4 w-4" />
                  Start Recall
                </>
              )}
            </button>
          </div>
        </div>

        {/* Footer */}
        <div className="mt-4 pt-4 border-t border-border">
          <p className="text-xs text-muted-foreground">
            This recall drill simulates emergency food safety response procedures. 
            Real recalls must be reported to FDA within 24 hours.
          </p>
        </div>
      </div>
    </>
  )
} 