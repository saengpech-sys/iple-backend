# Dockerfile for IPLE GoFiber Backend

# --- Build Stage ---
# ใช้ Go image ที่มี Alpine เป็นฐาน เพื่อให้ได้ build tools และขนาดไม่ใหญ่มาก
FROM golang:1.24.2 AS builder

# กำหนด Working Directory ใน container
WORKDIR /app

# Copy go.mod และ go.sum ก่อน เพื่อใช้ประโยชน์จาก Docker layer caching
# ถ้าไฟล์เหล่านี้ไม่เปลี่ยน ไม่ต้อง download dependencies ใหม่
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy โค้ดทั้งหมดเข้ามาใน container
COPY . .

# Build ตัวแอปพลิเคชัน
# -ldflags="-w -s" ลดขนาด binary โดยเอา debug symbols ออก
# CGO_ENABLED=0 สร้าง static binary (ไม่มี C dependency) สำคัญสำหรับ Alpine/Scratch
# -o /server กำหนด output path ของ binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /server ./cmd/server/main.go

# --- Final Stage ---
# ใช้ Base image ที่เล็กมากๆ เช่น Alpine หรือ Scratch (ถ้า static จริงๆ)
# Distroless ก็เป็นทางเลือกที่ดี: gcr.io/distroless/static-debian11
FROM alpine:latest
# FROM scratch # ใช้ได้ถ้า binary เป็น static และไม่มี dependency อื่นๆ นอกจาก libc (ถ้าใช้ CGO)

WORKDIR /app

# สร้าง User/Group ที่ไม่มีสิทธิ์ Root สำหรับรันแอป
# RUN addgroup -S appgroup && adduser -S appuser -G appgroup
# USER appuser

# Copy เฉพาะ binary ที่ build แล้วจาก stage ก่อนหน้า
COPY --from=builder /server /app/server

# (Optional) Copy ไฟล์ .env.example หรือไฟล์ config อื่นๆ ที่จำเป็นตอน Runtime
# COPY --from=builder /app/.env.example /app/.env.example

# (Optional) Copy migration files ถ้าต้องการรัน migration จาก container นี้ (ไม่แนะนำสำหรับ Production)
# COPY --from=builder /app/internal/database/migrations /app/internal/database/migrations

# ระบุ Port ที่แอปพลิเคชันจะ Listen (ควรตรงกับ PORT ใน Config)
EXPOSE 8080

# (Optional but Recommended) Health Check Endpoint
# ตรวจสอบให้แน่ใจว่ามี Endpoint /healthz ในแอปพลิเคชัน
HEALTHCHECK --interval=15s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/healthz || exit 1
# หมายเหตุ: Alpine ไม่มี curl/wget โดย default อาจต้องติดตั้ง (apk add --no-cache curl) หรือใช้ tool อื่น

# คำสั่งสำหรับรันแอปพลิเคชันเมื่อ Container เริ่มทำงาน
ENTRYPOINT ["/app/server"]