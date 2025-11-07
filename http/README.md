# API Test Files

이 디렉토리는 VSCode REST Client 확장을 사용한 API 테스트 파일을 포함합니다.

## REST Client 설치

VSCode에서 "REST Client" 확장을 설치하세요:
- Extension ID: `humao.rest-client`
- [마켓플레이스 링크](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)

## 사용 방법

1. `.http` 파일을 엽니다
2. 요청 위의 `Send Request` 링크를 클릭하거나 `Ctrl+Alt+R` (Mac: `Cmd+Alt+R`)를 누릅니다
3. 결과가 새 탭에 표시됩니다

## 파일 구조

- `variables.http` - 환경 변수 설정 (먼저 확인하세요!)
- `auth.http` - 회원가입 및 로그인
- `users.http` - 사용자 관리 (CRUD)
- `oidc.http` - OIDC 인증 (Logto)

## 테스트 순서

1. 서버 실행: `go run main.go`
2. `variables.http`에서 환경 설정 확인
3. `auth.http`에서 회원가입 및 로그인
4. 로그인 후 받은 토큰을 `variables.http`의 `@token` 변수에 복사
5. `users.http`에서 인증이 필요한 API 테스트

## 주의사항

- 실제 비밀번호는 커밋하지 마세요
- 토큰은 24시간 후 만료됩니다
- 포인트 업데이트는 누적 방식입니다 (기존 값 + 요청 값)
