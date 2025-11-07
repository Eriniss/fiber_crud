# Fiber CRUD API

사용자 인증 및 관리를 위한 REST API 서버

## 기술 스택

- **웹 프레임워크**: Fiber v3
- **데이터베이스**: SQLite3 + GORM
- **인증**: JWT + bcrypt
- **OIDC**: Logto (준비 완료)

## 주요 기능

### 보안 기능
- ✅ bcrypt를 사용한 안전한 비밀번호 해싱
- ✅ JWT 기반 인증 시스템
- ✅ JWT 검증 미들웨어
- ✅ 이메일 형식 검증
- ✅ 비밀번호 강도 검증 (최소 8자, 대소문자, 숫자, 특수문자 포함)
- ✅ 환경변수를 통한 보안 설정 관리
- ✅ CORS 설정

### 사용자 관리
- ✅ 회원가입 및 로그인
- ✅ 사용자 조회 (전체/개별)
- ✅ 사용자 정보 수정
- ✅ 사용자 삭제 (Soft Delete)
- ✅ 이메일 중복 검증
- ✅ 이메일 변경 방지 (보안)

### OIDC 준비 (Logto)
- ✅ OIDC 로그인 리다이렉트 엔드포인트
- ✅ OIDC 콜백 핸들러 구조
- ✅ OIDC 로그아웃 엔드포인트

## 설치 및 실행

### 1. 환경 설정

`.env.example` 파일을 `.env`로 복사하고 필요한 값을 설정하세요:

```bash
cp .env.example .env
```

필수 환경변수:
- `API_PORT`: 서버 포트 (기본값: 8080)
- `DATABASE_PATH`: SQLite 데이터베이스 파일 경로
- `JWT_SECRET_KEY`: JWT 서명용 비밀 키 (강력한 랜덤 문자열 사용)
- `ALLOWED_ORIGINS`: CORS 허용 오리진

선택 환경변수 (Logto OIDC 사용 시):
- `LOGTO_ENDPOINT`
- `LOGTO_APP_ID`
- `LOGTO_APP_SECRET`
- `LOGTO_REDIRECT_URI`

### 2. 의존성 설치

```bash
go mod download
```

### 3. 서버 실행

```bash
go run main.go
```

## 데이터베이스 구조

### Users 테이블

| 필드 | 타입 | 제약조건 | 설명 |
|------|------|----------|------|
| `id` | INTEGER | PRIMARY KEY, AUTO_INCREMENT | 사용자 고유 식별자 |
| `created_at` | DATETIME | NOT NULL | 계정 생성 시간 |
| `updated_at` | DATETIME | NOT NULL | 마지막 수정 시간 |
| `deleted_at` | DATETIME | NULL | 삭제 시간 (Soft Delete) |
| `email` | TEXT | UNIQUE, NOT NULL | 사용자 이메일 (로그인 ID, 변경 불가) |
| `password` | TEXT | NOT NULL | bcrypt 해시된 비밀번호 |
| `name` | TEXT | | 사용자 이름 |
| `group` | TEXT | | 사용자 그룹 (admin, user) |
| `gender` | TEXT | | 성별 (male, female) |
| `point` | INTEGER | DEFAULT 0 | 사용자 포인트 |

### User 모델 (Go)

```go
type User struct {
    gorm.Model           // ID, CreatedAt, UpdatedAt, DeletedAt 포함
    Email    string     `gorm:"unique" json:"email"`  // 고유 이메일
    Password string     `json:"password"`             // bcrypt 해시
    Name     string     `json:"name"`                 // 이름
    Group    string     `json:"group"`                // admin, user
    Gender   string     `json:"gender"`               // male, female
    Point    int        `json:"point"`                // 포인트
}
```

### 필드 상세 설명

#### `email` (이메일)
- 고유 식별자로 사용되며, 중복될 수 없습니다
- RFC 5322 표준에 따른 형식 검증이 수행됩니다
- 보안상 이유로 생성 후 변경이 불가능합니다
- 예시: `user@example.com`

#### `password` (비밀번호)
- bcrypt 알고리즘 (cost: 10)으로 해시되어 저장됩니다
- 평문 비밀번호는 절대 저장되지 않습니다
- 강도 검증 규칙:
  - 최소 8자 이상
  - 대문자 1개 이상 포함
  - 소문자 1개 이상 포함
  - 숫자 1개 이상 포함
  - 특수문자 1개 이상 포함
- API 응답에서 자동으로 제외됩니다

#### `name` (이름)
- 사용자의 표시 이름
- 자유롭게 변경 가능합니다

#### `group` (그룹)
- 사용자의 권한 그룹을 나타냅니다
- 가능한 값: `admin`, `user`
- 향후 RBAC (역할 기반 접근 제어)에 사용될 예정입니다

