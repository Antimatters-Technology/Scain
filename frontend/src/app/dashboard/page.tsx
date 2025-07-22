'use client';

import { Suspense } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { CardSensor } from '@/components/CardSensor';
import { TimelineLot } from '@/components/TimelineLot';
import { GraphViewport } from '@/components/GraphViewport';
import { RecallDrillModal } from '@/components/RecallDrillModal';
import { DashboardSkeleton } from '@/components/dashboard-skeleton';
import { CircleCheckBig, Clock, TriangleAlert, TrendingUp, MapPin } from 'lucide-react';

// Mock data for demonstration
const mockSensorData = [
  { id: 'LOT-2024-001', airTemp: 4.2, probeTemp: 3.8, humidity: 85, shockG: 0.1, timestamp: '2024-01-19T10:30:00Z', lat: 40.7128, lng: -74.006, status: 'good' },
  { id: 'LOT-2024-002', airTemp: 8.5, probeTemp: 7.2, humidity: 78, shockG: 0.3, timestamp: '2024-01-19T10:25:00Z', lat: 40.7589, lng: -73.9851, status: 'warning' },
  { id: 'LOT-2024-003', airTemp: 2.1, probeTemp: 1.9, humidity: 92, shockG: 0.05, timestamp: '2024-01-19T10:20:00Z', lat: 40.6892, lng: -74.0445, status: 'good' },
];

export default function DashboardPage() {
  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="sticky top-0 z-40 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4">
          <div className="flex h-16 items-center justify-between">
            <div className="flex items-center gap-4">
              <h1 className="text-2xl font-bold text-foreground">
                Scain <span className="text-primary">Dashboard</span>
              </h1>
            </div>
            <div className="flex items-center gap-4">
              <div className="hidden md:flex items-center gap-2 text-sm text-muted-foreground">
                <div className="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
                Live • FSMA §204 Compliant
              </div>
              <RecallDrillModal />
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main id="main-content" className="container mx-auto px-4 py-8">
        {/* KPI Metrics */}
        <div className="mb-8 grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Lots</CardTitle>
              <CircleCheckBig className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold font-display">3</div>
              <p className="text-xs text-muted-foreground">+2 from yesterday</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Active (24h)</CardTitle>
              <Clock className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold font-display text-primary">0</div>
              <p className="text-xs text-muted-foreground">Currently monitoring</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Warnings</CardTitle>
              <TriangleAlert className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold font-display text-amber-500">1</div>
              <p className="text-xs text-muted-foreground">Require attention</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Avg Temp</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold font-display">4.3°C</div>
              <p className="text-xs text-muted-foreground">Within range</p>
            </CardContent>
          </Card>
        </div>

        {/* Main Grid Layout */}
        <div className="grid gap-8 lg:grid-cols-3">
          {/* Left Column - Sensors */}
          <div className="lg:col-span-2">
            <h2 className="text-lg font-semibold mb-4">Live Sensors</h2>
            <Suspense fallback={<DashboardSkeleton />}>
              <div className="grid gap-4 sm:grid-cols-2">
                {/* Group sensors by lot */}
                <div className="space-y-4">
                  <h3 className="text-sm font-medium text-muted-foreground">LOT-2024-001</h3>
                  <div className="grid gap-2">
                    <CardSensor
                      title="Air Temp"
                      value={4.2}
                      unit="°C"
                      status="good"
                      icon="thermometer"
                    />
                    <CardSensor
                      title="Probe Temp"
                      value={3.8}
                      unit="°C"
                      status="good"
                      icon="thermometer"
                    />
                    <CardSensor
                      title="Humidity"
                      value={85}
                      unit="%"
                      status="good"
                      icon="droplets"
                    />
                    <CardSensor
                      title="Shock"
                      value={0.1}
                      unit="G"
                      status="good"
                      icon="zap"
                    />
                  </div>
                </div>

                <div className="space-y-4">
                  <h3 className="text-sm font-medium text-muted-foreground">LOT-2024-002</h3>
                  <div className="grid gap-2">
                    <CardSensor
                      title="Air Temp"
                      value={8.5}
                      unit="°C"
                      status="warning"
                      icon="thermometer"
                    />
                    <CardSensor
                      title="Probe Temp"
                      value={7.2}
                      unit="°C"
                      status="warning"
                      icon="thermometer"
                    />
                    <CardSensor
                      title="Humidity"
                      value={78}
                      unit="%"
                      status="good"
                      icon="droplets"
                    />
                    <CardSensor
                      title="Shock"
                      value={0.3}
                      unit="G"
                      status="warning"
                      icon="zap"
                    />
                  </div>
                </div>

                <div className="space-y-4">
                  <h3 className="text-sm font-medium text-muted-foreground">LOT-2024-003</h3>
                  <div className="grid gap-2">
                    <CardSensor
                      title="Air Temp"
                      value={2.1}
                      unit="°C"
                      status="good"
                      icon="thermometer"
                    />
                    <CardSensor
                      title="Probe Temp"
                      value={1.9}
                      unit="°C"
                      status="good"
                      icon="thermometer"
                    />
                    <CardSensor
                      title="Humidity"
                      value={92}
                      unit="%"
                      status="warning"
                      icon="droplets"
                    />
                    <CardSensor
                      title="Shock"
                      value={0.05}
                      unit="G"
                      status="good"
                      icon="zap"
                    />
                  </div>
                </div>
              </div>
            </Suspense>
          </div>

          {/* Right Column - Timeline & Graph */}
          <div className="space-y-8">
            {/* Recent Events Timeline */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MapPin className="h-4 w-4" />
                  Recent Events
                </CardTitle>
                <CardDescription>Latest EPCIS events from the supply chain</CardDescription>
              </CardHeader>
              <CardContent>
                <TimelineLot lotCode="LOT-2024-001" />
              </CardContent>
            </Card>

            {/* Temperature Trend Chart */}
            <Card>
              <CardHeader>
                <CardTitle>Temperature Trend</CardTitle>
                <CardDescription>24-hour temperature monitoring</CardDescription>
              </CardHeader>
              <CardContent>
                <GraphViewport data={mockSensorData} />
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Compliance Badge */}
        <div className="mt-12 rounded-lg border border-primary/20 bg-primary/5 p-6">
          <div className="flex items-start gap-4">
            <CircleCheckBig className="mt-1 h-5 w-5 text-primary" />
            <div>
              <h3 className="font-medium text-foreground">FSMA §204 Compliance Active</h3>
              <p className="mt-1 text-sm text-muted-foreground">
                All sensor data is being recorded with 2-second trace capability for FDA compliance. 
                Data retention: 7 years as required by the Food Traceability Rule.
              </p>
              <div className="mt-3 flex flex-wrap gap-2">
                <span className="inline-flex items-center rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary">
                  FSMA §204
                </span>
                <span className="inline-flex items-center rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary">
                  SFCR Part 5
                </span>
                <span className="inline-flex items-center rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary">
                  Hyperledger Fabric
                </span>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
} 