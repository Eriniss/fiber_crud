# Logto 통합 체크리스트

## 📋 Logto 콘솔 설정

### Application 설정
- [ ] Logto 콘솔 접속 (https://cloud.logto.io/beh25r)
- [ ] Traditional Web Application 생성
- [ ] Application 이름: `fiber_crud_api`
- [ ] App ID 복사
- [ ] App Secret 복사
- [ ] Endpoint (issuer) URL 복사

### Redirect URIs 설정
- [ ] Redirect URI 추가: `http://localhost:8080/oidc/callback`
- [ ] Redirect URI 추가: `http://localhost:3000/callback` (프론트엔드용)
- [ ] Post Logout Redirect URI: `http://localhost:3000`
- [ ] Post Logout Redirect URI: `http://localhost:8080`

### CORS 설정
- [ ] Allowed Origin 추가: `http://localhost:3000`
- [ ] Allowed Origin 추가: `http://localhost:8080`

### Roles 설정
- [ ] Admin Role 생성
  - Name: `admin`
  - Description: `관리자 권한`
- [ ] User Role 생성
  - Name: `user`
  - Description: `일반 사용자 권한`

### Permissions 설정 (선택사항)
- [ ] Admin 권한 설정
  - `users:read`
  - `users:write`
  - `users:delete`
  - `all:manage`
- [ ] User 권한 설정
  - `profile:read`
  - `profile:write`

---

## 🔧 프로젝트 설정

### 환경 변수 설정
- [ ] `.env` 파일 열기
- [ ] `LOGTO_ENDPOINT` 추가
- [ ] `LOGTO_APP_ID` 추가
- [ ] `LOGTO_APP_SECRET` 추가
- [ ] `LOGTO_REDIRECT_URI` 추가
- [ ] `LOGTO_POST_LOGOUT_REDIRECT_URI` 추가

### Go 패키지 설치
- [ ] `go get github.com/coreos/go-oidc/v3/oidc`
- [ ] `go get golang.org/x/oauth2`
- [ ] `go mod tidy`

### 데이터베이스 마이그레이션
- [ ] User 모델에 OIDC 필드 추가
  - `OIDCSubject string`
  - `OIDCProvider string`
  - `IsOIDCUser bool`
- [ ] 데이터베이스 마이그레이션 실행
- [ ] 기존 사용자 데이터 백업 (선택사항)

---

## 💻 코드 구현

### OIDC 핸들러 구현
- [ ] `handlers/oidc_handlers/config.go` 생성
  - OIDC Provider 설정
  - OAuth2 Config 초기화
- [ ] `handlers/oidc_handlers/login.go` 수정
  - 실제 Authorization URL 생성
  - State 토큰 생성 및 저장
- [ ] `handlers/oidc_handlers/callback.go` 수정
  - 인증 코드를 토큰으로 교환
  - ID Token 검증
  - 사용자 정보 추출
  - 로컬 User 생성/업데이트
  - JWT 토큰 발급
- [ ] `handlers/oidc_handlers/logout.go` 구현
  - Logto 로그아웃 URL로 리다이렉트

### User 모델 확장
- [ ] `models/user.go` 수정
  - OIDC 필드 추가
  - JSON 태그 설정
  - GORM 인덱스 설정

### 유틸리티 함수
- [ ] `utils/oidc.go` 생성
  - Role 매핑 함수
  - 사용자 정보 추출 함수
  - State 토큰 생성/검증 함수

### 미들웨어 확장
- [ ] `middleware/auth.go` 수정
  - OIDC 토큰 검증 추가 (선택사항)

---

## 🧪 테스트

### 수동 테스트
- [ ] 서버 실행: `go run main.go`
- [ ] OIDC 로그인 URL 요청
- [ ] 브라우저에서 Logto 로그인
- [ ] 콜백 처리 확인
- [ ] JWT 토큰 발급 확인
- [ ] 토큰으로 API 호출 테스트

### Role 테스트
- [ ] Admin Role 사용자 생성 및 로그인
- [ ] User Role 사용자 생성 및 로그인
- [ ] Role별 권한 확인

### 에러 케이스 테스트
- [ ] 잘못된 인증 코드
- [ ] 만료된 토큰
- [ ] 네트워크 오류 시나리오

---

## 📚 문서화

- [ ] API 문서에 OIDC 로그인 플로우 추가
- [ ] README.md 업데이트
- [ ] 환경 변수 문서 업데이트
- [ ] 통합 가이드 작성 완료

---

## 🚀 배포 준비

### 프로덕션 설정
- [ ] Logto Application에 프로덕션 Redirect URI 추가
- [ ] HTTPS 인증서 설정
- [ ] 환경 변수를 프로덕션 환경에 설정
- [ ] CORS 설정 업데이트

### 보안 체크
- [ ] OIDC Secret 안전하게 저장
- [ ] State 토큰 검증 구현
- [ ] CSRF 방어 구현
- [ ] Rate Limiting 적용

---

## ✅ 완료 확인

모든 항목이 체크되면:
- [ ] 전체 통합 테스트 실행
- [ ] 문서 최종 검토
- [ ] 팀원에게 가이드 공유
- [ ] 프로덕션 배포

---

## 📞 도움이 필요한 경우

- Logto 공식 문서: https://docs.logto.io/
- Logto Discord: https://discord.gg/logto
- GitHub Issues: https://github.com/logto-io/logto
