# Gommunity
- Golang으로 만드는 커뮤니티 백엔드 서버 Repository

## 프로젝트 구성

- **`cmd/main.go`**: Entry Point
- **`internal/controller`**: API Endpoint & HTTP 요청 처리
- **`internal/dto`**: DTO
- **`internal/middleware`**: Custom Middleware (ex. JWT...)
- **`internal/model`**: Prisma 스키마 및 Prisma Client 관련 코드
- **`internal/repository`**: DB 접근 관련 코드. 데이터 읽기 & 쓰기
- **`internal/service`**: 비즈니스 로직
- **`pkg/crpyt`**: 암호화 라이브러리 (SHA256, JWT, Base64, ...)
- **`pkg/exception`**: 에러 처리 라이브러리 
- **`pkg/utils`**: 기타 유틸성 라이브러리 (param validator, uuid generator, ...)

## 사용 라이브러리
- [fiber](https://github.com/gofiber/fiber) 
- [validator](https://github.com/go-playground/validator) 
- [prisma-client-go](https://github.com/steebchen/prisma-client-go)
- [golang-jwt](https://github.com/golang-jwt/jwt)
- [air](https://github.com/air-verse/air)