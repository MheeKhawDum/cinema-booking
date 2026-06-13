# 🎬 Cinema Ticket Booking System

โปรเจกต์นี้เป็น Take-Home Assignment สำหรับตำแหน่ง Full Stack Developer ครับ
เป็นระบบจองตั๋วหนังออนไลน์ที่รองรับการจองพร้อมกันหลายคนโดยไม่เกิด Double Booking

> ⚠️ นี่เป็นครั้งแรกที่ผมทำระบบที่มี Redis Lock, WebSocket และ Message Queue
> หลายส่วนอาจยังไม่ perfect แต่พยายามทำให้ครบและเข้าใจแต่ละส่วนจริงๆ ครับ

---

## 1. System Architecture Diagram

```
┌─────────────────────────────────────────────────────┐
│                   Browser (Vue 3)                   │
│         HTTP Request + WebSocket Connection         │
└──────────────┬────────────────────────┬─────────────┘
               │ REST API               │ WebSocket
               ▼                        ▼
┌──────────────────────────────────────────────────────┐
│              Backend API (Go + Gin)                  │
│   Auth │ Booking │ Movie │ Admin │ WebSocket Hub     │
└────┬───────┬──────────┬──────────┬────────────┬──────┘
     │       │          │          │            │
     ▼       ▼          ▼          ▼            ▼
  Google  MongoDB     Redis     RabbitMQ    WebSocket
  OAuth    (data)    (lock)    (queue)      Clients
                        │          │
                        │          ▼
                        │     Consumer Worker
                        │     (audit log + notify)
                        │
                   TTL 5 นาที
                   (auto expire)
```

ผมวาดเป็น ASCII เพราะยังไม่ได้ใช้ tool วาด diagram ครับ แต่ flow หลักๆ คือ

- Browser คุยกับ Backend ผ่าน REST API และ WebSocket
- Backend ใช้ Redis lock ที่นั่งก่อนจอง
- เมื่อจองสำเร็จ Backend ส่ง event ไป RabbitMQ
- Consumer รับ event ไปบันทึก log

---

## 2. Tech Stack Overview

| Layer | Technology | เหตุผลที่เลือก |
|-------|-----------|---------------|
| Backend | Go + Gin | โจทย์กำหนด + Goroutines จัดการ concurrent ได้ดี |
| Frontend | Vue 3 + Vite | โจทย์กำหนด + Reactivity เหมาะกับ real-time UI |
| Database | MongoDB | โจทย์กำหนด + Schema ยืดหยุ่นสำหรับ Audit Log |
| Cache/Lock | Redis | In-memory เร็วมาก + SETNX เป็น atomic operation |
| Realtime | WebSocket | Persistent connection ส่งข้อมูลได้ทั้งสองทาง |
| Queue | RabbitMQ | Async processing + message persistence |
| Auth | Google OAuth 2.0 + JWT | โจทย์กำหนด |
| Deploy | Docker + docker-compose | รันทั้งระบบด้วยคำสั่งเดียว |

---

## 3. Booking Flow อธิบายทีละ Step

```
User กดเลือกที่นั่ง A5
         │
         ▼
[1] Frontend ส่ง POST /api/bookings
    { showtime_id: "xxx", seat: "A5" }
         │
         ▼
[2] Backend เช็คก่อนว่าที่นั่งถูก BOOKED ใน MongoDB หรือยัง
    ถ้ามีแล้ว → return 409 Conflict
         │
         ▼
[3] Backend เรียก Redis SETNX
    key:   "lock:seat:{showtime_id}:A5"
    value: "{user_id}"
    TTL:   5 นาที
         │
    ┌────┴────┐
    │         │
  OK(1)    FAIL(0)
    │         │
    ▼         ▼
[4] สร้าง   return 409
 Booking   "seat is being
 PENDING   selected"
    │
    ▼
[5] Broadcast WebSocket
    { seat: "A5", status: "LOCKED" }
    → ทุก client ที่ดู showtime นี้เห็นทันที
         │
         ▼
[6] ตอบ Frontend กลับ
    { booking_id: "yyy", expires_at: "..." }
         │
         ▼
[7] User จ่ายเงิน (mock) ภายใน 5 นาที
    POST /api/bookings/{id}/payment
         │
    ┌────┴────────────┐
    │                 │
  จ่ายทัน         หมดเวลา
    │                 │
    ▼                 ▼
[8] MongoDB        MongoDB
  PENDING          PENDING
    →                →
  BOOKED          TIMEOUT
    │                 │
    ▼                 ▼
  DEL lock        DEL lock
  Broadcast       Broadcast
  "BOOKED"       "AVAILABLE"
    │
    ▼
[9] Publish event → RabbitMQ
    → Consumer บันทึก Audit Log
    → Mock notification
```

---

## 4. Redis Lock Strategy

### ทำไมต้องใช้ Redis Lock

ถ้าไม่มี Lock และมีคน 2 คนกดจองที่นั่งเดียวกันพร้อมกัน:

