# 개발 히스토리

## 1. 초기 설정
- main.go 파일 생성
- fiber/v3, gorm, sqlite3 의존성 패키지 설치

## 2. 로그인 DB 모델 구현

### SQLite Users

| id  | create_at | updated_at | deleted_at | email  | password | role |
| --- | --------- | ---------- | ---------- | ------ | -------- | ---- |
| int | date      | date       | date       | string | string   | int  |

- **email**: 이메일  
- **password**: 비밀번호  
- **role**: 유저의 권한 설정