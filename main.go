package main

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var ArrowWords = []string{"java", "python", "ruby"}

func CreateComment(c *gin.Context) (string, string) {
	// Извлекаем уникальный идентификатор из HTTP-заголовка
	uniqueID := c.GetHeader("X-Unique-ID")

	// Извлекаем текст комментария из тела запроса
	var request struct {
		CommentText string `json:"commentText"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return "", ""
	}

	// Приводим весь текст комментария к нижнему регистру
	commentText := strings.ToLower(request.CommentText)

	return uniqueID, commentText

}

// Псевдокод для проверки комментария
func verifyComment(commentText string, arrowWords []string) bool {
	// Создаем регулярное выражение, объединяя запрещенные слова через вертикальную черту |
	re := regexp.MustCompile(strings.Join(arrowWords, "|"))

	// Ищем все вхождения запрещенных слов в тексте комментария
	matches := re.FindAllString(commentText, -1)

	// Если найдены совпадения (запрещенные слова), возвращаем false
	if len(matches) > 0 {
		return false
	}

	// В противном случае, комментарий прошел проверку
	return true
}

func answer(c *gin.Context, uniqueID string, verified bool) {
	if verified {
		// Если комментарий прошел проверку, отправляем статус 200 и uniqueID
		c.JSON(http.StatusOK, gin.H{"uniqueID": uniqueID, "message": "Comment verified"})
	} else {
		// Если комментарий не прошел проверку, отправляем статус 400 и uniqueID
		c.JSON(http.StatusBadRequest, gin.H{"uniqueID": uniqueID, "error": "Comment verification failed"})
	}
}

func main() {

	router := gin.Default()

	// Маршрут для обработки POST-запросов
	router.POST("/verify", func(c *gin.Context) {
		uniqueID, commentText := CreateComment(c)
		verified := verifyComment(commentText, ArrowWords)
		answer(c, uniqueID, verified)
	})

	router.Run(":8081") // Порт вашего сервиса верификации
}
