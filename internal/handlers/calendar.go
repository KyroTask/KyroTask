package handlers

import (
	"net/http"
	"time"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/gin-gonic/gin"
)

type CalendarHandler struct{}

func NewCalendarHandler() *CalendarHandler {
	return &CalendarHandler{}
}

type CalendarEvent struct {
	ID         uint       `json:"id"`
	Title      string     `json:"title"`
	Date       time.Time  `json:"date"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	Type       string     `json:"type"` // task, habit_log, milestone
	Status     string     `json:"status,omitempty"`
	Priority   string     `json:"priority,omitempty"`
	Color      string     `json:"color,omitempty"`
	Project    string     `json:"project,omitempty"`
	ResourceID uint       `json:"resource_id"`
}

func (h *CalendarHandler) Get(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end query parameters are required (YYYY-MM-DD)"})
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}
	end = end.Add(24*time.Hour - time.Nanosecond) // Include the entire end day

	var events []CalendarEvent

	// 1. Tasks: match by due_date, start_time, or created_at falling in range
	var tasks []models.Task
	database.DB.Preload("Project").Where("user_id = ?", userID).
		Where(
			"(due_date >= ? AND due_date <= ?) OR "+
				"(start_time >= ? AND start_time <= ?) OR "+
				"(due_date IS NULL AND start_time IS NULL AND created_at >= ? AND created_at <= ?) OR "+
				"is_daily = ?",
			start, end,
			start, end,
			start, end,
			true,
		).Find(&tasks)

	// Deduplicate: track which task IDs we've added
	addedTaskIDs := make(map[uint]bool)

	for _, t := range tasks {
		if addedTaskIDs[t.ID] {
			continue
		}
		addedTaskIDs[t.ID] = true

		// Determine the event date: prefer due_date, then start_time, then created_at
		var eventDate time.Time
		if t.DueDate != nil {
			eventDate = *t.DueDate
		} else if t.StartTime != nil {
			eventDate = *t.StartTime
		} else {
			eventDate = t.CreatedAt
		}

		projectName := ""
		if t.Project != nil {
			projectName = t.Project.Name
		}

		// For daily tasks, generate one event per day in the range
		if t.IsDaily {
			current := start
			for current.Before(end) {
				events = append(events, CalendarEvent{
					ID:         t.ID,
					Title:      t.Title,
					Date:       current,
					EndDate:    t.EndTime,
					Type:       "task",
					Status:     t.Status,
					Priority:   t.Priority,
					Project:    projectName,
					ResourceID: t.ID,
				})
				current = current.AddDate(0, 0, 1)
			}
		} else {
			events = append(events, CalendarEvent{
				ID:         t.ID,
				Title:      t.Title,
				Date:       eventDate,
				EndDate:    t.EndTime,
				Type:       "task",
				Status:     t.Status,
				Priority:   t.Priority,
				Project:    projectName,
				ResourceID: t.ID,
			})
		}
	}

	// 2. Habit logs in range
	var habitLogs []models.HabitLog
	database.DB.Where("user_id = ? AND log_date >= ? AND log_date <= ?", userID, start, end).Find(&habitLogs)

	// 3. Scheduled habits (active habits that should appear on days in range)
	var activeHabits []models.Habit
	database.DB.Where("user_id = ? AND is_active = ?", userID, true).Find(&activeHabits)

	// Build a set of already-logged habit+date combos to avoid duplicates
	loggedSet := make(map[string]bool)
	for _, hl := range habitLogs {
		key := time.Time(hl.LogDate).Format("2006-01-02")
		loggedSet[key+"-"+string(rune(hl.HabitID))] = true
	}

	today := time.Now().Truncate(24 * time.Hour)

	for _, habit := range activeHabits {
		current := start
		for current.Before(end) {
			shouldShow := false
			dayOfWeek := int(current.Weekday()) // 0 = Sunday

			if habit.Frequency == "daily" {
				shouldShow = true
			} else if habit.Frequency == "weekly" && habit.ScheduledDays != "" {
				for _, ch := range habit.ScheduledDays {
					if ch == ',' {
						continue
					}
					if int(ch-'0') == dayOfWeek {
						shouldShow = true
						break
					}
				}
			}

			// ONLY show the habit if it is scheduled for TODAY, and is NOT yet completed
			if shouldShow && current.Equal(today) {
				dateStr := current.Format("2006-01-02")

				// Check if already logged
				isLogged := false
				for _, hl := range habitLogs {
					if hl.HabitID == habit.ID && hl.LogDate.Format("2006-01-02") == dateStr {
						isLogged = true
						break
					}
				}

				if !isLogged {
					events = append(events, CalendarEvent{
						ID:         habit.ID,
						Title:      habit.Name,
						Date:       current,
						Type:       "habit",
						Status:     "pending",
						Color:      "#f43f5e", // rose
						ResourceID: habit.ID,
					})
				}
			}
			current = current.AddDate(0, 0, 1)
		}
	}

	// 4. Milestones with due_date in range
	var milestones []models.Milestone
	database.DB.Joins("JOIN goals ON goals.id = milestones.goal_id").
		Where("goals.user_id = ? AND milestones.due_date >= ? AND milestones.due_date <= ?", userID, start, end).
		Find(&milestones)
	for _, m := range milestones {
		if m.DueDate != nil {
			events = append(events, CalendarEvent{
				ID:         m.ID,
				Title:      m.Title,
				Date:       *m.DueDate,
				Type:       "milestone",
				Status:     m.Status,
				Color:      "#8b5cf6", // violet
				ResourceID: m.GoalID,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"start":  startStr,
		"end":    endStr,
	})
}
