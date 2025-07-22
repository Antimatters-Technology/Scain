# Development Guide

This guide covers setting up and developing the Scain food traceability platform with a **primary focus on backend implementation**.

## 🚀 Quick Start

### Prerequisites

- **Node.js**: 20.0.0 or higher
- **npm**: 9.0.0 or higher
- **Git**: Latest version
- **PostgreSQL**: 15+ (for database)
- **Redis**: 7+ (for caching)

### Installation

```bash
# Clone the repository
git clone https://github.com/your-org/scain.git
cd scain

# Install dependencies for all packages
npm install

# Start development servers
npm run dev:full
```

### Verify Installation

- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8081/health

---

## 🛠️ Development Environment

### Project Structure

```
scain/
├── frontend/          # Next.js application ✅
│   ├── src/          # Source code
│   ├── components/   # React components
│   ├── public/       # Static assets
│   └── package.json  # Frontend dependencies
├── backend/          # Fastify API ⚠️ MINIMAL
│   ├── src/          # TypeScript source (5 files only)
│   │   ├── index.ts  # Basic server (63 lines)
│   │   ├── epcis/    # Schema definitions only
│   │   └── utils/    # Basic utilities
│   ├── dist/         # Compiled output
│   └── package.json  # Backend dependencies
├── docs/             # Documentation
├── package.json      # Workspace configuration
└── README.md         # Project overview
```

### Development Scripts

```bash
# Development
npm run dev:full      # Start both frontend and backend
npm run dev           # Frontend only
npm run dev:backend   # Backend only

# Building
npm run build         # Build both applications
npm run build:frontend # Frontend only
npm run build:backend # Backend only

# Production
npm start             # Start frontend production server
npm run start:backend # Start backend production server

# Testing
npm test              # Run all tests (BACKEND HAS NO TESTS)
npm run test:backend  # Backend tests only (EMPTY)

# Utilities
npm run lint          # Lint both applications
npm run clean         # Clean build artifacts
```

---

## ⚡ Backend Development

### Technology Stack

**Current Backend Stack:**
- **Language**: TypeScript (with Node.js runtime)
- **Framework**: Fastify (high-performance web framework)
- **Database**: PostgreSQL (planned)
- **Cache**: Redis (planned)
- **Validation**: Zod schemas
- **Logging**: Pino

### Why TypeScript + Fastify?

#### **TypeScript Benefits:**
- ✅ **Type Safety**: Catches errors at compile time
- ✅ **Better IDE Support**: IntelliSense, refactoring, debugging
- ✅ **Team Collaboration**: Self-documenting code
- ✅ **Enterprise Ready**: Large codebase management
- ✅ **EPCIS Compliance**: Complex schema validation needs strong typing

#### **Fastify Benefits:**
- ✅ **Performance**: One of the fastest Node.js frameworks
- ✅ **Low Overhead**: Minimal memory footprint
- ✅ **Plugin Ecosystem**: Rich middleware support
- ✅ **TypeScript Native**: Built with TypeScript in mind
- ✅ **JSON Schema**: Built-in request validation

### Current Backend State

**✅ Implemented (5 files, 230 lines):**
```typescript
// backend/src/index.ts - Basic server
import Fastify from 'fastify';

const fastify = Fastify({
  logger: false
});

fastify.get('/health', async () => {
  return { 
    status: 'ok',
    timestamp: new Date().toISOString(),
    version: '1.0.0'
  };
});
```

**❌ Missing (95% of backend):**
```
backend/src/                    # ❌ MISSING DIRECTORIES
├── api/                       # ❌ API routes
├── database/                  # ❌ Database layer
├── services/                  # ❌ Business logic
├── middleware/                # ❌ Authentication & validation
├── jobs/                      # ❌ Background processing
├── websockets/                # ❌ Real-time features
└── blockchain/                # ❌ Blockchain integration
```

---

## 🗄️ Database Setup

### PostgreSQL Installation

```bash
# macOS
brew install postgresql
brew services start postgresql

# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql

# Windows
# Download from https://www.postgresql.org/download/windows/
```

### Database Configuration

```bash
# Create database
createdb scain_dev

# Create user (optional)
createuser scain_user

# Set password
psql -d scain_dev -c "ALTER USER scain_user PASSWORD 'scain_password';"
```

### Environment Variables

```bash
# backend/.env
DATABASE_URL=postgresql://scain_user:scain_password@localhost:5432/scain_dev
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-super-secret-jwt-key
NODE_ENV=development
PORT=8081
HOST=0.0.0.0
LOG_LEVEL=debug
```

---

## 🔧 Backend Implementation Guide

### Phase 1: Core Infrastructure

#### 1. Database Models

