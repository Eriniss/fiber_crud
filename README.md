# 개발 히스토리

## 1. 초기 설정
- main.go 파일 생성
- fiber/v3, gorm, sqlite3 의존성 패키지 설치

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