#### `gender` (성별)
- 사용자의 성별 정보
- 가능한 값: `male`, `female`
- 선택 사항입니다

#### `point` (포인트)
- 사용자의 점수 또는 포인트를 나타냅니다
- 기본값: 0
- 업데이트 시 증가/감소 연산이 가능합니다

#### Soft Delete
- `deleted_at` 필드가 NULL이 아닌 경우 삭제된 것으로 간주됩니다
- 실제 데이터는 데이터베이스에 남아있어 복구가 가능합니다
- GORM이 자동으로 삭제된 레코드를 쿼리에서 제외합니다

## API 문서

### 기본 정보

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **인증 방식**: Bearer Token (JWT)

---

## 인증 API

### 1. 회원가입

새로운 사용자 계정을 생성합니다.

**Endpoint**
```
POST /auth/user
```

**인증**: 불필요 (Public)

**요청 본문**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "name": "홍길동",
  "group": "user",
  "gender": "male",
  "point": 0
}
```

**필수 필드**
- `email`: 유효한 이메일 형식
- `password`: 비밀번호 강도 조건 충족 필요

**선택 필드**
- `name`, `group`, `gender`, `point`

**성공 응답** (200 OK)
```json
{
  "message": "User created",
  "user": {
    "ID": 1,
    "CreatedAt": "2025-11-07T10:30:00Z",
    "UpdatedAt": "2025-11-07T10:30:00Z",
    "DeletedAt": null,
    "email": "user@example.com",
    "password": "",
    "name": "홍길동",
    "group": "user",
    "gender": "male",
    "point": 0
  }
}
```

**에러 응답**

400 Bad Request - 잘못된 입력
```json
{
  "error": "invalid email format"
}
```
또는
```json
{
  "error": "password must be at least 8 characters long"
}
```

409 Conflict - 이메일 중복
```json
{
  "error": "Email already exists"
}
```

---

### 2. 로그인

사용자 인증 후 JWT 토큰을 발급받습니다.

**Endpoint**
```
POST /auth/sign_in
```

**인증**: 불필요 (Public)

**요청 본문**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**성공 응답** (200 OK)
```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "홍길동",
    "group": "user",
    "gender": "male",
    "point": 0
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**JWT 토큰 정보**
- 유효기간: 24시간
- 포함 정보: `id`, `email`, `name`, `point`, `exp`

**에러 응답**

404 Not Found - 사용자 없음
```json
{
  "error": "User not found"
}
```

401 Unauthorized - 비밀번호 불일치
```json
{
  "error": "Invalid password"
}
```

---

## 사용자 관리 API

**인증 필요**: 모든 요청에 `Authorization: Bearer <token>` 헤더 필요

### 3. 모든 사용자 조회

등록된 모든 사용자 목록을 조회합니다.

**Endpoint**
```
GET /auth/users
```

**요청 헤더**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**성공 응답** (200 OK)
```json
{
  "users": [
    {
      "ID": 1,
      "CreatedAt": "2025-11-07T10:30:00Z",
      "UpdatedAt": "2025-11-07T10:30:00Z",
      "DeletedAt": null,
      "email": "user1@example.com",
      "password": "",
      "name": "홍길동",
      "group": "user",
      "gender": "male",
      "point": 100
    },
    {
      "ID": 2,
      "CreatedAt": "2025-11-07T11:00:00Z",
      "UpdatedAt": "2025-11-07T11:00:00Z",
      "DeletedAt": null,
      "email": "admin@example.com",
      "password": "",
      "name": "관리자",
      "group": "admin",
      "gender": "female",
      "point": 500
    }
  ]
}
```

**에러 응답**

401 Unauthorized - 토큰 없음 또는 유효하지 않음
```json
{
  "error": "Missing authorization header"
}
```
또는
```json
{
  "error": "Invalid or expired token"
}
```

---

### 4. 특정 사용자 조회

ID로 특정 사용자의 정보를 조회합니다.

**Endpoint**
```
GET /auth/user/:id
```

**URL 파라미터**
- `id`: 사용자 ID (정수)