```typescript
// backend/src/database/models/Device.ts
export interface Device {
  id: string;
  deviceId: string;
  type: 'ESP32' | 'ExpressLink' | 'LoRaWAN' | 'Tracker';
  location: string;
  status: 'active' | 'inactive' | 'error';
  metadata: Record<string, unknown>;
  createdAt: Date;
  updatedAt: Date;
}

// backend/src/database/models/Event.ts
export interface Event {
  id: string;
  eventType: 'ObjectEvent' | 'TransformationEvent' | 'AggregationEvent';
  eventTime: Date;
  deviceId: string;
  epcList: string[];
  sensorData: Record<string, unknown>;
  hash: string;
  blockchainTx?: string;
  createdAt: Date;
}
```

#### 2. Database Connection

```typescript
// backend/src/database/connection.ts
import { Pool } from 'pg';
import logger from '../utils/logger.js';

const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
  max: 20,
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 2000,
});

pool.on('error', (err) => {
  logger.error('Unexpected error on idle client', err);
  process.exit(-1);
});

export default pool;
```

#### 3. API Routes Structure

```typescript
// backend/src/api/routes/devices.ts
import { FastifyInstance } from 'fastify';
import { DeviceService } from '../../services/deviceService.js';

export default async function deviceRoutes(fastify: FastifyInstance) {
  const deviceService = new DeviceService();

  // GET /api/devices
  fastify.get('/api/devices', async (request, reply) => {
    try {
      const devices = await deviceService.getDevices();
      return { status: 'success', data: { devices } };
    } catch (error) {
      logger.error('Error fetching devices', error);
      return reply.status(500).send({ 
        status: 'error', 
        message: 'Internal server error' 
      });
    }
  });

  // POST /api/devices
  fastify.post('/api/devices', async (request, reply) => {
    try {
      const device = await deviceService.createDevice(request.body);
      return { status: 'success', data: { device } };
    } catch (error) {
      logger.error('Error creating device', error);
      return reply.status(400).send({ 
        status: 'error', 
        message: 'Invalid device data' 
      });
    }
  });
}
```

#### 4. Service Layer

```typescript
// backend/src/services/deviceService.ts
import pool from '../database/connection.js';
import { Device } from '../database/models/Device.js';
import logger from '../utils/logger.js';

export class DeviceService {
  async getDevices(): Promise<Device[]> {
    const result = await pool.query(
      'SELECT * FROM devices ORDER BY created_at DESC'
    );
    return result.rows;
  }

  async createDevice(deviceData: Partial<Device>): Promise<Device> {
    const result = await pool.query(
      `INSERT INTO devices (device_id, type, location, status, metadata) 
       VALUES ($1, $2, $3, $4, $5) 
       RETURNING *`,
      [
        deviceData.deviceId,
        deviceData.type,
        deviceData.location,
        deviceData.status || 'active',
        JSON.stringify(deviceData.metadata || {})
      ]
    );
    return result.rows[0];
  }

  async getDeviceById(id: string): Promise<Device | null> {
    const result = await pool.query(
      'SELECT * FROM devices WHERE id = $1',
      [id]
    );
    return result.rows[0] || null;
  }
}
```

---

## 🧪 Testing

### Current Testing Status
- **Tests**: ❌ No tests implemented
- **Coverage**: 0%
- **Test Framework**: Jest (configured but unused)

### Test Structure

```typescript
// backend/src/api/routes/__tests__/devices.test.ts
import { test } from 'tap';
import { build } from '../../../src/app.js';

test('devices API', async (t) => {
  const app = build();

  // Test GET /api/devices
  const response = await app.inject({
    method: 'GET',
    url: '/api/devices'
  });

  t.equal(response.statusCode, 200);
  t.same(JSON.parse(response.payload), {
    status: 'success',
    data: { devices: [] }
  });
});
```

### Running Tests

```bash
# Run all tests
npm run test:backend

# Run tests in watch mode
npm run test:watch --prefix backend

# Run tests with coverage
npm run test:coverage --prefix backend
```

---

## 🔐 Authentication

### JWT Implementation

```typescript
// backend/src/middleware/auth.ts
import { FastifyInstance } from 'fastify';
import jwt from 'jsonwebtoken';

export async function authPlugin(fastify: FastifyInstance) {
  await fastify.register(require('@fastify/jwt'), {
    secret: process.env.JWT_SECRET || 'default-secret'
  });

  fastify.addHook('onRequest', async (request, reply) => {
    try {
      if (request.url.startsWith('/api/')) {
        await request.jwtVerify();
      }
    } catch (err) {
      reply.send(err);
    }
  });
}
```

### User Management

