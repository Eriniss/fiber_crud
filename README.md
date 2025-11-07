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

## API 엔드포인트

### 인증 불필요 (Public)

| Method | Endpoint | 설명 |
|--------|----------|------|
| POST | `/auth/user` | 회원가입 |
| POST | `/auth/sign_in` | 로그인 |

### 인증 필요 (Protected)

모든 요청에 `Authorization: Bearer <token>` 헤더 필요

| Method | Endpoint | 설명 |
|--------|----------|------|
| GET | `/auth/users` | 모든 사용자 조회 |
| GET | `/auth/user/:id` | 특정 사용자 조회 |
| PUT | `/auth/user/:id` | 사용자 정보 수정 |
| DELETE | `/auth/user/:id` | 사용자 삭제 |

### OIDC (Logto)

| Method | Endpoint | 설명 |
|--------|----------|------|
| GET | `/oidc/login` | OIDC 로그인 시작 |
| GET | `/oidc/callback` | OIDC 콜백 |
| POST | `/oidc/logout` | OIDC 로그아웃 |

## 데이터 모델

### User

```go
type User struct {
    ID        uint      // 자동 생성
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
    Email     string    // 고유, 변경 불가
    Password  string    // bcrypt 해시
    Name      string
    Group     string    // admin, user
    Gender    string    // male, female
    Point     int
}
```

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

## 다음 개발 계획

- [ ] Logto SDK 통합 및 OIDC 완전 구현
- [ ] Rate limiting 미들웨어
- [ ] 로깅 시스템
- [ ] API 문서 자동화 (Swagger)
- [ ] 단위 테스트 작성
- [ ] 역할 기반 권한 제어 (RBAC)

## 라이선스

MIT
