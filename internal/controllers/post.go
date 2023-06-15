package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"otus-sonet/internal/models"
	"otus-sonet/internal/utils"
)

// PostCreate Обработка запроса на добавление нового поста
func PostCreate(c *gin.Context) {

	var post models.Post
	//парсим входные данные
	if err := c.ShouldBindJSON(&post); err != nil {
		c.AbortWithStatus(400)
		return
	}
	//проверяем что текст не пустой
	if post.Text == "" {
		c.AbortWithStatus(400)
		return
	}
	//берем id автора из авторизации
	userId, err := GetAuthUser(c)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}
	post.UserId = userId.Id
	//Сохроняем в БД
	postId, err := models.Posts.Create(post)
	if err != nil {
		utils.Code500(c, err.Error(), -5)
	}
	//Возвращаем успех
	c.JSON(200, gin.H{"id": postId})
}

// PostDelete Обработка запроса на удаление поста
func PostDelete(c *gin.Context) {
	//Получаем id поста из входных параметров
	postId := c.Param("id")
	if postId == "" {
		c.AbortWithStatus(400)
		return
	}
	//проверяем что валидные
	if _, err := uuid.Parse(postId); err != nil {
		c.AbortWithStatus(400)
		return
	}
	//берем id автора из авторизации
	userId, err := GetAuthUser(c)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}
	//Удаляем из БД
	if err := models.Posts.Delete(postId, userId.Id); err != nil {
		utils.Code500(c, err.Error(), -5)
		return
	}
	//Возвращаем успех
	c.AbortWithStatus(200)
}

// PostUpdate Обработка запроса обновления поста
func PostUpdate(c *gin.Context) {

	var post models.Post
	//парсим входные данные
	if err := c.ShouldBindJSON(&post); err != nil {
		c.AbortWithStatus(400)
		return
	}
	fmt.Println("post:", post)
	//проверяем что текст не пустой и id присутствует
	if post.Text == "" {
		c.AbortWithStatus(400)
		return
	}
	if post.Id == "" {
		c.AbortWithStatus(400)
		return
	}

	//проверяем что id валидный
	if _, err := uuid.Parse(post.Id); err != nil {
		c.AbortWithStatus(400)
		return
	}

	//берем id автора из авторизации
	userId, err := GetAuthUser(c)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}
	post.UserId = userId.Id
	//Сохроняем в БД
	err = models.Posts.Update(post)
	if err != nil {
		utils.Code500(c, err.Error(), -5)
	}
	//Возвращаем успех
	c.AbortWithStatus(200)
}

func PostDelete_(c *gin.Context) {

	id := c.Param("user_id")
	if id == "" {
		c.AbortWithStatus(400)
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		c.AbortWithStatus(400)
		return
	}

	userId, err := GetAuthUser(c)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	if err := models.Friends.Delete(id, userId.Id); err != nil {
		utils.Code500(c, err.Error(), -5)
		return
	}
	c.AbortWithStatus(200)
}

func PostGet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(400)
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		c.AbortWithStatus(400)
		return
	}

	post, err := models.Posts.Get(id)
	if err != nil {
		utils.Code500(c, err.Error(), -5)
		return
	}

	c.JSON(200, post)
}
