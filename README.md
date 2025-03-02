# 개발 히스토리

## 1. 초기 설정
- main.go 파일 생성
- fiber/v3, gorm, sqlite3 의존성 패키지 설치
- 범용적으로 사용되는 gorm 드라이버 사용

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
