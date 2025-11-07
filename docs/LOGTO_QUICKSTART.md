# Logto 빠른 시작 가이드

Logto 통합을 위한 가장 빠른 방법입니다.

## 🚀 5분 안에 시작하기

### Step 1: Logto 콘솔에서 정보 수집 (2분)

1. https://cloud.logto.io/beh25r 접속
2. **Applications** > **Create Application**
3. **Traditional Web** 선택
4. 다음 정보 복사:
   - App ID
   - App Secret
   - Endpoint (예: `https://beh25r.logto.app`)

### Step 2: Redirect URI 설정 (1분)

Logto Application 설정에서:
```
Redirect URIs:
- http://localhost:8080/oidc/callback

Post Logout Redirect URIs:
- http://localhost:3000
- http://localhost:8080

CORS Allowed Origins:
- http://localhost:3000
- http://localhost:8080
```

### Step 3: Roles 생성 (1분)

**User Management > Roles**에서:
1. `admin` Role 생성
2. `user` Role 생성

### Step 4: 환경 변수 설정 (1분)

`.env` 파일에 추가:
```env
LOGTO_ENDPOINT=https://beh25r.logto.app
LOGTO_APP_ID=your-app-id
LOGTO_APP_SECRET=your-app-secret
LOGTO_REDIRECT_URI=http://localhost:8080/oidc/callback
LOGTO_POST_LOGOUT_REDIRECT_URI=http://localhost:3000
```

### Step 5: 패키지 설치 (30초)

```bash
go get github.com/coreos/go-oidc/v3/oidc
go get golang.org/x/oauth2
go mod tidy
```

---

## ✅ 준비 완료!

이제 다음 문서를 참고하여 코드를 구현하세요:
- `LOGTO_IMPLEMENTATION_EXAMPLE.md` - 상세한 코드 예시
- `LOGTO_INTEGRATION.md` - 통합 가이드
- `LOGTO_CHECKLIST.md` - 체크리스트

---

## 🧪 첫 테스트

준비가 완료되면:

1. 서버 실행:
```bash
go run main.go
```

2. 브라우저에서 로그인 URL 얻기:
```bash
curl http://localhost:8080/oidc/login
```

3. 반환된 `auth_url`을 브라우저에서 열기

4. Logto에서 로그인

5. 콜백으로 리다이렉트되면 JWT 토큰 확인

---

## 📞 문제 해결

### OIDC Provider 초기화 실패
- `LOGTO_ENDPOINT`가 정확한지 확인
- 엔드포인트 끝에 `/` 없는지 확인
- 네트워크 연결 확인

### 콜백에서 에러
- Redirect URI가 Logto 콘솔에 정확히 등록되어 있는지 확인
- State 토큰 검증 확인
- 브라우저 콘솔에서 에러 메시지 확인

### Role이 올바르게 매핑되지 않음
- Logto에서 사용자에게 Role이 할당되어 있는지 확인
- `roles` scope가 요청에 포함되어 있는지 확인

---

## 🎯 다음 단계

1. Role 기반 권한 제어 구현
2. 프론트엔드 통합
3. Refresh Token 처리
4. 에러 처리 개선
5. 프로덕션 배포

자세한 내용은 다른 문서를 참고하세요!
