# Backend Architecture Flow

## 1. Request Lifecycle

How every HTTP request flows through the backend:

```mermaid
flowchart TD
    Client["🌐 Frontend / Telegram Bot"]
    
    subgraph "Gin Router"
        CORS["CORS Middleware"]
        GZIP["GZIP Compression"]
        Auth["Auth Middleware<br/>JWT Verification"]
        
        subgraph "Public Routes"
            Health["/health"]
            TGVerify["POST /auth/telegram/verify"]
            TGWidget["POST /auth/telegram/widget"]
            FBAuth["POST /auth/firebase"]
            Webhook["POST /telegram/webhook"]
        end
        
        subgraph "Protected Routes"
            Dashboard["GET /dashboard"]
            Calendar["GET /calendar"]
            Tasks["CRUD /tasks"]
            Projects["CRUD /projects"]
            Goals["CRUD /goals"]
            Habits["CRUD /habits"]
            Milestones["CRUD /milestones"]
            Activities["GET /activities"]
            Comments["CRUD /tasks/:id/comments"]
            Tags["CRUD /tags"]
            UserProfile["GET+PUT /user/me"]
        end
    end
    
    subgraph "Handlers Layer"
        AuthH["AuthHandler"]
        TelegramH["TelegramHandler"]
        DashboardH["DashboardHandler"]
        CalendarH["CalendarHandler"]
        TaskH["TaskHandler"]
        ProjectH["ProjectHandler"]
        GoalH["GoalHandler"]
        HabitH["HabitHandler"]
        MilestoneH["MilestoneHandler"]
        ActivityH["ActivityHandler"]
        CommentH["CommentHandler"]
        TagH["TagHandler"]
    end
    
    subgraph "Services Layer"
        TGService["TelegramService<br/>Bot API + JWT"]
        FBService["FirebaseAuthService<br/>Token Verification"]
        NotifService["NotificationService<br/>Telegram Messages"]
        GoalService["GoalService<br/>Progress Calculation"]
        ProjectService["ProjectService<br/>Progress Calculation"]
        HabitService["HabitService<br/>Streak Calculation"]
        Scheduler["Scheduler<br/>Background Jobs"]
    end
    
    subgraph "Data Layer"
        DB["SQLite / PostgreSQL<br/>via GORM"]
    end
    
    Client --> CORS --> GZIP --> Auth
    Auth --> Tasks & Projects & Goals & Habits & Milestones & Activities & Dashboard & Calendar & Comments & Tags & UserProfile
    
    Tasks --> TaskH --> NotifService & GoalService
    Projects --> ProjectH --> ProjectService
    Goals --> GoalH --> GoalService --> ProjectService
    Habits --> HabitH --> HabitService & NotifService
    Dashboard --> DashboardH
    Calendar --> CalendarH
    
    TaskH & ProjectH & GoalH & HabitH & MilestoneH & ActivityH & CommentH & TagH --> DB
    
    Scheduler --> NotifService --> TGService
```

## 2. Database Schema Relationships

```mermaid
erDiagram
    USER {
        uint id PK
        int64 telegram_id UK
        string firebase_uid UK
        string email UK
        string username
        string first_name
        string last_name
        string photo_url
    }
    
    PROJECT {
        uint id PK
        uint user_id FK
        string name
        string slug UK
        int progress
        string status
    }
    
    GOAL {
        uint id PK
        uint user_id FK
        uint project_id FK
        string title
        string slug UK
        int progress
        string status
        date target_date
    }
    
    MILESTONE {
        uint id PK
        uint goal_id FK
        string title
        string status
        date due_date
    }
    
    TASK {
        uint id PK
        uint user_id FK
        uint project_id FK
        uint goal_id FK
        uint milestone_id FK
        uint parent_id FK
        string title
        string slug UK
        string status
        string priority
        date due_date
    }
    
    HABIT {
        uint id PK
        uint user_id FK
        uint goal_id FK
        string name
        string frequency
        int current_streak
        int best_streak
    }
    
    HABIT_LOG {
        uint id PK
        uint habit_id FK
        uint user_id
        date log_date
        bool completed
    }
    
    COMMENT {
        uint id PK
        uint task_id FK
        uint user_id FK
        string content
    }
    
    TAG {
        uint id PK
        uint user_id FK
        string name
        string color
    }
    
    ACTIVITY {
        uint id PK
        uint user_id FK
        string resource_type
        uint resource_id
        string action
    }
    
    NOTIFICATION {
        uint id PK
        uint user_id FK
        string type
        uint related_id
        int telegram_message_id
        string status
    }
    
    USER ||--o{ PROJECT : "owns"
    USER ||--o{ GOAL : "owns"
    USER ||--o{ TASK : "owns"
    USER ||--o{ HABIT : "owns"
    USER ||--o{ ACTIVITY : "logs"
    USER ||--o{ NOTIFICATION : "receives"
    USER ||--o{ TAG : "owns"
    
    PROJECT ||--o{ GOAL : "contains"
    GOAL ||--o{ MILESTONE : "has"
    GOAL ||--o{ TASK : "linked to"
    GOAL ||--o{ HABIT : "linked to"
    MILESTONE ||--o{ TASK : "linked to"
    
    TASK ||--o{ TASK : "subtasks"
    TASK ||--o{ COMMENT : "has"
    TASK }o--o{ TAG : "tagged with"
    
    HABIT ||--o{ HABIT_LOG : "logged"
```

## 3. Auto-Progress Chain

When a Task is completed, a chain reaction updates everything:

```mermaid
flowchart LR
    A["✅ Task Completed"] --> B["GoalService<br/>Recalculate Progress"]
    B --> C["ProjectService<br/>Recalculate Progress"]
    B --> D{"Progress = 100%?"}
    D -->|Yes| E["🎉 Send Goal Completed<br/>Notification"]
    D -->|No| F["📊 Send Progress Update<br/>Notification"]
    
    G["✅ Habit Logged"] --> H["HabitService<br/>Calculate Streak"]
    H --> I{"Streak Milestone?"}
    I -->|7, 30, 100 days| J["🔥 Send Streak<br/>Notification"]
    I -->|No| K["Update DB Only"]
```

## 4. Scheduler Background Jobs

```mermaid
flowchart TD
    S["⏰ Scheduler<br/>Runs every 1 minute"]
    S --> T1["Check Tasks Due in 1 Hour"]
    S --> T2["Check Habits Not Logged<br/>at 9:00 PM"]
    S --> T3["Check Goals Deadline<br/>in 1 Week"]
    
    T1 --> N1["📋 Send Task Due Notification"]
    T2 --> N2["🌿 Send Habit Reminder"]
    T3 --> N3["🎯 Send Goal Deadline Alert"]
    
    N1 & N2 & N3 --> TG["Telegram Bot API<br/>sendMessage"]
```