**요청 예시**
```
GET /auth/user/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**성공 응답** (200 OK)
```json
{
  "user": {
    "ID": 1,
    "CreatedAt": "2025-11-07T10:30:00Z",
    "UpdatedAt": "2025-11-07T10:30:00Z",
    "DeletedAt": null,
    "email": "user@example.com",
    "password": "",
    "name": "홍길동",
    "group": "user",
    "gender": "male",
    "point": 100
  }
}
```

**에러 응답**

404 Not Found - 사용자 없음
```json
{
  "error": "User not found"
}
```

---

### 5. 사용자 정보 수정

사용자의 정보를 업데이트합니다.

**Endpoint**
```
PUT /auth/user/:id
```

**URL 파라미터**
- `id`: 사용자 ID (정수)

**요청 본문** (일부 필드만 전송 가능)
```json
{
  "name": "김철수",
  "group": "admin",
  "gender": "male",
  "point": 50,
  "password": "NewSecurePass456!"
}
```

**필드 수정 규칙**
- `email`: 수정 불가 (403 에러 반환)
- `password`: 변경 시 비밀번호 강도 검증 수행
- `point`: 기존 값에서 더하거나 빼기 (누적)
- 나머지 필드: 제공된 값으로 덮어쓰기

**포인트 업데이트 예시**
```json
{
  "point": 50
}
```
- 기존 포인트: 100
- 요청 포인트: 50
- 결과 포인트: 150 (100 + 50)

음수 값으로 포인트 차감 가능:
```json
{
  "point": -30
}
```
- 기존 포인트: 150
- 요청 포인트: -30
- 결과 포인트: 120 (150 - 30)

**성공 응답** (200 OK)
```json
{
  "message": "User updated",
  "user": {
    "ID": 1,
    "CreatedAt": "2025-11-07T10:30:00Z",
    "UpdatedAt": "2025-11-07T12:00:00Z",
    "DeletedAt": null,
    "email": "user@example.com",
    "password": "",
    "name": "김철수",
    "group": "admin",
    "gender": "male",
    "point": 150
  }
}
```

**에러 응답**

403 Forbidden - 이메일 변경 시도
```json
{
  "error": "Email cannot be changed"
}
```

400 Bad Request - 비밀번호 강도 미달
```json
{
  "error": "password must contain at least one uppercase letter"
}
```

404 Not Found - 사용자 없음
```json
{
  "error": "User not found"
}
```

---

### 6. 사용자 삭제

사용자를 Soft Delete 방식으로 삭제합니다.

**Endpoint**
```
DELETE /auth/user/:id
```

**URL 파라미터**
- `id`: 사용자 ID (정수)

**요청 예시**
```
DELETE /auth/user/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**성공 응답** (200 OK)
```json
{
  "message": "User deleted"
}
```

**에러 응답**

404 Not Found - 사용자 없음 또는 이미 삭제됨
```json
{
  "error": "User not found"
}
```

**참고**:
- Soft Delete 방식이므로 실제로 데이터가 삭제되지 않고 `deleted_at` 필드에 삭제 시간이 기록됩니다
- 삭제된 사용자는 일반 조회에서 제외됩니다
- 데이터베이스에서 직접 복구가 가능합니다

---

## OIDC API (Logto)

Logto를 통한 소셜 로그인 기능입니다. (구현 준비 완료, 실제 연동은 진행 중)

### 7. OIDC 로그인 시작

**Endpoint**
```
GET /oidc/login
```

**인증**: 불필요

**성공 응답** (200 OK)
```json
{
  "auth_url": "https://your-logto-domain.logto.app/oidc/auth?client_id=...",
  "message": "Redirect to this URL for OIDC login"
}
```

**사용 방법**
1. 클라이언트는 응답받은 `auth_url`로 사용자를 리다이렉트
2. 사용자가 Logto에서 로그인
3. Logto가 `/oidc/callback`으로 콜백

---

### 8. OIDC 콜백

Logto 인증 후 리다이렉트되는 엔드포인트입니다.

**Endpoint**
```
GET /oidc/callback?code=<authorization_code>
```

**쿼리 파라미터**
- `code`: Logto에서 발급한 인증 코드

**성공 응답** (200 OK)
```json
{
  "message": "OIDC callback received",
  "code": "auth_code_here",
  "note": "This endpoint needs to be implemented with Logto SDK"
}
```

**참고**: 실제 구현 시 이 엔드포인트에서:
1. 인증 코드를 액세스 토큰으로 교환
2. 액세스 토큰으로 사용자 정보 조회
3. 로컬 DB에서 사용자 조회 또는 생성
4. JWT 토큰 발급 및 반환

---

### 9. OIDC 로그아웃

**Endpoint**
```
POST /oidc/logout
```

**성공 응답** (200 OK)
```json
{
  "message": "Logout successful",
  "note": "Client should remove the JWT token"
}
```

---

## 에러 코드 정리

| 상태 코드 | 설명 | 예시 |
|----------|------|------|
| 200 | 성공 | 모든 성공적인 요청 |
| 400 | 잘못된 요청 | 입력 검증 실패, 잘못된 형식 |
| 401 | 인증 실패 | 토큰 없음, 토큰 만료, 잘못된 비밀번호 |
| 403 | 권한 없음 | 이메일 변경 시도 |
| 404 | 리소스 없음 | 사용자를 찾을 수 없음 |
| 409 | 충돌 | 이메일 중복 |
| 500 | 서버 오류 | 내부 서버 에러 |

---

## 인증 흐름

### 일반 로그인 흐름

