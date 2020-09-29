# Cloudflare-Bypass-Script
Node.js를 통한 Hcaptcha 챌린지 해결과 크롬을 이용한 Cloudflare Javascript 계산을 통한 HTTP 공격

## Setup
- 알맞은 chromedriver 설치
- hcaptcha-solver 패키지 사용
- Golang으로 컴파일 후 바이너리로 사용

## Warning
- 현재 Cloudflare의 패치로 작동 안함
- Cloudflare는 여러번의 Javascript challenge 보완으로 인해 Bot 감지 성능을 높임
- 이에 대한 추후 패치를 할 생각이 
