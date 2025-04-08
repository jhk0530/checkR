# checkR 📰🐾

Reddit에서 선택한 서브레딧의 최신 글을 주기적으로 수집해, RSS 피드로 변환해주는 자동화 도구입니다.  
GitHub Actions를 통해 매시간 자동으로 실행되며, 결과물은 GitHub Pages를 통해 누구나 구독할 수 있는 RSS로 배포됩니다.

---

## ✨ 기능 소개

- 선택한 서브레딧들(rstats, rprogramming 등)의 최신 글 수집
- 최근 **2시간 이내**에 작성된 글만 필터링
- RSS 2.0 형식으로 `docs/{subreddit}.xml` 생성
- GitHub Actions + cron 으로 **1시간마다 자동 업데이트**
- GitHub Pages를 통해 퍼블릭 RSS 주소 제공

---

## 📦 사용 기술

- [Go](https://golang.org/)
- [go-reddit](https://github.com/vartanbeno/go-reddit) Reddit API wrapper
- [joho/godotenv](https://github.com/joho/godotenv) `.env` 파일 로딩
- GitHub Actions (자동화 및 배포)

---

## 🌐 RSS 피드 주소

```
/feed subscribe jhk0530.github.io/checkR/rstats.xml
/feed subscribe jhk0530.github.io/checkR/rprogramming.xml
/feed subscribe jhk0530.github.io/checkR/RStudio.xml
```
