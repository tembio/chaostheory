# Kubernetes Manifests for Leaderboard System

## Files
- `rabbitmq-deployment.yaml`: RabbitMQ Deployment and Service
- `leaderboard-deployment.yaml`: Leaderboard app Deployment and Service
- `leaderboard-pvc.yaml`: PersistentVolumeClaim for leaderboard SQLite DB
- `mockeventgenerator-deployment.yaml`: Mock event generator Deployment (no Service because it does not expose an HTTP API)

## Usage
1. **Build and push leaderboard and mockeventgenerator Docker images** 
2. **Apply the manifests:**
   ```sh
   kubectl apply -f k8s/leaderboard-pvc.yaml
   kubectl apply -f k8s/rabbitmq-deployment.yaml
   kubectl apply -f k8s/leaderboard-deployment.yaml
   kubectl apply -f k8s/mockeventgenerator-deployment.yaml
   ```
3. **Access services:**
   - Leaderboard API: NodePort on your cluster (see `kubectl get svc leaderboard`)
   - RabbitMQ Management: NodePort on port 15672 (see `kubectl get svc rabbitmq`)
   - Mockeventgenerator: No Service 

## Notes
- The leaderboard DB is stored in a persistent volume (`leaderboard-pvc`).