```typescript
// backend/src/services/authService.ts
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';
import pool from '../database/connection.js';

export class AuthService {
  async registerUser(email: string, password: string) {
    const hashedPassword = await bcrypt.hash(password, 10);
    
    const result = await pool.query(
      'INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email',
      [email, hashedPassword]
    );
    
    return result.rows[0];
  }

  async loginUser(email: string, password: string) {
    const result = await pool.query(
      'SELECT * FROM users WHERE email = $1',
      [email]
    );
    
    const user = result.rows[0];
    if (!user) {
      throw new Error('User not found');
    }

    const isValid = await bcrypt.compare(password, user.password_hash);
    if (!isValid) {
      throw new Error('Invalid password');
    }

    const token = jwt.sign(
      { userId: user.id, email: user.email },
      process.env.JWT_SECRET!,
      { expiresIn: '24h' }
    );

    return { user: { id: user.id, email: user.email }, token };
  }
}
```

---

## 📊 EPCIS Processing

### Event Validation

```typescript
// backend/src/epcis/validator.ts
import { EpcisEventSchema, EpcisEvent } from './schemas/index.js';
import { computeHash } from '../utils/hash.js';

export class EpcisValidator {
  validateEvent(eventData: unknown): EpcisEvent {
    return EpcisEventSchema.parse(eventData);
  }

  computeEventHash(event: EpcisEvent): string {
    return computeHash(event);
  }

  async processEvent(eventData: unknown) {
    // Validate event
    const event = this.validateEvent(eventData);
    
    // Compute hash
    const hash = this.computeEventHash(event);
    
    // Store in database
    // Anchor to blockchain
    
    return { event, hash };
  }
}
```

---

## 🔄 Background Jobs

### Cron Job Setup

```typescript
// backend/src/jobs/anchorJob.ts
import cron from 'node-cron';
import { AnchorService } from '../services/anchorService.js';
import logger from '../utils/logger.js';

export class AnchorJob {
  private anchorService: AnchorService;

  constructor() {
    this.anchorService = new AnchorService();
  }

  start() {
    // Run every 5 minutes
    cron.schedule('*/5 * * * *', async () => {
      try {
        logger.info('Starting anchor job');
        await this.anchorService.processPendingAnchors();
        logger.info('Anchor job completed');
      } catch (error) {
        logger.error('Anchor job failed', error);
      }
    });
  }
}
```

---

## 🚀 Deployment

### Development Deployment

```bash
# Build backend
npm run build:backend

# Start production server
npm run start:backend
```

### Docker Deployment

```bash
# Build container
docker build -t scain-backend ./backend

# Run container
docker run -d -p 8081:8081 --name scain-backend scain-backend
```

---

## 🔧 Configuration

### TypeScript Configuration

```json
// backend/tsconfig.json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "moduleResolution": "node",
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "**/*.test.ts"]
}
```

### ESLint Configuration

```javascript
// backend/eslint.config.js
export default [
  {
    files: ['src/**/*.ts'],
    languageOptions: {
      parser: '@typescript-eslint/parser',
      parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module',
      },
    },
    plugins: {
      '@typescript-eslint': '@typescript-eslint/eslint-plugin',
    },
    rules: {
      '@typescript-eslint/no-unused-vars': 'error',
      '@typescript-eslint/explicit-function-return-type': 'warn',
      '@typescript-eslint/no-explicit-any': 'warn',
    },
  },
];
```

---

## 🤝 Contributing

### Backend Development Workflow

1. **Setup Database**
   ```bash
   # Install PostgreSQL
   brew install postgresql
   
   # Create database
   createdb scain_dev
   
   # Run migrations (when implemented)
   npm run migrate
   ```

2. **Implement API Endpoints**
   ```bash
   # Create new route file
   touch backend/src/api/routes/devices.ts
   
   # Add to server
   # Test endpoint
   curl http://localhost:8081/api/devices
   ```

3. **Add Tests**
   ```bash
   # Create test file
   touch backend/src/api/routes/__tests__/devices.test.ts
   
   # Run tests
   npm run test:backend
   ```

### Code Style Guidelines

```typescript
// Use interfaces for object shapes
interface Device {
  id: string;
  deviceId: string;
  type: DeviceType;
  location: string;
  createdAt: Date;
}

// Use enums for constants
enum DeviceType {
  ESP32 = 'ESP32',
  EXPRESS_LINK = 'ExpressLink',
  LORA_WAN = 'LoRaWAN',
  TRACKER = 'Tracker'
}

// Use type aliases for unions
type ApiResponse<T> = {
  data: T;
  status: 'success' | 'error';
  message?: string;
};
```

---

## 📚 Resources

### Documentation
- [Fastify Documentation](https://www.fastify.io/docs/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Jest Testing Framework](https://jestjs.io/docs/)

### Tools
- [VS Code](https://code.visualstudio.com/) - Recommended IDE
- [Postman](https://www.postman.com/) - API testing
- [pgAdmin](https://www.pgadmin.org/) - PostgreSQL administration
- [Redis Commander](https://github.com/joeferner/redis-commander) - Redis management

---

**Last Updated**: July 2025  
**Backend Completion**: 5%  
**Priority**: Core API implementation  
**Status**: Development - Infrastructure needed 