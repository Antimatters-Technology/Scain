# Deployment Guide

This guide covers deploying the Scain food traceability platform to various environments.

## 🚀 Quick Deployment

### Development Deployment

```bash
# Build both applications
npm run build

# Start production servers
npm start              # Frontend on port 3000
npm run start:backend  # Backend on port 8081
```

### Docker Deployment

```bash
# Build containers
docker build -t scain-frontend ./frontend
docker build -t scain-backend ./backend

# Run containers
docker run -d -p 3000:3000 --name scain-frontend scain-frontend
docker run -d -p 8081:8081 --name scain-backend scain-backend
```

## 🐳 Docker Configuration

### Frontend Dockerfile

```dockerfile
# frontend/Dockerfile
FROM node:18-alpine AS base

# Install dependencies only when needed
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Install dependencies based on the preferred package manager
COPY package.json package-lock.json* ./
RUN npm ci --only=production

# Rebuild the source code only when needed
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry during the build.
ENV NEXT_TELEMETRY_DISABLED 1

RUN npm run build

# Production image, copy all the files and run next
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public

# Set the correct permission for prerender cache
RUN mkdir .next
RUN chown nextjs:nodejs .next

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

CMD ["node", "server.js"]
```

### Backend Dockerfile

```dockerfile
# backend/Dockerfile
FROM node:18-alpine AS base

# Install dependencies only when needed
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Install dependencies based on the preferred package manager
COPY package.json package-lock.json* ./
RUN npm ci --only=production

# Rebuild the source code only when needed
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

RUN npm run build

# Production image, copy all the files and run the app
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nodejs

COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

USER nodejs

EXPOSE 8081

ENV PORT 8081
ENV HOST "0.0.0.0"

CMD ["node", "dist/index.js"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - NEXT_PUBLIC_API_URL=http://localhost:8081
    depends_on:
      - backend

  backend:
    build: ./backend
    ports:
      - "8081:8081"
    environment:
      - NODE_ENV=production
      - PORT=8081
      - HOST=0.0.0.0
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=scain
      - POSTGRES_USER=scain
      - POSTGRES_PASSWORD=scain_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

## ☁️ Cloud Deployment

### AWS Deployment

#### Using AWS ECS

```bash
# Create ECR repositories
aws ecr create-repository --repository-name scain-frontend
aws ecr create-repository --repository-name scain-backend

# Build and push images
docker build -t scain-frontend ./frontend
docker build -t scain-backend ./backend

aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

docker tag scain-frontend:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/scain-frontend:latest
docker tag scain-backend:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/scain-backend:latest

docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/scain-frontend:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/scain-backend:latest
```

#### Using AWS App Runner

```yaml
# apprunner.yaml
version: 1.0
runtime: nodejs18
build:
  commands:
    build:
      - npm install
      - npm run build
run:
  runtime-version: 18
  command: npm start
  network:
    port: 3000
    env: PORT
```

### Google Cloud Platform

#### Using Cloud Run

```bash
# Build and deploy to Cloud Run
gcloud builds submit --tag gcr.io/PROJECT_ID/scain-frontend ./frontend
gcloud builds submit --tag gcr.io/PROJECT_ID/scain-backend ./backend

gcloud run deploy scain-frontend --image gcr.io/PROJECT_ID/scain-frontend --platform managed --region us-central1 --allow-unauthenticated
gcloud run deploy scain-backend --image gcr.io/PROJECT_ID/scain-backend --platform managed --region us-central1 --allow-unauthenticated
```

### Azure

#### Using Azure Container Instances

```bash
# Build and push to Azure Container Registry
az acr build --registry myregistry --image scain-frontend ./frontend
az acr build --registry myregistry --image scain-backend ./backend

# Deploy to Container Instances
az container create \
  --resource-group myResourceGroup \
  --name scain-frontend \
  --image myregistry.azurecr.io/scain-frontend:latest \
  --dns-name-label scain-frontend \
  --ports 3000

az container create \
  --resource-group myResourceGroup \
  --name scain-backend \
  --image myregistry.azurecr.io/scain-backend:latest \
  --dns-name-label scain-backend \
  --ports 8081
```

## 🔧 Environment Configuration

### Production Environment Variables

```bash
# Frontend (.env.production)
NEXT_PUBLIC_API_URL=https://api.scain.com
NEXT_PUBLIC_APP_NAME=Scain
NEXT_PUBLIC_GA_ID=G-XXXXXXXXXX

# Backend (.env.production)
NODE_ENV=production
PORT=8081
HOST=0.0.0.0
DATABASE_URL=postgresql://user:pass@host:5432/scain
REDIS_URL=redis://host:6379
JWT_SECRET=your-super-secret-jwt-key
CORS_ORIGIN=https://scain.com
```

### Development Environment Variables

```bash
# Frontend (.env.development)
NEXT_PUBLIC_API_URL=http://localhost:8081
NEXT_PUBLIC_APP_NAME=Scain Dev

