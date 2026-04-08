# SigNoz LDAP Kurulum Rehberi

## Gereksinimler

- Docker & Docker Compose **veya** Kubernetes (k3s/k8s)
- LDAP / Active Directory sunucusuna ağ erişimi
- ClickHouse (SigNoz stack'inin parçası)

## Hızlı Kurulum (Docker Compose)

### 1. Ortam Değişkenlerini Ayarla

```bash
cd deploy/docker
cp .env.ldap.example .env.ldap
```

`.env.ldap` dosyasını düzenleyin:

```env
# Admin hesabı
SIGNOZ_ADMIN_EMAIL=admin@setyazilim.com
SIGNOZ_ADMIN_PASSWORD=SetAdmin2024!@
SIGNOZ_ORG_NAME=SETYazilim
SIGNOZ_JWT_SECRET=guclu-bir-secret-buraya

# LDAP / Active Directory
LDAP_HOST=172.16.1.172
LDAP_PORT=389
LDAP_USE_TLS=false
LDAP_AD_DOMAINS="SETYAZILIM","SETSOFTWARE"
LDAP_EMAIL_DOMAIN=setyazilim.com
```

### 2. Başlat

```bash
docker compose -f docker-compose.yaml -f docker-compose.ldap.yaml --env-file .env.ldap up -d
```

### 3. Doğrulama

```bash
# API çalışıyor mu?
curl http://localhost:8080/api/v1/version

# LDAP yapılandırıldı mı?
curl "http://localhost:8080/api/v2/sessions/context?email=kullanici@setyazilim.com&ref=http://localhost:8080"
# Beklenen: "provider":"ldap"
```

### 4. Giriş

Tarayıcıdan `http://SUNUCU_IP:8080` adresine gidin:
- **Admin**: `admin@setyazilim.com` / `SetAdmin2024!@`
- **LDAP Kullanıcı**: `kullaniciadi@setyazilim.com` + AD şifresi

---

## Kubernetes Kurulumu

### 1. LDAP İmajını Oluştur

```bash
# Repo kök dizininde
docker build -t signoz-ldap:latest -f cmd/community/Dockerfile .
```

k3s kullanıyorsanız:
```bash
docker save signoz-ldap:latest | sudo k3s ctr images import -
```

### 2. Secret/ConfigMap Değerlerini Düzenle

`deploy/kubernetes/signoz-ldap.yaml` içindeki Secret ve ConfigMap'i düzenleyin:

```yaml
# Secret
stringData:
  jwt-secret: "guclu-bir-secret"
  admin-password: "SetAdmin2024!@"
  clickhouse-dsn: "tcp://signoz-clickhouse:9000/?username=admin&password=SIFRE"

# ConfigMap
data:
  admin-email: "admin@setyazilim.com"
  ldap-host: "172.16.1.172"
  ldap-port: "389"
  ldap-ad-domains: '"SETYAZILIM","SETSOFTWARE"'
  ldap-email-domain: "setyazilim.com"
```

### 3. Deploy

```bash
kubectl apply -f deploy/kubernetes/signoz-ldap.yaml -n signoz
```

### 4. LDAP Auth Domain Oluştur

Pod sağlıklı olduğunda:

```bash
kubectl apply -f deploy/kubernetes/signoz-ldap-init-job.yaml -n signoz
```

### 5. Erişim

```bash
# NodePort üzerinden
http://SUNUCU_IP:30690
```

---

## LDAP Yapılandırma Parametreleri

| Parametre | Açıklama | Varsayılan |
|-----------|----------|------------|
| `LDAP_HOST` | AD sunucu IP veya hostname | - |
| `LDAP_PORT` | LDAP portu | `389` |
| `LDAP_USE_TLS` | LDAPS (TLS) kullan | `false` |
| `LDAP_INSECURE_SKIP_VERIFY` | TLS sertifika doğrulamasını atla | `false` |
| `LDAP_AD_DOMAINS` | AD domain isimleri (bind için) | - |
| `LDAP_EMAIL_DOMAIN` | Email domain (tetikleyici) | - |
| `LDAP_BASE_DN` | Kullanıcı arama Base DN | (boş = sadece bind) |
| `LDAP_USER_SEARCH_FILTER` | Kullanıcı arama filtresi | `(sAMAccountName=%s)` |

## Nasıl Çalışır

1. Kullanıcı `kullanici@setyazilim.com` + şifre ile giriş yapar
2. Backend email domain'e bakarak (`setyazilim.com`) LDAP auth domain'i bulur
3. `SETYAZILIM\kullanici` ve `SETSOFTWARE\kullanici` şeklinde AD'ye bind dener
4. Bind başarılıysa, kullanıcı yerel DB'de yoksa otomatik oluşturulur (Viewer rolü)
5. JWT token döner, kullanıcı giriş yapar

## Sorun Giderme

```bash
# Backend loglarını kontrol et
docker logs signoz 2>&1 | grep ldap

# LDAP bağlantısını test et
docker exec signoz sh -c "timeout 5 cat < /dev/tcp/172.16.1.172/389 && echo OK || echo FAIL"

# Auth domain'leri listele (token gerekir)
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/api/v1/domains
```