```
1. 클라이언트 → POST /auth/sign_in
   {email, password}

2. 서버 → 사용자 검증
   - 이메일로 사용자 조회
   - bcrypt로 비밀번호 검증

3. 서버 → JWT 토큰 생성
   - 24시간 유효기간
   - 사용자 정보 포함

4. 서버 → 클라이언트
   {user, token}

5. 클라이언트 → 이후 모든 요청
   Authorization: Bearer <token>
```

### OIDC 로그인 흐름 (예정)

```
1. 클라이언트 → GET /oidc/login

2. 서버 → Logto 인증 URL 반환

3. 클라이언트 → Logto 로그인 페이지로 리다이렉트

4. 사용자 → Logto에서 로그인

5. Logto → GET /oidc/callback?code=xxx

6. 서버 → 코드를 토큰으로 교환
   → 사용자 정보 조회
   → 로컬 DB 생성/조회
   → JWT 발급

7. 서버 → 클라이언트
   {user, token}
```

---

## 보안 고려사항

### 구현된 보안 기능
1. **비밀번호 보안**: bcrypt (cost: 10) 사용
2. **JWT 보안**: 환경변수에서 시크릿 키 관리, 24시간 만료
3. **입력 검증**: 이메일 형식 및 비밀번호 강도 검증
4. **CORS**: 환경변수로 허용 오리진 관리
5. **필드 보호**: 이메일 변경 방지
6. **응답 보안**: 비밀번호 필드 응답에서 제거

### 권장사항
- 프로덕션 환경에서는 반드시 강력한 `JWT_SECRET_KEY` 사용
- HTTPS 사용 권장
- Rate limiting 구현 고려
- 로그 모니터링 구현 고려

---

## 프로젝트 구조

```
fiber_crud/
├── config/           # 애플리케이션 설정
├── database/         # DB 초기화 및 연결
├── handlers/
│   ├── user_handlers/  # 사용자 관련 핸들러
│   └── oidc_handlers/  # OIDC 관련 핸들러
├── middleware/       # JWT 인증 등 미들웨어
├── models/          # 데이터 모델
├── routes/          # API 라우팅
├── utils/           # 유틸리티 (JWT, 해싱, 검증)
├── main.go          # 애플리케이션 진입점
├── .env.example     # 환경변수 예제
└── README.md
```

## 개발 히스토리

### v2.0 (최신)
- ✅ 블로그 기능 제거 (사용자 관리에 집중)
- ✅ SHA-256에서 bcrypt로 마이그레이션
- ✅ JWT Secret Key 환경변수화
- ✅ 입력 검증 로직 추가
- ✅ JWT 인증 미들웨어 구현
- ✅ CORS 설정 개선
- ✅ Logto OIDC 연동 준비

### v1.0
- 기본 CRUD 구현
- JWT 토큰 기반 로그인
- SQLite + GORM 연동

## Logto OIDC 통합 가이드

Logto를 통한 소셜 로그인 및 사용자 관리를 구현하려면 다음 문서를 참고하세요:

### 📚 통합 문서
- **[빠른 시작 가이드](docs/LOGTO_QUICKSTART.md)** - 5분 안에 시작하기
- **[통합 가이드](docs/LOGTO_INTEGRATION.md)** - 전체 통합 프로세스
- **[구현 예시](docs/LOGTO_IMPLEMENTATION_EXAMPLE.md)** - 상세한 코드 예시
- **[체크리스트](docs/LOGTO_CHECKLIST.md)** - 단계별 체크리스트

### 🎯 통합 요약

1. **Logto 콘솔 설정**
   - Application 생성 (Traditional Web)
   - Redirect URIs 설정
   - Roles 생성 (admin, user)

2. **환경 변수 설정**
   ```env
   LOGTO_ENDPOINT=https://beh25r.logto.app
   LOGTO_APP_ID=your-app-id
   LOGTO_APP_SECRET=your-app-secret
   LOGTO_REDIRECT_URI=http://localhost:8080/oidc/callback
   ```

3. **패키지 설치**
   ```bash
   go get github.com/coreos/go-oidc/v3/oidc
   go get golang.org/x/oauth2
   ```

4. **코드 구현**
   - User 모델 확장 (OIDC 필드 추가)
   - OIDC 핸들러 구현
   - Role 매핑 로직 구현

자세한 내용은 [docs/LOGTO_QUICKSTART.md](docs/LOGTO_QUICKSTART.md)를 참고하세요.

---

## 다음 개발 계획

- [ ] Logto SDK 통합 및 OIDC 완전 구현 → **[가이드 준비 완료](docs/LOGTO_INTEGRATION.md)**
- [ ] Rate limiting 미들웨어
- [ ] 로깅 시스템
- [ ] API 문서 자동화 (Swagger)
- [ ] 단위 테스트 작성
- [ ] 역할 기반 권한 제어 (RBAC)

## 라이선스

MIT
