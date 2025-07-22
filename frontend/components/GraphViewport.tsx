'use client'

import * as React from 'react'
import { Line } from 'react-chartjs-2'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js'
import { cn } from '@/lib/utils'

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

interface SensorData {
  id: string
  airTemp: number
  probeTemp: number
  humidity: number
  shockG: number
  timestamp: string
  lat: number
  lng: number
  status: string
}

interface GraphViewportProps {
  data: SensorData[]
  className?: string
}

export function GraphViewport({ data, className }: GraphViewportProps) {
  // Generate mock historical data for the chart
  const generateTimePoints = () => {
    const points = []
    const now = new Date()
    for (let i = 23; i >= 0; i--) {
      const time = new Date(now.getTime() - i * 60 * 60 * 1000)
      points.push(time.toLocaleTimeString('en-US', { 
        hour: '2-digit', 
        minute: '2-digit',
        hour12: false 
      }))
    }
    return points
  }

  const generateMockData = (baseTemp: number, variance: number = 1) => {
    return Array.from({ length: 24 }, () => 
      baseTemp + (Math.random() - 0.5) * variance * 2
    )
  }

  const timeLabels = generateTimePoints()
  const avgProbeTemp = data.reduce((acc, sensor) => acc + sensor.probeTemp, 0) / data.length
  
  const chartData = {
    labels: timeLabels,
    datasets: [
      {
        label: 'Probe Temperature',
        data: generateMockData(avgProbeTemp, 0.5),
        borderColor: 'rgb(16, 185, 129)', // primary color
        backgroundColor: 'rgba(16, 185, 129, 0.1)',
        borderWidth: 2,
        fill: true,
        tension: 0.4,
        pointBackgroundColor: 'rgb(16, 185, 129)',
        pointBorderColor: '#ffffff',
        pointBorderWidth: 2,
        pointRadius: 4,
        pointHoverRadius: 6,
      },
      {
        label: 'Air Temperature',
        data: generateMockData(avgProbeTemp + 0.5, 0.8),
        borderColor: 'rgb(245, 158, 11)', // accent color
        backgroundColor: 'rgba(245, 158, 11, 0.1)',
        borderWidth: 2,
        fill: false,
        tension: 0.4,
        pointBackgroundColor: 'rgb(245, 158, 11)',
        pointBorderColor: '#ffffff',
        pointBorderWidth: 2,
        pointRadius: 3,
        pointHoverRadius: 5,
      }
    ],
  }

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top' as const,
        labels: {
          boxWidth: 12,
          boxHeight: 12,
          usePointStyle: true,
          font: {
            size: 12,
            family: 'Inter',
          },
          color: 'rgb(107, 114, 128)', // muted-foreground
        },
      },
      tooltip: {
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
        titleColor: 'rgb(17, 24, 39)',
        bodyColor: 'rgb(75, 85, 99)',
        borderColor: 'rgb(229, 231, 235)',
        borderWidth: 1,
        cornerRadius: 8,
        displayColors: true,
        intersect: false,
        mode: 'index' as const,
        callbacks: {
          label: function(context: any) {
            return `${context.dataset.label}: ${context.parsed.y.toFixed(1)}°C`
          }
        }
      },
    },
    scales: {
      x: {
        grid: {
          color: 'rgba(229, 231, 235, 0.3)',
          drawBorder: false,
        },
        ticks: {
          color: 'rgb(107, 114, 128)',
          font: {
            size: 11,
            family: 'Inter',
          },
          maxTicksLimit: 8,
        },
      },
      y: {
        beginAtZero: false,
        grid: {
          color: 'rgba(229, 231, 235, 0.3)',
          drawBorder: false,
        },
        ticks: {
          color: 'rgb(107, 114, 128)',
          font: {
            size: 11,
            family: 'Inter',
          },
          callback: function(value: any) {
            return `${value}°C`
          }
        },
      },
    },
    interaction: {
      intersect: false,
      mode: 'index' as const,
    },
    elements: {
      point: {
        hoverBackgroundColor: 'rgb(16, 185, 129)',
      }
    }
  }

  return (
    <div className={cn('relative', className)}>
      {/* Chart Container */}
      <div className="h-64 w-full">
        <Line data={chartData} options={options} />
      </div>
      
      {/* Map Fallback */}
      {!process.env.NEXT_PUBLIC_MAPBOX_TOKEN && (
        <div className="mt-4 rounded-lg border-2 border-dashed border-muted-foreground/25 p-4">
          <div className="flex items-center justify-center space-y-2 text-center">
            <div>
              <p className="text-sm text-muted-foreground">
                Map view unavailable
              </p>
              <p className="text-xs text-muted-foreground">
                Add MAPBOX_TOKEN to enable location tracking
              </p>
            </div>
          </div>
        </div>
      )}
      
      {/* Data Summary */}
      <div className="mt-4 grid grid-cols-3 gap-4 text-center">
        <div>
          <p className="text-lg font-bold font-display text-foreground">
            {data.length}
          </p>
          <p className="text-xs text-muted-foreground">Active Sensors</p>
        </div>
        <div>
          <p className="text-lg font-bold font-display text-primary">
            {avgProbeTemp.toFixed(1)}°C
          </p>
          <p className="text-xs text-muted-foreground">Avg Temperature</p>
        </div>
        <div>
          <p className="text-lg font-bold font-display text-green-600">
            100%
          </p>
          <p className="text-xs text-muted-foreground">Data Integrity</p>
        </div>
      </div>
    </div>
  )
} 