```
User A: check A5 → AVAILABLE  ← ตรงนี้มีช่องว่าง!
User B: check A5 → AVAILABLE  ← ทั้งคู่เห็น AVAILABLE
User A: INSERT booking A5 ✅
User B: INSERT booking A5 ✅  ← Double Booking!
```

Redis SETNX แก้ปัญหานี้เพราะเป็น **atomic operation** ทำใน step เดียว ไม่มีช่องว่าง

### วิธีที่ใช้

```
1. SETNX lock:seat:{showtime_id}:{seat} {user_id} EX 300
   - ถ้า key ยังไม่มี → set สำเร็จ → ได้ lock (return 1)
   - ถ้ามีอยู่แล้ว  → set ไม่ได้ → ไม่ได้ lock (return 0)

2. TTL 5 นาที (300 วินาที)
   - ถ้า user ไม่จ่ายเงินใน 5 นาที Redis expire อัตโนมัติ
   - ป้องกัน deadlock กรณี server crash

3. Release Lock ใช้ Lua Script
   - เช็คว่าเป็นเจ้าของ lock ก่อนค่อยลบ
   - ป้องกัน user A ลบ lock ของ user B โดยไม่ตั้งใจ
```

### Trade-offs ที่รู้

- ใช้ Redis single node ถ้า Redis crash lock จะหาย
- Production ควรใช้ Redlock algorithm (หลาย Redis node)
- แต่สำหรับโปรเจกต์นี้ single node เพียงพอครับ

---

## 5. Message Queue ใช้ทำอะไร

เลือกใช้ **RabbitMQ** เพราะ setup ง่ายกว่า Kafka และเพียงพอสำหรับ use case นี้

### Use Cases ที่ implement

**1. Booking Success Event**
```
Backend (Producer)
  → publish "booking.events" queue
  → { booking_id, user_id, seat, showtime_id }

Consumer Worker
  → รับ event
  → บันทึก Audit Log ลง MongoDB
  → Mock notification (print log)
```

**2. ทำไมไม่ทำตรงๆ ใน API handler**

ถ้าทำตรงๆ: Backend ต้องรอ log เสร็จ → response ช้า

ใช้ Queue: Backend ส่ง event แล้วตอบ user ทันที → Worker ทำ background
ทำให้ latency ของ API ต่ำกว่า

---

## 6. วิธีรันระบบ

### Requirements
- Docker Desktop
- ไฟล์ `.env` ที่ root folder (ดู `.env.example`)

### ขั้นตอน

```bash
# 1. clone project
git clone <repo-url>
cd cinema-booking

# 2. ตั้งค่า environment
cp .env.example .env
# แก้ค่าใน .env โดยเฉพาะ GOOGLE_CLIENT_ID และ GOOGLE_CLIENT_SECRET

# 3. รันทั้งระบบ
docker compose up --build

# รอจนเห็น log:
# ✅ MongoDB connected
# ✅ Redis connected
# 🚀 Server running on :8080
```

### URLs หลังรัน

| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| Backend API | http://localhost:8080 |
| Health Check | http://localhost:8080/health |
| RabbitMQ Dashboard | http://localhost:15672 (rabbit/rabbitpass) |

### ทำให้ตัวเองเป็น Admin

```bash
docker exec -it cinema-booking-mongo-1 \
  mongosh "mongodb://admin:secret@localhost:27017" \
  --authenticationDatabase admin \
  --eval 'use cinema; db.users.updateOne({email:"your@email.com"},{$set:{role:"ADMIN"}})'
```

---

## 7. Assumptions & Trade-offs

### สิ่งที่ตัดสินใจเองระหว่างทำ

**Payment เป็น Mock**
ไม่ได้ต่อ payment gateway จริง เพราะโจทย์ระบุว่าไม่คาดหวัง
กดปุ่มแล้วสำเร็จเลยครับ

**Seat Layout คงที่ 10x10**
กำหนดเป็น rows A-J, cols 1-10 รวม 100 ที่นั่งต่อรอบ
ยังไม่ได้ทำให้ configure ได้ตาม hall

**WebSocket ไม่มี Auth**
ตอนนี้ใครก็ connect WebSocket ได้ถ้ารู้ showtime_id
Production ควรเพิ่ม token ใน query param และ verify ก่อน

**Timeout Watcher ทุก 30 วินาที**
ตั้ง ticker ทุก 30 วินาที ซึ่งหมายความว่า booking อาจ timeout ช้าไป 30 วินาที
ยอมรับได้สำหรับ use case นี้ครับ

**RabbitMQ Retry 10 ครั้ง**
เพราะ Docker container start ไม่พร้อมกัน RabbitMQ บางครั้ง start ช้ากว่า Backend
เลยใส่ retry logic เพื่อป้องกัน crash ตอนเริ่มต้น

### สิ่งที่อยากทำเพิ่มถ้ามีเวลา

- [ ] Unit tests สำหรับ Redis Lock logic
- [ ] WebSocket authentication
- [ ] Redlock สำหรับ production
- [ ] Rate limiting ป้องกัน API spam
- [ ] Email notification จริง (ตอนนี้แค่ print log)