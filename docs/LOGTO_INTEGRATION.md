# Logto 통합 가이드

## 개요

이 문서는 Fiber CRUD API와 Logto OIDC를 통합하는 방법을 설명합니다.

## 사전 준비사항

### 1. Logto 콘솔 설정

#### Application 생성
1. Logto 콘솔 접속: https://cloud.logto.io/beh25r
2. **Applications** > **Create Application**
3. **Traditional Web** 선택
4. Application 이름: `fiber_crud_api`

#### Redirect URIs 설정
```
개발: http://localhost:8080/oidc/callback
프론트: http://localhost:3000/callback
```

#### CORS 설정
```
http://localhost:3000
http://localhost:8080
```

### 2. Role 설정

Logto 콘솔에서 두 개의 Role 생성:

**Admin Role:**
- Name: `admin`
- Description: 관리자 권한
- Permissions: users:*, all:manage

**User Role:**
- Name: `user`
- Description: 일반 사용자
- Permissions: profile:read, profile:write

### 3. 환경 변수 수집

Logto 콘솔의 Application Details에서 확인:
- App ID
- App Secret
- Endpoint (issuer)

## 통합 단계

### Step 1: 환경 변수 설정

`.env` 파일에 Logto 정보 추가:

```env
# 기존 설정
API_PORT=8080
DATABASE_PATH=./database.db
JWT_SECRET_KEY=your-secret-key
ALLOWED_ORIGINS=http://localhost:3000

# Logto OIDC 설정
LOGTO_ENDPOINT=https://beh25r.logto.app
LOGTO_APP_ID=<your-app-id>
LOGTO_APP_SECRET=<your-app-secret>
LOGTO_REDIRECT_URI=http://localhost:8080/oidc/callback
LOGTO_POST_LOGOUT_REDIRECT_URI=http://localhost:3000
```

### Step 2: Go 패키지 설치

```bash
# OIDC 클라이언트 라이브러리
go get github.com/coreos/go-oidc/v3/oidc
go get golang.org/x/oauth2
```

### Step 3: 코드 구현

다음 파일들을 수정/생성해야 합니다:

#### 1. `handlers/oidc_handlers/login.go` - 로그인 URL 생성
#### 2. `handlers/oidc_handlers/callback.go` - 인증 코드 처리
#### 3. `handlers/oidc_handlers/token.go` - 토큰 교환 및 검증
#### 4. `models/user.go` - OIDC Subject 필드 추가

### Step 4: User 모델 확장

OIDC 연동을 위해 User 모델에 필드 추가:

```go
type User struct {
    gorm.Model
    Email    string `gorm:"unique" json:"email"`
    Password string `json:"password"`
    Name     string `json:"name"`
    Group    string `json:"group"` // admin, user
    Gender   string `json:"gender"`
    Point    int    `json:"point"`

    // OIDC 필드 추가
    OIDCSubject  string `gorm:"index" json:"-"`        // Logto user ID
    OIDCProvider string `json:"-"`                      // "logto"
    IsOIDCUser   bool   `json:"is_oidc_user"`          // OIDC 사용자 여부
}
```

### Step 5: Role 매핑 전략

Logto Role → Local Group 매핑:

```go
// Logto에서 받은 role을 로컬 group으로 변환
func mapLogtoRoleToGroup(logtoRole string) string {
    roleMap := map[string]string{
        "admin": "admin",
        "user":  "user",
    }

    if group, ok := roleMap[logtoRole]; ok {
        return group
    }
    return "user" // 기본값
}
```

## 통합 후 테스트 순서

### 1. 로그인 플로우 테스트
```bash
# 1. 로그인 URL 얻기
GET http://localhost:8080/oidc/login

# 2. 브라우저에서 해당 URL 접속
# 3. Logto에서 로그인
# 4. 콜백으로 리다이렉트되며 JWT 토큰 받음
```

### 2. 역할 확인
```bash
# 로그인 후 자신의 정보 확인
GET http://localhost:8080/auth/user/:id
Authorization: Bearer <token>

# group 필드가 "admin" 또는 "user"인지 확인
```

### 3. 권한 테스트
- Admin: 모든 사용자 조회, 수정, 삭제 가능
- User: 자신의 정보만 조회, 수정 가능

## 통합 시 고려사항

### 1. 기존 사용자와의 호환성
- 기존 이메일/비밀번호 로그인 유지
- OIDC 로그인은 추가 옵션으로 제공
- 같은 이메일로 두 방식 병행 가능하도록 설계

### 2. Role 동기화
- Logto에서 Role이 변경되면 다음 로그인 시 반영
- 또는 주기적으로 동기화하는 백그라운드 작업 구현

### 3. 보안
- OIDC 토큰은 짧은 유효기간 (1시간)
- Refresh Token으로 자동 갱신
- JWT Secret과 OIDC Secret 분리 관리

### 4. 에러 처리
- Logto 서버 다운 시 fallback
- 네트워크 오류 처리
- 잘못된 토큰 처리

## 다음 단계

1. ✅ Logto 콘솔 설정 완료
2. ⬜ 환경 변수 설정
3. ⬜ Go OIDC 패키지 설치
4. ⬜ OIDC 핸들러 구현
5. ⬜ User 모델 확장
6. ⬜ Role 기반 권한 제어 구현
7. ⬜ 통합 테스트

## 참고 자료

- [Logto 공식 문서](https://docs.logto.io/)
- [go-oidc 라이브러리](https://github.com/coreos/go-oidc)
- [OAuth 2.0 스펙](https://oauth.net/2/)
- [OIDC 스펙](https://openid.net/connect/)
