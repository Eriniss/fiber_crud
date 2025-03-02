# 개발 계획

## 1. 사용 스택

- fiber/v3
- sqlite3
- gorm driver

## 2. 파일 구조

- multi-repo 형식 사용
- 데이터베이스의 콜렉션을 기준으로 모듈화
- ex: user, blog 등

## 3. 기능 명세

- 로그인
    - User 역할별 차등 기능 구현(0, 1, 2)
    - 비밀번호 변경 기능 구현
    - 이메일 인증 기능 구현
    - Google oauth 2.0을 구현한 소셜 로그인 구현

- 블로그
    - CRUD 기능 구현
    - user role(0, 1)일 경우 모든 포스트글 제어
    - user role(2)일 경우 자신이 작성한 포스트글만 제어


# 개발자 노트

## 1. 초기 설정
- main.go 파일 생성
- fiber/v3, gorm, sqlite3 의존성 패키지 설치
- 범용적으로 사용되는 gorm 드라이버 사용
- .env 파일을 사용하기 위해 godotenv 패키지 추가

## 2. 로그인(User) DB 모델 구현

### SQLite Users

| id  | create_at | updated_at | deleted_at | email  | password | role |
| --- | --------- | ---------- | ---------- | ------ | -------- | ---- |
| int | date      | date       | date       | string | string   | int  |

- **email**: 이메일  
- **password**: 비밀번호  
- **role**: 유저의 권한 설정

- gorm.Open() 메서드를 사용하여 "database.db" 파일을 오픈. 없을 시 자동 생성
- gorm.DB.AutoMigrate() 메서드를 사용하여 *User 모델을 매핑
- database 모듈의 InitDatabase() 함수가 초기 데이터베이스 생성

## 3. 핸들러(미들웨어) 작성

### 목표 handlers 기능

- 모든 User 조회(R)
- 단일 User 조회(R)
- 단일 User 생성(C)
- 단일 User 수정(U)
- 단일 User 삭제(D)

### handlers 패키지 생성

- handlers 내 user를 컨트롤 할 수 있는 미들웨어 작성
- fiber.Ctx가 기본적으로 사용되며 앞서 작성한 gorm.DB 패키지를 이용하여 sqlite3 데이터베이스 제어 가능
- fiber.Ctx 내 JSON 메서드를 사용하면 json 형태로 데이터 입/출력 가능
- 단일 User 생성 API 에서 email 중복 불가 기능 추가
- 단일 User 생성 API 에서 crypto를 이용하여 password의 sha-256 해시 저장 기능 추가
- 단일 User 생성 API 에서 salt 기능 추가. salt는 .env 파일에서 해당 값을 가져와 생성

## 4. 라우트 작성

### 생성한 핸들러를 API 엔드포인트와 연결

- routes 패키지 생성
- "/auth" API 그룹 생성
- 해당 라우트 내에 user handlers 연동하여 테스트
- API 테스트에 ThunderClient 사용
