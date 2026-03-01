# Telegram Notification System Specification

This document outlines the notification system for the Telegram Task Manager. It covers notification types, workflows, and the UI design for Telegram messages.

## 1. Notification Types

### 📋 Task Notifications

| Type              | Trigger                          | Content                             |
| :---------------- | :------------------------------- | :---------------------------------- |
| **Due Soon**      | 1 hour before `due_date`         | Reminder of upcoming task deadline. |
| **Overdue**       | At `due_date` (if not completed) | Alert that a task is now overdue.   |
| **Daily Summary** | Every morning (e.g., 8:00 AM)    | List of tasks due today.            |

### 🔄 Habit Notifications

| Type                 | Trigger                                  | Content                                   |
| :------------------- | :--------------------------------------- | :---------------------------------------- |
| **Daily Reminder**   | User-set time or default (e.g., 9:00 PM) | Reminder to log daily habits.             |
| **Weekly Review**    | Sunday evening                           | Summary of habit completion for the week. |
| **Streak Milestone** | When streak reaches 7, 30, 100 days      | Celebration of habit consistency.         |

### 🎯 Goal Notifications

| Type                  | Trigger                     | Content                                  |
| :-------------------- | :-------------------------- | :--------------------------------------- |
| **Progress Update**   | When progress % increases   | Encouragement based on goal advancement. |
| **Target Date Alert** | 1 week before `target_date` | Reminder of goal deadline.               |
| **Goal Completed**    | When progress reaches 100%  | Celebration of goal achievement.         |

---

## 2. Notification Workflows

### 📤 Sending Workflow

1. **Event Triggered**: An event occurs in the system (e.g., task created, time reached).
2. **Notification Queued**: The system checks if the user has a linked `telegram_id`.
3. **Template Selection**: The appropriate Telegram message template is selected.
4. **Message Sent**: The message is sent via the Telegram Bot API using the `telegram_id`.
5. **Log Created**: The notification is logged in the database for history/retry.

### 🗑️ Deletion/Clear Workflow

- **Automatic Removal**: When a user clicks the "Done" button on a notification (Task or Habit), the system will:
  1. Mark the item as completed in the database.
  2. **Delete the message** from the Telegram chat immediately using the `deleteMessage` API.
- **Auto-Sync**: If a task is completed via the web app, the system will attempt to find and delete any active Telegram notifications for that task.

---

## 3. Telegram UI Design ("Beast" Mode)

### � Task Notification

```text
🔔 *Task Reminder*
━━━━━━━━━━━━━━━━━━
📝 *Task:* [Task Title]
📅 *Due:* [Due Date/Time]
🏷️ *Project:* [Project Name]

[Open in App](https://app.url/tasks/[slug])
━━━━━━━━━━━━━━━━━━
```

**Buttons:**
`[ ✅ Done & Remove ]`

### 🌿 Daily Habit Check-in

```text
✨ *Daily Habit Check-in*
━━━━━━━━━━━━━━━━━━
Time to log your daily habits!

✅ *[Habit Name]* - [Streak: 5 days 🔥]

━━━━━━━━━━━━━━━━━━
```

**Buttons:**
`[ ✅ Done & Remove ]`

### 📅 Weekly Habit Review

```text
🗓️ *Weekly Habit Review*
━━━━━━━━━━━━━━━━━━
*Habit:* [Habit Name]

*Tasks for this week:*
• [Task 1]
• [Task 2]
• [Task 3]

━━━━━━━━━━━━━━━━━━
```

> [!NOTE]
> Tasks listed in the Weekly Review do **not** have individual "Done" buttons. Only the Habit itself can be marked as done.

**Buttons:**
`[ ✅ Weekly Done & Remove ]`

### 🏆 Goal Milestone

```text
🎉 *Goal Progress Update!*
━━━━━━━━━━━━━━━━━━
Your goal *[Goal Title]* is now:

▓▓▓▓▓▓▓▓░░░░░░░░ *75% Complete*

━━━━━━━━━━━━━━━━━━
```

**Buttons:**
`[ 🗑️ Dismiss ]`

---

## 4. Technical Requirements

### Backend (Go)

- **Telegram Bot API Integration**: Use `telebot` or direct HTTP calls to `https://api.telegram.org/bot<token>/sendMessage`.
- **Scheduler**: Implement a background worker (e.g., using `gocron` or a simple ticker) to check for upcoming deadlines.
- **Notification Model**:
  ```go
  type Notification struct {
      ID          uint
      UserID      uint
      Type        string // task_due, habit_reminder, etc.
      Payload     string // JSON data for the message
      SentAt      time.Time
      Status      string // sent, failed, read
      TelegramMsgID int // To allow updating/deleting messages
  }
  ```

### Frontend (Mini-App)

- **Settings Page**: Allow users to toggle specific notification types and set reminder times.
- **Deep Linking**: Ensure buttons in Telegram open the correct page in the Mini-App.
