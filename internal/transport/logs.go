package transport

import (
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h Handler) getLogById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	log, err := h.logsService.GetById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, string(log.Data))
}
