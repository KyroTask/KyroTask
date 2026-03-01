package models

type AnalyticsDashboardData struct {
	Tasks     TaskStats     `json:"tasks"`
	Habits    HabitStats    `json:"habits"`
	Pomodoros PomodoroStats `json:"pomodoros"`
	Goals     GoalStats     `json:"goals"`
}

type TaskStats struct {
	TotalCompleted int `json:"total_completed"`
	TotalPending   int `json:"total_pending"`
}

type HabitStats struct {
	ActiveHabits  int `json:"active_habits"`
	TotalLogs     int `json:"total_logs"`
	HighestStreak int `json:"highest_streak"`
}

type PomodoroStats struct {
	TotalSessionsCompleted int    `json:"total_sessions_completed"`
	TotalFocusMinutes      int    `json:"total_focus_minutes"`
	CurrentLevel           int    `json:"current_level"`
	CurrentPhase           string `json:"current_phase"`
}

type GoalStats struct {
	ActiveGoals    int `json:"active_goals"`
	CompletedGoals int `json:"completed_goals"`
}