# Backend (.env.development)
NODE_ENV=development
PORT=8081
HOST=0.0.0.0
DATABASE_URL=postgresql://scain:scain@localhost:5432/scain_dev
REDIS_URL=redis://localhost:6379
LOG_LEVEL=debug
```

## 📊 Monitoring & Logging

### Application Monitoring

```typescript
// backend/src/utils/monitoring.ts
import { FastifyInstance } from 'fastify';

export function setupMonitoring(fastify: FastifyInstance) {
  // Health check endpoint
  fastify.get('/health', async (request, reply) => {
    return {
      status: 'ok',
      timestamp: new Date().toISOString(),
      version: process.env.npm_package_version || '1.0.0',
      uptime: process.uptime(),
      memory: process.memoryUsage()
    };
  });

  // Metrics endpoint
  fastify.get('/metrics', async (request, reply) => {
    return {
      requests: fastify.metrics.requests,
      errors: fastify.metrics.errors,
      responseTime: fastify.metrics.responseTime
    };
  });
}
```

### Logging Configuration

```typescript
// backend/src/utils/logger.ts
import pino from 'pino';

export const logger = pino({
  level: process.env.LOG_LEVEL || 'info',
  transport: {
    target: 'pino-pretty',
    options: {
      colorize: true,
      translateTime: 'SYS:standard',
      ignore: 'pid,hostname'
    }
  }
});
```

## 🔒 Security Configuration

### SSL/TLS Setup

```nginx
# nginx.conf
server {
    listen 443 ssl http2;
    server_name scain.com;

    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /api {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### CORS Configuration

```typescript
// backend/src/plugins/cors.ts
import { FastifyInstance } from 'fastify';

export async function corsPlugin(fastify: FastifyInstance) {
  await fastify.register(require('@fastify/cors'), {
    origin: process.env.CORS_ORIGIN || 'http://localhost:3000',
    credentials: true,
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS']
  });
}
```

## 🚀 CI/CD Pipeline

### GitHub Actions

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
      
      - run: npm ci
      - run: npm run lint
      - run: npm test
      - run: npm run build

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v1
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1
      
      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v1
      
      - name: Build and push Docker images
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
        run: |
          docker build -t $ECR_REGISTRY/scain-frontend:latest ./frontend
          docker build -t $ECR_REGISTRY/scain-backend:latest ./backend
          docker push $ECR_REGISTRY/scain-frontend:latest
          docker push $ECR_REGISTRY/scain-backend:latest
      
      - name: Deploy to ECS
        run: |
          aws ecs update-service --cluster scain-cluster --service scain-frontend --force-new-deployment
          aws ecs update-service --cluster scain-cluster --service scain-backend --force-new-deployment
```

## 📋 Deployment Checklist

### Pre-Deployment
- [ ] All tests passing
- [ ] Code review completed
- [ ] Security scan completed
- [ ] Performance testing completed
- [ ] Database migrations ready
- [ ] Environment variables configured
- [ ] SSL certificates installed
- [ ] Monitoring configured

### Deployment
- [ ] Backup current production data
- [ ] Deploy backend first
- [ ] Run database migrations
- [ ] Deploy frontend
- [ ] Verify health checks
- [ ] Test critical user flows
- [ ] Monitor error rates
- [ ] Update DNS if needed

### Post-Deployment
- [ ] Monitor application performance
- [ ] Check error logs
- [ ] Verify all integrations
- [ ] Update documentation
- [ ] Notify stakeholders
- [ ] Schedule rollback if needed

## 🔄 Rollback Procedures

### Quick Rollback

```bash
# Rollback to previous Docker image
docker tag scain-frontend:previous scain-frontend:latest
docker tag scain-backend:previous scain-backend:latest

# Restart services
docker-compose restart frontend backend
```

### Database Rollback

```bash
# Restore from backup
pg_restore -d scain_production backup_file.sql

# Or rollback migrations
npm run migrate:rollback
```

## 📊 Performance Optimization

### Frontend Optimization

```typescript
// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
  images: {
    domains: ['localhost'],
  },
  compress: true,
  poweredByHeader: false,
  generateEtags: false,
  onDemandEntries: {
    maxInactiveAge: 25 * 1000,
    pagesBufferLength: 2,
  },
}

module.exports = nextConfig
```

### Backend Optimization

```typescript
// backend/src/index.ts
import fastify from 'fastify';
import compression from '@fastify/compress';

const app = fastify({
  logger: true,
  trustProxy: true,
  bodyLimit: 1048576, // 1MB
});

await app.register(compression, { threshold: 1024 });
```

---

**Last Updated**: July 2025  
**Deployment Version**: 1.0.0 